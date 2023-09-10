package telegram

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type botLogger struct {
	base           *zap.Logger
	chatID         int64
	chatTitle      string
	chatType       string
	messageID      int
	messageContent string
	userID         int64
	userName       string
}

func (l *botLogger) log(logLevel zapcore.Level, message string) {
	l.base.Log(
		logLevel,
		message,
		zap.Int64("chat_id", l.chatID),
		zap.String("chat_name", l.chatTitle),
		zap.String("chat_type", l.chatType),
		zap.Int("message_id", l.messageID),
		zap.String("message_content", l.messageContent),
		zap.Int64("user_id", l.userID),
		zap.String("user_name", l.userName),
	)
}

func (l *botLogger) Error(message string) {
	l.log(zap.ErrorLevel, message)
}

func (l *botLogger) Errorf(format string, args ...any) {
	l.log(zap.ErrorLevel, fmt.Sprintf(format, args...))
}

func (l *botLogger) Warn(message string) {
	l.log(zap.WarnLevel, message)
}

func (l *botLogger) Warnf(format string, args ...any) {
	l.log(zap.WarnLevel, fmt.Sprintf(format, args...))
}

func (l *botLogger) Info(message string) {
	l.log(zap.InfoLevel, message)
}

func (l *botLogger) Infof(format string, args ...any) {
	l.log(zap.InfoLevel, fmt.Sprintf(format, args...))
}

func (l *botLogger) Debug(message string) {
	l.log(zap.DebugLevel, message)
}

func (l *botLogger) Debugf(format string, args ...any) {
	l.log(zap.DebugLevel, fmt.Sprintf(format, args...))
}
