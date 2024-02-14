use log::{debug, error, info};
use serde::Deserialize;
use std::fs::File;
use std::io::Read;
use std::path::Path;
use teloxide::prelude::*;
use teloxide::types::Message;

#[derive(Deserialize, Debug, PartialEq)]
struct Config {
    bot_token: String,
    pg_data_source: String,
}

fn load_config() -> Config {
    let path = Path::new("config.local.yaml");
    let _display = path.display();

    let mut file = match File::open(&path) {
        Ok(file) => file,
        Err(why) => panic!("could not open config file {_display}: {:?}", why),
    };

    let mut content = String::new();
    match file.read_to_string(&mut content) {
        Ok(_) => (),
        Err(why) => panic!("could not read config file {_display}: {:?}", why),
    };

    let config = match serde_yaml::from_str::<Config>(content.as_str()) {
        Ok(c) => c,
        Err(why) => panic!("could not deserialize config file {_display}: {:?}", why),
    };

    assert_ne!("", config.bot_token, "bot_token can not be empty");
    assert_ne!("", config.pg_data_source, "pg_data_source can not be empty");

    config
}

#[tokio::main]
async fn main() {
    let version = option_env!("CARGO_PKG_VERSION").unwrap_or("0.1.0-alpha.1");

    println!("========================");
    println!("Reminder Bot");
    println!();
    println!("Version: {version}");
    println!("========================");

    pretty_env_logger::init();

    if version.contains("alpha") || version.contains("beta") || version.contains("rc") {
        log::set_max_level(log::LevelFilter::Trace)
    } else {
        log::set_max_level(log::LevelFilter::Info)
    }

    let config = load_config();
    let bot = Bot::new(config.bot_token);
    teloxide::repl(bot, |bot: Bot, msg: Message| async move {
        bot.send_message(
            msg.chat.id,
            format!("your message is: {}", msg.text().unwrap_or("no text")),
        )
        .await?;
        Ok(())
    })
    .await;
}
