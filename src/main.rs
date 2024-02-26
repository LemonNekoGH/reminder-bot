mod commands;
mod config;
mod models;
mod schema;


use diesel::prelude::*;
use dotenvy::dotenv;
use models::{get_all_reminders};

use std::sync::{Arc, Mutex};
use teloxide::{prelude::*};

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

    let db = Arc::new(Mutex::new(PgConnection::establish(db_url).unwrap()));

    let db_cloned_2 = db.clone();
    let bot_token_cloned_2 = bot_token.clone();

    tokio::join!(
        create_cron_tasks(bot_token_cloned_2, db_cloned_2),
        commands::Command::repl(
            Bot::new(bot_token),
            move |bot: Bot, msg: Message, cmd: commands::Command| {
                log::trace!("creating command processor");
                commands::process_cmd(bot, msg, cmd, db.clone())
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
