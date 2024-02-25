mod config;
mod schema;
mod models;

use std::str::FromStr;
use chrono::Utc;
use diesel::prelude::*;
use dotenvy::dotenv;
use models::{get_all_reminders, save_new_reminder};
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

    let (bot_token, db_url) = &config::load_config();

    let _version = option_env!("CARGO_PKG_VERSION").unwrap_or("0.1.0-alpha.1");

    println!("========================");
    println!("Reminder Bot");
    println!();
    println!("Version: {_version}");
    println!("========================");

    let db = Arc::new(Mutex::new(PgConnection::establish(&db_url).unwrap()));

    let db_cloned_2 = db.clone();
    let bot_token_cloned_2 = bot_token.clone();

    tokio::join!(
        create_cron_tasks(bot_token_cloned_2, db_cloned_2),
        Command::repl(
            Bot::new(bot_token),
            move |bot: Bot, msg: Message, cmd: Command| {
                log::trace!("creating command processor");
                process_cmd(bot, msg, cmd, db.clone())
            }
        )
    );
}

// TODO: create task when reminder create
async fn create_cron_tasks(bot_token: String, db: Arc<Mutex<PgConnection>>) {
    log::trace!("creating cron tasks");
    let scheduler = tokio_cron_scheduler::JobScheduler::new().await.unwrap();

    let reminders = get_all_reminders(db).expect("cannot get reminders");
    for reminder in reminders {
        let bot_token_cloned_1 = bot_token.clone();
        let _ = scheduler
            .add(
                tokio_cron_scheduler::Job::new_async(
                    reminder.cron_exp.as_str(),
                    move |uuid, mut l| {
                        let chat_id_cloned_1 = reminder.chat_id;
                        let chat_id_cloned_2 = reminder.chat_id;
                        let content_cloned = reminder.content.clone();
                        let bot_token_cloned_2 = bot_token_cloned_1.clone();
                        Box::pin(async move {
                            let bot = Bot::new(bot_token_cloned_2.clone());
                            // Query the next execution time for this job
                            let next_tick = l.next_tick_for_job(uuid).await;
                            match next_tick {
                                Ok(Some(ts)) => {
                                    let send_message = bot
                                        .send_message(ChatId(chat_id_cloned_1), content_cloned)
                                        .await;
                                    if send_message.is_err() {
                                        log::error!(
                                            "could send notice time for {} in {}",
                                            reminder.owner,
                                            chat_id_cloned_2
                                        )
                                    } else {
                                        log::info!(
                                            "sent notice time for {} in {}, next time is: {:?}",
                                            reminder.owner,
                                            chat_id_cloned_2,
                                            ts
                                        )
                                    }
                                }
                                _ => {
                                    log::error!(
                                        "could not get next notice time for {} in {}",
                                        reminder.owner,
                                        chat_id_cloned_2
                                    )
                                }
                            }
                        })
                    },
                )
                .unwrap(),
            )
            .await;
        log::trace!(
            "initialized reminder for {} in {}",
            reminder.owner,
            reminder.chat_id
        );
    }

    scheduler.start().await.expect("failed to start cron tasks");
    tokio::time::sleep(std::time::Duration::from_secs(100)).await;
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

    let parse_result = match cron::Schedule::from_str(exp) {
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
    for date_time in parse_result.upcoming(Utc).take(5) {
        next_5_runs.push_str(date_time.to_string().as_str());
        next_5_runs.push('\n');
    }

    bot.send_message(msg.chat.id, next_5_runs).await
}
