package logger

import (
	"go.uber.org/zap"
)

var logger *zap.Logger

func init() {
	logger, err := zap.NewProduction()
	if err != nil {
		logger.Fatal("initializing log failed", zap.Error(err))
	}
}

func NewLogger() *zap.Logger {
	return logger
}
