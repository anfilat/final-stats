package logger

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/anfilat/final-stats/internal/symo"
)

func New(logLevel string) (symo.Logger, error) {
	log := logrus.New()

	result := logger{
		logger: log,
	}

	if logLevel != "" {
		level, err := logrus.ParseLevel(logLevel)
		if err != nil {
			return result, fmt.Errorf("failed to parse log level: %w", err)
		}
		log.SetLevel(level)
	}

	return result, nil
}

type logger struct {
	logger *logrus.Logger
}

func (l logger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

func (l logger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l logger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l logger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}
