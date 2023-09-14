using Microsoft.EntityFrameworkCore;

namespace ReminderBot.Bot.Database;

public partial class PersistenceContext : DbContext
{
    public PersistenceContext() : base() {}
}
