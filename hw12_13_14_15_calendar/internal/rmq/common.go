package rmq

import (
	"fmt"

	"github.com/streadway/amqp"
)

type RMQ struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	uri       string
	queueName string
}

func New(uri, queue string) *RMQ {
	return &RMQ{
		uri:       uri,
		queueName: queue,
	}
}

func (r *RMQ) Connect() (err error) {
	r.conn, err = amqp.Dial(r.uri)
	if err != nil {
		return fmt.Errorf("dial: %s", err)
	}

	r.channel, err = r.conn.Channel()
	if err != nil {
		return fmt.Errorf("channel: %s", err)
	}

	_, err = r.channel.QueueDeclare(
		r.queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("queue declare: %s", err)
	}

	return
}

func (r *RMQ) Close() {
	_ = r.channel.Close()
	_ = r.conn.Close()
}
