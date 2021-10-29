package main

import (
	"log"

	"github.com/serdyanuk/microtask/config"
	"github.com/serdyanuk/microtask/internal/files"
	"github.com/serdyanuk/microtask/internal/rabbitmq"
	"github.com/serdyanuk/microtask/pkg/imgmanager"
)

func main() {
	cfg := config.Get()
	publisher, err := rabbitmq.NewProcessingPublisher(&cfg.Rabbitmq)
	if err != nil {
		log.Fatal(err)
	}

	imgm, err := imgmanager.New(cfg.Storage.FilesDir)
	if err != nil {
		log.Fatal(err)
	}

	s := files.NewApiServer(cfg.FilesService, imgm, publisher)
	log.Fatal(s.Run())
}
