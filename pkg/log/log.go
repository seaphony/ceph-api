package log

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

const (
	TraceID    = "trace_id"
	Username   = "usename"
	httpPath   = "http_path"
	grpcMethod = "grpc_method"
	httpMethod = "http_method"
	httpQuery  = "http_query"
)

type Config struct {
	Json  bool   `yaml:"json"`
	Level string `yaml:"level"`
}

func GetLogger(cfg Config) zerolog.Logger {
	logger := CreateLogger(cfg)
	zerolog.DefaultContextLogger = &logger
	return logger
}

func CreateLogger(cfg Config) zerolog.Logger {
	var logger zerolog.Logger
	if cfg.Json {
		logger = zerolog.New(os.Stdout)
	} else {
		logger = zerolog.New(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		})
	}
	lvl, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		lvl = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(lvl)
	l := logger.With().Caller().Timestamp()

	logger = l.Logger()
	return logger
}
