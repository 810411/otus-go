package rmq

import "github.com/streadway/amqp"

type Producer struct {
	*RMQ
}

func NewProducer(uri, queue string) *Producer {
	return &Producer{
		New(uri, queue),
	}
}

func (p *Producer) Publish(body []byte) error {
	return p.channel.Publish(
		"",
		p.queueName,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
}
