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
	"errors"
	"fmt"
	redisv9 "github.com/redis/go-redis/v9"
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
				// TODO ？ 关注操作，为什么要重新缓存视频信息？只是因为视频信息中包含关注信息嘛
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
					redisResult, err := redis.RedisClient.Get(ctx, key).Result()
					if err != nil && !errors.Is(err, redisv9.Nil) {
						return &relationPb.RelationActionResponse{
							StatusCode: 1,
							StatusMsg:  "获取视频流失败",
						}, err
					}
					if !errors.Is(err, redisv9.Nil) {
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
								return &relationPb.RelationActionResponse{
									StatusCode: 1,
									StatusMsg:  "解析视频流失败",
								}, err
							}
							err = redis.RedisClient.Set(ctx, key, updatedVideo, time.Hour).Err()
							if err != nil {
								return &relationPb.RelationActionResponse{
									StatusCode: 1,
									StatusMsg:  "更新视频流失败",
								}, err
							}
						}
					}
				}
				return &relationPb.RelationActionResponse{
					StatusCode: 0,
					StatusMsg:  "关注成功！",
				}, nil
			} else {
				return &relationPb.RelationActionResponse{
					StatusCode: 1,
					StatusMsg:  "关注失败，不可重复关注！",
				}, nil
			}
		} else {
			return &relationPb.RelationActionResponse{
				StatusCode: 1,
				StatusMsg:  "关注失败！",
			}, err
		}

	} else if actionType == 2 {
		// 取消关注
		Relation := BuildRelationModel(int(toUserId), fromUserId)
		if _, err := mysql.NewRelationDao(ctx).CancelFollow(&Relation); err == nil {
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
				redisResult, err := redis.RedisClient.Get(ctx, key).Result()
				if err != nil && err != redisv9.Nil {
					return &relationPb.RelationActionResponse{
						StatusCode: 1,
						StatusMsg:  "获取视频流失败",
					}, err
				}
				if err != redisv9.Nil {
					var video videoPb.Video
					err = json.Unmarshal([]byte(redisResult), &video)
					if err != nil {
						return &relationPb.RelationActionResponse{
							StatusCode: 1,
							StatusMsg:  "解析视频流失败",
						}, err
					}
					if video.Author.Id == toUserId {
						video.Author.IsFollow = false
						updatedVideo, err := json.Marshal(&video)
						if err != nil {
							return &relationPb.RelationActionResponse{
								StatusCode: 1,
								StatusMsg:  "解析视频流失败",
							}, err
						}
						err = redis.RedisClient.Set(ctx, key, updatedVideo, time.Hour).Err()
						if err != nil {
							return &relationPb.RelationActionResponse{
								StatusCode: 1,
								StatusMsg:  "更新视频流失败",
							}, err
						}
					}
				}
			}
			return &relationPb.RelationActionResponse{
				StatusCode: 0,
				StatusMsg:  "取消关注成功！",
			}, nil
		}
		return &relationPb.RelationActionResponse{
			StatusCode: 1,
			StatusMsg:  "取消关注失败！",
		}, nil

	}
	return &relationPb.RelationActionResponse{
		StatusCode: 1,
		StatusMsg:  "操作失败！",
	}, nil
}

func (r *RelationService) ListFollowRelation(ctx context.Context, req *relationPb.RelationFollowRequest) (*relationPb.RelationFollowResponse, error) {
	userId := req.UserId
	token := req.Token
	follows, err := mysql.NewRelationDao(ctx).FindAllFollow(int(userId))
	if err != nil {
		return &relationPb.RelationFollowResponse{
			StatusCode: 1,
			StatusMsg:  "获取关注列表失败！",
		}, err
	}

	// 初始化，对空切片进行操作不会导致panic
	res := &relationPb.RelationFollowResponse{
		StatusCode: 0,
		StatusMsg:  "获取关注列表成功！",
		UserList:   make([]*relationPb.User, 0),
	}

	for _, follow := range follows {
		user, _ := mysql.NewRelationDao(ctx).FindUserById(follow.UserID)
		if user.ID == 0 {
			return &relationPb.RelationFollowResponse{
				StatusCode: 1,
				StatusMsg:  "关注的用户不存在",
			}, nil
		}
		res.UserList = append(res.UserList, BuildUserPbModel(ctx, user, token))
	}
	return &relationPb.RelationFollowResponse{
		StatusCode: 0,
		StatusMsg:  "获取关注列表成功！",
		UserList:   res.UserList,
	}, nil
}

func (r *RelationService) ListFollowerRelation(ctx context.Context, req *relationPb.RelationFollowerRequest) (*relationPb.RelationFollowerResponse, error) {
	userId := req.UserId
	token := req.Token
	follows, err := mysql.NewRelationDao(ctx).FindAllFollower(int(userId))
	if err != nil {
		return &relationPb.RelationFollowerResponse{
			StatusCode: 1,
			StatusMsg:  "获取粉丝列表失败！",
		}, err
	}

	res := &relationPb.RelationFollowerResponse{
		StatusCode: 0,
		StatusMsg:  "获取粉丝列表成功！",
		UserList:   make([]*relationPb.User, 0),
	}

	for _, follow := range follows {
		user, _ := mysql.NewRelationDao(ctx).FindUserById(follow.FollowedUserID)
		if user.ID == 0 {
			return &relationPb.RelationFollowerResponse{
				StatusCode: 1,
				StatusMsg:  "用户不存在！",
			}, err
		}
		res.UserList = append(res.UserList, BuildUserPbModel(ctx, user, token))
	}
	return res, nil
}

func (r *RelationService) ListFriendRelation(ctx context.Context, req *relationPb.RelationFriendRequest) (*relationPb.RelationFriendResponse, error) {
	userId := req.UserId
	token := req.Token
	follows, err := mysql.NewRelationDao(ctx).FindAllFollower(int(userId))
	if err != nil {
		return &relationPb.RelationFriendResponse{
			StatusCode: 1,
			StatusMsg:  "获取好友列表失败！",
		}, err
	}

	res := &relationPb.RelationFriendResponse{
		StatusCode: 0,
		StatusMsg:  "获取好友列表成功！",
		UserList:   make([]*relationPb.User, 0),
	}

	for _, follow := range follows {
		count, _ := mysql.NewRelationDao(ctx).GetFriendCount(follow.FollowedUserID, int(userId))
		if count == 1 {
			user, _ := mysql.NewRelationDao(ctx).FindUserById(follow.FollowedUserID)
			if user.ID == 0 {
				return &relationPb.RelationFriendResponse{
					StatusCode: 1,
					StatusMsg:  "用户不存在！",
				}, err
			}
			res.UserList = append(res.UserList, BuildUserPbModel(ctx, user, token))
		} else {
			continue
		}
	}
	return res, nil
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
