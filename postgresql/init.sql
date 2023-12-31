CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS ig_user (
    id uuid DEFAULT uuid_generate_v4(),
    name VARCHAR(255),
    PRIMARY KEY (id)
);
CREATE TABLE IF NOT EXISTS ig_media (
    id uuid DEFAULT uuid_generate_v4(),
    profile_id uuid REFERENCES ig_user(id),
    media_type VARCHAR(20) NOT NULL,
    media_url VARCHAR(255) NOT NULL,
    caption TEXT NOT NULL,
    is_comment_enabled BOOLEAN NOT NULL,
    ig_created_at TIMESTAMP,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    PRIMARY KEY (id)
);
CREATE TABLE IF NOT EXISTS ig_tag (
    id uuid DEFAULT uuid_generate_v4(),
    tag VARCHAR(255) NOT NULL,
    PRIMARY KEY (id)
);
CREATE TABLE IF NOT EXISTS ig_media_tag (
    id uuid DEFAULT uuid_generate_v4(),
    ig_media_id uuid REFERENCES ig_media(id),
    ig_tag_id uuid REFERENCES ig_tag(id),
    PRIMARY KEY (id)
);