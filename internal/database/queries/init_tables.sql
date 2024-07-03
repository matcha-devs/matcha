-- Copyright (c) 2024 Seoyoung Cho, Ali A. Shah, Carlos Andres Cotera Jurado.

CREATE TABLE IF NOT EXISTS users
(
    id         BIGINT(20) UNSIGNED     NOT NULL AUTO_INCREMENT PRIMARY KEY,
    firstname  VARCHAR(255)            NOT NULL,
    middlename VARCHAR(255),
    lastname   VARCHAR(255)            NOT NULL,
    email      VARCHAR(255)            NOT NULL UNIQUE,
    password   VARCHAR(255)            NOT NULL,
    birthdate VARCHAR(255) NOT NULL,
    created_on timestamp DEFAULT NOW() NOT NULL
);

CREATE TABLE IF NOT EXISTS openid
(
    id         BIGINT(20) UNSIGNED NOT NULL PRIMARY KEY,
    created_on timestamp DEFAULT NOW()
);

-- TODO(@everyone): database design
CREATE TABLE IF NOT EXISTS asset_class_aggregations
(
    id                BIGINT(20) UNSIGNED  NOT NULL AUTO_INCREMENT PRIMARY KEY,
    cash              BIGINT(20) DEFAULT 0 NOT NULL,
    stocks            BIGINT(20) DEFAULT 0 NOT NULL,
    credit_card       BIGINT(20) DEFAULT 0 NOT NULL,
    other_loan        BIGINT(20) DEFAULT 0 NOT NULL,
    retirement_cash   BIGINT(20) DEFAULT 0 NOT NULL,
    retirement_stocks BIGINT(20) DEFAULT 0 NOT NULL,
    real_estate       BIGINT(20) DEFAULT 0 NOT NULL,
    other_property    BIGINT(20) DEFAULT 0 NOT NULL
);

CREATE TABLE IF NOT EXISTS transactions
(
    id                   BIGINT(20) UNSIGNED                                                  NOT NULL AUTO_INCREMENT,
    user_id              BIGINT(20) UNSIGNED                                                  NOT NULL,
    financial_account_id INT UNSIGNED                                                         NOT NULL,
    amount               BIGINT(20)                                                           NOT NULL,
    type                 ENUM ('RESTAURANTS', 'BILLS', 'HOUSING', 'GROCERY', 'TRAVEL', 'ETC') NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS financial_accounts
(
    id             INT UNSIGNED                               NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id        BIGINT(20) UNSIGNED                        NOT NULL,
    institution_id INT UNSIGNED                               NOT NULL,
    asset_class    ENUM ('CASH', 'STOCKS', 'CREDIT_CARD', 'OTHER_LOAN', 'RETIREMENT_CASH',
        'RETIREMENT_STOCKS', 'REAL_ESTATE', 'OTHER_PROPERTY') NOT NULL,
    name           VARCHAR(255)                               NOT NULL,
    net_value      BIGINT(20) DEFAULT 0                       NOT NULL
);

CREATE TABLE IF NOT EXISTS institutions
(
    id   INT          NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);