package script

import (
	"Douyin/app/video/service"
	"Douyin/consts"
	"Douyin/idl/video/videoPb"
	"Douyin/mq"
	"context"
	"encoding/json"
)

// 同步视频操作对象
type SyncVideo struct {
}

// 对外 api
func VideoCreateSync(ctx context.Context) {
	Sync := new(SyncVideo)
	err := Sync.SyncVideoCreate(ctx, consts.CreateVideoQueue)
	if err != nil {
		return
	}
}

func Video2RedisSync(ctx context.Context) {
	Sync := new(SyncVideo)
	err := Sync.SyncVideo2Redis(ctx, consts.Video2RedisQueue)
	if err != nil {
		return
	}
}

func (s *SyncVideo) SyncVideoCreate(ctx context.Context, queueName string) error {
	// 获得消息队列
	msg, err := mq.ConsumeMessage(ctx, queueName)
	if err != nil {
		return err
	}
	var forever chan struct{}
	go func() {
		for d := range msg {
			// 落库
			var req *videoPb.PublishRequest
			err = json.Unmarshal(d.Body, &req)
			if err != nil {
				return
			}
			err = service.VideoMQ2DB(ctx, req)
			if err != nil {
				return
			}
			d.Ack(false)
		}
	}()
	<-forever
	return nil
}

func (s *SyncVideo) SyncVideo2Redis(ctx context.Context, queueName string) error {
	msg, err := mq.ConsumeMessage(ctx, queueName)
	if err != nil {
		return err
	}
	var forever chan struct{}
	go func() {
		for d := range msg {
			// 落库
			var req *videoPb.Video
			err = json.Unmarshal(d.Body, &req)
			if err != nil {
				return
			}
			err = service.VideoMQ2Redis(ctx, req)
			if err != nil {
				return
			}
			d.Ack(false)
		}
	}()
	<-forever
	return nil
}
