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

// 注册
func (u *User) Register(ctx context.Context, userRegisterRequest *user.UserRegisterRequest,
	userRegisterResponse *user.UserRegisterResponse) error {
	userRegister := &model.User{
		ID:        0,
		UserName:  userRegisterRequest.UserName,
		Pwd:       userRegisterRequest.Pwd,
		FirstName: userRegisterRequest.FirstName,
	}
	_, err := u.UserDataService.AddUser(userRegister)
	if err != nil {
		return err
	}
	userRegisterResponse.Message = "OK"
	return nil
}

// 登陆
func (u *User) Login(ctx context.Context, userLogin *user.UserLoginRequest,
	loginResponse *user.UserLoginResponse) error {
	isOK, err := u.UserDataService.CheckPwd(userLogin.UserName, userLogin.Pwd)
	if err != nil {
		return err
	}
	loginResponse.IsSuccess = isOK
	return nil
}

// 查询用户信息
func (u *User) GetUserInfo(ctx context.Context, userInfoRequest *user.UserInfoRequest,
	userInfoResponse *user.UserInfoResponse) error {
	userInfo, err := u.UserDataService.FindUserByName(userInfoRequest.UserName)
	if err != nil {
		return err
	}
	userInfoResponse = UserForResponse(userInfo)
	return nil
}

func UserForResponse(userModel *model.User) *user.UserInfoResponse {
	response := &user.UserInfoResponse{}
	response.UserName = userModel.UserName
	response.FirstName = userModel.FirstName
	response.UserId = userModel.ID
	return response
}
