package rpc

import (
	"Douyin/config"
	"Douyin/discovery"
	"Douyin/idl/comment/commentPb"
	"Douyin/idl/favorite/favoritePb"
	"Douyin/idl/message/messagePb"
	"Douyin/idl/relation/relationPb"
	"Douyin/idl/testsrv/testsrvPb"
	"Douyin/idl/user/userPb"
	"Douyin/idl/video/videoPb"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
)

var (
	Register   *discovery.Resolver
	ctx        context.Context
	CancelFunc context.CancelFunc

	// TODO 这里进行 rpc 客户端全局声明
	TestClient     testsrvPb.TestServiceClient
	UserClient     userPb.UserServiceClient
	VideoClient    videoPb.VideoServiceClient
	RelationClient relationPb.RelationServiceClient
	FavoriteClient favoritePb.FavoriteServiceClient
	CommentClient  commentPb.CommentServiceClient
	MessageClient  messagePb.MessageServiceClient
)

// Init 初始化 RPC 客户端连接
func InitRpc() {
	EtcdAddress := fmt.Sprintf("%s:%d", config.Config.EtcdConfig.EtcdHost, config.Config.EtcdConfig.EtcdPort)
	fmt.Println(EtcdAddress)
	Register = discovery.NewResolver([]string{EtcdAddress}, logrus.New())
	resolver.Register(Register)
	ctx, CancelFunc = context.WithTimeout(context.Background(), 3*time.Second)
	defer CancelFunc()

	defer Register.Close()
	// TODO 这里进行 rpc客户端 初始化
	initClient(config.Config.Domain.TestServiceDomain, &TestClient)
	initClient(config.Config.Domain.UserServiceDomain, &UserClient)
	initClient(config.Config.Domain.VideoServiceDomain, &VideoClient)
	initClient(config.Config.Domain.RelationServiceDomain, &RelationClient)
	initClient(config.Config.Domain.FavoriteServiceDomain, &FavoriteClient)
	initClient(config.Config.Domain.CommentServiceDomain, &CommentClient)
	initClient(config.Config.Domain.MessageServiceDomain, &MessageClient)
}

// TODO 初始化客户端
func initClient(serviceName string, client interface{}) {
	conn, err := connectServer(serviceName)

	if err != nil {
		panic(err)
	}

	// TODO 这里添加不同服务端的客户端变量初始化
	switch c := client.(type) {
	case *testsrvPb.TestServiceClient:
		*c = testsrvPb.NewTestServiceClient(conn)
	case *userPb.UserServiceClient:
		*c = userPb.NewUserServiceClient(conn)
	case *videoPb.VideoServiceClient:
		*c = videoPb.NewVideoServiceClient(conn)
	case *relationPb.RelationServiceClient:
		*c = relationPb.NewRelationServiceClient(conn)
	case *favoritePb.FavoriteServiceClient:
		*c = favoritePb.NewFavoriteServiceClient(conn)
	case *commentPb.CommentServiceClient:
		*c = commentPb.NewCommentServiceClient(conn)
	case *messagePb.MessageServiceClient:
		*c = messagePb.NewMessageServiceClient(conn)
	default:
		panic("unsupported client type")
	}
}

// 初始化客户端连接服务端
func connectServer(serviceName string) (conn *grpc.ClientConn, err error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	}
	addr := fmt.Sprintf("%s:///%s", Register.Scheme(), serviceName)
	log.Printf("connectServer addr: %s", addr)
	// 调试信息
	log.Printf("Resolving address: %s", addr)

	// TODO 建立 gRPC 连接：使用上下文控制超时  【已弃用】（暂时使用）
	conn, err = grpc.Dial(addr, opts...)
	if err != nil {
		log.Printf("连接 %s 失败: %v", addr, err)
	}

	if err != nil {
		log.Printf("Failed to connect to %s: %v", addr, err)
	}
	return
}
