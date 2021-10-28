package config

import (
	"log"
	"sync"

	"github.com/spf13/viper"
)

var (
	config Config
	once   sync.Once
)

type Config struct {
	FilesService
	ProcessingService
	Rabbitmq
	Storage
}

// FilesService config
type FilesService struct {
	Port          string
	FileSizeLimit string
}

// ProcessingService config
type ProcessingService struct {
}

// Rabbitmq config
type Rabbitmq struct {
	URL string
}

// Storage config
type Storage struct {
	FilesDir string
}

func Get() *Config {
	once.Do(func() {
		viper.SetConfigName("config")
		viper.AddConfigPath("./config")
		viper.SetConfigType("yaml")

		if err := viper.ReadInConfig(); err != nil {
			log.Fatal(err)
		}
		if err := viper.Unmarshal(&config); err != nil {
			log.Fatal(err)
		}
	})

	return &config
}
