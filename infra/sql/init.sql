-- User Profile
CREATE TABLE user_profiles (
    user_id VARCHAR(255) PRIMARY KEY,
    display_name VARCHAR(255) NOT NULL,
    picture_url VARCHAR(255),
    access_token VARCHAR(255) NOT NULL,
    refresh_token VARCHAR(255) NOT NULL
);
