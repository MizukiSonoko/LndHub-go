package logger

import (
	"go.uber.org/zap"
)

var logger, _ = zap.NewProduction()

func NewLogger() *zap.Logger {
	return logger
}
