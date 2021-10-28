package main

import (
	"log"

	"github.com/serdyanuk/microtask/config"
	"github.com/serdyanuk/microtask/internal/files"
	"github.com/serdyanuk/microtask/pkg/imgmanager"
)

func main() {
	cfg := config.Get()
	imgm, err := imgmanager.New(cfg.Storage.FilesDir)
	if err != nil {
		log.Fatal(err)
	}

	s := files.NewApiServer(cfg.FilesService, imgm)
	log.Fatal(s.Run())
}
