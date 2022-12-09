package main

import (
	"fmt"
	common "github.com/HuangQinTang/micro_shop_common"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	microOpentracing "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/micro_shop/api/handler"
	go_api "github.com/micro_shop/api/proto/api"
	go_micro_service_user "github.com/micro_shop/api/proto/user"
	"github.com/opentracing/opentracing-go"
)

const (
	apiServName = "go.micro.api"
	apiServHost = "0.0.0.0:8086"
	// consul
	consulHost = "127.0.0.1"
	consulPort = int64(8500)
	// jaeger
	traceServ = "127.0.0.1:6831"
)

func main() {
	//注册中心
	consulRegister := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			fmt.Sprintf("%s:%d", consulHost, consulPort),
		}
	})

	//链路追踪
	t, io, err := common.NewTracer(apiServName, traceServ)
	if err != nil {
		log.Fatal(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	service := micro.NewService(
		micro.Name(apiServName),
		micro.Version("latest"),
		// 暴露的服务地址
		micro.Address(apiServHost),
		// 注册中心
		micro.Registry(consulRegister),
		// 链路追踪
		micro.WrapClient(microOpentracing.NewClientWrapper(opentracing.GlobalTracer())),
		micro.Metadata(map[string]string{"protocol": "http"}),
	)
	service.Init()

	userServ := go_micro_service_user.NewUserService(common.UserServName, service.Client())

	if err = go_api.RegisterUserApiHandler(service.Server(), &handler.UserApi{UserService: userServ}); err != nil {
		log.Fatal(err)
	}

	if err = service.Run(); err != nil {
		log.Fatal(err)
	}
}
