use crate::establish_connection;
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
#[diesel(table_name = crate::schema::reminders)]
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

#[test]
fn test_save_new_reminder() {
    use crate::schema::reminders::dsl::reminders;
    use crate::schema::reminders::id;
    // save
    let db = Arc::new(Mutex::new(establish_connection()));
    let save_result = save_new_reminder(
        1234,
        4321,
        "empty_string".to_string(),
        "* * * * *".to_string(),
        db.clone(),
    );
    assert!(save_result.is_ok());

    let saved = save_result.unwrap();

    // query
    let mut db_binding = db.lock().unwrap();
    let db_unlocked = db_binding.deref_mut();
    let select_result = reminders
        .filter(id.eq(saved.id))
        .select(Reminders::as_select())
        .load(db_unlocked);
    assert!(select_result.is_ok());

    let result_set = select_result.unwrap();
    assert_eq!(1, result_set.len());

    let result = &result_set[0];
    assert_eq!(1234, result.chat_id);
    assert_eq!(4321, result.owner);
    assert_eq!("empty_string", result.content);
    assert_eq!("* * * * *", result.cron_exp);

    // clean up
    let clean_up_result = diesel::delete(reminders.filter(id.eq(saved.id))).execute(db_unlocked);
    assert!(clean_up_result.is_ok());
}
