
DROP TABLE IF EXISTS game_records;

CREATE TABLE game_records (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    p1 INTEGER,
    p2 INTEGER,
    date_recorded DATE,
    moves TEXT,
    PRIMARY KEY (p1, p2, date_recorded),
    FOREIGN KEY (p1) REFERENCES users(id),
    FOREIGN KEY (p2) REFERENCES users(id)
);
