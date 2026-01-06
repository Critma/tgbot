CREATE TABLE reminders (
    id SERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(telegram_id),
    message TEXT NOT NULL,
    scheduled_time TIMESTAMPTZ NOT NULL,
    repeat_interval INTERVAL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);