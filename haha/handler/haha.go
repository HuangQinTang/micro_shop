package handler

import (
	"context"
	"encoding/json"
	log "github.com/micro/go-micro/v2/logger"

	"haha/client"
	"github.com/micro/go-micro/v2/errors"
	api "github.com/micro/go-micro/v2/api/proto"
	haha "path/to/service/proto/haha"
)

type Haha struct{}

func extractValue(pair *api.Pair) string {
	if pair == nil {
		return ""
	}
	if len(pair.Values) == 0 {
		return ""
	}
	return pair.Values[0]
}

// Haha.Call is called by the API as /haha/call with post body {"name": "foo"}
func (e *Haha) Call(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Info("Received Haha.Call request")

	// extract the client from the context
	hahaClient, ok := client.HahaFromContext(ctx)
	if !ok {
		return errors.InternalServerError("go.micro.api.haha.haha.call", "haha client not found")
	}

	// make request
	response, err := hahaClient.Call(ctx, &haha.Request{
		Name: extractValue(req.Post["name"]),
	})
	if err != nil {
		return errors.InternalServerError("go.micro.api.haha.haha.call", err.Error())
	}

	b, _ := json.Marshal(response)

	rsp.StatusCode = 200
	rsp.Body = string(b)

	return nil
}
