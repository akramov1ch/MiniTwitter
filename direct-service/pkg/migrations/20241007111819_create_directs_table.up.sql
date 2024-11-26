CREATE TABLE IF NOT EXISTS directs (
    id SERIAL PRIMARY KEY,
    sender_id INT NOT NULL,
    receiver_id INT NOT NULL,
    tweet_id INT NOT NULL,
    text TEXT NOT NULL,
    media TEXT[] NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
);