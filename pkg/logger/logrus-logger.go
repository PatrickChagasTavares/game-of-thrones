package logger

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
)

type loggerImpl struct {
	log *logrus.Logger
}

func NewLogrusLogger() Logger {
	log := logrus.New()
	log.SetOutput(os.Stdout)
	log.SetFormatter(&logrus.JSONFormatter{})
	return &loggerImpl{log}
}

func (l *loggerImpl) Fatal(arg ...any) {
	l.log.Fatal(arg...)
}

func (l *loggerImpl) Info(arg ...any) {
	l.log.Info(arg...)
}

func (l *loggerImpl) InfoContext(ctx context.Context, arg ...any) {
	l.log.WithContext(ctx).Info(arg...)
}

func (l *loggerImpl) Error(arg ...any) {
	l.log.Error(arg...)
}

func (l *loggerImpl) ErrorContext(ctx context.Context, arg ...any) {
	l.log.WithContext(ctx).Error(arg...)
}
