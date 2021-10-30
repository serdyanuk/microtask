package main

import (
	"github.com/serdyanuk/microtask/config"
	"github.com/serdyanuk/microtask/internal/processing"
	"github.com/serdyanuk/microtask/internal/rabbitmq"
	"github.com/serdyanuk/microtask/pkg/imgmanager"
	"github.com/serdyanuk/microtask/pkg/logger"
)

func main() {
	cfg := config.Get()
	logger := logger.Get()

	consumer, err := rabbitmq.NewProcessingConsumer(&cfg.Rabbitmq, logger)
	if err != nil {
		logger.Fatal(err)
	}

	imgm := imgmanager.New(cfg.Storage.FilesDir)

	o := processing.NewOptimizer(cfg.ProcessingService, consumer, imgm, logger)
	logger.Fatal(o.Run())
}
