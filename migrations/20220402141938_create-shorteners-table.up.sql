CREATE TABLE IF NOT EXISTS shorteners(
    id uuid,
    url TEXT NOT NULL,
    shorten_url TEXT NOT NULL,
    expired_at TEXT
);
