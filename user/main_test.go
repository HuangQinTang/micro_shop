package main

import (
	"context"
	"fmt"
	"github.com/HuangQinTang/micro_shop/common"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	microOpentracing "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	go_micro_service_user "github.com/micro_shop/user/proto/user"
	"github.com/opentracing/opentracing-go"
	"testing"
)

func Test_User(testing *testing.T) {
	// 注册中心
	consul := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			fmt.Sprintf("%s:%d", consulHost, consulPort),
		}
	})

	userClientName := "go.micro.client.user"

	// 链路追踪
	t, io, err := common.NewTracer(userClientName, traceServ)
	if err != nil {
		panic(err.Error())
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	// 配置客户端
	service := micro.NewService(
		micro.Name(userClientName),
		micro.Version("latest"),
		micro.Address(userServHost),
		// 添加注册中心
		micro.Registry(consul),
		//绑定链路追踪
		micro.WrapClient(microOpentracing.NewClientWrapper(opentracing.GlobalTracer())),
	)
	//service.Init()

	userService := go_micro_service_user.NewUserService(userServName, service.Client())
	res, err := userService.GetUserInfo(context.TODO(), &go_micro_service_user.UserInfoReq{
		UserName: "admin",
	})
	if err != nil {
		fmt.Println("接口调不通", err.Error())
	}
	fmt.Printf("%#v\n", res)
	fmt.Println("---")
	fmt.Println(res)
}
