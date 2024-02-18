use crate::models::{NewReminder, Reminders};
use diesel::query_builder::Query;
use diesel::{prelude::*, Connection, PgConnection, RunQueryDsl};
use std::env;

pub mod models;
pub mod schema;

pub fn establish_connection() -> PgConnection {
    let db_url = env::var("DB_URL").expect("env variable DB_URL must be set");
    PgConnection::establish(&db_url).unwrap_or_else(|_| panic!("error connecting to {}", db_url))
}

pub fn save_new_reminder(
    chat_id: i64,
    owner: i64,
    content: String,
    cron_exp: String,
) -> QueryResult<Reminders> {
    use crate::schema::reminders;
    let new_reminder = NewReminder {
        chat_id: &chat_id,
        owner: &owner,
        content: &content,
        cron_exp: &cron_exp,
    };
    let pg_connection = &mut establish_connection();
    diesel::insert_into(reminders::table)
        .values(&new_reminder)
        .returning(Reminders::as_returning())
        .get_result(pg_connection)
}
