
-- +migrate Up

# 添加两步验证秘钥 
alter table user add column two_verify_secret varchar(255);