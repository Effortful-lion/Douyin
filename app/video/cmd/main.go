package main

import (
	"Douyin/app/video/dao/mysql"
	"Douyin/app/video/dao/redis"
	"Douyin/app/video/script"
	"Douyin/app/video/service"
	"Douyin/config"
	"Douyin/discovery"
	"Douyin/idl/video/videoPb"
	"Douyin/mq"
	"context"
	"fmt"
	"log"
	"net"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	config.InitConfig()
	mysql.InitMysql()
	redis.InitRedis()
	mq.InitRabbitMQ()
	loadingScript()

	// 注册本服务
	registerService()
}

func loadingScript() {
	ctx := context.Background()
	go script.VideoCreateSync(ctx)
	go script.Video2RedisSync(ctx)
}

type gRPCRegisterConfig struct {
	Addr         string
	RegisterFunc func(g *grpc.Server)
}

func registerService() {
	// 注册服务到 etcd
	registerEtcdService()
	// 注册服务到 grpc 服务器
	registerGrpcService()
}

// 注册本服务到 grpc 服务器
func registerGrpcService() {
	// 创建 gRPC 服务器, 设置最大接收和发送消息大小为 32MB
	// maxMsgSize := 32 * 1024 * 1024 // 32MB
	maxMsgSize := 32 * 1024 * 1024
	s := grpc.NewServer(
		grpc.MaxRecvMsgSize(maxMsgSize),
		grpc.MaxSendMsgSize(maxMsgSize),
	)
	defer s.Stop()

	// TODO 不同： 向 grpc服务器 执行服务注册
	cfg := gRPCRegisterConfig{
		Addr: config.Config.ServiceConfig.VideoServiceAddress,
		RegisterFunc: func(g *grpc.Server) {
			videoPb.RegisterVideoServiceServer(g, service.NewVideoService())
		},
	}
	cfg.RegisterFunc(s)

	// 监听端口，监听地址：服务注册地址
	lis, err := net.Listen("tcp", cfg.Addr)
	if err != nil {
		log.Println("cannot listen")
	}
	// 启动 gRPC 服务器，开启监听端口的处理程序

	log.Printf("grpc server started as: %s \n", cfg.Addr)
	err = s.Serve(lis)
	if err != nil {
		log.Println("server started error", err)
		return
	}

	// TODO s 服务器资源优雅关闭
}

// 注册服务到 etcd
func registerEtcdService() {
	// 获取 etcd 地址并创建 etcd 服务注册器
	etcd_addr := fmt.Sprintf("%s:%d", config.Config.EtcdConfig.EtcdHost, config.Config.EtcdConfig.EtcdPort)
	// 创建 etcd 注册器
	r := discovery.NewRegister([]string{etcd_addr}, logrus.New())
	defer r.Stop()

	// TODO 不同：构造服务节点信息
	info := discovery.Server{
		Name: config.Config.Domain.VideoServiceDomain,
		Addr: config.Config.ServiceConfig.VideoServiceAddress,
	}
	logrus.Println(info)

	// 注册 服务到 etcd
	_, err := r.Register(info, 2)
	if err != nil {
		logrus.Fatalln(err)
	}
}
