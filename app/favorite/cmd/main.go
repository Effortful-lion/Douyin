package main

import (
	"Douyin/app/favorite/dao/mysql"
	"Douyin/app/favorite/dao/redis"
	"Douyin/app/favorite/script"
	"Douyin/app/favorite/service"
	"Douyin/config"
	"Douyin/discovery"
	"Douyin/idl/favorite/favoritePb"
	"Douyin/mq"
	"context"
	"fmt"
	"log"
	"net"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	// 初始化配置
	config.InitConfig()
	// 初始化 mysql 服务
	mysql.InitMysql()
	// 初始化 redis 服务
	redis.InitRedis()
	// 初始化 mq
	mq.InitRabbitMQ()
	// TODO	loadScript() 作用： 同步落库
	loadingScript()

	// 注册服务
	registerService()

}

func registerService() {
	// 注册服务到 etcd
	registerEtcdService()
	// 注册服务到 grpc 服务器
	registerGrpcService()
}

type gRPCRegisterConfig struct {
	Addr         string
	RegisterFunc func(g *grpc.Server)
}

// 注册本服务到 grpc 服务器
func registerGrpcService() {
	// 创建 gRPC 服务器
	s := grpc.NewServer()
	defer s.Stop()

	// TODO 不同： 向 grpc服务器 执行服务注册
	cfg := gRPCRegisterConfig{
		Addr: config.Config.ServiceConfig.FavoriteServiceAddress,
		RegisterFunc: func(g *grpc.Server) {
			favoritePb.RegisterFavoriteServiceServer(g, service.NewFavoriteService())
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
	// 获取 etcd 地址
	etcd_addr := fmt.Sprintf("%s:%d", config.Config.EtcdConfig.EtcdHost, config.Config.EtcdConfig.EtcdPort)
	// 创建 etcd 注册器
	r := discovery.NewRegister([]string{etcd_addr}, logrus.New())
	defer r.Stop()

	// TODO 不同：构造服务节点信息
	info := discovery.Server{
		Name: config.Config.Domain.FavoriteServiceDomain,
		Addr: config.Config.ServiceConfig.FavoriteServiceAddress,
	}
	logrus.Println(info)

	// 注册 服务到 etcd
	_, err := r.Register(info, 2)
	if err != nil {
		logrus.Fatalln(err)
	}
}

func loadingScript() {
	ctx := context.Background()
	go script.FavoriteCreateSync(ctx)
	go script.FavoriteDeleteSync(ctx)
}
