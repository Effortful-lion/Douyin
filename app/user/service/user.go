package service

import (
	"Douyin/app/user/dao/mysql"
	"Douyin/config"
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
	token := in.Token
	uid := int(in.UserId)

	user, _ := mysql.NewUserDao(ctx).FindUserById(uid)
	if user.ID == 0 {
		return &userPb.UserInfoResponse{
			StatusCode: 1,
			StatusMsg:  "用户不存在",
			User:       nil,
		}, nil
	}

	User := BuildUserPbModel(ctx, user, token)
	return &userPb.UserInfoResponse{
		StatusCode: 0,
		StatusMsg:  "获取用户信息成功",
		User:       User,
	}, nil
}

func (s *UserService) Register(ctx context.Context, in *userPb.UserRequest) (*userPb.UserResponse, error) {
	username := in.Username
	password := in.Password
	if len(username) > 32 || len(password) > 32 {
		return nil, errors.New("用户名或密码长度不合法")
	}
	// 数据库查询
	if user, _ := mysql.NewUserDao(ctx).FindUserByUserName(username); user.ID != 0 {
		return &userPb.UserResponse{
			StatusCode: 1,
			StatusMsg:  "用户名已存在",
		}, nil
	}

	// 数据库插入
	user := BuildUserModel(username, password)
	if id, err := mysql.NewUserDao(ctx).CreateUser(&user); err != nil {
		return &userPb.UserResponse{
			StatusCode: 1,
			StatusMsg:  "注册失败",
		}, err
	} else {
		token := util.GenerateToken(&user, 0)
		return &userPb.UserResponse{
			StatusCode: 0,
			StatusMsg:  "注册成功",
			UserId:     int64(id),
			Token:      token,
		}, nil
	}
}

func (s *UserService) Login(ctx context.Context, req *userPb.UserRequest) (*userPb.UserResponse, error) {
	user, err := mysql.NewUserDao(ctx).FindUserByUserName(req.Username)
	if err != nil {
		return &userPb.UserResponse{
			StatusCode: 1,
			StatusMsg:  "用户名或密码错误",
		}, err
	}
	if user.ID == 0 || util.Md5(req.Password) != user.Password {
		return &userPb.UserResponse{
			StatusCode: 1,
			StatusMsg:  "用户名或密码错误",
		}, err
	}
	token := util.GenerateToken(user, 0)
	uid := int64(user.ID)
	return &userPb.UserResponse{
		StatusCode: 0,
		StatusMsg:  "登录成功",
		UserId:     uid,
		Token:      token,
	}, nil
}

func BuildUserModel(username string, password string) model.User {
	config.InitConfig()
	avatar := config.Config.User.Avatar
	background := config.Config.User.Background
	signature := config.Config.User.Signature
	return model.User{
		Username:        username,
		Password:        util.Md5(password),
		Avatar:          avatar,
		BackgroundImage: background,
		Signature:       signature,
	}
}
func BuildUserPbModel(ctx context.Context, user *model.User, token string) *userPb.User {
	uid := int(user.ID)
	FollowCount, _ := mysql.NewUserDao(ctx).GetFollowCount(uid)
	FollowerCount, _ := mysql.NewUserDao(ctx).GetFollowerCount(uid)
	WorkCount, _ := mysql.NewUserDao(ctx).GetWorkCount(uid)
	FavoriteCount, _ := mysql.NewUserDao(ctx).GetFavoriteCount(uid)
	TotalFavorited, _ := mysql.NewUserDao(ctx).GetTotalFavorited(uid)
	IsFollow, _ := mysql.NewUserDao(ctx).GetIsFollowed(uid, token)
	return &userPb.User{
		Id:              int64(uid),
		Name:            user.Username,
		Avatar:          user.Avatar,
		BackgroundImage: user.BackgroundImage,
		Signature:       user.Signature,
		FollowCount:     FollowCount,
		FollowerCount:   FollowerCount,
		WorkCount:       WorkCount,
		FavoriteCount:   FavoriteCount,
		TotalFavorited:  TotalFavorited,
		IsFollow:        IsFollow,
	}
}
