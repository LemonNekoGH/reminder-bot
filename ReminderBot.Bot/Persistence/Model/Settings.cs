using System.ComponentModel.DataAnnotations.Schema;

using Microsoft.EntityFrameworkCore;

namespace ReminderBot.Bot.Persistence.Model;

[PrimaryKey("ChatId")]
[Table("settings", Schema = "reminder_bot")]
public record Settings
{
    [Column("chat_id")]
    public long ChatId { get; set; }

    [Column("allow_all_users")]
    public bool AllowAllUsers { get; set; }

    public static async Task<bool> IsAllowAllUsersAsync(long chatId)
    {
        await using var ctx = new PersistenceContext();
        Settings? settings = await ctx.Settings.FirstOrDefaultAsync(x => x.ChatId == chatId);
        return settings?.AllowAllUsers ?? false;
    }
}
