using System.Text.Json.Serialization;

namespace ReminderBot.Bot.Handler;

public sealed record class CallbackQueryData
{
    [JsonPropertyName("chat_id")]
    public int ChatId { get; set; }
    [JsonPropertyName("chat_type")]
    public int ChatType { get; set; }
    [JsonPropertyName("chat_name")]
    public string ChatName { get; set; } = "";
    [JsonPropertyName("message_id")]
    public int MessageId { get; set; }
    [JsonPropertyName("message")]
    public string Message { get; set; } = "";
    [JsonPropertyName("user_id")]
    public int UserId { get; set; }
    [JsonPropertyName("username")]
    public string Username { get; set; } = "";
    [JsonPropertyName("operation_type")]
    public int OperationType { get; set; }
    [JsonPropertyName("remind_item_id")]
    public string RemindItemId { get; set; } = "";
}

