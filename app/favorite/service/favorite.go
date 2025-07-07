package service

import (
	"Douyin/app/favorite/dao/mysql"
	"Douyin/app/favorite/dao/redis"
	"Douyin/consts"
	"Douyin/idl/favorite/favoritePb"
	"Douyin/model"
	"Douyin/mq"
	"Douyin/util"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	redisv9 "github.com/redis/go-redis/v9"
)

type FavoriteService struct {
	favoritePb.UnimplementedFavoriteServiceServer
}

func NewFavoriteService() *FavoriteService {
	return &FavoriteService{}
}

func (f *FavoriteService) FavoriteAction(ctx context.Context, req *favoritePb.FavoriteActionRequest) (*favoritePb.FavoriteActionResponse, error) {
	actionType := req.ActionType
	vid := req.VideoId
	uid, _ := util.GetUserIdFromToken(req.Token)

	body, _ := json.Marshal(&req)

	// 点赞
	if actionType == 1 {
		// 不能重复点赞
		isFavorite, _ := mysql.NewFavoriteDao(ctx).GetIsFavoriteByUserIdAndVid(int64(uid), vid)

		if isFavorite {
			return &favoritePb.FavoriteActionResponse{
				StatusCode: 1,
				StatusMsg:  "不能重复点赞",
			}, nil
		}

		//修改redis
		key := fmt.Sprintf("%d", vid)
		redisResult, err := redis.RedisClient.Get(ctx, key).Result()
		if err != nil && err != redisv9.Nil {
			return &favoritePb.FavoriteActionResponse{
				StatusCode: 1,
				StatusMsg:  "点赞失败",
			}, err
		}

		if err != redisv9.Nil { // 在redis中找到了视频信息
			var video favoritePb.Video
			err = json.Unmarshal([]byte(redisResult), &video)
			if err != nil {
				return &favoritePb.FavoriteActionResponse{
					StatusCode: 1,
					StatusMsg:  "点赞失败",
				}, err
			}

			video.FavoriteCount += 1

			videoJson, _ := json.Marshal(&video)
			redis.RedisClient.Set(ctx, key, videoJson, time.Hour)
		}

		// 加入消息队列
		err = mq.SendMessage2MQ(body, consts.CreateFavorite2MQ)
		if err != nil {
			return &favoritePb.FavoriteActionResponse{
				StatusCode: 1,
				StatusMsg:  "点赞失败",
			}, err
		}
	}

	// 取消点赞
	if actionType == 2 {
		//修改redis
		key := fmt.Sprintf("%d", vid)
		redisResult, err := redis.RedisClient.Get(ctx, key).Result()
		if err != nil && !errors.Is(err, redisv9.Nil) {
			return &favoritePb.FavoriteActionResponse{
				StatusCode: 1,
				StatusMsg:  "取消点赞失败",
			}, err
		}

		if !errors.Is(err, redisv9.Nil) {
			// 在redis中找到了视频信息
			var video favoritePb.Video
			err = json.Unmarshal([]byte(redisResult), &video)
			if err != nil {
				return &favoritePb.FavoriteActionResponse{
					StatusCode: 1,
					StatusMsg:  "取消点赞失败",
				}, err
			}

			video.FavoriteCount -= 1

			videoJson, _ := json.Marshal(&video)
			redis.RedisClient.Set(ctx, key, videoJson, time.Hour)
		}

		// 加入消息队列
		err = mq.SendMessage2MQ(body, consts.DeleteFavorite2MQ)
		if err != nil {
			return &favoritePb.FavoriteActionResponse{
				StatusCode: 1,
				StatusMsg:  "取消点赞失败",
			}, err
		}
	}

	return &favoritePb.FavoriteActionResponse{
		StatusCode: 0,
		StatusMsg:  "操作成功",
	}, nil
}

func (f *FavoriteService) FavoriteList(ctx context.Context, req *favoritePb.FavoriteListRequest) (*favoritePb.FavoriteListResponse, error) {
	uid := req.UserId

	favorites, err := mysql.NewFavoriteDao(ctx).GetFavoriteListByUserId(uid)

	if err != nil {
		return &favoritePb.FavoriteListResponse{
			StatusCode: 1,
			StatusMsg:  "获取喜欢列表失败",
		}, err
	}

	res := &favoritePb.FavoriteListResponse{
		StatusCode: 0,
		StatusMsg:  "获取喜欢列表成功",
		VideoList:  make([]*favoritePb.Video, 0),
	}

	res.StatusCode = 0
	res.StatusMsg = "获取喜欢列表成功"

	for _, favorite := range favorites {
		vid := favorite.VideoID
		res.VideoList = append(res.VideoList, BuildVideoPbModelByVid(ctx, vid))
	}
	return res, nil
}

func BuildVideoPbModelByVid(ctx context.Context, vid int) *favoritePb.Video {
	FavoriteCount, _ := mysql.NewFavoriteDao(ctx).GetFavoriteCount(vid)
	PlayUrl, _ := mysql.NewFavoriteDao(ctx).GetPlayUrlByVid(vid)
	CoverUrl, _ := mysql.NewFavoriteDao(ctx).GetCoverUrlByVid(vid)

	return &favoritePb.Video{
		Id:            int64(vid),
		PlayUrl:       PlayUrl,
		CoverUrl:      CoverUrl,
		FavoriteCount: FavoriteCount,
	}
}

func FavoriteMQ2DB(ctx context.Context, req *favoritePb.FavoriteActionRequest) error {
	token := req.Token
	videoId := req.VideoId
	actionType := req.ActionType

	// 解析token
	user, _ := util.GetUserFromToken(token)

	// 点赞
	if actionType == 1 {
		//加入redis
		//加入校队列

		favorite := model.Favorite{
			UserID:  int(user.ID), // uint to int
			VideoID: int(videoId), // int64 to int
		}
		if err := mysql.NewFavoriteDao(ctx).CreateFavorite(&favorite); err != nil {
			return err
		}

		return nil
	}

	// 取消点赞
	favorite := model.Favorite{
		UserID:  int(user.ID), // uint to int
		VideoID: int(videoId), // int64 to int
	}

	if err := mysql.NewFavoriteDao(ctx).DeleteFavorite(&favorite); err != nil {
		return err
	}

	return nil
}
