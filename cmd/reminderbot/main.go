package main

import (
	"fmt"
	"runtime"

	"github.com/lemonnekogh/reminderbot/config"
	"github.com/lemonnekogh/reminderbot/internal/bot"
	"github.com/lemonnekogh/reminderbot/internal/database"
	"github.com/lemonnekogh/reminderbot/internal/logger"
	"go.uber.org/fx"
)

func printBanner() {
	fmt.Println("=============================")
	fmt.Println("ReminderBot:    " + config.Version)
	fmt.Println()
	fmt.Println("Go              " + runtime.Version())
	fmt.Printf("Platform:       %v_%v\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("Number of CPUs: %v\n", runtime.NumCPU())
	fmt.Println("=============================")
}

func main() {
	printBanner()
	app := fx.New(
		fx.Provide(config.NewConfig()),
		fx.Provide(logger.NewLoggerService()),
		fx.Provide(database.NewDatabaseService()),
		fx.Provide(bot.NewTelegramBot()),
		fx.Invoke(bot.Run()),
	)
	app.Run()
}
