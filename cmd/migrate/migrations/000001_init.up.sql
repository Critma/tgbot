CREATE TABLE users (
    telegram_id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ DEFAULT NOW()
);