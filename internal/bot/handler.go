package bot

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/lemonnekogh/reminderbot/ent/operations"
	"github.com/lemonnekogh/reminderbot/ent/reminds"
	"github.com/lemonnekogh/reminderbot/ent/schema"
	"github.com/lemonnekogh/reminderbot/pkg/telegram"
)

func handleSetName(c *telegram.Context) bool {
	chat := c.Chat()

	reminds, err := c.GetDatabaseClient().Reminds.Query().Where(
		reminds.Owner(c.Sender().ID),
	).All(context.Background())
	if err != nil {
		_ = c.NewMessage(chat.ID).WithError(err).Send("出现了一些错误")
		return false
	}

	if len(reminds) == 0 {
		_ = c.NewMessage(chat.ID).Send("你没有设置任何提醒项")
		return false
	}

	inlineKeyboardBtn := make([][]tgbotapi.InlineKeyboardButton, 0)
	for _, remind := range reminds {
		inlineKeyboardBtn = append(inlineKeyboardBtn, []tgbotapi.InlineKeyboardButton{{
			Text: remind.Name,
		}})
	}

	_ = c.NewMessage(chat.ID).WithReplyMarkup(tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: inlineKeyboardBtn,
	}).WithReply(c.Message().MessageID).Send("请选择要修改名称的提醒项")

	return false
}

func handleNewReminds(c *telegram.Context) bool {
	msg := c.NewMessage(c.Chat().ID).Send("请输入提醒项名称")
	if msg.MessageID == 0 {
		return false
	}

	_, err := c.GetDatabaseClient().
		Operations.
		Create().
		SetID(uuid.New()).
		SetOperator(c.Sender().ID).
		SetMessageID(msg.MessageID).
		SetType(int(schema.EnumOperationTypeCreate)).
		Save(context.Background())
	if err != nil {
		c.Logger().Error("创建操作项时失败：" + err.Error())
		// 数据库操作失败，修改消息
		err2 := c.EditMessage(msg.MessageID, "出现了一些问题，请重试")
		if err2 != nil {
			c.Logger().Error("修改消息时失败：" + err2.Error())
		}
	}

	return false
}

func handleSetContent(c *telegram.Context) bool {
	return false
}

func handleSetPeriod(c *telegram.Context) bool {
	return false
}

func handleShowAll(c *telegram.Context) bool {
	return false
}

func handleDelete(c *telegram.Context) bool {
	return false
}

func handleSettings(c *telegram.Context) bool {
	return false
}

func handleCallbackQuery(c *telegram.Context) bool {
	query := c.CallbackQuery()
	if query == nil {
		c.Logger().Info("没有 Callback Query")
		return false
	}

	// 检查点击按钮的人是否与发送消息的人一致
	if query.Message.From.ID != query.From.ID {
		c.Logger().Info("用户在点击消息按钮时被拒绝：无权操作")
		_ = c.NewMessage(c.Chat().ID).Send("无权操作")

		return false
	}

	return true
}

func handleReply(c *telegram.Context) bool {
	// 没有被回复的消息，不处理
	if c.Message().ReplyToMessage == nil {
		c.Logger().Info("没有被回复的消息")
		return false
	}

	// 查询被回复的消息正在进行的操作
	replyToMessageID := c.Message().ReplyToMessage.MessageID
	operation, err := c.GetDatabaseClient().Operations.Query().
		Where(operations.MessageID(replyToMessageID)).
		First(context.Background())

	if err != nil {
		c.Logger().Info("用户回复正在进行的操作时被拒绝：数据库错误 " + err.Error())
		_ = c.NewMessage(c.Chat().ID).WithError(err).Send("出了一些问题，请尝试重新回复上一条消息")

		return false
	}

	if operation.Completed {
		c.Logger().Info("用户回复正在进行的操作时被拒绝：操作已经完成")
		_ = c.NewMessage(c.Chat().ID).Send("此操作已经完成")

		return false
	}

	if operation.Operator != c.Sender().ID {
		c.Logger().Info("用户回复正在进行的操作时被拒绝：无权操作")
		_ = c.NewMessage(c.Chat().ID).Send("无权操作")
	}

	// TODO: 完成操作
	c.Logger().Infof("操作「%s」收到回复：%s", schema.OperationType(operation.Type), c.Message().Text)

	return false
}
