package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Message struct {
	bot     *tgbotapi.BotAPI
	err     error
	logger  *botLogger
	chatID  int64
	options []func(msg *tgbotapi.MessageConfig)
}

func (m *Message) WithError(err error) *Message {
	m.options = append(m.options, func(msg *tgbotapi.MessageConfig) {
		m.err = err
	})

	return m
}

func (m *Message) WithReply(replyTo int) *Message {
	m.options = append(m.options, func(msg *tgbotapi.MessageConfig) {
		msg.ReplyToMessageID = replyTo
	})

	return m
}

func (m *Message) WithReplyMarkup(replyMarkup any) *Message {
	m.options = append(m.options, func(msg *tgbotapi.MessageConfig) {
		msg.ReplyMarkup = replyMarkup
	})

	return m
}

func (m *Message) Send(message string) (returnedMsg tgbotapi.Message) {
	if m.err != nil {
		m.logger.Warnf("遇到错误: %v", m.err)
	}

	msg := tgbotapi.NewMessage(m.chatID, message)

	// apply options
	for _, option := range m.options {
		option(&msg)
	}

	returnedMsg, err := m.bot.Send(msg)
	if err != nil {
		m.logger.Warnf("发送「%s」时遇到错误: %v", message, err)
	}

	return
}
