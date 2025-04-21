-- +migrate Up
ALTER TABLE `user`
	ADD COLUMN `introduction` VARCHAR(255)   NOT NULL DEFAULT ''    COMMENT '个人介绍'  AFTER `sex`;