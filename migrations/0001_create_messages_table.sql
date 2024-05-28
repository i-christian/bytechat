-- Create messages table to be stored in postgresql
CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    room TEXT NOT NULL,
    text TEXT NOT NULL,
    user TEXT NOT NULL,
    date TIMESTAMPTZ NOT NULL
);
