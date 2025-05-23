DROP TABLE IF EXISTS users;

CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL,
    session_token TEXT NOT NULL,
    csrf_token TEXT NOT NULL,
    provider TEXT NOT NULL
);
