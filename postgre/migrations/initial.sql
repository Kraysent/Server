CREATE TABLE users 
(
    id BIGSERIAL PRIMARY KEY,
    login TEXT NOT NULL,
    salt TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    description TEXT DEFAULT 'test'
);