package main

import (
	"log"

	"github.com/serdyanuk/microtask/config"
	"github.com/serdyanuk/microtask/internal/processing"
	"github.com/serdyanuk/microtask/internal/rabbitmq"
	"github.com/serdyanuk/microtask/pkg/imgmanager"
)

func main() {
	cfg := config.Get()
	consumer, err := rabbitmq.NewProcessingConsumer(&cfg.Rabbitmq)
	if err != nil {
		log.Fatal(err)
	}

	imgm, err := imgmanager.New(cfg.Storage.FilesDir)
	if err != nil {
		log.Fatal(err)
	}

	o := processing.NewOptimizer(cfg.ProcessingService, consumer, imgm)
	log.Fatal(o.Run())
}
