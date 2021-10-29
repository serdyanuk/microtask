package rabbitmq

import (
	"github.com/serdyanuk/microtask/config"
	"github.com/serdyanuk/microtask/pkg/imgmanager"
	"github.com/streadway/amqp"
)

type ProcessingPublisher struct {
	cfg   *config.Rabbitmq
	conn  *amqp.Connection
	ch    *amqp.Channel
	queue amqp.Queue
}

func NewProcessingPublisher(cfg *config.Rabbitmq) (*ProcessingPublisher, error) {
	conn, err := connect(cfg)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	queue, err := ch.QueueDeclare(
		cfg.QueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}
	return &ProcessingPublisher{
		cfg:   cfg,
		conn:  conn,
		ch:    ch,
		queue: queue,
	}, nil
}

func (p *ProcessingPublisher) Publish(msg *imgmanager.ImageStat) error {
	return p.ch.Publish(
		"",           // exchange
		p.queue.Name, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        msg.MustMarshal(),
		})
}
