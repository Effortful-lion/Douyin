package rpc

import (
	"Douyin/idl/user/userPb"
	"context"
)

// 这里负责实际调用 user 模块的 rpc 方法

func UserRegister(ctx context.Context, req *userPb.UserRequest) (res *userPb.UserResponse, err error) {
	res, err = UserClient.Register(ctx, req)
	if err != nil {
		return
	}
	return
}

func UserLogin(ctx context.Context, req *userPb.UserRequest) (res *userPb.UserResponse, err error) {
	res, err = UserClient.Login(ctx, req)
	if err != nil {
		return
	}
	return
}

func UserInfo(ctx context.Context, req *userPb.UserInfoRequest) (res *userPb.UserInfoResponse, err error) {
	res, err = UserClient.UserInfo(ctx, req)
	if err != nil {
		return
	}
	return
}
