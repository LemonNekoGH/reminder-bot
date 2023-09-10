package bot

import (
	"errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/lemonnekogh/reminderbot/config"
	"github.com/lemonnekogh/reminderbot/internal/database"
	"github.com/lemonnekogh/reminderbot/internal/logger"
	"github.com/lemonnekogh/reminderbot/pkg/telegram"
	"go.uber.org/fx"
)

type TelegramBotParams struct {
	fx.In
	Lifecycle fx.Lifecycle

	Config   *config.Config
	DBClient *database.DatabaseService
	Logger   *logger.LoggerService
}

type TelegramBotService struct {
	bot *telegram.TelegramBot
}

func NewTelegramBot() func(TelegramBotParams) (*TelegramBotService, error) {
	return func(tbp TelegramBotParams) (*TelegramBotService, error) {
		if tbp.Config.TelegramBotToken == "" {
			return nil, errors.New("telegram bot token must bot be empty")
		}

		telegramBotAPI, err := tgbotapi.NewBotAPI(tbp.Config.TelegramBotToken)
		if err != nil {
			return nil, err
		}

		botService := telegram.NewTelegramBot(
			telegramBotAPI,
			tbp.Logger,
			tbp.DBClient,
		)

		botService.HandleCommand("start", checkIsAllowUser, handleNewReminds)
		botService.HandleCommand("set_name", handleSetName)
		botService.HandleCommand("set_content", handleSetContent)
		botService.HandleCommand("set_period", handleSetPeriod)
		botService.HandleCommand("show_all", checkIsFromAdmin, handleShowAll)
		botService.HandleCommand("delete", handleDelete)
		botService.HandleCommand("settings", checkIsFromAdmin, handleSettings)
		botService.Handle(handleCallbackQuery)
		botService.Handle(handleReply)

		tbp.Lifecycle.Append(fx.Hook{
			OnStop: botService.Stop,
		})

		return &TelegramBotService{bot: botService}, nil
	}
}

func Run() func(*TelegramBotService) {
	return func(tbs *TelegramBotService) {
		tbs.bot.Run()
	}
}
