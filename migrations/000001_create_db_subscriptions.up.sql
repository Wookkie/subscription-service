CREATE TABLE subscriptions (
    id TEXT PRIMARY KEY,
    service_name TEXT NOT NULL,
    price INT NOT NULL,
    user_id TEXT NOT NULL,
    start_date TIMESTAMP,
    end_date TIMESTAMP
);
