using System.ComponentModel.DataAnnotations.Schema;

using Microsoft.EntityFrameworkCore;

namespace ReminderBot.Bot.Persistence.Model;

public enum OperationType
{
    Unknown,
    Create,
    SetName,
    SetContent,
    SetPeriod,
}

[PrimaryKey("MessageId")]
[Table("operations", Schema = "reminder_bot")]
public record class Operations
{
    [Column("id")]
    public int MessageId { get; set; }

    [Column("operator_user")]
    public long OperatorUser { get; set; }

    [Column("completed")]
    public bool Completed { get; set; }

    [Column("remind_item_id", TypeName = "text")]
    public string? RemindItemId { get; set; } = null;

    [Column("type")]
    public OperationType OperationType { get; set; }

    public static async Task<Operations?> GetOperationByMessageIdAsync(int messageId)
    {
        await using var ctx = new PersistenceContext();
        return await ctx.Operations.Where(it => it.MessageId == messageId).FirstOrDefaultAsync();
    }

    public static async Task CreateOperationAsync(int messageId, long userId, OperationType type, string? remindItemId = null)
    {
        await using var ctx = new PersistenceContext();
        await ctx.Operations.AddAsync(new Operations
        {
            MessageId = messageId,
            OperatorUser = userId,
            OperationType = type,
            RemindItemId = remindItemId
        });
        await ctx.SaveChangesAsync();
    }
}
