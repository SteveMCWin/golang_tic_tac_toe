
DROP TABLE IF EXISTS users;

CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL,
    email TEXT NOT NULL,
    avatar_url TEXT NOT NULL,
    provider TEXT NOT NULL,
    games_played INTEGER NOT NULL,
    games_won INTEGER NOT NULL,
    elo INTEGER NOT NULL
);
