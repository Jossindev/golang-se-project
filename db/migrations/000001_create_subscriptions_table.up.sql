CREATE TABLE IF NOT EXISTS subscriptions (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    city VARCHAR(255) NOT NULL,
    frequency VARCHAR(50) NOT NULL,
    confirmed BOOLEAN DEFAULT FALSE,
    token VARCHAR(255),
    last_sent TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    UNIQUE(email, city)
    );

CREATE INDEX IF NOT EXISTS idx_subscriptions_token ON subscriptions(token);