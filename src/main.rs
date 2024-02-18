use chrono::Utc;
use serde::Deserialize;
use std::fmt::format;
use std::fs::File;
use std::io::Read;
use std::path::Path;
use std::str::FromStr;
use teloxide::{prelude::*, utils::command::BotCommands};

#[derive(Deserialize, Debug, PartialEq)]
struct Config {
    bot_token: String,
    pg_data_source: String,
}

#[derive(BotCommands, Clone)]
#[command(
    rename_rule = "snake_case",
    description = "欢迎使用 ReminderBot，以下是可以使用的命令："
)]
enum Command {
    #[command(description = "显示所有可用命令")]
    Help,
    #[command(description = "允许或禁止非管理员使用 ReminderBot")]
    SetAllowAllMember(bool),
    #[command(description = "修改 ReminderBot 使用的时区")]
    SetTimezone(String),
    #[command(description = "新建提醒事项")]
    NewReminder(String),
    #[command(description = "删除提醒事项")]
    DeleteReminder(i32),
    #[command(description = "显示所有提醒项")]
    ListReminders,
    #[command(description = "显示 ReminderBot 相关信息")]
    About,
}

const ENV_NAME_BOT_TOKEN: &str = "TG_BOT_TOKEN";
const ENV_NAME_DB_DATA_SOURCE: &str = "DB_DATA_SOURCE";

fn load_config() -> (String, String) {
    log::trace!("loading configuration from environment variables...");

    if dotenv().is_err() {
        log::info!("not .env file found, skipped");
    }

    let bot_token = env::var(ENV_NAME_BOT_TOKEN)
        .expect(format!("environment variable {ENV_NAME_BOT_TOKEN} not found").as_str());
    assert_ne!(
        "", bot_token,
        "environment variable {ENV_NAME_BOT_TOKEN} cannot be empty"
    );

    let db_data_source = env::var(ENV_NAME_DB_DATA_SOURCE)
        .expect(format!("environment variable {ENV_NAME_DB_DATA_SOURCE} not found").as_str());
    assert_ne!(
        "", db_data_source,
        "environment variable {ENV_NAME_DB_DATA_SOURCE} cannot be empty"
    );

    log::trace!("configuration loaded");

    (bot_token, db_data_source)
}

#[tokio::main]
async fn main() {
    pretty_env_logger::init();

    let config = load_config();
    log::debug!("configuration loaded");

    let _version = option_env!("CARGO_PKG_VERSION").unwrap_or("0.1.0-alpha.1");

    println!("========================");
    println!("Reminder Bot");
    println!();
    println!("Version: {_version}");
    println!("========================");

    let bot = Bot::new(config.bot_token);
    Command::repl(bot, process_cmd).await;

    log::info!("bot started")
}

async fn process_cmd(bot: Bot, msg: Message, cmd: Command) -> ResponseResult<()> {
    match cmd {
        Command::Help => {
            bot.send_message(msg.chat.id, Command::descriptions().to_string())
                .await?;
        }
        Command::NewReminder(arg) => {
            let blank_index = arg.find(" ").unwrap_or(0);
            let (content, cron_exp) = arg.split_at(blank_index);

            // TODO: check if administrator
            log::debug!(
                "received new reminder request from {}, cron expression '{cron_exp}'",
                msg.chat.id
            );

            log::trace!("parsing cron expression");
            let mut parse_result = match cron_parser::parse(cron_exp, &Utc::now()) {
                Ok(result) => result,
                Err(why) => {
                    log::debug!("cron expression not valid, {:?}", why);
                    bot.send_message(
                        msg.chat.id,
                        format!("{cron_exp} 不是一个有效的 cron 表达式"),
                    )
                    .await?;

                    return Ok(());
                }
            };

            log::trace!("get next 5th runs time");
            let mut next_5_runs =
                format!("好的，你的提醒内容是「{content}」，将来 5 次提醒时间如下：\n");
            next_5_runs.push_str(parse_result.to_string().as_str());
            next_5_runs.push_str("\n");
            for _ in 0..5 {
                parse_result = cron_parser::parse(cron_exp, &parse_result).unwrap();
                next_5_runs.push_str(parse_result.to_string().as_str());
                next_5_runs.push_str("\n");
            }
            bot.send_message(msg.chat.id, next_5_runs).await?;
        }
        Command::DeleteReminder(id) => {
            bot.send_message(msg.chat.id, "command not implemented")
                .await?;
        }
        Command::ListReminders => {
            bot.send_message(msg.chat.id, "command not implemented")
                .await?;
        }
        Command::SetAllowAllMember(allow) => {
            bot.send_message(msg.chat.id, "command not implemented")
                .await?;
        }
        Command::About => {
            bot.send_message(msg.chat.id, "command not implemented")
                .await?;
        }
        Command::SetTimezone(timezone) => {
            bot.send_message(msg.chat.id, "command not implemented")
                .await?;
        }
    }

    Ok(())
}
