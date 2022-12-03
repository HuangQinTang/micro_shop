package handler

import (
	"context"
	"user/domain/model"
	"user/domain/service"
	user "user/proto/user"
)

type User struct {
	UserDataService service.IUserDataService
}

//注册
func (u *User) Register(ctx context.Context, userRegisterRequest *user.UserRegisterReq, userRegisterResponse *user.UserRegisterRes) error {
	userRegister := &model.User{
		UserName:     userRegisterRequest.UserName,
		FirstName:    userRegisterRequest.FirstName,
		HashPassword: userRegisterRequest.Pwd,
	}
	_, err := u.UserDataService.AddUser(userRegister)
	if err != nil {
		return err
	}
	userRegisterResponse.Message = "添加成功"
	return nil
}

//登录
func (u *User) Login(ctx context.Context, userLogin *user.UserLoginReq, loginResponse *user.UserLoginRes) error {
	isOk, err := u.UserDataService.CheckPwd(userLogin.UserName, userLogin.Pwd)
	if err != nil {
		return err
	}
	loginResponse.IsSuccess = isOk
	return nil
}

//查询用户信息
func (u *User) GetUserInfo(ctx context.Context, userInfoRequest *user.UserInfoReq, userInfoResponse *user.UserInfoRes) error {
	userInfo, err := u.UserDataService.FindUserByName(userInfoRequest.UserName)
	if err != nil {
		return err
	}
	userInfoResponse = UserForResponse(userInfo)
	return nil
}

//类型转化
func UserForResponse(userModel *model.User) *user.UserInfoRes {
	response := &user.UserInfoRes{}
	response.UserName = userModel.UserName
	response.FirstName = userModel.FirstName
	response.UserId = userModel.ID
	return response
}