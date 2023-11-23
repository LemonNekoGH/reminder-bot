using NLog;

using ReminderBot.Bot.Extensions;
using ReminderBot.Bot.Persistence.Model;
using ReminderBot.Bot.Types;

using Telegram.Bot;
using Telegram.Bot.Types;

namespace ReminderBot.Bot.Handler;

public class CommandHandler(ITelegramBotClient client) : Handler<Command>(client)
{
    public async void HandleNewRemindAsync(Command cmd)
    {
        // if (cmd.Name != "") return;

        Logger logger = this.InitLogger(cmd);
        logger.Info(cmd.Name);

        bool isFromGroup = cmd.Message.Chat.IsGroup();
        bool isAllowAllUsers = await Settings.IsAllowAllUsersAsync(cmd.Message.Chat.Id);
        bool isFromAdmin = await cmd.Message.IsFromAdminAsync(this.BotClient);

        if (isFromGroup && !isAllowAllUsers && !isFromAdmin)
        {
            await this.BotClient.SendTextMessageAsync(cmd.Message.Chat.Id, "无权使用");
            logger.Warn("无权使用");
            return;
        }

        Message message = await this.BotClient.SendTextMessageAsync(cmd.Message.Chat.Id, "请将提醒名称回复至此消息");
        await Operations.CreateOperationAsync(message.MessageId, cmd.Message.FromUserId(), OperationType.Create);
    }

    public async void HandleSetNameAsync(Command cmd)
    {

    }

    public async void HandleSetContentAsync(Command cmd) { }
    public async void HandleSetPeriodAsync(Command cmd) { }
    public async void HandleShowAllAsync(Command cmd) { }
    public async void HandleDeleteAsync(Command cmd)
    {

    }
    public async void HandleHelpAsync(Command cmd) { }
}
