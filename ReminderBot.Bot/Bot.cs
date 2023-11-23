using NLog;

using ReminderBot.Bot.Extensions;
using ReminderBot.Bot.Types;

using Telegram.Bot;
using Telegram.Bot.Types;
using Telegram.Bot.Types.Enums;

using NotSupportedException = System.NotSupportedException;

namespace ReminderBot.Bot;

public class Bot(ITelegramBotClient client)
{
    public delegate void CommandReceivedHandler(Command ctx);
    public delegate void CommonMessageReceivedHandler(Message ctx);
    public delegate void CallbackQueryReceivedHandler(CallbackQuery ctx);

    public event CommandReceivedHandler? CommandReceived;
    public event CommonMessageReceivedHandler? CommonMessageReceived;
    public event CallbackQueryReceivedHandler? CallbackQueryReceived;

    private readonly Logger logger = LogManager.GetCurrentClassLogger();

    public void Start(CancellationToken cts)
    {
        client.Timeout = TimeSpan.FromSeconds(10);
        client.TestApiAsync(cts).Wait(cts);

        Task.Run(async () =>
        {
            while (!cts.IsCancellationRequested)
            {
                Update[] updates = await client.GetUpdatesAsync(allowedUpdates: default, cancellationToken: cts);
                this.logger.Debug($"Update received, count: {updates.Length}");
                foreach (Update update in updates)
                    this.HandleUpdate(update);
            }
        }, cts);

        this.logger.Info("Bot started");
    }

    private void HandleUpdate(Update update)
    {
        switch (update.Type)
        {
            case UpdateType.Message when update.IsCommand():
                {
                    this.CommandReceived?.Invoke(new Command(update.Message!, update.Command()));
                    return;
                }
            case UpdateType.Message:
                {
                    this.CommonMessageReceived?.Invoke(update.Message!);
                    return;
                }
            case UpdateType.CallbackQuery:
                {
                    this.CallbackQueryReceived?.Invoke(new CallbackQuery());

                    return;
                }
            default:
                throw new NotSupportedException($"{update.Type} is not support");
        }
    }
}
