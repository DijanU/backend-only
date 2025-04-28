CREATE TABLE series (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    ranking INT UNIQUE NOT NULL,
    title TEXT NOT NULL,
    status TEXT,
    lws_episodes INT DEFAULT 0,
    t_episodes INT DEFAULT 0
);