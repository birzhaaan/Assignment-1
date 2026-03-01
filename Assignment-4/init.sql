DROP TABLE IF EXISTS users;

CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       name TEXT NOT NULL,
                       email TEXT UNIQUE NOT NULL
);

INSERT INTO users (name, email)
VALUES ('Test User', 'test@example.com');