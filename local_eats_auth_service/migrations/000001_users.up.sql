CREATE TABLE users
(
    id                  UUID PRIMARY KEY,
    user_name            VARCHAR(50) UNIQUE  NOT NULL,
    email               VARCHAR(100) UNIQUE NOT NULL,
    password_hash       VARCHAR(255)        NOT NULL,
    full_name           VARCHAR(100)        NOT NULL,
    user_type           VARCHAR(20)         NOT NULL,
    address             TEXT,
    phone_number        VARCHAR(20),
    bio                 TEXT,
    specialties         TEXT[],
    years_of_experience INTEGER,
    is_verified         BOOLEAN                  DEFAULT false,
    created_at          TIMESTAMP  DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP  DEFAULT CURRENT_TIMESTAMP,
    deleted_at          TIMESTAMP
)