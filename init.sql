-- Copyright (c) 2024 Seoyoung Cho, Ali A. Shah, Carlos Andres Cotera Jurado.
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS openid (
    id INT PRIMARY KEY
);

INSERT IGNORE INTO users (username, email, password) VALUES 
('Ancient One', 'ancientone@gmail.com', 'pw1'),
('Alice', 'alice@example.com', 'pw2'),
('Bob', 'bob@example.com', 'pw3'),
('Charlie', 'charlie@example.com', 'pw4');
