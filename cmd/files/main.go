package main

import (
	"fmt"
	"log"

	"github.com/serdyanuk/microtask/config"
	"github.com/serdyanuk/microtask/internal/files"
	"github.com/serdyanuk/microtask/pkg/imgmanager"
)

func main() {
	fmt.Println(1)
	cfg := config.Get()
	imgm := imgmanager.New(cfg.Storage.FilesDir)

	s := files.NewApiServer(cfg.FilesService, imgm)
	log.Fatal(s.Run())
}
