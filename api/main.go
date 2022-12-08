package main

import (
	"HuangQinTang/micro_shop/api/proto/api"
	go_micro_service_cart "HuangQinTang/micro_shop/api/proto/cart"
	"fmt"
	"github.com/HuangQinTang/micro_shop/common"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	microOpentracing "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/opentracing/opentracing-go"
)

const (
	apiServName = "go.micro.api"
	apiServHost = "0.0.0.1:8608"
	// consul
	consulHost = "127.0.0.1"
	consulPort = int64(8500)
	// jaeger
	traceServ = "127.0.0.1:6831"
)

func main() {
	//注册中心
	consul := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			fmt.Sprintf("%s:%d", consulHost, consulPort),
		}
	})

	//链路追踪
	t, io, err := common.NewTracer(apiServName, traceServ)
	if err != nil {
		log.Error(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	service := micro.NewService(
		micro.Name(apiServName),
		micro.Version("latest"),
		// 暴露的服务地址
		micro.Address(apiServHost),
		// 注册中心
		micro.Registry(consul),
		// 链路追踪
		micro.WrapHandler(microOpentracing.NewHandlerWrapper(opentracing.GlobalTracer())),
	)
	service.Init()

	go_micro_service_cart.NewCartService("go.micro.service.cart", service.Client())

	api.RegisterApiHandler()
}
