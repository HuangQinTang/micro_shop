package main

import (
	"fmt"
	"github.com/HuangQinTang/micro_shop/common"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	microOpentracing "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/opentracing/opentracing-go"
	"user/domain/repository"
	"user/domain/service"
	"user/handler"
	user "user/proto/user"
)

const (
	// consul
	consulHost   = "127.0.0.1"
	consulPort   = int64(8500)
	consulPrefix = "/micro/config"
	// rpc
	userServHost = "127.0.0.1:8600"
	userServName = "go.micro.srv.user"
	// jaeger
	traceServ = "127.0.0.1:6831"
)

func main() {
	// 配置中心
	consulConfig, err := common.GetConsulConfig(consulHost, consulPort, consulPrefix)
	if err != nil {
		panic("连接配置中心失败")
	}

	//获取数据库配置
	mysqlInfo := common.GetMysqlFromConsul(consulConfig, "mysql")
	// 连接数据库
	db, err := gorm.Open("mysql", mysqlInfo.User+":"+mysqlInfo.Pwd+"@/"+mysqlInfo.Database+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	db.SingularTable(true) //true严格匹配表面，默认false,复数映射表面

	// 注册中心
	consulRegister := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			fmt.Sprintf("%s:%d", consulHost, consulPort),
		}
	})

	// 链路追踪
	t, io, err := common.NewTracer(userServName, traceServ)
	if err != nil {
		log.Fatal(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t) //设置全局的tracer

	// 添加服务配置
	srv := micro.NewService(
		micro.Name(userServName),
		micro.Version("latest"),
		micro.Address(userServHost),
		// 添加注册中心
		micro.Registry(consulRegister),
		// 绑定链路追踪
		micro.WrapHandler(microOpentracing.NewHandlerWrapper(opentracing.GlobalTracer())),
	)
	// 服务初始化
	srv.Init()

	// 注入mysql连接进repository层（传统的dao层，封装mysql操作）
	rp := repository.NewUserRepository(db)
	rp.InitTable() //创建数据表，只執行一次

	// 创建服务实列（传统的service层，处理业务逻辑）
	userDataService := service.NewUserDataService(rp)
	// 将服务实列注册Handler（Handler相当于Controller）
	user.RegisterUserHandler(srv.Server(), &handler.User{UserDataService: userDataService})

	// Register Struct as Subscriber
	//micro.RegisterSubscriber(userServName, srv.Server(), new(subscriber.User))

	// Run srv
	if err = srv.Run(); err != nil {
		log.Fatal(err)
	}
}
