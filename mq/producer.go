package mq

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"time"
)

func SendMessage2MQ(body []byte, queueName string) error {
	// 将 json 格式的 req 数据发送到消息队列
	log.Println("Starting to send message...")
	// 检查连接是否有效
	if RabbitMQ == nil || RabbitMQ.IsClosed() {
		log.Println("RabbitMQ connection is nil or closed")
		return fmt.Errorf("RabbitMQ connection is invalid")
	}

	log.Println("Setting confirm mode...")
	if err := VideoMessageChannel.Confirm(false); err != nil {
		log.Printf("Failed to set confirm mode: %v", err)
		return err
	}

	log.Println("Declaring queue...")
	q, err := VideoMessageChannel.QueueDeclare(
		queueName,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,
	)
	if err != nil {
		log.Printf("Failed to declare queue: %v", err)
		return err
	}

	log.Println("Setting up publish confirmation...")
	// 增大确认通道的缓冲区
	confirms := VideoMessageChannel.NotifyPublish(make(chan amqp.Confirmation, 10))

	log.Println("Publishing message...")
	// TODO 这里 body 太大了，消息队列无法传递（报错）
	err = VideoMessageChannel.Publish(
		"",     // default exchange
		q.Name, // routing key
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
		},
	)
	if err != nil {
		log.Printf("Failed to publish message: %v", err)
		return err
	}

	log.Println("Waiting for confirmation...")
	select {
	case confirm := <-confirms:
		if !confirm.Ack {
			log.Printf("Message not acknowledged, deliveryTag=%d", confirm.DeliveryTag)
			return fmt.Errorf("消息未被确认，deliveryTag=%d", confirm.DeliveryTag)
		}
		log.Println("Message acknowledged successfully")
	case <-time.After(10 * time.Second): // 适当延长超时时间
		log.Println("Message publish confirmation timed out")
		return fmt.Errorf("消息发布确认超时")
	}

	log.Println("消息发送成功")
	return nil
}
