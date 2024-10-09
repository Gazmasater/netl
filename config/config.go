package config

import (
	"netl/pkg/logger"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func SetLogLevel() {
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "debug"
	}

	var level zapcore.Level
	switch logLevel {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	default:
		level = zap.DebugLevel
	}

	logger.SetLogger(logger.New(zap.NewAtomicLevelAt(level)))
}
