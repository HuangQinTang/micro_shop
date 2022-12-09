package handler

import (
	"context"
	"encoding/json"
	"errors"
	common "github.com/HuangQinTang/micro_shop_common"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro_shop/api/proto/api"
	userServ "github.com/micro_shop/api/proto/user"
)

type UserApi struct {
	UserService userServ.UserService
}

// Register userapi/register
func (u *UserApi) Register(ctx context.Context, req *go_api.Request, res *go_api.Response) error {
	log.Info("路径:【", req.Path, "】 请求方式：【", req.Method, "】 请求头:【", req.Header, "】 请求参数:【", req.Post, "】")

	userName, firstName, pwd := extractValue(req.Post["user_name"]), extractValue(req.Post["first_name"]), extractValue(req.Post["pwd"])
	if userName == "" {
		return errors.New("user_name 不能为空")
	}
	if pwd == "" {
		return errors.New("pwd 不能为空")
	}

	resp, err := u.UserService.Register(context.Background(), &userServ.UserRegisterReq{
		UserName:  userName,
		FirstName: firstName,
		Pwd:       pwd,
	})
	if err != nil {
		return err
	}
	respByte, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	res.StatusCode = 200
	res.TraceId = common.WithTrace(ctx)
	res.Body = string(respByte)
	return nil
}

func (u *UserApi) Login(ctx context.Context, req *go_api.Request, res *go_api.Response) error {
	log.Info("路径:【", req.Path, "】 请求方式：【", req.Method, "】 请求头:【", req.Header, "】 请求参数:【", req.Post, "】")

	userName, pwd := extractValue(req.Post["user_name"]), extractValue(req.Post["pwd"])
	if userName == "" {
		return errors.New("user_name 不能为空")
	}
	if pwd == "" {
		return errors.New("pwd 不能为空")
	}

	resp, err := u.UserService.Login(context.Background(), &userServ.UserLoginReq{
		UserName: userName,
		Pwd:      pwd,
	})
	if err != nil {
		return err
	}
	respByte, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	res.StatusCode = 200
	res.TraceId = common.WithTrace(ctx)
	res.Body = string(respByte)
	return nil
}

func (u *UserApi) GetUserInfo(ctx context.Context, req *go_api.Request, res *go_api.Response) error {
	log.Info("路径:【", req.Path, "】 请求方式：【", req.Method, "】 请求头:【", req.Header, "】 请求参数:【", req.Post, "】")

	userName := extractValue(req.Get["user_name"])
	if userName == "" {
		return errors.New("user_name 不能为空")
	}

	resp, err := u.UserService.GetUserInfo(context.Background(), &userServ.UserInfoReq{UserName: userName})
	if err != nil {
		return err
	}
	respByte, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	res.StatusCode = 200
	res.TraceId = common.WithTrace(ctx)
	res.Body = string(respByte)
	return nil
}

func extractValue(pair *go_api.Pair) string {
	if pair == nil {
		return ""
	}
	if len(pair.Values) == 0 {
		return ""
	}
	return pair.Values[0]
}