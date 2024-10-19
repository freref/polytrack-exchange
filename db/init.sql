CREATE TABLE
    IF NOT EXISTS tracks (
        id SERIAL PRIMARY KEY,
        title VARCHAR(255) NOT NULL,
        track_description VARCHAR(255) NOT NULL,
        code TEXT NOT NULL,
        vote integer NOT NULL
    );

CREATE TABLE
    IF NOT EXISTS users (
        username VARCHAR(25) PRIMARY KEY,
        email VARCHAR(50),
        password_hash VARCHAR(60)
    );