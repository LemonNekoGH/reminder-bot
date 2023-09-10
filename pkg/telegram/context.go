package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/lemonnekogh/reminderbot/ent"
	"github.com/lemonnekogh/reminderbot/internal/database"
)

type Context struct {
	bot    *TelegramBot
	item   tgbotapi.Update
	db     *database.DatabaseService
	logger *botLogger
}

// returns raw update info from telegram.
func (c *Context) Raw() tgbotapi.Update {
	return c.item
}

func (c *Context) CallbackQuery() *tgbotapi.CallbackQuery {
	return c.item.CallbackQuery
}

func (c *Context) Chat() *tgbotapi.Chat {
	return c.item.Message.Chat
}

func (c *Context) Sender() *tgbotapi.User {
	return c.item.Message.From
}

func (c *Context) Message() *tgbotapi.Message {
	return c.item.Message
}

func (c *Context) IsFromAdmin() (bool, error) {
	members, err := c.bot.bot.GetChatAdministrators(tgbotapi.ChatAdministratorsConfig{
		ChatConfig: tgbotapi.ChatConfig{
			ChatID: c.item.Message.Chat.ID,
		},
	})

	if err != nil {
		return false, err
	}

	for _, member := range members {
		if member.User.ID == c.item.Message.From.ID {
			return true, nil
		}
	}

	return false, nil
}

func (c *Context) GetDatabaseClient() *ent.Client {
	return c.db.Client
}

func (c *Context) NewMessage(chatID int64) *Message {
	return &Message{
		bot:    c.bot.bot,
		logger: c.logger,
		chatID: chatID,
	}
}

func (c *Context) DeleteMessage(messageID int) error {
	cfg := tgbotapi.NewDeleteMessage(c.Chat().ID, messageID)
	_, err := c.bot.bot.Send(cfg)

	return err
}

func (c *Context) EditMessage(messageID int, text string) error {
	cfg := tgbotapi.NewEditMessageText(c.Chat().ID, messageID, text)
	_, err := c.bot.bot.Send(cfg)

	return err
}

func (c *Context) Logger() *botLogger {
	return c.logger
}
