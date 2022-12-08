package client

import (
	"context"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/server"
	api "path/to/service/proto/api"
)

type apiKey struct {}

// FromContext retrieves the client from the Context
func ApiFromContext(ctx context.Context) (api.ApiService, bool) {
	c, ok := ctx.Value(apiKey{}).(api.ApiService)
	return c, ok
}

// Client returns a wrapper for the ApiClient
func ApiWrapper(service micro.Service) server.HandlerWrapper {
	client := api.NewApiService("go.micro.service.template", service.Client())

	return func(fn server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			ctx = context.WithValue(ctx, apiKey{}, client)
			return fn(ctx, req, rsp)
		}
	}
}
