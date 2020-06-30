package rmq

import (
	"github.com/streadway/amqp"
)

type Consumer struct {
	*RMQ
}

func NewConsumer(uri, queue string) *Consumer {
	return &Consumer{
		New(uri, queue),
	}
}

func (c *Consumer) Consume() (<-chan amqp.Delivery, error) {
	return c.channel.Consume(
		c.queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
}
