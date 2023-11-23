using NLog;

using ReminderBot.Bot.Extensions;
using ReminderBot.Bot.Types;

using Telegram.Bot;
using Telegram.Bot.Types;
using Telegram.Bot.Types.Enums;

namespace ReminderBot.Bot.Handler;

public abstract class Handler<T>(ITelegramBotClient client)
{
    protected ITelegramBotClient BotClient { get; private set; } = client;

    protected Logger InitLogger(T update)
    {
        Logger logger = LogManager.GetCurrentClassLogger();

        if (update is Message message)
        {
            logger = logger.WithProperties(new Dictionary<string, object>()
            {
                {"ChatID", message.Chat.Id},
                {"ChatType", message.Chat.Type},
                {"Username", message.FromUserName()},
                {"UserID", message.FromUserId()},
                {"MessageID", message.MessageId},
                {"MessageContent", message.Text()}

            });

            return logger;
        }

        if (update is CallbackQuery callbackQuery)
        {
            long messageId = default;
            long chatId = default;
            ChatType chatType = ChatType.Private;
            string messageContent = "";
            if (callbackQuery.Message != null)
            {
                messageId = callbackQuery.Message.MessageId;
                chatId = callbackQuery.Message.Chat.Id;
                chatType = callbackQuery.Message.Chat.Type;
                messageContent = callbackQuery.Message.Text();
            }

            logger = logger.WithProperties(new Dictionary<string, object>()
            {
                {"ChatID", chatId},
                {"ChatType", chatType},
                {"FallbackChatType", callbackQuery.Message == null},
                {"Username", callbackQuery.FromUserName()},
                {"UserID", callbackQuery.FromUserId()},
                {"MessageID", messageId},
                {"CallbackQueryID", callbackQuery.Id},
                {"MessageContent", messageContent}
            });

            return logger;
        }

        if (update is Command cmd)
        {
            logger = logger.WithProperties(new Dictionary<string, object>()
            {
                {"ChatID", cmd.Message.Chat.Id},
                {"ChatType", cmd.Message.Chat.Type},
                {"Username", cmd.Message.FromUserName()},
                {"UserID", cmd.Message.FromUserId()},
                {"MessageID", cmd.Message.MessageId},
                {"MessageContent", cmd.Message.Text()}
            });

            return logger;
        }

        throw new NotSupportedException("Not supported for type: " + typeof(T));
    }
}
