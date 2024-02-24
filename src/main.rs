mod schema;

use chrono::Utc;
use diesel::prelude::*;
use dotenvy::dotenv;
use reminder_bot::models::save_new_reminder;
use reminder_bot::{establish_connection, parse_cron_exp};
use std::env;
use std::sync::{Arc, Mutex};
use teloxide::{prelude::*, utils::command::BotCommands, RequestError};

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

#[tokio::main]
async fn main() {
    pretty_env_logger::init();
    if dotenv().is_err() {
        log::info!(".env file not exists, skipped loading env variables");
    }

    let _version = option_env!("CARGO_PKG_VERSION").unwrap_or("0.1.0-alpha.1");

    println!("========================");
    println!("Reminder Bot");
    println!();
    println!("Version: {_version}");
    println!("========================");

    let db = Arc::new(Mutex::new(establish_connection()));
    let bot = Bot::new(env::var("TG_BOT_TOKEN").expect("env variable TG_BOT_TOKEN must be set"));
    // `move` is important
    Command::repl(bot, move |bot: Bot, msg: Message, cmd: Command| {
        process_cmd(bot, msg, cmd, db.clone())
    })
    .await
}

async fn process_cmd(
    bot: Bot,
    msg: Message,
    cmd: Command,
    db: Arc<Mutex<PgConnection>>,
) -> Result<(), RequestError> {
    match cmd {
        Command::Help => {
            bot.send_message(msg.chat.id, Command::descriptions().to_string())
                .await?
        }
        Command::NewReminder(arg) => process_new_reminder(bot, msg, arg, db).await?,
        Command::DeleteReminder(_id) => {
            bot.send_message(msg.chat.id, "command not implemented")
                .await?
        }
        Command::ListReminders => {
            bot.send_message(msg.chat.id, "command not implemented")
                .await?
        }
        Command::SetAllowAllMember(_allow) => {
            bot.send_message(msg.chat.id, "command not implemented")
                .await?
        }
        Command::About => {
            bot.send_message(msg.chat.id, "command not implemented")
                .await?
        }
        Command::SetTimezone(_timezone) => {
            bot.send_message(msg.chat.id, "command not implemented")
                .await?
        }
    };

    Ok(())
}

async fn process_new_reminder(
    bot: Bot,
    msg: Message,
    arg: String,
    db: Arc<Mutex<PgConnection>>,
) -> Result<Message, RequestError> {
    let blank_index = arg.find(' ').unwrap_or(0);
    let (c, exp) = arg.split_at(blank_index);

    // TODO: check if administrator
    log::info!(
        "received new reminder request from {}, cron expression '{exp}'",
        msg.chat.id
    );

    log::trace!("parsing cron expression");

    let mut parse_result = match parse_cron_exp(exp, &Utc::now()) {
        Ok(result) => result,
        Err(why) => {
            log::trace!("cron expression not valid, {:?}", why);

            return bot
                .send_message(msg.chat.id, format!("{exp} 不是一个有效的 cron 表达式"))
                .await;
        }
    };

    log::trace!("inserting data to database");

    let result = save_new_reminder(
        msg.chat.id.0,
        msg.chat.id.0,
        String::from(c),
        String::from(exp),
        db,
    );
    if result.is_err() {
        log::warn!(
            "error when insert reminder data: {:?}",
            result.err().unwrap()
        );
        return bot
            .send_message(msg.chat.id, "出现了一些问题，请稍后再试")
            .await;
    }

    log::trace!("get next 5th runs time");
    let mut next_5_runs = "好的，你的提醒项已经保存，将来 5 次提醒时间如下：\n".to_string();
    next_5_runs.push_str(parse_result.to_string().as_str());
    next_5_runs.push('\n');
    for _ in 0..5 {
        parse_result = cron_parser::parse(exp, &parse_result).unwrap();
        next_5_runs.push_str(parse_result.to_string().as_str());
        next_5_runs.push('\n');
    }
    bot.send_message(msg.chat.id, next_5_runs).await
}
