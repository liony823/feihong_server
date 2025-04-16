-- +migrate Up

-- 用户密保表
create table `user_security` (
  id bigint not null primary key auto_increment,
  uid varchar(40) not null default '',
  question varchar(200) not null default '',
  answer varchar(200) not null default '',
  created_at timeStamp     not null DEFAULT CURRENT_TIMESTAMP, -- 创建时间
  updated_at timeStamp     not null DEFAULT CURRENT_TIMESTAMP  -- 更新时间
);

CREATE UNIQUE INDEX `user_security_uid_index` on `user_security` (`uid`);