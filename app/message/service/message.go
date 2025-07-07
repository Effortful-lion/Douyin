package service

import (
	"Douyin/app/message/dao/mysql"
	"Douyin/app/message/dao/redis"
	"Douyin/idl/message/messagePb"
	"Douyin/model"
	"Douyin/util"
	"context"
	"encoding/json"
	"errors"
	redisv9 "github.com/redis/go-redis/v9"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
	"time"
)

type MessageService struct {
	messagePb.UnimplementedMessageServiceServer
}

func NewMessageService() *MessageService {
	return &MessageService{}
}

func (m *MessageService) ChatMessage(ctx context.Context, req *messagePb.MessageChatRequest) (*messagePb.MessageChatResponse, error) {
	token := req.Token
	toUserId := req.ToUserId

	fromUserId, err := util.GetUserIdFromToken(token)
	if err != nil {
		return &messagePb.MessageChatResponse{
			StatusCode: 1,
			StatusMsg:  "获取用户ID失败！",
		}, status.Errorf(codes.Internal, "获取用户ID失败：%v", err)
	}

	if int64(fromUserId) == toUserId {
		return &messagePb.MessageChatResponse{
			StatusCode: 1,
			StatusMsg:  "不能查看和自己的聊天记录！",
		}, nil
	}

	// 构建 Redis 键
	var redisKey string
	if strconv.Itoa(fromUserId) < strconv.Itoa(int(toUserId)) {
		redisKey = "chat_messages:" + strconv.Itoa(fromUserId) + ":" + strconv.Itoa(int(toUserId))
	} else {
		redisKey = "chat_messages:" + strconv.Itoa(int(toUserId)) + ":" + strconv.Itoa(fromUserId)
	}

	// 尝试从 Redis 缓存中获取数据
	redisResult, err := redis.RedisClient.Get(ctx, redisKey).Result()
	if err != nil && errors.Is(err, redisv9.Nil) {
		return &messagePb.MessageChatResponse{
			StatusCode: 1,
			StatusMsg:  "获取聊天记录失败！",
		}, nil
	}

	res := &messagePb.MessageChatResponse{
		StatusCode:  0,
		StatusMsg:   "获取聊天记录成功",
		MessageList: make([]*messagePb.Message, 0),
	}

	if redisResult != "" {
		// 如果缓存中存在数据，则解码并直接返回缓存数据
		err = json.Unmarshal([]byte(redisResult), &res.MessageList)
		if err != nil {
			return &messagePb.MessageChatResponse{
				StatusCode: 1,
				StatusMsg:  "获取聊天记录失败！",
			}, nil
		}

		return &messagePb.MessageChatResponse{
			StatusCode:  0,
			StatusMsg:   "获取聊天记录成功！",
			MessageList: res.MessageList,
		}, nil
	}

	// 缓存中不存在数据，则从数据库获取聊天记录
	messages, err := mysql.NewMessageDao(ctx).FindAllMessages(int64(fromUserId), toUserId)
	if err != nil {
		return &messagePb.MessageChatResponse{
			StatusCode: 1,
			StatusMsg:  "获取聊天记录失败！",
		}, nil
	}

	for _, message := range messages {
		res.MessageList = append(res.MessageList, BuildMessagePbModel(message))
	}

	// 将结果存入 Redis 缓存
	jsonBytes, err := json.Marshal(&res.MessageList)
	if err != nil {
		return &messagePb.MessageChatResponse{
			StatusCode: 1,
			StatusMsg:  "获取聊天记录失败！",
		}, nil
	}

	err = redis.RedisClient.Set(ctx, redisKey, string(jsonBytes), time.Hour).Err()
	if err != nil {
		return &messagePb.MessageChatResponse{
			StatusCode: 1,
			StatusMsg:  "获取聊天记录失败！",
		}, nil
	}

	return &messagePb.MessageChatResponse{
		StatusCode:  0,
		StatusMsg:   "获取聊天记录成功！",
		MessageList: res.MessageList,
	}, nil
}

