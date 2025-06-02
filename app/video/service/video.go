package service

import (
	"Douyin/app/video/dao/mysql"
	"Douyin/consts"
	"Douyin/idl/video/videoPb"
	"Douyin/model"
	"Douyin/mq"
	"Douyin/util/resp"
	"context"
	"encoding/json"
)

type VideoService struct {
	videoPb.UnimplementedVideoServiceServer
}

func NewVideoService() *VideoService {
	return &VideoService{}
}

func (s *VideoService) Publish(c context.Context, req *videoPb.PublishRequest) (*videoPb.PublishResponse, error) {
	//加入消息队列
	body, _ := json.Marshal(&req)
	err := mq.SendMessage2MQ(body, consts.CreateVideoQueue)
	if err != nil {
		return &videoPb.PublishResponse{StatusCode: int32(resp.CodePublishFail), StatusMsg: resp.CodePublishFail.Msg()}, nil
	}
	return &videoPb.PublishResponse{StatusCode: int32(resp.CodePublishSuccess), StatusMsg: resp.CodePublishSuccess.Msg()}, nil
}

func (s *VideoService) PublishList(ctx context.Context, req *videoPb.PublishListRequest) (*videoPb.PublishListResponse, error) {
	token := req.Token
	uid := int(req.UserId)

	videos, err := mysql.NewVideoDao(ctx).GetVideoListByUserId(uid)
	if err != nil {
		return &videoPb.PublishListResponse{StatusCode: int32(resp.CodePublishListFail), StatusMsg: resp.CodePublishListFail.Msg()}, nil
	}
	var videoList []*videoPb.Video
	for _, video := range videos {
		videoList = append(videoList, BuildVideoPbModel(ctx, video, token))
	}
	return &videoPb.PublishListResponse{StatusCode: int32(resp.CodePublishListSuccess), StatusMsg: resp.CodePublishListSuccess.Msg(), VideoList: videoList}, nil
}

func BuildVideoPbModel(ctx context.Context, video *model.Video, token string) *videoPb.Video {
	Author, _ := mysql.NewVideoDao(ctx).FindUser(video)
	vid := int(video.ID)
	FavoriteCount, _ := mysql.NewVideoDao(ctx).GetFavoriteCount(vid)
	CommentCount, _ := mysql.NewVideoDao(ctx).GetCommentCount(vid)
	IsFavorite, _ := mysql.NewVideoDao(ctx).GetIsFavorite(vid, token)
	return &videoPb.Video{
		Id:            int64(vid),
		PlayUrl:       video.PlayUrl,
		CoverUrl:      video.CoverUrl,
		Title:         video.Title,
		FavoriteCount: FavoriteCount,
		CommentCount:  CommentCount,
		IsFavorite:    IsFavorite,
		Author:        BuildUserPbModel(ctx, Author, token),
	}
}

func BuildUserPbModel(ctx context.Context, user *model.User, token string) *videoPb.User {
	uid := int(user.ID)
	FollowCount, _ := mysql.NewVideoDao(ctx).GetFollowCount(uid)
	FollowerCount, _ := mysql.NewVideoDao(ctx).GetFollowerCount(uid)
	WorkCount, _ := mysql.NewVideoDao(ctx).GetWorkCount(uid)
	FavoriteCount, _ := mysql.NewVideoDao(ctx).GetFavoriteCount(uid)
	TotalFavorited, _ := mysql.NewVideoDao(ctx).GetTotalFavorited(uid)
	IsFollow, _ := mysql.NewVideoDao(ctx).GetIsFollowed(uid, token)
	return &videoPb.User{
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
