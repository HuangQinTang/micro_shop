package main

import (
	"context"
	"fmt"
	"github.com/HuangQinTang/micro_shop/common"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	microOpentracing "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	go_micro_service_product "github.com/micro_shop/product/proto/product"
	"github.com/opentracing/opentracing-go"
	"log"
	"testing"
)

const (
	productClientName = "go.micro.client.product"
)

func Test_Product(testing *testing.T) {
	//注册中心
	consulRegistor := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			fmt.Sprintf("%s:%d", consulHost, consulPort),
		}
	})
	//链路追踪
	t, io, err := common.NewTracer(productClientName, traceServ)
	if err != nil {
		log.Fatal(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	service := micro.NewService(
		micro.Name(productClientName),
		micro.Version("latest"),
		micro.Address(productServHost),
		//添加注册中心
		micro.Registry(consulRegistor),
		//绑定链路追踪
		micro.WrapClient(microOpentracing.NewClientWrapper(opentracing.GlobalTracer())),
	)

	productService := go_micro_service_product.NewProductService(productServName, service.Client())

	productAdd := &go_micro_service_product.ProductInfo{
		ProductName:        "imooc",
		ProductSku:         "cap",
		ProductPrice:       1.1,
		ProductDescription: "imooc-cap",
		ProductCategoryId:  1,
		ProductImage: []*go_micro_service_product.ProductImage{
			{
				ImageName: "cap-image",
				ImageCode: "capimage01",
				ImageUrl:  "capimage01",
			},
			{
				ImageName: "cap-image02",
				ImageCode: "capimage02",
				ImageUrl:  "capimage02",
			},
		},
		ProductSize: []*go_micro_service_product.ProductSize{
			{
				SizeName: "cap-size",
				SizeCode: "cap-size-code",
			},
		},
		ProductSeo: &go_micro_service_product.ProductSeo{
			SeoTitle:       "cap-seo",
			SeoKeywords:    "cap-seo",
			SeoDescription: "cap-seo",
			SeoCode:        "cap-seo",
		},
	}
	response, err := productService.AddProduct(context.TODO(), productAdd)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response)
}
