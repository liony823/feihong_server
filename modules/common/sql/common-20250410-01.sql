-- +migrate Up


CREATE TABLE `sms_config` (
    `id`  INT   NOT NULL    AUTO_INCREMENT ,
	`key` VARCHAR(255)   NOT NULL      ,
	`options` JSON   NOT NULL     COMMENT '配置信息' ,
    `created_at` TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP    COMMENT '创建时间' ,
    `updated_at` TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP    COMMENT '更新时间',
	PRIMARY KEY  (`id`, `key`)  
) COMMENT='短信商配置';

CREATE TABLE `fs_config` (
    `id`  INT   NOT NULL    AUTO_INCREMENT ,
    `title` VARCHAR(255)   NULL   COMMENT '标题',
	`key` VARCHAR(255)   NOT NULL      ,
	`options` JSON   NOT NULL     COMMENT '配置信息' ,
    `created_at` TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP    COMMENT '创建时间' ,
    `updated_at` TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP    COMMENT '更新时间',
	PRIMARY KEY  (`id`, `key`)  
) COMMENT='资源文件存储配置';

CREATE TABLE `operation_log` (
	`id` INT   NOT NULL    AUTO_INCREMENT  ,
	`uid` VARCHAR(40)   NOT NULL     COMMENT '用户id' ,
	`username` VARCHAR(40)   NOT NULL     COMMENT '用户名称' ,
	`method` VARCHAR(40)   NULL     COMMENT '操作方法' ,
	`path` VARCHAR(255)    NULL     COMMENT '请求路径' ,
	`ip` VARCHAR(255)    NULL     COMMENT '操作ip' ,
	`payload` TEXT    NULL     COMMENT '传递参数' ,
    `response` TEXT    NULL     COMMENT '返回参数',
	`status` SMALLINT   NOT NULL     COMMENT '操作状态',
	`duration` INT   NOT NULL     COMMENT '执行耗时',
    `error_msg` TEXT    NULL     COMMENT '错误信息' ,
	`created_at` TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP    COMMENT '创建时间' ,
    `updated_at` TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP    COMMENT '更新时间',
	PRIMARY KEY  (`id`)  
) COMMENT='系统操作日志';
