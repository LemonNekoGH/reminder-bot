package bot

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/lemonnekogh/reminderbot/ent/operations"
	"github.com/lemonnekogh/reminderbot/pkg/telegram"
)

func onNewRemindsReply(c *telegram.Context, text string, replyToMessageID int) {
	chatID := c.Chat().ID
	tx, err := c.GetDatabaseClient().BeginTx(context.Background(), nil)

	if err != nil {
		c.NewMessage(chatID).WithError(err).WithReply(c.Message().MessageID).Send("出现了一些错误")
		return
	}
	// 创建新的提醒项
	remind, err := tx.Reminds.Create().
		SetID(uuid.New()).
		SetChatID(chatID).
		SetName(text).
		SetOwner(c.Sender().ID).
		Save(context.Background())
	if err != nil {
		c.NewMessage(chatID).WithError(err).WithReply(c.Message().MessageID).Send("出现了一些错误")
		return
	}
	// 设置操作项为完成并且成功
	_, err = tx.Operations.Update().
		Where(operations.MessageID(replyToMessageID)).
		SetCompleted(true).
		SetSuccess(true).
		SetRemindID(remind.ID).
		Save(context.Background())
	if err != nil {
		err2 := tx.Rollback()
		if err2 != nil {
			c.Logger().Error("创建新提醒项，回滚时出错：" + err2.Error())
		}

		c.NewMessage(chatID).WithError(err).WithReply(c.Message().MessageID).Send("出现了一些错误")

		return
	}
	// 发送成功消息
	msg := c.NewMessage(chatID).WithReply(c.Message().MessageID).Send(
		fmt.Sprintf("提醒项创建成功\nID: `%s`\n请使用 `/set_content` 设置提醒项内容，使用 `/set_period` 设置提醒间隔，提醒项将会在设置完成后生效", remind.ID),
	)
	if msg.MessageID == 0 {
		err2 := tx.Rollback()
		if err2 != nil {
			c.Logger().Error("创建新提醒项，回滚时出错：" + err2.Error())
		}

		return
	}

	c.Logger().Info("新提醒项创建完成：" + remind.ID.String())
}

func onSetNameReply(c *telegram.Context, text string, replyToMessageID int) {
	log.Fatalf("not implemented")
}

func onSetCronReply(c *telegram.Context, text string, replyToMessageID int) {
	log.Fatalf("not implemented")
}

func onSetContentReply(c *telegram.Context, text string, replyToMessageID int) {
	log.Fatalf("not implemented")
}
