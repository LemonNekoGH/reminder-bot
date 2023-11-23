namespace ReminderBot.Bot;

using DotNetEnv;

public static class Config
{
    public static string BotToken { get; private set; } = "";
    public static string DbSourceString { get; private set; } = "";

    private const string botTokenEnvVar = "REMINDER_BOT_TOKEN";
    private const string dbSourceStringEnvVar = "REMINDER_BOT_DATABASE_SOURCE";

    public static void LoadConfig()
    {
        var env = Env.Load(options: LoadOptions.NoEnvVars()).ToDictionary();
        BotToken = env.GetValueOrDefault(botTokenEnvVar, "");
        if (BotToken == "")
        {
            BotToken = Environment.GetEnvironmentVariable(botTokenEnvVar) ??
                throw new EnvVariableNotFoundException("You need provide token for bot.", botTokenEnvVar);
        }

        DbSourceString = env.GetValueOrDefault(dbSourceStringEnvVar, "");
        if (DbSourceString == "")
        {
            DbSourceString = Environment.GetEnvironmentVariable(dbSourceStringEnvVar) ??
                throw new EnvVariableNotFoundException("You need provide database connection url for data persistence.", dbSourceStringEnvVar);
        }
    }
}
