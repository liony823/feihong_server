-- +migrate Up

CREATE TABLE `test`.`sign_config` (
	`id` INT   NOT NULL    AUTO_INCREMENT  ,
	`sign_type` VARCHAR(20)   NOT NULL     COMMENT '签到模式' ,
	`random_min` INT   NOT NULL     COMMENT '随机发放最小金额，单位分' ,
	`random_max` INT   NOT NULL     COMMENT '随机发放最大金额，单位分' ,
	`rule` TEXT   NULL     COMMENT '签到规则说明' ,
	`daily_signin` INT   NOT NULL     COMMENT '日签金额' ,
	`continue_signin` JSON   NULL     COMMENT '连签配置' ,
	`create_at` TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP    COMMENT '创建时间' ,
	PRIMARY KEY  (`id`)  
) COMMENT='签到配置';