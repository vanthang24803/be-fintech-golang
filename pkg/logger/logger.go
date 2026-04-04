package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

func InitLogger() {
	env := os.Getenv("APP_ENV")

	var config zap.Config
	if env == "production" {
		config = zap.NewProductionConfig()
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		// Development config uses colored, human-readable console output
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	config.InitialFields = map[string]any{}

	var err error
	Log, err = config.Build(zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}
}
