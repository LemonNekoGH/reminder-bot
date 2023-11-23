using NLog;

using ReminderBot.Bot.Persistence;

using Telegram.Bot;
using Telegram.Bot.Types;

namespace ReminderBot.Bot.Handler;

public class MessageHandler(ITelegramBotClient client) : Handler<Message>(client)
{
    public async void HandleMessageAsync(Message msg)
    {
        await using var ctx = new PersistenceContext();

    }
}

