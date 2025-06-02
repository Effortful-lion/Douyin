package rpc

import (
	"Douyin/idl/relation/relationPb"
	"context"
)

func ActionRelation(ctx context.Context, req *relationPb.RelationActionRequest) (res *relationPb.RelationActionResponse, err error) {
	res, err = RelationClient.ActionRelation(ctx, req)
	if err != nil {
		return
	}
	return
}

func ListFollowRelation(ctx context.Context, req *relationPb.RelationFollowRequest) (res *relationPb.RelationFollowResponse, err error) {
	res, err = RelationClient.ListFollowRelation(ctx, req)
	if err != nil {
		return
	}
	return
}

func ListFollowerRelation(ctx context.Context, req *relationPb.RelationFollowerRequest) (res *relationPb.RelationFollowerResponse, err error) {
	res, err = RelationClient.ListFollowerRelation(ctx, req)
	if err != nil {
		return
	}
	return
}

func ListFriendRelation(ctx context.Context, req *relationPb.RelationFriendRequest) (res *relationPb.RelationFriendResponse, err error) {
	res, err = RelationClient.ListFriendRelation(ctx, req)
	if err != nil {
		return
	}
	return
}
