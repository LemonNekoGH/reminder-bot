-- Your SQL goes here
CREATE TABLE reminders(
    "id" uuid not null primary key default gen_random_uuid(),
    "chat_id" int8 not null,
    "owner" int8 not null,
    "content" varchar not null,
    "cron_exp" varchar not null,
    "create_at" timestamp not null default now(),
    "deleted_at" timestamp
);
