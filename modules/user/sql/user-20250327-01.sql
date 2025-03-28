
-- +migrate Up

# 添加两步验证开关 
alter table user add column two_verify_on smallint not null default 0;