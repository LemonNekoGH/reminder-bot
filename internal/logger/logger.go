package logger

import (
	"context"
	"log"
	"strings"

	"github.com/lemonnekogh/reminderbot/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type LoggerServiceParams struct {
	fx.In
	Lifecycle fx.Lifecycle
}

type LoggerService struct {
	Logger *zap.Logger
}

func shouldEnableDebugLogging() bool {
	return strings.Contains(config.Version, "alpha") ||
		strings.Contains(config.Version, "beta") ||
		strings.Contains(config.Version, "rc") ||
		strings.HasPrefix(config.Version, "0.")
}

func NewLoggerService() func(params LoggerServiceParams) *LoggerService {
	return func(params LoggerServiceParams) *LoggerService {
		var (
			logger *zap.Logger
			err    error
		)

		if shouldEnableDebugLogging() {
			logger, err = zap.NewDevelopment()
			if err != nil {
				log.Fatalln("Failed to create logger: " + err.Error())
			}
		} else {
			logger, err = zap.NewProduction()
			if err != nil {
				log.Fatalln("Failed to create logger: " + err.Error())
			}
		}

		params.Lifecycle.Append(fx.Hook{
			OnStop: func(ctx context.Context) error {
				return logger.Sync()
			},
		})

		return &LoggerService{
			Logger: logger,
		}
	}
}
