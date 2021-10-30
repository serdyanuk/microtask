package rabbitmq

import (
	"github.com/serdyanuk/microtask/config"
	"github.com/serdyanuk/microtask/pkg/logger"
	"github.com/streadway/amqp"
)

type ProcessingConsumer struct {
	conn   *amqp.Connection
	cfg    *config.Rabbitmq
	logger *logger.Logger
}

func NewProcessingConsumer(cfg *config.Rabbitmq, logger *logger.Logger) (*ProcessingConsumer, error) {
	conn, err := connect(cfg, logger)
	if err != nil {
		return nil, err
	}
	return &ProcessingConsumer{
		conn:   conn,
		cfg:    cfg,
		logger: logger,
	}, nil
}

func (c *ProcessingConsumer) CreateChannel() (ch *amqp.Channel, err error) {
	ch, err = c.conn.Channel()
	if err != nil {
		return nil, err
	}
	return ch, nil
}

func (c *ProcessingConsumer) QueueDeclare(ch *amqp.Channel) (queueName string, err error) {
	queue, err := ch.QueueDeclare(
		c.cfg.QueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return "", err
	}
	return queue.Name, nil
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
