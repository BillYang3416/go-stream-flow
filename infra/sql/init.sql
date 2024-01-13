-- User Profile
CREATE TABLE user_profiles (
    user_id SERIAL PRIMARY KEY,
    display_name VARCHAR(255) NOT NULL,
    picture_url VARCHAR(255)
);

-- Oauth Details
CREATE TABLE oauth_details (
    oauth_id VARCHAR(255) PRIMARY KEY,
    user_id INT NOT NULL REFERENCES user_profiles(user_id),
    provider VARCHAR(255) NOT NULL,
    access_token VARCHAR(255) NOT NULL,
    refresh_token VARCHAR(255) NOT NULL
);

-- Credentials
CREATE TABLE user_credentials (
    credential_id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES user_profiles(user_id),
    username VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL
);

-- User File
CREATE TABLE user_uploaded_files (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    size BIGINT NOT NULL,
    content BYTEA NOT NULL,
    user_id int NOT NULL REFERENCES user_profiles(user_id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    email_sent BOOLEAN NOT NULL,
    email_sent_at TIMESTAMPTZ,
    email_recipient VARCHAR(255),
    error_message TEXT
);
