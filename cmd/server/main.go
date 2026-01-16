package main

import (
	"fmt"
	"time"

	"github.com/0xsj/overwatch-pkg/config"
	"github.com/0xsj/overwatch-pkg/errors"
)

type Config struct {
	Env           string        `env:"APP_ENV" default:"development"`
	Host          string        `env:"HOST" default:"0.0.0.0"`
	Port          int           `env:"PORT" default:"8085"`
	Timeout       time.Duration `env:"TIMEOUT" default:"30s"`
	LogLevel      string        `env:"LOG_LEVEL" default:"info"`
	BatchSize     int           `env:"BATCH_SIZE" default:"100"`
	MaxRetries    int           `env:"MAX_RETRIES" default:"3"`
	FlushInterval time.Duration `env:"FLUSH_INTERVAL" default:"10s"`
}

func main() {
	var cfg Config
	if err := config.Load(&cfg); err != nil {
		panic(errors.Internal(err).WithOperation("Ingest.LoadConfig"))
	}

	fmt.Println("overwatch-ingest")
	fmt.Printf("Environment: %s\n", cfg.Env)
	fmt.Printf("Listening on: %s:%d\n", cfg.Host, cfg.Port)
	fmt.Printf("Log Level: %s\n", cfg.LogLevel)
	fmt.Printf("Batch Size: %d\n", cfg.BatchSize)
	fmt.Printf("Max Retries: %d\n", cfg.MaxRetries)
	fmt.Printf("Flush Interval: %s\n", cfg.FlushInterval)
}
