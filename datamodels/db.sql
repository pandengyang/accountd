CREATE DATABASE `demo` CHARACTER SET utf8mb4;
GRANT ALL on `demo`.* TO `demo`@`localhost` IDENTIFIED BY 'demo';

DROP TABLE IF EXISTS `account`;
CREATE TABLE `account` (
	`id` BIGINT(11) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '用户 ID',
	`nickname` VARCHAR(64) NOT NULL COMMENT '昵称',
	`phone` VARCHAR(128) NOT NULL COMMENT '手机',
	`password` VARCHAR(128) NOT NULL COMMENT '密码',
	`salt` VARCHAR(128) DEFAULT '' COMMENT '盐',
	`state` ENUM ('A', 'F') NOT NULL DEFAULT 'A' COMMENT '激活/冻结',
	`created_at` BIGINT(20) NOT NULL COMMENT '创建时间',
	PRIMARY KEY(`id`),
	UNIQUE KEY (`nickname`),
	UNIQUE KEY (`phone`)
) ENGINE=INNODB DEFAULT CHARSET=utf8 COMMENT '账户信息';

CREATE TABLE `refresh_token` (
	`token` VARCHAR(64) NOT NULL COMMENT 'Token',
	`expires_at` BIGINT(20) NOT NULL COMMENT '过期时间',
	UNIQUE KEY (`token`)
) ENGINE=INNODB DEFAULT CHARSET=utf8 COMMENT '刷新 Token';
