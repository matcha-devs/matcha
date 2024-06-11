-- Copyright (c) 2024 Seoyoung Cho, Ali A. Shah, Carlos Andres Cotera Jurado.
CREATE TABLE IF NOT EXISTS users
(
    id       INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    email    VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created  timestamp DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS openid
(
    id INT PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS accounts
(
    id   INT PRIMARY KEY,
    bank VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS transactions
(
    id       INT PRIMARY KEY,
    category VARCHAR(255) NOT NULL,
    amount   INT          NOT NULL,
    bank     VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS assets
(
    id   INT PRIMARY KEY,
    type VARCHAR(255) NOT NULL
);