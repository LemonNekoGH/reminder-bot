using System.Text.Json;

using Telegram.Bot;
using Telegram.Bot.Types;
using Telegram.Bot.Types.Enums;

namespace ReminderBot.Bot.Extensions;

public static class MessageExtension
{
    public static string Text(this Message message) => message.Text ?? "";

    public static long FromUserId(this Message message) => message.From?.Id ?? 0;

    public static string FromUserName(this Message message) => message.From == null ? "unknown" : message.From.Username ?? "unknown";

    public static async Task<bool> IsFromAdminAsync(this Message message, ITelegramBotClient bot)
    {
        if (message.FromUserId() == 0)
        {
            return false;
        }

        ChatMember[] administrators = await bot.GetChatAdministratorsAsync(message.Chat.Id);
        return administrators.Any(admin => message.FromUserId() == admin.User.Id);

    }
}

public static class CallbackQueryExtension
{
    public static long FromUserId(this CallbackQuery cb) => cb.From.Id;

    public static string FromUserName(this CallbackQuery cb) => cb.From.Username ?? "unknown";

    public static T? GetData<T>(this CallbackQuery cb) => cb.Data == null ? default : JsonSerializer.Deserialize<T>(cb.Data);
}

public static class UpdateExtension
{
    public static bool IsCommand(this Update update) => update.Message is { Text: not null } &&
        update.Message.Text.StartsWith("/");

    public static string Command(this Update update)
    {
        if (update.Message?.Text == null)
        {
            return "";
        }

        string cmd = update.Message.Text;
        cmd = cmd.Split(' ')[0];
        cmd = cmd.Split('@')[0];

        return cmd[1..];
    }
}

public static class ChatExtension
{
    public static bool IsGroup(this Chat c) => c.Type is ChatType.Group or ChatType.Supergroup;
}
