package main

import (
	"Douyin/app/message/dao/mysql"
	"Douyin/app/message/dao/redis"
	"Douyin/app/message/service"
	"Douyin/config"
	"Douyin/discovery"
	"Douyin/idl/message/messagePb"
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	config.InitConfig()
	mysql.InitMysql()
	redis.InitRedis()

	// 服务注册
	registerService()
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
	s := grpc.NewServer()
	defer s.Stop()

	// TODO 不同： 向 grpc服务器 执行服务注册
	cfg := gRPCRegisterConfig{
		Addr: config.Config.ServiceConfig.MessageServiceAddress,
		RegisterFunc: func(g *grpc.Server) {
			messagePb.RegisterMessageServiceServer(g, service.NewMessageService())
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
		Name: config.Config.Domain.MessageServiceDomain,
		Addr: config.Config.ServiceConfig.MessageServiceAddress,
	}
	logrus.Println(info)

	// 注册 服务到 etcd
	_, err := r.Register(info, 2)
	if err != nil {
		logrus.Fatalln(err)
	}
}
