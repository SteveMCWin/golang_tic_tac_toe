
DROP TABLE IF EXISTS game_records;

CREATE TABLE game_records (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    p1 INTEGER,
    p2 INTEGER,
    date_recorded DATE,
    moves TEXT,
    winner TEXT check (winner in ('x', 'o', '_')), -- 'x', 'o' or defs.BOARD_TIE
    FOREIGN KEY (p1) REFERENCES users(id),
    FOREIGN KEY (p2) REFERENCES users(id)
);
