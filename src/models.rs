use crate::schema::reminders;
use diesel::prelude::*;
use std::ops::DerefMut;
use std::sync::{Arc, Mutex};

// #[derive(Queryable, Selectable)]
// #[diesel(table_name = crate::schema::settings)]
// #[diesel(check_for_backend(diesel::pg::Pg))]
// pub struct Settings {
//     pub id: i64, // chat id
//     pub allow_all_members: bool,
//     pub timezone: String,
// }

#[derive(Queryable, Selectable)]
#[diesel(table_name = crate::schema::reminders)]
#[diesel(check_for_backend(diesel::pg::Pg))]
pub struct Reminders {
    pub id: uuid::Uuid,
    pub chat_id: i64,
    pub owner: i64,
    pub content: String,
    pub cron_exp: String,
}

#[derive(Insertable)]
#[diesel(table_name = reminders)]
pub struct NewReminder<'a> {
    pub chat_id: &'a i64,
    pub owner: &'a i64,
    pub content: &'a String,
    pub cron_exp: &'a String,
}

pub fn save_new_reminder(
    chat_id: i64,
    owner: i64,
    content: String,
    cron_exp: String,
    db: Arc<Mutex<PgConnection>>,
) -> QueryResult<Reminders> {
    let new_reminder = NewReminder {
        chat_id: &chat_id,
        owner: &owner,
        content: &content,
        cron_exp: &cron_exp,
    };
    let mut binding = db.lock().unwrap();
    let db = binding.deref_mut();
    diesel::insert_into(reminders::table)
        .values(&new_reminder)
        .returning(Reminders::as_returning())
        .get_result(db)
}
