CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name text NOT NULL,
    CONSTRAINT username_unique UNIQUE (name)
);

CREATE TABLE replays (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    file_path text NOT NULL,
    owner INTEGER NOT NULL,
    FOREIGN KEY (owner) REFERENCES users (id)
);
