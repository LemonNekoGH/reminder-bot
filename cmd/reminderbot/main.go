package main

import (
	"github.com/lemonnekogh/reminderbot/config"
	"github.com/lemonnekogh/reminderbot/internal/bot"
	"github.com/lemonnekogh/reminderbot/internal/database"
	"github.com/lemonnekogh/reminderbot/internal/logger"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		fx.Provide(config.NewConfig()),
		fx.Provide(logger.NewLoggerService()),
		fx.Provide(database.NewDatabaseService()),
		fx.Provide(bot.NewTelegramBot()),
		fx.Invoke(bot.Run()),
	)
	app.Run()
}
