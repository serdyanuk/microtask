package main

import (
	"github.com/serdyanuk/microtask/config"
	"github.com/serdyanuk/microtask/internal/files"
	"github.com/serdyanuk/microtask/internal/rabbitmq"
	"github.com/serdyanuk/microtask/pkg/imgmanager"
	"github.com/serdyanuk/microtask/pkg/logger"
)

func main() {
	cfg := config.Get()
	logger := logger.Get()

	publisher, err := rabbitmq.NewProcessingPublisher(&cfg.Rabbitmq, logger)
	if err != nil {
		logger.Fatal(err)
	}

	imgm := imgmanager.New(cfg.Storage.FilesDir)

	s := files.NewApiServer(cfg.FilesService, imgm, publisher, logger)
	s.Run()
}
