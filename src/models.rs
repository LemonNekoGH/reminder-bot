use crate::schema::reminders;
use diesel::prelude::*;

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
