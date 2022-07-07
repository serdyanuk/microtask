package processing

import (
	"encoding/json"

	"github.com/serdyanuk/microtask/config"
	"github.com/serdyanuk/microtask/internal/rabbitmq"
	"github.com/serdyanuk/microtask/pkg/imgmanager"
	"github.com/serdyanuk/microtask/pkg/logger"
	"github.com/streadway/amqp"
)

// Optimizer represents processing service.
type Optimizer struct {
	cfg      config.ProcessingService
	consumer *rabbitmq.ProcessingConsumer
	resizer  imgmanager.Resizer
	logger   *logger.Logger
}

// NewOptimizer is used to create processing service.
func NewOptimizer(cfg config.ProcessingService, consumer *rabbitmq.ProcessingConsumer, resizer imgmanager.Resizer, logger *logger.Logger) *Optimizer {
	return &Optimizer{
		cfg:      cfg,
		consumer: consumer,
		resizer:  resizer,
		logger:   logger,
	}
}

func (o *Optimizer) Run() error {
	ch, err := o.consumer.CreateChannel()
	if err != nil {
		return err
	}
	defer ch.Close()

	queueName, err := o.consumer.QueueDeclare(ch)
	if err != nil {
		return err
	}

	msgs, err := o.consumer.Consume(ch, queueName)
	if err != nil {
		return err
	}

	for i := 0; i < o.cfg.WorkerPoolSize; i++ {
		go o.reader(msgs)
	}

	o.logger.Info("Processing service is ready to receive messages")

	forever := make(chan bool)
	<-forever
	return nil
}

func (o *Optimizer) reader(msgs <-chan amqp.Delivery) {
	for m := range msgs {
		msg := imgmanager.ImageStat{}
		err := json.Unmarshal(m.Body, &msg)
		if err != nil {
			o.logger.Error(err)
			continue
		}

		o.logger.Infof("receive message: %s", msg)

		stat, err := o.resizer.LoadAndResize(msg.ID, o.cfg.ResizePower)
		if err != nil {
			o.logger.Error(err)
			continue
		}

		o.logger.Infof("image resizing success: %s", stat)
	}
}
