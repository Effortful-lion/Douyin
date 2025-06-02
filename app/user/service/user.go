package service

import (
	"Douyin/idl/user/userPb"
	"Douyin/model"
	"Douyin/util"
	"context"
	"errors"
)

type UserService struct {
	// 注入（继承），实现真正的业务逻辑
	userPb.UnimplementedUserServiceServer
}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) UserInfo(ctx context.Context, in *userPb.UserInfoRequest) (*userPb.UserInfoResponse, error) {
	//token := in.Token
	//uid := int(in.UserId)

	//user, _ := dao.NewUserDao(ctx).FindUserById(uid)
	//if user.ID == 0 {
	//	UserInfoResponseData(res, 1, "用户不存在")
	//	return nil
	//}
	//
	//User := BuildUserPbModel(ctx, user, token)
	//UserInfoResponseData(res, 0, "获取用户信息成功", User)
	return &userPb.UserInfoResponse{
		StatusCode: 0,
		StatusMsg:  "获取用户信息成功",
		User: &userPb.User{
			Id:            1,
			Name:          "test",
			FollowCount:   10,
			FollowerCount: 10,
			IsFollow:      true,
		},
	}, nil
}

func (s *UserService) Register(ctx context.Context, in *userPb.UserRequest) (*userPb.UserResponse, error) {
	username := in.Username
	password := in.Password
	if len(username) > 32 || len(password) > 32 {
		return nil, errors.New("用户名或密码长度不合法")
	}
	// 数据库查询
	// if user, _ := dao.NewUserDao(ctx).FindUserByUserName(username); user.ID != 0 {
	// 	UserResponseData(res, 1, "用户名已存在")
	// 	return nil
	// }

	// 数据库插入
	// user := BuildUserModel(username, password)
	// if id, err := dao.NewUserDao(ctx).CreateUser(&user); err != nil {
	// 	UserResponseData(res, 1, "注册失败")
	// 	return err
	// } else {
	// 	token := util.GenerateToken(&user, 0)
	// 	return UserResponseData(res, 0, "注册成功", id, token)
	// }
	return &userPb.UserResponse{
		StatusCode: 0,
		StatusMsg:  "注册成功",
		UserId:     1,
		Token:      util.GenerateToken(&model.User{ID: 1}, 0),
	}, nil
}

func (s *UserService) Login(context.Context, *userPb.UserRequest) (*userPb.UserResponse, error) {
	//user, err := dao.NewUserDao(ctx).FindUserByUserName(req.Username)

	//if user.ID == 0 || util.Md5(req.Password) != user.Password {
	//	return nil, errors.New("用户名或密码错误")
	//}
	//token := util.GenerateToken(user, 0)
	//uid := int64(user.ID)
	return &userPb.UserResponse{
		StatusCode: 0,
		StatusMsg:  "登录成功",
		UserId:     1,
		Token:      util.GenerateToken(&model.User{ID: 1}, 0),
	}, nil
}
