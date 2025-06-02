package service

import (
	"Douyin/idl/testsrv/testsrvPb"
	"context"
)

type TestService struct {
	testsrvPb.UnimplementedTestServiceServer
}

func NewTestService() *TestService {
	return &TestService{}
}

func (s *TestService) Test(ctx context.Context, req *testsrvPb.TestRequest) (*testsrvPb.TestResponse, error) {
	return &testsrvPb.TestResponse{
		Message: "Hello, " + req.Name,
	}, nil
	// 后面可能还会涉及到数据库的操作
}
