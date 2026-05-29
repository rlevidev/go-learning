CREATE TABLE urls (
    id SERIAL PRIMARY KEY,
    url_original TEXT NOT NULL,
    short_code VARCHAR(10) UNIQUE NOT NULL,
    access_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
