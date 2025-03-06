BEGIN;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Core Movie Tables
CREATE TABLE movies (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    duration INT NOT NULL CHECK (duration > 0),
    watch_url VARCHAR(512) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE genres (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE movie_genres (
    movie_id UUID REFERENCES movies(id) ON DELETE CASCADE,
    genre_id UUID REFERENCES genres(id) ON DELETE CASCADE,
    PRIMARY KEY (movie_id, genre_id)
);

CREATE TABLE artists (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE movie_artists (
    movie_id UUID REFERENCES movies(id) ON DELETE CASCADE,
    artist_id UUID REFERENCES artists(id) ON DELETE CASCADE,
    PRIMARY KEY (movie_id, artist_id)
);

-- User Management
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(20) UNIQUE NOT NULL
);

INSERT INTO roles (name) VALUES ('admin'), ('user');

CREATE TABLE user_roles (
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    role_id INT REFERENCES roles(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, role_id)
);

-- Viewership Tracking
CREATE TABLE views (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    movie_id UUID NOT NULL REFERENCES movies(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Voting System
CREATE TABLE votes (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    movie_id UUID NOT NULL REFERENCES movies(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, movie_id)
);

-- Indexing Strategy
CREATE INDEX idx_movies_title ON movies(title);
CREATE INDEX idx_movies_created ON movies(created_at);
CREATE INDEX idx_views_movie ON views using btree(movie_id);
CREATE INDEX idx_votes_user ON votes using btree(user_id);
CREATE INDEX idx_votes_movie ON votes using btree(movie_id);
CREATE INDEX idx_user_roles ON user_roles using btree(user_id);

COMMIT;