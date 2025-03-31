package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewZapConfig() zap.Config {
	var logConfig zap.Config
	Level := ""

	level, err := zapcore.ParseLevel(Level)
	if err != nil {
		level = zapcore.InfoLevel
	}
	logConfig.Level = zap.NewAtomicLevelAt(level)

	// Настройка формата вывода времени
	logConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	return logConfig
}
