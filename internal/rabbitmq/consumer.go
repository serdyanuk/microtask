package rabbitmq

import (
	"github.com/serdyanuk/microtask/config"
	"github.com/streadway/amqp"
)

type ProcessingConsumer struct {
	conn *amqp.Connection
	cfg  *config.Rabbitmq
}

func NewProcessingConsumer(cfg *config.Rabbitmq) (*ProcessingConsumer, error) {
	conn, err := connect(cfg)
	if err != nil {
		return nil, err
	}
	return &ProcessingConsumer{
		conn: conn,
		cfg:  cfg,
	}, nil
}

func (c *ProcessingConsumer) CreateChannel() (ch *amqp.Channel, q string, err error) {
	ch, err = c.conn.Channel()
	if err != nil {
		return nil, "", err
	}
	queue, err := ch.QueueDeclare(
		c.cfg.QueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, "", err
	}
	return ch, queue.Name, nil
}

func (c *ProcessingConsumer) Consume(ch *amqp.Channel, queueName string) (<-chan amqp.Delivery, error) {
	msgs, err := ch.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return nil, err
	}

	return msgs, nil
}
