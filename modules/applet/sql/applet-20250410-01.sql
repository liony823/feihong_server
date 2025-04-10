-- +migrate Up

CREATE TABLE `applet_config` (
	`id` INT   NOT NULL    AUTO_INCREMENT  ,
	`display_name` VARCHAR(40)   NOT NULL     COMMENT '小程序名称' ,
	`key` VARCHAR(40)   NOT NULL     COMMENT '唯一值' ,
	`icon` VARCHAR(225)   NOT NULL     COMMENT 'icon' ,
	`link` VARCHAR(225)   NOT NULL     COMMENT '路径' ,
	`priority` INT   NOT NULL     COMMENT '优先级' ,
	`status` SMALLINT   NOT NULL     COMMENT '状态' ,
	`is_default` SMALLINT   NOT NULL     COMMENT '是否是默认' ,
	`created_at` TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP    COMMENT '创建时间' ,
	`updated_at` TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP    COMMENT '更新时间',
	PRIMARY KEY  (`id`,`key`)  
) COMMENT='小程序配置';