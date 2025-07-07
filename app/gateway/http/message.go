package http

import (
	"Douyin/app/gateway/rpc"
	"Douyin/idl/message/messagePb"
	"Douyin/util/resp"
	"github.com/gin-gonic/gin"
	"strconv"
)

type MessageHandler struct {
}

func NewMessageHandler() *MessageHandler {
	return &MessageHandler{}
}

func (m *MessageHandler) ActionMessageHandler(ctx *gin.Context) {
	var req messagePb.MessageActionRequest
	req.Token = ctx.Query("token")
	toUserId, _ := strconv.Atoi(ctx.Query("to_user_id"))
	req.ToUserId = int64(toUserId)
	actionType, _ := strconv.Atoi(ctx.Query("action_type"))
	req.ActionType = int32(actionType)
	req.Content = ctx.Query("content")
	var res *messagePb.MessageActionResponse
	//hystrix.ConfigureCommand("ActionMessage", wrapper.ActionMessageFuseConfig)
	//err := hystrix.Do("ActionMessage", func() (err error) {
	//	res, err = rpc.ActionMessage(ctx, &req)
	//	return err
	//}, func(err error) error {
	//	return err
	//})
	res, err := rpc.ActionMessage(ctx, &req)
	if err != nil {
		resp.ResponseError(ctx, resp.CodeError, err)
		return
	}
	resp.ResponseSuccess(ctx, res)
}

func (m *MessageHandler) ChatMessageHandler(ctx *gin.Context) {
	var req messagePb.MessageChatRequest
	req.Token = ctx.Query("token")
	toUserId, _ := strconv.Atoi(ctx.Query("to_user_id"))
	req.ToUserId = int64(toUserId)
	var res *messagePb.MessageChatResponse
	//hystrix.ConfigureCommand("ChatMessage", wrapper.ChatMessageFuseConfig)
	//err := hystrix.Do("ChatMessage", func() (err error) {
	//	res, err = rpc.ChatMessage(ctx, &req)
	//	return err
	//}, func(err error) error {
	//	return err
	//})
	res, err := rpc.ChatMessage(ctx, &req)
	if err != nil {
		resp.ResponseError(ctx, resp.CodeError, err)
		return
	}
	resp.ResponseSuccess(ctx, res)
}
