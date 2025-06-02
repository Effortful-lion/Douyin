package http

import (
	"Douyin/app/gateway/rpc"
	"Douyin/idl/user/userPb"
	"Douyin/util/resp"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	// 依赖注入
	// userService *service.UserService
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (u *UserHandler) UserInfoHandler(ctx *gin.Context) {
	var req userPb.UserInfoRequest
	uid, _ := strconv.Atoi(ctx.Query("user_id"))
	req.UserId = int64(uid)
	req.Token = ctx.Query("token")
	var res *userPb.UserInfoResponse
	//hystrix.ConfigureCommand("UserInfo", wrapper.UserInfoFuseConfig)
	//err := hystrix.Do("UserInfo", func() (err error) {
	//	res, err = rpc.UserInfo(ctx, &req)
	//	return err
	//}, func(err error) error {
	//	return err
	//})
	//if err != nil {
	//	ctx.JSON(http.StatusInternalServerError, util.FailRequest(err.Error()))
	//	return
	//}
	res, err := rpc.UserInfo(ctx, &req)
	if err != nil {
		resp.ResponseError(ctx, resp.CodeError, err)
		return
	}
	resp.ResponseSuccess(ctx, res)
}

func (u *UserHandler) RegisterHandler(c *gin.Context) {
	// handler 负责：
	//1. 参数基本校验
	var req *userPb.UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.ResponseError(c, resp.CodeInvalidParams, err)
		return
	}
	var res *userPb.UserResponse
	//2. 调用 user 模块，进行数据查询和校验
	res, err := rpc.UserRegister(c, req)
	if err != nil {
		resp.ResponseError(c, resp.CodeError, err)
		return
	}
	//3. 组装响应数据
	resp.ResponseSuccess(c, res)
}

func (u *UserHandler) LoginHandler(c *gin.Context) {
	var req *userPb.UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.ResponseError(c, resp.CodeInvalidParams, err)
		return
	}
	var res *userPb.UserResponse
	res, err := rpc.UserLogin(c, req)
	if err != nil {
		resp.ResponseError(c, resp.CodeError, err)
		return
	}
	resp.ResponseSuccess(c, res)
}
