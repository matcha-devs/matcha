-- Copyright (c) 2024 Seoyoung Cho, Ali A. Shah, Carlos Andres Cotera Jurado.

CREATE TABLE IF NOT EXISTS users
(
    id         BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    username   VARCHAR(255)        NOT NULL UNIQUE,
    email      VARCHAR(255)        NOT NULL UNIQUE,
    password   VARCHAR(255)        NOT NULL,
    created_on timestamp DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS openid
(
    id         BIGINT(20) UNSIGNED NOT NULL PRIMARY KEY,
    created_on timestamp DEFAULT NOW()
);

-- TODO(@everyone): database design
-- CREATE TABLE IF NOT EXISTS assets
-- (
--     id       INT UNSIGNED PRIMARY KEY,
--     asset_id INT UNSIGNED,
--     user_id  BIGINT(20) UNSIGNED,
--     balance  BIGINT(20) UNSIGNED
-- );
--
-- CREATE TABLE IF NOT EXISTS liquid
-- (
--     id         INT UNSIGNED PRIMARY KEY,
--     name       VARCHAR(255) NOT NULL,
--     created_at timestamp DEFAULT NOW()
-- );
--
-- CREATE TABLE IF NOT EXISTS solid
-- (
--     id         INT UNSIGNED PRIMARY KEY,
--     name       VARCHAR(255) NOT NULL,
--     created_at timestamp DEFAULT NOW()
-- );
--
-- CREATE TABLE IF NOT EXISTS transactions
-- (
--     id         INT UNSIGNED PRIMARY KEY,
--     category   VARCHAR(255) NOT NULL,
--     amount     INT          NOT NULL,
--     bank       VARCHAR(255) NOT NULL,
--     created_at timestamp DEFAULT NOW()
-- );
