package main

import (
	"fmt"
	"github.com/HuangQinTang/micro_shop/common"
	"github.com/jinzhu/gorm"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	ratelimit "github.com/micro/go-plugins/wrapper/ratelimiter/uber/v2"
	microOpentracing "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/micro_shop/category/domain/repository"
	"github.com/micro_shop/category/domain/service"
	"github.com/micro_shop/category/handler"
	"github.com/opentracing/opentracing-go"

	_ "github.com/jinzhu/gorm/dialects/mysql"

	category "github.com/micro_shop/category/proto/category"
)

const (
	QPS = 100
	// consul
	consulHost   = "127.0.0.1"
	consulPort   = int64(8500)
	consulPrefix = "/micro/config"
	// rpc
	categoryServHost = "127.0.0.1:8601"
	categoryServName = "go.micro.service.category"
	// jaeger
	traceServ = "127.0.0.1:6831"
)

func main() {
	// 配置中心
	consulConfig, err := common.GetConsulConfig(consulHost, consulPort, consulPrefix)
	if err != nil {
		log.Error(err)
	}
	// 获取mysql配置,路径中不带前缀
	mysqlInfo := common.GetMysqlFromConsul(consulConfig, "mysql")

	// 连接数据库
	db, err := gorm.Open("mysql", mysqlInfo.User+":"+mysqlInfo.Pwd+"@/"+mysqlInfo.Database+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Error(err)
	}
	defer db.Close()
	db.SingularTable(true)       // 禁止复表
	db.DB().SetMaxIdleConns(10)  // 空闲连接
	db.DB().SetMaxOpenConns(100) //最大连接数

	// 注册中心
	consulRegistry := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			fmt.Sprintf("%s:%d", consulHost, consulPort),
		}
	})
	// 链路追踪
	t, io, err := common.NewTracer(categoryServName, traceServ)
	if err != nil {
		log.Fatal(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t) //设置全局的tracer

	// New Service
	srv := micro.NewService(
		micro.Name(categoryServName),
		micro.Version("latest"),
		// 这里设置地址和需要暴露的端口
		micro.Address(categoryServHost),
		// 添加consul 作为注册中心
		micro.Registry(consulRegistry),
		// 绑定链路追踪
		micro.WrapHandler(microOpentracing.NewHandlerWrapper(opentracing.GlobalTracer())),
		//添加限流
		micro.WrapHandler(ratelimit.NewHandlerWrapper(QPS)),
	)
	srv.Init()

	rp := repository.NewCategoryRepository(db)
	rp.InitTable()

	categoryDataService := service.NewCategoryDataService(rp)
	err = category.RegisterCategoryHandler(srv.Server(), &handler.Category{CategoryDataService: categoryDataService})
	if err != nil {
		log.Error(err)
	}

	// Run service
	if err = srv.Run(); err != nil {
		log.Fatal(err)
	}
}
