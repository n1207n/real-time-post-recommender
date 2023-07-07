CREATE TABLE posts (
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    body TEXT,
    votes INTEGER,
    timestamp TIMESTAMP
);
