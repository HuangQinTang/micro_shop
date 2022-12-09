package main

import (
	"fmt"
	common "github.com/HuangQinTang/micro_shop_common"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	ratelimit "github.com/micro/go-plugins/wrapper/ratelimiter/uber/v2"
	microOpentracing "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/micro_shop/cart/domain/repository"
	"github.com/micro_shop/cart/domain/service"
	"github.com/micro_shop/cart/handler"
	cart "github.com/micro_shop/cart/proto/cart"
	"github.com/opentracing/opentracing-go"
)

const (
	// consul
	consulHost   = "127.0.0.1"
	consulPort   = int64(8500)
	consulPrefix = "/micro/config"
	// rpc
	cartServHost = "127.0.0.1:8601"
	cartServName = common.CartServName
	// jaeger
	traceServ = "127.0.0.1:6831"
	QPS       = 100
)

func main() {
	//配置中心
	consulConfig, err := common.GetConsulConfig(consulHost, consulPort, consulPrefix)
	if err != nil {
		log.Error(err)
	}
	//注册中心
	consul := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			fmt.Sprintf("%s:%d", consulHost, consulPort),
		}
	})

	//链路追踪
	t, io, err := common.NewTracer(cartServName, traceServ)
	if err != nil {
		log.Error(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	//数据库连接
	mysqlInfo := common.GetMysqlFromConsul(consulConfig, "mysql")
	//创建数据库连接
	db, err := gorm.Open("mysql", mysqlInfo.User+":"+mysqlInfo.Pwd+"@/"+mysqlInfo.Database+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Error(err)
	}
	defer db.Close()
	//禁止副表
	db.SingularTable(true)

	//第一次初始化
	err = repository.NewCartRepository(db).InitTable()
	if err != nil {
		log.Error(err)
	}

	// New Service
	srv := micro.NewService(
		micro.Name(cartServName),
		micro.Version("latest"),
		//暴露的服务地址
		micro.Address(cartServHost),
		//注册中心
		micro.Registry(consul),
		//链路追踪
		micro.WrapHandler(microOpentracing.NewHandlerWrapper(opentracing.GlobalTracer())),
		//添加限流
		micro.WrapHandler(ratelimit.NewHandlerWrapper(QPS)),
	)

	// Initialise service
	srv.Init()

	cartDataService := service.NewCartDataService(repository.NewCartRepository(db))

	// Register Handler
	cart.RegisterCartHandler(srv.Server(), &handler.Cart{CartDataService: cartDataService})

	// Run service
	if err = srv.Run(); err != nil {
		log.Fatal(err)
	}
}
