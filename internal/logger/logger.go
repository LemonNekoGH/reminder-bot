package logger

import (
	"context"

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

func NewLoggerService() func(params LoggerServiceParams) *LoggerService {
	return func(params LoggerServiceParams) *LoggerService {
		logger, _ := zap.NewProduction()

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
