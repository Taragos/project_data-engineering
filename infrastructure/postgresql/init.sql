CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS ig_profiles (
    id uuid DEFAULT uuid_generate_v4(),
    name VARCHAR(255),
    PRIMARY KEY (id)
);
CREATE TABLE IF NOT EXISTS ig_media (
    id uuid DEFAULT uuid_generate_v4(),
    profile_id uuid REFERENCES ig_profiles(id),
    media_type VARCHAR(20) NOT NULL,
    media_url VARCHAR(255) NOT NULL,
    caption TEXT NOT NULL,
    is_comment_enabled BOOLEAN NOT NULL,
    ig_created_at TIMESTAMP,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    PRIMARY KEY (id)
)