using System.Runtime.InteropServices;

using ReminderBot.Bot;

internal class Program
{
    private static string GetOS()
    {
        if (OperatingSystem.IsMacOS())
            return "Darwin";
        else if (OperatingSystem.IsWindows())
            return "Windows";
        else if (OperatingSystem.IsLinux())
            return "Linux";

        return "Unknown";
    }

    private static void Main(string[] args)
    {
        Console.WriteLine("=======================");
        Console.WriteLine("Reminder Bot");
        Console.WriteLine();
        Console.WriteLine($"Version:          {BotVersion.Version}");
        Console.WriteLine($"Operating System: {GetOS()}");
        Console.WriteLine($"Architecture:     {RuntimeInformation.OSArchitecture}");
        Console.WriteLine($"ProcessorCount:   {Environment.ProcessorCount}");
        Console.WriteLine("=======================");

        Config.LoadConfig();

        Console.WriteLine("Hello, World!");
    }
}
