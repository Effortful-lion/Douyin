package main

import (
	"Douyin/app/gateway/router"
	"Douyin/app/gateway/rpc"
	"Douyin/config"
	"Douyin/util"
	"fmt"
	"log"
	"net/http"
	"time"
)

// 网关：接收所有请求，转发服务

func main() {
	log.Println("网关服务启动...")
	// 初始化配置
	config.InitConfig()
	// 初始化 rpc 服务 和 rpc 客户端连接（后面的网关请求处理中需要使用）
	rpc.InitRpc()
	// 启动网关服务
	go startListen()
	// 阻塞主进程，保持服务运行
	select {}
}

// 启动 http 服务监听
func startListen() {
	// 加入熔断 TODO main太臃肿了
	// wrapper.NewServiceWrapper(userServiceName)
	// wrapper.NewServiceWrapper(taskServiceName)

	// 初始化gin路由
	r := router.NewRouter()
	// 配置http服务
	server := &http.Server{
		Addr:           config.Config.GinConfig.AppHost + ":" + config.Config.GinConfig.AppPort,
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	// 启动http服务：目前是 4000 端口
	if err := server.ListenAndServe(); err != nil {
		fmt.Println("gateway启动失败, err: ", err)
	}

	// 启动一个协程，优雅关闭http服务
	go func() {
		// 优雅关闭
		util.GracefullyShutdown(server)
	}()
}
