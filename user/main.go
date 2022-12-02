package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"user/domain/repository"
	"user/domain/service"
	"user/handler"
	user "user/proto/user"
)

func main() {
	// 服务参数设置
	srv := micro.NewService(
		micro.Name("go.micro.srv.user"),
		micro.Version("latest"),
	)

	// 服务初始化
	srv.Init()
	// 连接数据库
	db, err := gorm.Open("mysql", "root:root@/microshop?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	db.SingularTable(true) //true严格匹配表面，默认false,复数映射表面

	rp := repository.NewUserRepository(db) //注入mysql连接进repository层
	rp.InitTable()                         //创建数据表，只執行一次

	// 创建服务实列
	userDataService := service.NewUserDataService(rp)

	// 将服务注册Handler
	user.RegisterUserHandler(srv.Server(), &handler.User{UserDataService: userDataService})

	// Register Struct as Subscriber
	//micro.RegisterSubscriber("go.micro.srv.user", srv.Server(), new(subscriber.User))

	// Run srv
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
