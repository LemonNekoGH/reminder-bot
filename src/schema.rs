// @generated automatically by Diesel CLI.

diesel::table! {
    reminders (id) {
        id -> Uuid,
        chat_id -> Int8,
        owner -> Int8,
        content -> Varchar,
        cron_exp -> Varchar,
        create_at -> Timestamp,
        deleted_at -> Nullable<Timestamp>,
    }
}
