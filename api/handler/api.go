package handler

import (
	"context"
	"encoding/json"
	log "github.com/micro/go-micro/v2/logger"

	"HuangQinTang/micro_shop/api/client"
	api "github.com/micro/go-micro/v2/api/proto"
	"github.com/micro/go-micro/v2/errors"
	api "path/to/service/proto/api"
)

type Api struct{}

func extractValue(pair *api.Pair) string {
	if pair == nil {
		return ""
	}
	if len(pair.Values) == 0 {
		return ""
	}
	return pair.Values[0]
}

// Api.Call is called by the API as /api/call with post body {"name": "foo"}
func (e *Api) Call(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Info("Received Api.Call request")

	// extract the client from the context
	apiClient, ok := client.ApiFromContext(ctx)
	if !ok {
		return errors.InternalServerError("go.micro.api.api.api.call", "api client not found")
	}

	// make request
	response, err := apiClient.Call(ctx, &api.Request{
		Name: extractValue(req.Post["name"]),
	})
	if err != nil {
		return errors.InternalServerError("go.micro.api.api.api.call", err.Error())
	}

	b, _ := json.Marshal(response)

	rsp.StatusCode = 200
	rsp.Body = string(b)

	return nil
}
