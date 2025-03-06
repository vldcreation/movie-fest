CREATE TABLE genre_views (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    genre_id UUID NOT NULL REFERENCES genres(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_genre_views ON genre_views using btree(genre_id);