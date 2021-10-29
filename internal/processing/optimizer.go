package processing

import (
	"encoding/json"

	"github.com/serdyanuk/microtask/config"
	"github.com/serdyanuk/microtask/internal/rabbitmq"
	"github.com/serdyanuk/microtask/pkg/imgmanager"
	"github.com/serdyanuk/microtask/pkg/logger"
)

type Optimizer struct {
	cfg      config.ProcessingService
	consumer *rabbitmq.ProcessingConsumer
	resizer  imgmanager.Resizer
	logger   *logger.Logger
}

func NewOptimizer(cfg config.ProcessingService, consumer *rabbitmq.ProcessingConsumer, resizer imgmanager.Resizer, logger *logger.Logger) *Optimizer {
	return &Optimizer{
		cfg:      cfg,
		consumer: consumer,
		resizer:  resizer,
		logger:   logger,
	}
}

func (o *Optimizer) Run() error {
	ch, queueName, err := o.consumer.CreateChannel()
	if err != nil {
		return nil
	}
	defer ch.Close()

	msgs, err := o.consumer.Consume(ch, queueName)
	if err != nil {
		return err
	}
	forever := make(chan bool)
	go func() {
		for m := range msgs {
			msg := imgmanager.ImageStat{}
			err = json.Unmarshal(m.Body, &msg)
			if err != nil {
				return
			}

			o.logger.Infof("receive message: id=%s x=%d y=%d size=%d", msg.ID, msg.Width, msg.Height, msg.Size)

			stat, err := o.resizer.LoadAndResize(msg.ID, o.cfg.ResizePower)
			if err != nil {
				o.logger.Error(err)
				continue
			}

			o.logger.Infof("image resizing success: id=%s x=%d y=%d size=%d", stat.ID, stat.Width, stat.Height, stat.Size)
		}
	}()
	<-forever
	return nil
}
