namespace ReminderBot.Bot;

using DotNetEnv;

public class Config
{
    public static string BotToken { get; private set; } = "";
    public static string DbConnectionURL { get; private set; } = "";

    private const string BotTokenEnvVar = "REMINDER_BOT_TOKEN";
    private const string DbConnectionURLEnvVar = "REMINDER_BOT_DATABASE_URL";

    public static void LoadConfig()
    {
        var env = Env.Load(options: LoadOptions.NoEnvVars()).ToDictionary();
        BotToken = env.GetValueOrDefault(BotTokenEnvVar, "");
        if (BotToken == "")
        {
            BotToken = Environment.GetEnvironmentVariable(BotTokenEnvVar) ??
                throw new EnvVariableNotFoundException("You need provide token for bot.", BotTokenEnvVar);
        }

        DbConnectionURL = env.GetValueOrDefault(DbConnectionURLEnvVar, "");
        if (DbConnectionURL == "")
        {
            DbConnectionURL = Environment.GetEnvironmentVariable(DbConnectionURLEnvVar) ??
                throw new EnvVariableNotFoundException("You need provide database connection url for data persistence.", DbConnectionURLEnvVar);
        }
    }
}
