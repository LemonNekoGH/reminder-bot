package bot

import (
	"context"

	"github.com/lemonnekogh/reminderbot/ent/settings"
	"github.com/lemonnekogh/reminderbot/pkg/telegram"
)

func checkIsAllowUser(c *telegram.Context) bool {
	chat := c.Chat()
	// 不是群组，不需要进行权限控制
	if !chat.IsGroup() {
		return true
	}

	s, err := c.GetDatabaseClient().Settings.Query().
		Where(settings.ChatID(chat.ID)).
		First(context.Background())

	if err != nil {
		_ = c.NewMessage(chat.ID).WithError(err).Send("出现了一些错误")
		return false
	}

	if !s.AllowAllUser {
		// 不允许所有人设置提醒项，检查是否来自管理员
		ok := checkIsFromAdmin(c)
		return ok
	}

	return true
}

func checkIsFromAdmin(c *telegram.Context) bool {
	chat := c.Chat()
	if !chat.IsGroup() {
		return true
	}

	isFromAdmin, err := c.IsFromAdmin()
	if err != nil {
		_ = c.NewMessage(chat.ID).WithError(err).Send("出现了一些错误")
		return false
	}

	if !isFromAdmin {
		_ = c.NewMessage(chat.ID).Send("无权使用")
		return false
	}

	return true
}
