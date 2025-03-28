-- +migrate Up

create table `sys_menu_user` (
    `id`           integer      not null primary key AUTO_INCREMENT,
    `menus`        text  comment '菜单ID列表',
    `uid`          varchar(40)      not null default '' comment '用户ID',
    `created_at`   timestamp   not null default CURRENT_TIMESTAMP comment '创建时间',
    `updated_at`   timestamp   not null default CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP comment '更新时间',
    unique key `uk_menus_uid` (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统菜单用户表';


