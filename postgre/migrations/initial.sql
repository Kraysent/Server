CREATE TABLE users 
(
    id BIGSERIAL PRIMARY KEY,
    login TEXT NOT NULL,
    salt BIGSERIAL NOT NULL,
    password_hash TEXT NOT NULL,
    description TEXT DEFAULT 'test'
);