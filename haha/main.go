package main

import (
	log "github.com/micro/go-micro/v2/logger"

	"github.com/micro/go-micro/v2"
	"haha/handler"
	"haha/client"

	haha "haha/proto/haha"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.api.haha"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init(
		// create wrap for the Haha service client
		micro.WrapHandler(client.HahaWrapper(service)),
	)

	// Register Handler
	haha.RegisterHahaHandler(service.Server(), new(handler.Haha))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
