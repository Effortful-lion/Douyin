package http

import (
	"Douyin/app/gateway/rpc"
	"Douyin/idl/favorite/favoritePb"
	"Douyin/util/resp"
	"github.com/gin-gonic/gin"
	"strconv"
)

type FavoriteHandler struct{}

func NewFavoriteHandler() *FavoriteHandler {
	return &FavoriteHandler{}
}

func (f *FavoriteHandler) FavoriteActionHandler(ctx *gin.Context) {
	var req favoritePb.FavoriteActionRequest
	req.Token = ctx.Query("token")
	vid, _ := strconv.Atoi(ctx.Query("video_id"))
	req.VideoId = int64(vid)
	ActionType, _ := strconv.Atoi(ctx.Query("action_type"))
	req.ActionType = int32(ActionType)

	var res *favoritePb.FavoriteActionResponse

	//hystrix.ConfigureCommand("FavoriteAction", wrapper.FavoriteActionFuseConfig)
	//err := hystrix.Do("FavoriteAction", func() (err error) {
	//	res, err = rpc.FavoriteAction(ctx, &req)
	//	if err != nil {
	//		return err
	//	}
	//	return err
	//}, func(err error) error {
	//	return err
	//})

	res, err := rpc.FavoriteAction(ctx, &req)

	if err != nil {
		resp.ResponseError(ctx, resp.CodeError, err)
		return
	}

	resp.ResponseSuccess(ctx, res)
}

func (f *FavoriteHandler) FavoriteListHandler(ctx *gin.Context) {
	var req favoritePb.FavoriteListRequest
	uid, _ := strconv.Atoi(ctx.Query("user_id"))
	req.UserId = int64(uid)
	req.Token = ctx.Query("token")

	var res *favoritePb.FavoriteListResponse

	//hystrix.ConfigureCommand("FavoriteList", wrapper.FavoriteListFuseConfig)
	//err := hystrix.Do("FavoriteList", func() (err error) {
	//	res, err = rpc.FavoriteList(ctx, &req)
	//	if err != nil {
	//		return err
	//	}
	//	return err
	//}, func(err error) error {
	//	return err
	//})

	res, err := rpc.FavoriteList(ctx, &req)

	if err != nil {
		resp.ResponseError(ctx, resp.CodeError, err)
		return
	}

	resp.ResponseSuccess(ctx, res)
}
