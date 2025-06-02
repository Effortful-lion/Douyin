# 基本介绍
基于 字节跳动暑期训练营 项目 - 简版抖音
基于 grpc + gorm + gin 实现
基于 mysql + redis 数据库
基于 rabbitMQ 消息队列
基于 阿里云 oss
基于 etcd 服务注册与发现、配置中心

# 项目结构
```go
├─app                           # 应用程序的主要代码
│  ├─comment                    # comment模块
│  │  ├─cmd                     # 微服务启动入口
│  │  ├─dao                     # 访问数据库相关代码
│  │  └─service                 # 业务逻辑处理
│  ├─favorite                   # favorite模块
│  │  ├─cmd                     # 微服务启动入口
│  │  ├─dao                     # 访问数据库相关代码
│  │  ├─script                  # 有关消息队列的脚本
│  │  └─service                 # 业务逻辑处理
│  ├─gateway                    # 网关模块
│  │  ├─cmd                     # 启动网关入口
│  │  ├─http                    # HTTP相关配置
│  │  ├─middleware              # 中间件
│  │  ├─router                  # 路由配置
│  │  ├─rpc                     # 远程过程调用配置
│  │  └─wrapper                 # 包装器(熔断和链路追踪)
│  ├─message                    # message模块
│  │  ├─cmd                     # 微服务启动入口
│  │  ├─dao                     # 访问数据库相关代码
│  │  └─service                 # 业务逻辑处理
│  ├─relation                   # relation模块
│  │  ├─cmd                     # 微服务启动入口
│  │  ├─dao                     # 访问数据库相关代码
│  │  └─service                 # 业务逻辑处理
│  ├─user                       # suer模块
│  │  ├─cmd                     # 微服务启动入口
│  │  ├─dao                     # 访问数据库相关代码
│  │  └─service                 # 业务逻辑处理
│  └─video                      # video模块
│      ├─cmd                     # 微服务启动入口
│      ├─dao                     # 访问数据库相关代码
│      ├─script                  # 有关消息队列的脚本
│      ├─service                 # 业务逻辑处理
│      └─tmp                     # 临时文件
├─config                        # 配置文件
├─consts                        # 常量定义
├─idl                           # 接口定义语言文件
│  ├─comment                    # comment模块的IDL文件
│  │  └─commentPb               # comment模块的Protocol Buffers文件
│  ├─favorite                   # favorite模块的IDL文件
│  │  └─favoritePb              # favorite模块的Protocol Buffers文件
│  ├─message                    # message模块的IDL文件
│  │  └─messagePb               # message模块的Protocol Buffers文件
│  ├─relation                   # relation模块的IDL文件
│  │  └─relationPb              # relation模块的Protocol Buffers文件
│  ├─user                       # user模块的IDL文件
│  │  └─userPb                  # user模块的Protocol Buffers文件
│  └─video                      # video模块的IDL文件
│      └─videoPb                # video模块的Protocol Buffers文件
├─model                         # 数据模型
├─mq                            # 消息队列相关代码
├─test                          # 测试代码
└─util                          # 工具函数和类

```
# 项目启动准备环境：
1. go 1.23.0
2. mysql 8.4.0
3. etcd 3.5.17
4. rabbitmq 4.1.0
5. redis 7.2.0
6. ffmpeg
7. grpc 开发环境