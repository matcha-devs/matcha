--  Copyright (c) 2024 Seoyoung Cho, Ali A. Shah, Carlos Andres Cotera Jurado.
CREATE TABLE `userdb`.`users` (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL
);

INSERT INTO users (username, email, password) VALUES 
('Ancient One', 'ancientone@gmail.com', 'pw1'),
('Alice', 'alice@example.com', 'pw2'),
('Bob', 'bob@example.com', 'pw3'),
('Charlie', 'charlie@example.com', 'pw4');
