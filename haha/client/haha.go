package client

import (
	"context"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/server"
	haha "path/to/service/proto/haha"
)

type hahaKey struct {}

// FromContext retrieves the client from the Context
func HahaFromContext(ctx context.Context) (haha.HahaService, bool) {
	c, ok := ctx.Value(hahaKey{}).(haha.HahaService)
	return c, ok
}

// Client returns a wrapper for the HahaClient
func HahaWrapper(service micro.Service) server.HandlerWrapper {
	client := haha.NewHahaService("go.micro.service.template", service.Client())

	return func(fn server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			ctx = context.WithValue(ctx, hahaKey{}, client)
			return fn(ctx, req, rsp)
		}
	}
}
