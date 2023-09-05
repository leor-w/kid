package rabbitmq

import (
	"fmt"

	"github.com/streadway/amqp"
)

type Option func(*Options)

// Producer 生产者
type Producer interface {
	Publish() interface{}
}

// Consumer 消费者
type Consumer interface {
	Consumer(msg string) error
}

type RabbitMQ struct {
	conn    *amqp.Connection
	ch      *amqp.Channel
	options *Options
}

func (mq *RabbitMQ) connect() error {
	connUrl := fmt.Sprintf("amqp://%s:%s@%s:%d", mq.options.User, mq.options.Pwd, mq.options.Host, mq.options.Port)
	var err error
	mq.conn, err = amqp.Dial(connUrl)
	if err != nil {
		return err
	}
	mq.ch, err = mq.conn.Channel()
	if err != nil {
		return err
	}
	return nil
}

func (mq *RabbitMQ) Destroy() error {
	if err := mq.ch.Close(); err != nil {
		return fmt.Errorf("关闭channel失败: %w", err)
	}
	if err := mq.conn.Close(); err != nil {
		return fmt.Errorf("关闭connection失败: %w", err)
	}
	return nil
}
