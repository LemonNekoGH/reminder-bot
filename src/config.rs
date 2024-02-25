use std::env;

pub fn load_config() -> (String, String) {
    let bot_token = env::var("TG_BOT_TOKEN").expect("env variable TG_BOT_TOKEN must be set");
    assert!(
        !bot_token.is_empty(),
        "env variable TG_BOT_TOKEN cannot be empty"
    );
    let db_url = env::var("DB_URL").expect("env variable DB_URL must be set");
    assert!(!db_url.is_empty(), "env variable DB_URL cannot be empty");
    (bot_token, db_url)
}
