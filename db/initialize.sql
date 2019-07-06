CREATE USER `attendance_rec`@`%`;
GRANT INSERT,SELECT,UPDATE,DELETE ON `attendance_rec_db`.* TO `attendance_rec`@`%`;

CREATE DATABASE IF NOT EXISTS `attendance_rec_db`;

CREATE TABLE IF NOT EXISTS `attendance_rec_db`.`data` (
	`id`                INT          NOT NULL AUTO_INCREMENT,
	`user_id`           INT          NOT NULL,
	`date`              DATE         NOT NULL,
	`reason`            TEXT         NOT NULL,
	PRIMARY KEY (`id`)
);

CREATE TABLE IF NOT EXISTS `attendance_rec_db`.`users` (
	`id`                INT          NOT NULL AUTO_INCREMENT,
	`user_name`         CHAR(18)     NOT NULL,
	PRIMARY KEY (`id`)
);
