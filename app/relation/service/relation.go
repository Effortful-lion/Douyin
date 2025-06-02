package service

import (
	"Douyin/app/relation/dao/mysql"
	"Douyin/app/relation/dao/redis"
	"Douyin/idl/relation/relationPb"
	"Douyin/idl/video/videoPb"
	"Douyin/model"
	"Douyin/util"
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type RelationService struct {
	relationPb.UnimplementedRelationServiceServer
}

func NewRelationService() *RelationService {
	return &RelationService{}
}

func (r *RelationService) ActionRelation(ctx context.Context, req *relationPb.RelationActionRequest) (*relationPb.RelationActionResponse, error) {
	token := req.Token
	toUserId := req.ToUserId
	actionType := req.ActionType

	fromUserId, err := util.GetUserIdFromToken(token)
	if err != nil {
		return &relationPb.RelationActionResponse{
			StatusCode: 1,
			StatusMsg:  "token解析失败",
		}, err
	}

	if int64(fromUserId) == toUserId {
		return &relationPb.RelationActionResponse{
			StatusCode: 1,
			StatusMsg:  "不能关注自己",
		}, nil
	}
	if actionType == 1 {
		// 关注
		user, _ := mysql.NewRelationDao(ctx).FindUserById(int(toUserId))
		if user.ID == 0 {
			return &relationPb.RelationActionResponse{
				StatusCode: 1,
				StatusMsg:  "关注的用户不存在",
			}, nil
		}

		Relation := BuildRelationModel(int(toUserId), fromUserId)
		if RowsAffected, err := mysql.NewRelationDao(ctx).AddFollow(&Relation); err == nil {
			if RowsAffected > 0 {
				pattern := fmt.Sprintf("%d:*", fromUserId)
				// 使用Keys命令获取所有键
				keys, err := redis.RedisClient.Keys(ctx, pattern).Result()
				if err != nil {
					return &relationPb.RelationActionResponse{
						StatusCode: 1,
						StatusMsg:  "获取视频流失败",
					}, err
				}
				//从缓存取对应的视频
				for _, key := range keys {
					redisResult, err := dao.RedisClient.Get(ctx, key).Result()
					if err != nil && err != redis.Nil {
						return &relationPb.RelationActionResponse{
							StatusCode: 1,
							StatusMsg:  "获取视频流失败",
						}, err
					}
					if err != redis.Nil {
						var video videoPb.Video
						err = json.Unmarshal([]byte(redisResult), &video)
						if err != nil {
							return &relationPb.RelationActionResponse{
								StatusCode: 1,
								StatusMsg:  "解析视频流失败",
							}, err
						}
						if video.Author.Id == toUserId {
							video.Author.IsFollow = true
							updatedVideo, err := json.Marshal(&video)
							if err != nil {
								RelationActionResponseData(res, 1, "解析视频流失败")
								return err
							}
							err = dao.RedisClient.Set(ctx, key, updatedVideo, time.Hour).Err()
							if err != nil {
								RelationActionResponseData(res, 1, "更新视频流失败")
								return err
							}
						}
					}
				}
				RelationActionResponseData(res, 0, "关注成功!")
				return nil
			} else {
				RelationActionResponseData(res, 1, "已经关注过该用户了！")
				return nil
			}
		} else {
			RelationActionResponseData(res, 1, "关注失败!")
			return err
		}

	} else if actionType == 2 {
		// 取消关注
		Relation := BuildRelationModel(int(toUserId), fromUserId)
		if _, err := dao.NewRelationDao(ctx).CancelFollow(&Relation); err == nil {
			pattern := fmt.Sprintf("%d:*", fromUserId)
			// 使用Keys命令获取所有键
			keys, err := dao.RedisClient.Keys(ctx, pattern).Result()
			if err != nil {
				RelationActionResponseData(res, 1, "获取视频流失败")
				return err
			}
			//从缓存取对应的视频
			for _, key := range keys {
				redisResult, err := dao.RedisClient.Get(ctx, key).Result()
				if err != nil && err != redis.Nil {
					RelationActionResponseData(res, 1, "获取视频流失败")
					return err
				}
				if err != redis.Nil {
					var video videoPb.Video
					err = json.Unmarshal([]byte(redisResult), &video)
					if err != nil {
						RelationActionResponseData(res, 1, "解析视频流失败")
						return err
					}
					if video.Author.Id == toUserId {
						video.Author.IsFollow = false
						updatedVideo, err := json.Marshal(&video)
						if err != nil {
							RelationActionResponseData(res, 1, "解析视频流失败")
							return err
						}
						err = dao.RedisClient.Set(ctx, key, updatedVideo, time.Hour).Err()
						if err != nil {
							RelationActionResponseData(res, 1, "更新视频流失败")
							return err
						}
					}
				}
			}
			RelationActionResponseData(res, 0, "取消关注成功！")
			return nil
		}
		RelationActionResponseData(res, 1, "取消关注失败！")
		return err

	}
	RelationActionResponseData(res, 1, "检查参数信息！")
	return nil
}

func (r *RelationService) ListFollowRelation(context.Context, *relationPb.RelationFollowRequest) (*relationPb.RelationFollowResponse, error) {

}

func (r *RelationService) ListFollowerRelation(context.Context, *relationPb.RelationFollowerRequest) (*relationPb.RelationFollowerResponse, error) {

}

func (r *RelationService) ListFriendRelation(context.Context, *relationPb.RelationFriendRequest) (*relationPb.RelationFriendResponse, error) {

}

func BuildRelationModel(userID int, followedUserID int) model.Follow {
	return model.Follow{
		UserID:         userID,
		FollowedUserID: followedUserID,
		CreatedAt:      time.Now(),
	}
}

func BuildUserPbModel(ctx context.Context, user *model.User, token string) *relationPb.User {
	uid := int(user.ID)
	FollowCount, _ := mysql.NewRelationDao(ctx).GetFollowCount(uid)
	FollowerCount, _ := mysql.NewRelationDao(ctx).GetFollowerCount(uid)
	WorkCount, _ := mysql.NewRelationDao(ctx).GetWorkCount(uid)
	FavoriteCount, _ := mysql.NewRelationDao(ctx).GetFavoriteCount(uid)
	TotalFavorited, _ := mysql.NewRelationDao(ctx).GetTotalFavorited(uid)
	IsFollow, _ := mysql.NewRelationDao(ctx).GetIsFollowed(uid, token)
	return &relationPb.User{
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
