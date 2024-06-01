--  Copyright (c) 2024 Seoyoung Cho and Carlos Andres Cotera Jurado.

CREATE TABLE `userdb`.`user`
(
    `ID`       INT         NOT NULL AUTO_INCREMENT,
    `Name`     VARCHAR(45) NOT NULL,
    `Email`    VARCHAR(45) NOT NULL,
    `Password` VARCHAR(45) NOT NULL,
    PRIMARY KEY (`ID`),
    UNIQUE INDEX `ID_UNIQUE` (`ID` ASC) VISIBLE
);