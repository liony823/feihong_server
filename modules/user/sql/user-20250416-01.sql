-- +migrate Up

ALTER TABLE `user` ADD COLUMN `updated_at_username` TIMESTAMP NULL COMMENT '用户名最后的修改时间';
ALTER TABLE `user` ADD COLUMN `search_by_username` smallint NOT NULL DEFAULT 1 COMMENT '是否可以通过用户名搜索到本人';