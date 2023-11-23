using Telegram.Bot.Types;

namespace ReminderBot.Bot.Types;

public class Command(Message message, string name)
{
    public Message Message { get; private set; } = message;
    public string Name { get; private set; } = name;
}
