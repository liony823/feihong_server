-- +migrate Up

ALTER TABLE `app_config`
    CHANGE COLUMN sigle_ip_register_limit_in_12hour sigle_ip_register_limit_in12hour smallint not null default 0 COMMENT '单IP12小时注册限制数, 0为不限制';




