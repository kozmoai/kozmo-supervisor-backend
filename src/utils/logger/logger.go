package logger

import (
	"github.com/caarlos0/env"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.SugaredLogger

type LogConfig struct {
	KOZMO_LOG_LEVEL int `env:"KOZMO_LOG_LEVEL" envDefault:"0"`
}

func init() {
	cfg := &LogConfig{}
	err := env.Parse(cfg)
	if err != nil {
		return
	}

	logConfig := zap.NewProductionConfig()
	logConfig.Level = zap.NewAtomicLevelAt(zapcore.Level(cfg.KOZMO_LOG_LEVEL))
	baseLogger, err := logConfig.Build()
	if err != nil {
		panic("failed to create the default logger: " + err.Error())
	}
	logger = baseLogger.Sugar()
}

func NewSugardLogger() *zap.SugaredLogger {
	return logger
}
