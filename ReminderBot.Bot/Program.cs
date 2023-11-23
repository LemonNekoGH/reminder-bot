using System.Runtime.InteropServices;

using Microsoft.EntityFrameworkCore;

using NLog;
using NLog.Config;
using NLog.Layouts;

using ReminderBot.Bot.Handler;
using ReminderBot.Bot.Persistence;

using Telegram.Bot;

namespace ReminderBot.Bot;

internal static class Program
{
    private static readonly CancellationTokenSource cts = new();

    private static string GetOs()
    {
        if (OperatingSystem.IsMacOS())
            return "Darwin";

        if (OperatingSystem.IsWindows())
            return "Windows";
        return OperatingSystem.IsLinux() ? "Linux" : "Unknown";
    }

    private static void InitialLogger(ISetupLoadConfigurationBuilder builder)
    {
        LogLevel minLevel = LogLevel.Info;
        if (BotVersion.Version.Contains("alpha")
        || BotVersion.Version.Contains("beta")
        || BotVersion.Version.Contains("rc")
        || BotVersion.Version.StartsWith("0."))
            minLevel = LogLevel.Debug;


        builder.ForLogger()
            .FilterMinLevel(minLevel)
            .WriteToConsole(layout: Layout.FromString("${longdate}|${level:uppercase=true}|${logger}|${message:withexception=true}|${event-properties:item=ID}"));
    }

    private static void Main()
    {
        Console.WriteLine("=======================");
        Console.WriteLine("Reminder Bot");
        Console.WriteLine("");
        Console.WriteLine($"Version:          {BotVersion.Version}");
        Console.WriteLine($"Operating System: {GetOs()}");
        Console.WriteLine($"Architecture:     {RuntimeInformation.OSArchitecture}");
        Console.WriteLine($"ProcessorCount:   {Environment.ProcessorCount}");
        Console.WriteLine("=======================");

        Config.LoadConfig();

        LogManager.Setup().LoadConfiguration(InitialLogger);

        var dbContext = new PersistenceContext();
        dbContext.Database.MigrateAsync().Wait();
        dbContext.Dispose();

        var botClient = new TelegramBotClient(Config.BotToken);

        var bot = new Bot(botClient);

        var cmdHandler = new CommandHandler(botClient);

        bot.CommandReceived += cmdHandler.HandleNewRemindAsync;

        bot.Start(cts.Token);

        // suspend this program
        var waitForStop = new TaskCompletionSource<bool>();
        Console.CancelKeyPress += (sender, args) =>
        {
            cts.Cancel();
            args.Cancel = true;

            bot.CommandReceived -= cmdHandler.HandleNewRemindAsync;

            waitForStop.TrySetResult(true);
        };
        waitForStop.Task.Wait();
    }
}
