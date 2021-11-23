package logger

import (
	"go.uber.org/zap"
)

type Logger struct {
	*zap.SugaredLogger
}

func New(structuredLogging bool) Logger {
	if structuredLogging {
		logger, _ := zap.NewProduction()
		return Logger{logger.Sugar()}
	} else {
		logger, _ := zap.NewDevelopment()
		return Logger{logger.Sugar()}
	}
}
