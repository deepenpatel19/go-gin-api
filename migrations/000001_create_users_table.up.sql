CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    first_name varchar(50),
    last_name varchar(50),
    email varchar(150) unique not null,
    username varchar(150) unique not null,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NULL,
    is_first_time boolean NOT NULL DEFAULT FALSE
);
