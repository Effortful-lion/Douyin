package http

import (
	"Douyin/app/gateway/rpc"
	"Douyin/idl/relation/relationPb"
	"Douyin/util/resp"
	"github.com/gin-gonic/gin"
	"strconv"
)

type RelationHandler struct{}

func NewRelationHandler() *RelationHandler {
	return &RelationHandler{}
}

func (r *RelationHandler) ListFriendRelationHandler(ctx *gin.Context) {
	var req relationPb.RelationFriendRequest
	uid, _ := strconv.Atoi(ctx.Query("user_id"))
	req.UserId = int64(uid)
	req.Token = ctx.Query("token")
	var res *relationPb.RelationFriendResponse
	//hystrix.ConfigureCommand("ListFriendRelation", wrapper.ListFriendRelationFuseConfig)
	//err := hystrix.Do("ListFriendRelation", func() (err error) {
	//	res, err = rpc.ListFriendRelation(ctx, &req)
	//	return err
	//}, func(err error) error {
	//	return err
	//})
	res, err := rpc.ListFriendRelation(ctx, &req)
	if err != nil {
		resp.ResponseError(ctx, resp.CodeError, err)
		return
	}
	resp.ResponseSuccess(ctx, res)
}

func (r *RelationHandler) ListFollowerRelationHandler(ctx *gin.Context) {
	var req relationPb.RelationFollowerRequest
	uid, _ := strconv.Atoi(ctx.Query("user_id"))
	req.UserId = int64(uid)
	req.Token = ctx.Query("token")
	var res *relationPb.RelationFollowerResponse
	//hystrix.ConfigureCommand("ListFollowerRelation", wrapper.ListFollowerRelationFuseConfig)
	//err := hystrix.Do("ListFollowerRelation", func() (err error) {
	//	res, err = rpc.ListFollowerRelation(ctx, &req)
	//	return err
	//}, func(err error) error {
	//	return err
	//})
	res, err := rpc.ListFollowerRelation(ctx, &req)
	if err != nil {
		resp.ResponseError(ctx, resp.CodeError, err)
		return
	}
	resp.ResponseSuccess(ctx, res)
}

func (r *RelationHandler) ListFollowRelationHandler(ctx *gin.Context) {
	var req relationPb.RelationFollowRequest
	uid, _ := strconv.Atoi(ctx.Query("user_id"))
	req.UserId = int64(uid)
	req.Token = ctx.Query("token")
	var res *relationPb.RelationFollowResponse
	//hystrix.ConfigureCommand("ListFollowRelation", wrapper.ListFollowRelationFuseConfig)
	//err := hystrix.Do("ListFollowRelation", func() (err error) {
	//	res, err = rpc.ListFollowRelation(ctx, &req)
	//	return err
	//}, func(err error) error {
	//	return err
	//})
	res, err := rpc.ListFollowRelation(ctx, &req)
	if err != nil {
		resp.ResponseError(ctx, resp.CodeError, err)
		return
	}
	resp.ResponseSuccess(ctx, res)
}

func (r *RelationHandler) ActionRelationHandler(ctx *gin.Context) {
	var req relationPb.RelationActionRequest
	req.Token = ctx.Query("token")
	toUserId, _ := strconv.Atoi(ctx.Query("to_user_id"))
	req.ToUserId = int64(toUserId)
	actionType, _ := strconv.Atoi(ctx.Query("action_type"))
	req.ActionType = int32(actionType)
	var res *relationPb.RelationActionResponse
	//hystrix.ConfigureCommand("ActionRelation", wrapper.ActionRelationFuseConfig)
	//err := hystrix.Do("ActionRelation", func() (err error) {
	//	res, err = rpc.ActionRelation(ctx, &req)
	//	return err
	//}, func(err error) error {
	//	return err
	//})
	res, err := rpc.ActionRelation(ctx, &req)
	if err != nil {
		resp.ResponseError(ctx, resp.CodeError, err)
		return
	}
	resp.ResponseSuccess(ctx, res)
}
