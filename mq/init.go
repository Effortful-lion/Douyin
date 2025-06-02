package mq

import (
	"Douyin/config"
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

var RabbitMQ *amqp.Connection

// 准备几个全局的通道
var (
	VideoMessageChannel *amqp.Channel
	//UserMessageChannel  *amqp.Channel
)

func InitRabbitMQ() {
	config.InitConfig()
	connString := fmt.Sprintf("%s://%s:%s@%s:%d/",
		config.Config.RabbitMQConfig.RabbitMQ,
		config.Config.RabbitMQConfig.RabbitMQUser,
		config.Config.RabbitMQConfig.RabbitMQPassword,
		config.Config.RabbitMQConfig.RabbitMQHost,
		config.Config.RabbitMQConfig.RabbitMQPort,
	)
	conn, err := amqp.Dial(connString)
	if err != nil {
		panic(err)
	}
	RabbitMQ = conn

	// 设置一个全局的默认的通道，以便后面进行消息的发布
	VideoMessageChannel, err = conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}

	log.Println("Video channel created successfully")
}
