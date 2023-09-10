package telegram

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/lemonnekogh/reminderbot/internal/database"
	"github.com/lemonnekogh/reminderbot/internal/logger"
	"github.com/nekomeowww/xo/exp/channelx"
	"go.uber.org/zap"
)

type TelegramBot struct {
	handlers []UpdateHandler
	bot      *tgbotapi.BotAPI
	logger   *logger.LoggerService
	db       *database.DatabaseService
	puller   *channelx.Puller[tgbotapi.Update]
}

func (b *TelegramBot) Handle(handler ...HandlerFunc) {
	b.handlers = append(b.handlers, UpdateHandler{handlers: handler})
}

func (b *TelegramBot) HandleCommand(cmd string, handler ...HandlerFunc) {
	b.handlers = append(b.handlers, UpdateHandler{command: cmd, handlers: handler})
}

func (b *TelegramBot) Run() {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updateChannel := b.bot.GetUpdatesChan(updateConfig)
	// create puller
	puller := channelx.NewPuller[tgbotapi.Update]()
	puller.WithNotifyChannel(updateChannel)
	puller.WithHandler(b.handleUpdate)
	puller.StartPull(context.Background())
	b.puller = puller

	b.logger.Logger.Info("Telegram bot started")
}

func (b *TelegramBot) Stop(context context.Context) error {
	return b.puller.StopPull(context)
}

func handleUpdate(ctx *Context, update tgbotapi.Update, handler UpdateHandler) {
	for _, h := range handler.handlers {
		ok := h(ctx)
		if !ok {
			break
		}
	}
}

func (b *TelegramBot) handleUpdate(update tgbotapi.Update) {
	msg := update.Message
	if msg == nil {
		b.logger.Logger.Error("空消息", zap.Int("update_id", update.UpdateID))
		return
	}

	ctx := &Context{
		bot:  b,
		item: update,
		db:   b.db,
		logger: &botLogger{
			base:           b.logger.Logger,
			chatID:         msg.Chat.ID,
			chatTitle:      msg.Chat.Title,
			userID:         msg.From.ID,
			messageContent: msg.Text,
		},
	}

	for _, handler := range b.handlers {
		if msg.IsCommand() {
			if handler.command == msg.Command() {
				b.logger.Logger.Info("收到命令：" + msg.Command())
				handleUpdate(ctx, update, handler)
			}

			continue
		}

		if handler.command == "" {
			handleUpdate(ctx, update, handler)
		}
	}
}

func NewTelegramBot(
	bot *tgbotapi.BotAPI,
	logger *logger.LoggerService,
	db *database.DatabaseService,
) *TelegramBot {
	return &TelegramBot{
		bot:    bot,
		logger: logger,
		db:     db,
	}
}
