CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8),
    title TEXT,
    comment TEXT,
    repeat VARCHAR(128)
);
CREATE INDEX scheduler_date ON scheduler (date);