package logger

import "go.uber.org/zap"

var log *zap.Logger

type Logger struct {
	*zap.Logger
}

func InitLogger() {

	logConfig := NewZapConfig()
	logConfig = zap.NewProductionConfig()

	var err error

	if log, err = logConfig.Build(zap.AddCallerSkip(1)); err != nil {
		panic("failed to initialize logger: " + err.Error())
	}

	defer log.Sync()
}

func GetLogger() *zap.Logger {
	if log == nil {
		panic("logger not initialized")
	}
	return log
}
