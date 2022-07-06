CREATE TABLE countries (
    id SERIAL PRIMARY KEY,
    code TEXT NOT NULL
);
CREATE TABLE cities (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    country_id BIGINT REFERENCES countries(id)
);
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    login TEXT NOT NULL UNIQUE,
    salt BIGINT NOT NULL,
    password_hash TEXT NOT NULL,
    description TEXT DEFAULT 'test',
    city_id BIGINT REFERENCES cities(id),
    registration_date TIMESTAMP DEFAULT NOW()
);
CREATE TABLE tokens (
    id SERIAL,
    user_id BIGINT REFERENCES users(id),
    value TEXT NOT NULL UNIQUE,
    start_date TIMESTAMP DEFAULT NOW(),
    expiration_date TIMESTAMP NOT NULL
);