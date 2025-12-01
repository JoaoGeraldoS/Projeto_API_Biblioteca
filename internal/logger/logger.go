package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var instance *zap.Logger

func NewLogger(loggerEnv string) *zap.Logger {
	level := zap.NewAtomicLevel()

	var isProduction bool

	if loggerEnv == "production" || loggerEnv == "prod" {
		isProduction = true
	}

	if isProduction {
		level.SetLevel(zap.InfoLevel)
	} else {
		level.SetLevel(zap.DebugLevel)
	}

	cfg := zap.Config{
		Level:            level,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	if isProduction {
		cfg.Encoding = "json"
		cfg.EncoderConfig = zap.NewProductionEncoderConfig()
		cfg.EncoderConfig.TimeKey = "timestamp"
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		cfg.Encoding = "console"
		cfg.EncoderConfig = zap.NewDevelopmentEncoderConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	l, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	instance = l
	return instance

}
