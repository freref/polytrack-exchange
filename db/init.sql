CREATE TABLE
    IF NOT EXISTS track (
        id SERIAL PRIMARY KEY,
        title VARCHAR(255) NOT NULL,
        code TEXT NOT NULL,
        upvote integer NOT NULL,
        downvote integer NOT NULL
    );