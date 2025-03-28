-- +migrate Up


create table `sys_menu` (
    `id`           integer      not null primary key AUTO_INCREMENT,
    `key`          varchar(40) not null default '' comment '菜单key',
    `path`         varchar(40) not null default '' comment '菜单路径',
    `name`         varchar(40) not null default '' comment '菜单名称',
    `icon`         varchar(40) not null default '' comment '菜单图标',
    `sort`         bigint      not null default 0  comment '菜单排序',
    `parent_key`   varchar(40) not null default '' comment '父菜单',
    `layout`       boolean     not null default false comment '是否是布局',
    `hidden_in_menu` boolean   not null default false comment '是否在菜单中隐藏',
    `redirect`     varchar(40) not null default '' comment '重定向路径',
    `component`    varchar(40) not null default '' comment '组件路径',
    `created_at`   timestamp   not null default CURRENT_TIMESTAMP comment '创建时间',
    `updated_at`   timestamp   not null default CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP comment '更新时间',
    unique key `uk_key` (`key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统菜单表';

-- 插入菜单数据
INSERT INTO `sys_menu` (`key`, `path`, `name`, `icon`, `sort`, `parent_key`, `layout`, `hidden_in_menu`, `redirect`, `component`) VALUES
('dashboard', '/dashboard', '仪表盘', 'dashboard', 0, '', true, false, '', ''),
('dashboard-analysis', '/dashboard/analysis', '概览页', '', 0, 'dashboard', false, false, '', './dashboard/analysis'),
('dashboard-monitor', '/dashboard/monitor', '监控页', '', 0, 'dashboard', false, false, '', './dashboard/monitor'),
('system', '/system', '系统管理', 'crown', 10, '', true, false, '', ''),
('system-menu', '/system/menu', '菜单管理', '', 2, 'system', false, false, '', './system/menu'),
('system-admin', '/system/admin', '账户管理', '', 3, 'system', false, false, '', './system/admin'),
('system-client', '/chat/client', '项目设置', '', 4, 'system', false, false, '', './system/client'),
('system-sms', '/chat/sms', '短信配置', '', 5, 'system', false, false, '', './system/sms'),
('system-bucket', '/chat/bucket', '存储桶配置', '', 6, 'system', false, false, '', './system/bucket'),
('system-sign', '/chat/sign', '签到配置', '', 7, 'system', false, false, '', './system/sign'),
('system-ua_pp', '/chat/ua_pp', '隐私政策', '', 8, 'system', false, false, '', './system/ua_pp'),
('system-log', '/system/log', '操作日志', '', 0, 'system', false, false, '', './system/log'),
('chat-applet', '/chat/applet', '小程序管理', '', 3, 'chat', false, false, '', './system/applet'),
('chat', '/chat', '业务系统', 'appstore', 1, '', false, false, '', ''),
('im', '/im', 'IM系统', 'message', 3, '', false, false, '', ''),
('im-user', '/im/user', '用户管理', '', 1, 'im', false, false, '', './im/user'),
('im-group', '/im/group', '群组管理', '', 2, 'im', false, false, '', './im/group'),
('im-message', '/im/message', '消息管理', '', 2, 'im', false, false, '', './im/message');