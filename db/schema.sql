CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE users (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
    email varchar(250) UNIQUE NOT NULL,
    username varchar(50) UNIQUE NOT NULL,
    password TEXT NOT NULL,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE images (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
    user_id uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    sent_image_id uuid UNIQUE NOT NULL,
    received_image_id uuid UNIQUE NOT NULL,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP
);