func (m *MessageService) ActionMessage(ctx context.Context, req *messagePb.MessageActionRequest) (*messagePb.MessageActionResponse, error) {
	token := req.Token
	fromUserID, err := util.GetUserIdFromToken(token)
	if err != nil {
		return &messagePb.MessageActionResponse{
			StatusCode: 1,
			StatusMsg:  "发送消息失败！",
		}, status.Errorf(codes.Internal, "获取用户ID失败：%v", err)
	}

	actionType := req.ActionType
	if actionType == 1 {
		toUserID := req.ToUserId
		content := req.Content

		message := BuildMessageModel(fromUserID, int(toUserID), content)
		id, err := mysql.NewMessageDao(ctx).CreateMessage(&message)
		if err != nil {
			return &messagePb.MessageActionResponse{
				StatusCode: 1,
				StatusMsg:  "发送消息失败！",
			}, nil
		}
		message.ID = uint(id)

		// 构建 Redis 键
		var redisKey string
		if strconv.Itoa(fromUserID) < strconv.Itoa(int(toUserID)) {
			redisKey = "chat_messages:" + strconv.Itoa(fromUserID) + ":" + strconv.Itoa(int(toUserID))
		} else {
			redisKey = "chat_messages:" + strconv.Itoa(int(toUserID)) + ":" + strconv.Itoa(fromUserID)
		}

		// 尝试从 Redis 缓存中获取数据
		redisResult, err := redis.RedisClient.Get(ctx, redisKey).Result()
		if err != nil && errors.Is(err, redisv9.Nil) {
			return &messagePb.MessageActionResponse{
				StatusCode: 1,
				StatusMsg:  "操作失败！",
			}, nil
		}

		var messageList []*messagePb.Message
		// 如果缓存中存在数据，则解码并合并到 messageList 中
		if redisResult != "" {
			// 解码 Redis 结果
			err = json.Unmarshal([]byte(redisResult), &messageList)
			if err != nil {
				return &messagePb.MessageActionResponse{
					StatusCode: 1,
					StatusMsg:  "操作失败！",
				}, nil
			}
			messageList = append(messageList, BuildMessagePbModel(&message))
		} else {
			// 如果缓存中不存在数据，则创建新的 messageList 切片
			messageList = []*messagePb.Message{}
		}

		// 将结果存入 Redis 缓存
		jsonBytes, err := json.Marshal(&messageList)
		if err != nil {
			return &messagePb.MessageActionResponse{
				StatusCode: 1,
				StatusMsg:  "操作失败！",
			}, nil
		}

		err = redis.RedisClient.Set(ctx, redisKey, string(jsonBytes), time.Hour).Err()
		if err != nil {
			return &messagePb.MessageActionResponse{
				StatusCode: 1,
				StatusMsg:  "操作失败！",
			}, nil
		}

		return &messagePb.MessageActionResponse{
			StatusCode: 0,
			StatusMsg:  "发送消息成功！",
		}, nil

	} else {
		return &messagePb.MessageActionResponse{
			StatusCode: 1,
			StatusMsg:  "非发送消息操作！",
		}, nil
	}
}

func BuildMessageModel(fromUserID int, toUserID int, content string) model.Message {
	return model.Message{
		FromUserID: fromUserID,
		ToUserID:   toUserID,
		Content:    content,
		CreatedAt:  time.Now(),
	}
}

func BuildMessagePbModel(message *model.Message) *messagePb.Message {
	return &messagePb.Message{
		Id:         int64(message.ID),
		FromUserId: int64(message.FromUserID),
		ToUserId:   int64(message.ToUserID),
		Content:    message.Content,
		//CreateTime: message.CreatedAt.Format("2006-01-02 15:04:05"),
		//CreateTime: time.ParseInLocation("2006-01-02 15:04:05", message.CreatedAt, time.Local),
		CreateTime: "0",
	}
}
