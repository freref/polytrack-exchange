CREATE TABLE
    IF NOT EXISTS tracks (
        id SERIAL PRIMARY KEY,
        title VARCHAR(255) NOT NULL,
        track_description VARCHAR(1000) NOT NULL,
        code TEXT NOT NULL
    );

CREATE TABLE
    IF NOT EXISTS users (
        username VARCHAR(25) PRIMARY KEY,
        email VARCHAR(50),
        password_hash VARCHAR(60)
    );

CREATE TABLE
    IF NOT EXISTS votes (
        id SERIAL PRIMARY KEY,
        username VARCHAR(25) NOT NULL,
        track_id INTEGER NOT NULL,
        vote INTEGER NOT NULL,
        CONSTRAINT fk_user FOREIGN KEY (username) REFERENCES users (username) ON DELETE CASCADE,
        CONSTRAINT fk_track FOREIGN KEY (track_id) REFERENCES tracks (id) ON DELETE CASCADE,
        UNIQUE (username, track_id)
    );