package service

import (
	"Douyin/app/video/dao/mysql"
	"Douyin/app/video/dao/redis"
	"Douyin/consts"
	"Douyin/idl/video/videoPb"
	"Douyin/model"
	"Douyin/mq"
	"Douyin/util"
	"Douyin/util/resp"
	"context"
	"encoding/json"
	"fmt"
	redisv9 "github.com/redis/go-redis/v9"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"os"
	"time"
)

type VideoService struct {
	videoPb.UnimplementedVideoServiceServer
}

func NewVideoService() *VideoService {
	return &VideoService{}
}

func (s *VideoService) Publish(c context.Context, req *videoPb.PublishRequest) (*videoPb.PublishResponse, error) {
	// 将 包含 data 视频文件字节数据 的 req 请求数据 转码 为 json 格式
	body, _ := json.Marshal(&req)
	// 将 json 格式的 req 数据 发送到 mq 队列中
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

func BuildVideoModel(uid int, VideoUrl string, coverUrl string, title string) model.Video {
	return model.Video{
		AuthorID: uid,
		PlayUrl:  VideoUrl,
		CoverUrl: coverUrl,
		Title:    title,
	}
}

func VideoMQ2DB(ctx context.Context, req *videoPb.PublishRequest) error {
	token := req.Token
	data := req.Data
	title := req.Title
	uid, _ := util.GetUserIdFromToken(token)
	VideoUrl, _ := util.UploadVideo(data)
	imgPath := util.VideoGetNetImgCount(1, VideoUrl)
	coverUrl := util.UploadJPG(imgPath, VideoUrl)
	err := os.Remove(imgPath)
	if err != nil {
		return err
	}
	video := BuildVideoModel(uid, VideoUrl, coverUrl, title)
	//将视频存入数据库
	if err := mysql.NewVideoDao(ctx).CreateVideo(&video); err != nil {
		return err
	}
	//将视频存入缓存
	var videoCache *videoPb.Video
	videoCache = BuildVideoPbModel(ctx, &video, token)
	videoJson, _ := json.Marshal(&videoCache)
	redis.RedisClient.Set(ctx, fmt.Sprintf("%d", video.ID), videoJson, time.Hour)
	return nil
}

func VideoMQ2Redis(ctx context.Context, req *videoPb.Video) error {
	videoJson, _ := json.Marshal(&req)
	redis.RedisClient.Set(ctx, fmt.Sprintf("%d", req.Id), videoJson, time.Hour)
	return nil
}

func (v *VideoService) Feed(ctx context.Context, req *videoPb.FeedRequest) (*videoPb.FeedResponse, error) {

	latestTimeStamp := time.Now().Unix()
	latestTime := time.Unix(latestTimeStamp, 0)
	token := req.Token

	// 使用Keys命令获取所有键
	keys, err := redis.RedisClient.Keys(ctx, "*").Result()
	if err != nil {
		return &videoPb.FeedResponse{
			StatusCode: int32(resp.CodeVideoFeedFail),
			StatusMsg:  status.Errorf(codes.Internal, "获取视频流失败").Error(),
		}, err
	}
	keys = util.SortKeys(keys)
	var videoList []*videoPb.Video

	//从缓存取对应的视频
	for _, key := range keys {
		// 尝试从 Redis 缓存中获取数据
		redisResult, err := redis.RedisClient.Get(ctx, key).Result()
		if err != nil && err != redisv9.Nil {
			return &videoPb.FeedResponse{
				StatusCode: int32(resp.CodeVideoFeedFail),
				StatusMsg:  status.Errorf(codes.Internal, "获取视频流失败").Error(),
			}, err
		}
		if err != redisv9.Nil {
			var video videoPb.Video
			err = json.Unmarshal([]byte(redisResult), &video)
			if err != nil {
				return &videoPb.FeedResponse{
					StatusCode: int32(resp.CodeVideoFeedFail),
					StatusMsg:  status.Errorf(codes.Internal, "获取视频流失败").Error(),
				}, err
			}
			if token == "" {
				video.IsFavorite = false
				video.Author.IsFollow = false
			} else {
				video.IsFavorite, _ = mysql.NewVideoDao(ctx).GetIsFavorite(int(video.Id), token)
				video.Author.IsFollow, _ = mysql.NewVideoDao(ctx).GetIsFollowed(int(video.Author.Id), token)
			}
			videoList = append(videoList, &video)
		}
	}
	if len(keys) == 30 {
		return &videoPb.FeedResponse{
			StatusCode: int32(resp.CodeVideoFeedSuccess),
			StatusMsg:  resp.CodeVideoFeedSuccess.Msg(),
			VideoList:  videoList,
			NextTime:   latestTimeStamp,
		}, nil
	}

	//从数据库取对应的视频
	videos, err := mysql.NewVideoDao(ctx).GetVideoListByLatestTime(latestTime, util.StringArray2IntArray(keys), 30-len(keys))
	if err != nil {
		return &videoPb.FeedResponse{
			StatusCode: int32(resp.CodeVideoFeedFail),
			StatusMsg:  status.Errorf(codes.Internal, "获取视频流失败").Error(),
		}, err
	}
	var nextTime int64
	if len(videos) != 0 {
		nextTime = videos[len(videos)-1].CreatedAt.Unix()
	}
	for _, video := range videos {
		videoPbModel := BuildVideoPbModel(ctx, video, token)
		videoList = append(videoList, videoPbModel)
		//将视频存入缓存，加入消息队列
		body, _ := json.Marshal(&videoPbModel)
		err := mq.SendMessage2MQ(body, consts.Video2RedisQueue)
		if err != nil {
			return &videoPb.FeedResponse{
				StatusCode: int32(resp.CodeVideoFeedFail),
				StatusMsg:  status.Errorf(codes.Internal, "视频加入消息队列失败").Error(),
			}, nil
		}
	}
	return &videoPb.FeedResponse{
		StatusCode: int32(resp.CodeVideoFeedSuccess),
		StatusMsg:  resp.CodeVideoFeedSuccess.Msg(),
		VideoList:  videoList,
		NextTime:   nextTime,
	}, nil
}
