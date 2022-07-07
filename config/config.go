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

// Config contains all configuration.
type Config struct {
	FilesService
	ProcessingService
	Rabbitmq
	Storage
}

// FilesService represents files service config.
type FilesService struct {
	Port          string
	FileSizeLimit string
}

// ProcessingService represents processing service config.
type ProcessingService struct {
	ResizePower    uint8
	WorkerPoolSize int
}

// ValidateConfig is used for validation processing service config.
func (cfg ProcessingService) ValidateConfig() error {
	if cfg.WorkerPoolSize <= 0 {
		return fmt.Errorf("processingService: workerPoolSize option must be greater than 0, got %d", cfg.WorkerPoolSize)
	}
	return nil
}

// Rabbitmq represnets a Rabbitmq config.
type Rabbitmq struct {
	Host      string
	User      string
	Password  string
	Addr      string
	QueueName string
}

// Storage represents a storage config
type Storage struct {
	FilesDir string
}

// Get is used to getting config once.
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
