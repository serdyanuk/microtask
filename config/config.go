package config

import (
	"fmt"
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
	ResizePower    uint8
	WorkerPoolSize int
}

func (cfg ProcessingService) ValidateConfig() error {
	if cfg.WorkerPoolSize <= 0 {
		return fmt.Errorf("processingService: workerPoolSize option must be greater than 0, got %d", cfg.WorkerPoolSize)
	}
	return nil
}

// Rabbitmq config
type Rabbitmq struct {
	Host      string
	User      string
	Password  string
	Addr      string
	QueueName string
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
