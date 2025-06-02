package rpc

import (
	"Douyin/idl/video/videoPb"
	"context"
)

func Feed(ctx context.Context, req *videoPb.FeedRequest) (res *videoPb.FeedResponse, err error) {
	res, err = VideoClient.Feed(ctx, req)
	if err != nil {
		return
	}
	return
}

func Publish(ctx context.Context, req *videoPb.PublishRequest) (res *videoPb.PublishResponse, err error) {
	res, err = VideoClient.Publish(ctx, req)
	if err != nil {
		return
	}
	return
}

func PublishList(ctx context.Context, req *videoPb.PublishListRequest) (res *videoPb.PublishListResponse, err error) {
	res, err = VideoClient.PublishList(ctx, req)
	if err != nil {
		return
	}
	return
}