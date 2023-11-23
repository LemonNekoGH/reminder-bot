using Microsoft.EntityFrameworkCore.Migrations;

using Npgsql.EntityFrameworkCore.PostgreSQL.Metadata;

#nullable disable

namespace ReminderBot.Bot.Migrations
{
    /// <inheritdoc />
    public partial class Init : Migration
    {
        /// <inheritdoc />
        protected override void Up(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.EnsureSchema(
                name: "reminder_bot");

            migrationBuilder.CreateTable(
                name: "operations",
                schema: "reminder_bot",
                columns: table => new
                {
                    id = table.Column<int>(type: "integer", nullable: false)
                        .Annotation("Npgsql:ValueGenerationStrategy", NpgsqlValueGenerationStrategy.IdentityByDefaultColumn),
                    operator_user = table.Column<long>(type: "bigint", nullable: false),
                    completed = table.Column<bool>(type: "boolean", nullable: false),
                    remind_item_id = table.Column<string>(type: "text", nullable: true),
                    type = table.Column<int>(type: "integer", nullable: false)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_operations", x => x.id);
                });

            migrationBuilder.CreateTable(
                name: "remind_items",
                schema: "reminder_bot",
                columns: table => new
                {
                    id = table.Column<string>(type: "text", nullable: false),
                    owner = table.Column<long>(type: "bigint", nullable: false),
                    period = table.Column<string>(type: "text", nullable: true),
                    name = table.Column<string>(type: "text", nullable: false),
                    content = table.Column<string>(type: "text", nullable: true),
                    chat_id = table.Column<long>(type: "bigint", nullable: false)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_remind_items", x => x.id);
                });

            migrationBuilder.CreateTable(
                name: "settings",
                schema: "reminder_bot",
                columns: table => new
                {
                    chat_id = table.Column<long>(type: "bigint", nullable: false)
                        .Annotation("Npgsql:ValueGenerationStrategy", NpgsqlValueGenerationStrategy.IdentityByDefaultColumn),
                    allow_all_users = table.Column<bool>(type: "boolean", nullable: false)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_settings", x => x.chat_id);
                });
        }

        /// <inheritdoc />
        protected override void Down(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropTable(
                name: "operations",
                schema: "reminder_bot");

            migrationBuilder.DropTable(
                name: "remind_items",
                schema: "reminder_bot");

            migrationBuilder.DropTable(
                name: "settings",
                schema: "reminder_bot");
        }
    }
}
