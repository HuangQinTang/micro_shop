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
	"github.com/micro_shop/product/domain/repository"
	"github.com/micro_shop/product/domain/service"
	"github.com/micro_shop/product/handler"
	"github.com/opentracing/opentracing-go"

	product "github.com/micro_shop/product/proto/product"
)

const (
	QPS = 100
	// consul
	consulHost   = "127.0.0.1"
	consulPort   = int64(8500)
	consulPrefix = "/micro/config"
	// rpc
	productServHost = "127.0.0.1:8602"
	productServName = common.ProductServName
	// jaeger
	traceServ = "127.0.0.1:6831"
)

func main() {
	//配置中心
	consulConfig, err := common.GetConsulConfig(consulHost, consulPort, consulPrefix)
	if err != nil {
		log.Fatal(err)
	}

	//数据库设置
	mysqlInfo := common.GetMysqlFromConsul(consulConfig, "mysql")
	db, err := gorm.Open("mysql", mysqlInfo.User+":"+mysqlInfo.Pwd+"@/"+mysqlInfo.Database+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.SingularTable(true)

	rp := repository.NewProductRepository(db)
	rp.InitTable()               // 初始化数据表，只执行一次即可
	db.DB().SetMaxIdleConns(10)  // 空闲连接
	db.DB().SetMaxOpenConns(100) //最大连接数

	//注册中心
	consulRegister := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			fmt.Sprintf("%s:%d", consulHost, consulPort),
		}
	})

	//链路追踪
	t, io, err := common.NewTracer(productServName, traceServ)
	if err != nil {
		log.Fatal(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	// 商品dao层
	productDataService := service.NewProductDataService(rp)

	// 添加服务配置
	srv := micro.NewService(
		micro.Name(productServName),
		micro.Version("latest"),
		micro.Address(productServHost),
		//添加注册中心
		micro.Registry(consulRegister),
		//绑定链路追踪
		micro.WrapHandler(microOpentracing.NewHandlerWrapper(opentracing.GlobalTracer())),
		//添加限流
		micro.WrapHandler(ratelimit.NewHandlerWrapper(QPS)),
	)

	// Initialise service
	srv.Init()

	// Register Handler
	product.RegisterProductHandler(srv.Server(), &handler.Product{ProductDataService: productDataService})

	// Run service
	if err = srv.Run(); err != nil {
		log.Fatal(err)
	}
}
