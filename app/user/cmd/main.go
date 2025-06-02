package main

import (
	"Douyin/app/gateway/rpc"
	"Douyin/app/user/dao/mysql"
	"Douyin/app/user/service"
	"Douyin/config"
	"Douyin/discovery"
	"Douyin/idl/user/userPb"
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	log.Println("用户服务启动...")
	// 初始化配置
	// 初始化配置
	config.InitConfig()
	// 初始化 rpc 服务 和 rpc 客户端连接（后面的网关请求处理中需要使用）
	rpc.InitRpc()
	// 初始化 mysql 服务
	mysql.InitMysql()
	// 注册本服务
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
	// 创建 gRPC 服务器
	s := grpc.NewServer()
	defer s.Stop()

	// 向 grpc服务器 执行服务注册
	cfg := gRPCRegisterConfig{
		Addr: config.Config.ServiceConfig.UserServiceAddress,
		RegisterFunc: func(g *grpc.Server) {
			userPb.RegisterUserServiceServer(g, service.NewUserService())
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

	// 构造服务节点信息
	info := discovery.Server{
		Name: config.Config.Domain.UserServiceDomain,
		Addr: config.Config.ServiceConfig.UserServiceAddress,
	}
	logrus.Println(info)

	// 注册 服务到 etcd
	_, err := r.Register(info, 2)
	if err != nil {
		logrus.Fatalln(err)
	}
}
