package logger

import (
	"go.uber.org/zap"
	"strings"
)

var logger *zap.Logger

func Init() error {
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		return err
	}
	return nil
}

func Debug(messages ...string) {
	logger.Debug(strings.Join(messages, ""))
}

func Info(messages ...string) {
	logger.Info(strings.Join(messages, ""))
}

func Warn(messages ...string) {
	logger.Warn(strings.Join(messages, ""))
}

func Error(messages ...string) {
	logger.Error(strings.Join(messages, ""))
}

func Fatal(messages ...string) {
	logger.Fatal(strings.Join(messages, ""))
}
