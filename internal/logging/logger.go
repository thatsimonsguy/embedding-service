package logging

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func Init() {
	config := zap.NewProductionConfig()

	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.LevelKey = "level"
	config.EncoderConfig.MessageKey = "message"
	config.EncoderConfig.CallerKey = "caller"
	config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	// Optional: Log level from env
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "debug" {
		config.Level.SetLevel(zap.DebugLevel)
	}

	l, err := config.Build()
	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}

	Logger = l
}
