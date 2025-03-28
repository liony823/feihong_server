
-- +migrate Up

# 修改两步验证秘钥字段
alter table user modify column two_verify_secret varchar(255) not null default '';