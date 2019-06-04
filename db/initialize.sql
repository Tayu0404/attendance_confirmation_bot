CREATE USER `attendance_rec`@`%`;
GRANT INSERT,SELECT,UPDATE,DELETE ON `attendance_rec_db`.* TO `attendance_rec`@`%`;

CREATE DATABASE IF NOT EXISTS `attendance_rec_db`;

CREATE TABLE IF NOT EXISTS `attendance_rec_db`.`user` (
    `uuid`              CHAR(18)     NOT NULL,
    `data`              DATE         NOT NULL,
	`reason`            TEXT         NOT NULL,

    PRIMARY KEY ( `uuid` )
);
