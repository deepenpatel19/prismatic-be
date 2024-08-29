package logger

import (
	"os"

	"github.com/deepenpatel19/prismatic-be/core"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func createLogger() *zap.Logger {
	environment := core.Config.Environment
	var level zapcore.Level
	config := zap.NewProductionEncoderConfig()
	var consoleEncoder zapcore.Encoder
	if environment == "local" || environment == "" {
		config = zap.NewDevelopmentEncoderConfig()
		level = -1 // DEBUG
		consoleEncoder = zapcore.NewConsoleEncoder(config)
	} else {
		level = 0 // INFO
		consoleEncoder = zapcore.NewJSONEncoder(config)
	}
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	defaultLogLevel := zap.NewAtomicLevelAt(level)
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), defaultLogLevel),
	)
	return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
}

func LoggerInit() {
	Logger = createLogger()
	Logger.Info("Logger configured Successfully")
}
