package http

import (
	"Douyin/app/gateway/rpc"
	"Douyin/idl/video/videoPb"
	"Douyin/util/resp"
	"bytes"
	"io"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PublishHandler struct{}

func NewPublishHandler() *PublishHandler {
	return &PublishHandler{}
}

func (p *PublishHandler) PublishListHandler(ctx *gin.Context) {
	var req videoPb.PublishListRequest
	uid, _ := strconv.Atoi(ctx.Query("user_id"))
	req.UserId = int64(uid)
	req.Token = ctx.Query("token")
	var res *videoPb.PublishListResponse
	//hystrix.ConfigureCommand("PublishList", wrapper.PublishListFuseConfig)
	//err := hystrix.Do("PublishList", func() (err error) {
	//	res, err = rpc.PublishList(ctx, &req)
	//	return err
	//}, func(err error) error {
	//	return err
	//})
	res, err := rpc.PublishList(ctx, &req)
	if err != nil {
		resp.ResponseError(ctx, resp.CodeInvalidParams, err)
		return
	}
	resp.ResponseSuccess(ctx, res)
}

func (p *PublishHandler) PublishHandler(ctx *gin.Context) {
	var req videoPb.PublishRequest
	req.Title = ctx.PostForm("title")
	req.Token = ctx.PostForm("token")
	//将获得的文件转为[]byte类型
	data, err := ctx.FormFile("data")
	if err != nil {
		resp.ResponseError(ctx, resp.CodeInvalidParams, err)
		return
	}
	file, err := data.Open()
	if err != nil {
		resp.ResponseError(ctx, resp.CodeInvalidParams, err)
		return
	}
	defer file.Close()
	// 使用缓冲区逐块读取文件内容并写入 req.Data
	var buffer bytes.Buffer
	_, err = io.Copy(&buffer, file)
	if err != nil {
		resp.ResponseError(ctx, resp.CodeInvalidParams, err)
		return
	}
	req.Data = buffer.Bytes()
	var res *videoPb.PublishResponse
	// hystrix.ConfigureCommand("Publish", wrapper.PublishFuseConfig)
	// err = hystrix.Do("Publish", func() (err error) {
	// 	res, err = rpc.Publish(ctx, &req)
	// 	return err
	// }, func(err error) error {
	// 	return err
	// })
	res, err = rpc.Publish(ctx, &req)
	if err != nil {
		resp.ResponseError(ctx, resp.CodeInvalidParams, err)
		return
	}
	resp.ResponseSuccess(ctx, res)
}
