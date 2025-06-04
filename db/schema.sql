CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE users (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
    email varchar(250) UNIQUE NOT NULL,
    name varchar(50) NOT NULL,
    username varchar(50) UNIQUE NOT NULL,
    password TEXT NOT NULL,
    is_active boolean DEFAULT TRUE,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE images (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
    user_id uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    sent_image_id uuid UNIQUE NOT NULL,
    received_image_id uuid UNIQUE NOT NULL,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE consumers (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
    username varchar(255) UNIQUE,
    custom_id varchar(255) UNIQUE,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE oauth2_credentials (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
    name varchar(255) NOT NULL,
    client_id varchar(255) UNIQUE NOT NULL,
    client_secret varchar(255) NOT NULL,
    redirect_uris text[], -- Array de URLs
    consumer_id uuid NOT NULL REFERENCES consumers (id),
    created_at timestamp DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE oauth2_tokens (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
    access_token varchar(512) UNIQUE NOT NULL,
    refresh_token varchar(512) UNIQUE,
    token_type varchar(50) DEFAULT 'bearer',
    expires_in integer,
    scope text,
    authenticated_userid varchar(255),
    credential_id uuid NOT NULL REFERENCES oauth2_credentials (id),
    created_at bigint
);

CREATE TABLE authorization_codes (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
    code varchar(255) UNIQUE NOT NULL DEFAULT encode(gen_random_bytes(32), 'hex'),
    client_id varchar(255) NOT NULL,
    user_id uuid NOT NULL REFERENCES users (id),
    redirect_uri text NOT NULL,
    scopes text[],
    code_challenge varchar(255),
    code_challenge_method varchar(50),
    is_used boolean DEFAULT FALSE,
    expires_at timestamp NOT NULL,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP
);

