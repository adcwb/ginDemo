package wechat

import (
	"fmt"
	"ginDemo/global"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

// SendQueue 将信息发送到队列中
func SendQueue(QueueName, message string) error {
	// 3. 声明消息要发送到的队列
	q, err := global.RabbitMQChannel.QueueDeclare(
		QueueName, // 队列名称 kfQueue
		false,     // 是否持久化
		false,     // 是否自动删除
		false,     // 是否独占
		false,     // 是否等待消费者
		nil,       // 额外参数
	)
	if err != nil {
		zap.L().Error("声明消息要发送到的队列失败", zap.Error(err))
	}

	err = global.RabbitMQChannel.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	return err
}

// PullQueue 收取消息
func PullQueue(QueueName string) ([]byte, error) {
	q, err := global.RabbitMQChannel.QueueDeclare(
		QueueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	if err != nil {
		zap.L().Error("声明消息要发送到的队列失败", zap.Error(err))
	}

	// 从队列中获取消息
	msg, ok, err := global.RabbitMQChannel.Get(q.Name, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get a message: %s", err)
	}
	if !ok {
		return nil, fmt.Errorf("no message in the queue")
	}

	// 返回读取到的数据
	return msg.Body, nil
}
