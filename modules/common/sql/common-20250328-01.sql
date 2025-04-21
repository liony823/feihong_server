-- +migrate Up


create table `sys_menu` (
    `id`           integer     not null primary key AUTO_INCREMENT,
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
INSERT INTO `sys_menu` (`key`, `path`, `name`, `icon`, `sort`, `parent_key`, `layout`, `hidden_in_menu`, `redirect`, `component`, `created_at`, `updated_at`) VALUES('dashboard', '/dashboard', '仪表盘', 'dashboard', 0, '', 1, 0, '', '', NOW(), NOW());
INSERT INTO `sys_menu` (`key`, `path`, `name`, `icon`, `sort`, `parent_key`, `layout`, `hidden_in_menu`, `redirect`, `component`, `created_at`, `updated_at`) VALUES('dashboard-analysis', '/dashboard/analysis', '概览页', '', 0, 'dashboard', 0, 0, '', './dashboard/analysis', NOW(), NOW());
INSERT INTO `sys_menu` (`key`, `path`, `name`, `icon`, `sort`, `parent_key`, `layout`, `hidden_in_menu`, `redirect`, `component`, `created_at`, `updated_at`) VALUES('dashboard-monitor', '/dashboard/monitor', '监控页', '', 0, 'dashboard', 0, 0, '', './dashboard/monitor', NOW(), NOW());
INSERT INTO `sys_menu` (`key`, `path`, `name`, `icon`, `sort`, `parent_key`, `layout`, `hidden_in_menu`, `redirect`, `component`, `created_at`, `updated_at`) VALUES('system', '/system', '系统管理', 'crown', 10, '', 1, 0, '', '', NOW(), NOW());
INSERT INTO `sys_menu` (`key`, `path`, `name`, `icon`, `sort`, `parent_key`, `layout`, `hidden_in_menu`, `redirect`, `component`, `created_at`, `updated_at`) VALUES('system-menu', '/system/menu', '菜单管理', '', 2, 'system', 0, 0, '', './system/menu', NOW(), NOW());
INSERT INTO `sys_menu` (`key`, `path`, `name`, `icon`, `sort`, `parent_key`, `layout`, `hidden_in_menu`, `redirect`, `component`, `created_at`, `updated_at`) VALUES('system-admin', '/system/admin', '账户管理', '', 3, 'system', 0, 0, '', './system/admin', NOW(), NOW());
INSERT INTO `sys_menu` (`key`, `path`, `name`, `icon`, `sort`, `parent_key`, `layout`, `hidden_in_menu`, `redirect`, `component`, `created_at`, `updated_at`) VALUES('chat-client', '/chat/client', '项目设置', '', 4, 'chat', 0, 0, '', './chat/client', NOW(), NOW());
INSERT INTO `sys_menu` (`key`, `path`, `name`, `icon`, `sort`, `parent_key`, `layout`, `hidden_in_menu`, `redirect`, `component`, `created_at`, `updated_at`) VALUES('chat-bucket', '/chat/bucket', '存储桶配置', '', 6, 'chat', 0, 0, '', './chat/bucket', NOW(), NOW());
INSERT INTO `sys_menu` (`key`, `path`, `name`, `icon`, `sort`, `parent_key`, `layout`, `hidden_in_menu`, `redirect`, `component`, `created_at`, `updated_at`) VALUES('chat-sign', '/chat/sign', '签到配置', '', 7, 'chat', 0, 0, '', './chat/sign', NOW(), NOW());
INSERT INTO `sys_menu` (`key`, `path`, `name`, `icon`, `sort`, `parent_key`, `layout`, `hidden_in_menu`, `redirect`, `component`, `created_at`, `updated_at`) VALUES('chat-ua_pp', '/chat/ua_pp', '隐私政策', '', 8, 'chat', 0, 0, '', './chat/ua_pp', NOW(), NOW());
INSERT INTO `sys_menu` (`key`, `path`, `name`, `icon`, `sort`, `parent_key`, `layout`, `hidden_in_menu`, `redirect`, `component`, `created_at`, `updated_at`) VALUES('chat-log', '/system/log', '操作日志', '', 0, 'system', 0, 0, '', './system/log', NOW(), NOW());
INSERT INTO `sys_menu` (`key`, `path`, `name`, `icon`, `sort`, `parent_key`, `layout`, `hidden_in_menu`, `redirect`, `component`, `created_at`, `updated_at`) VALUES('chat-applet', '/chat/applet', '小程序管理', '', 3, 'chat', 0, 0, '', './chat/applet', NOW(), NOW());
INSERT INTO `sys_menu` (`key`, `path`, `name`, `icon`, `sort`, `parent_key`, `layout`, `hidden_in_menu`, `redirect`, `component`, `created_at`, `updated_at`) VALUES('chat', '/chat', '业务系统', 'appstore', 1, '', 0, 0, '', '', NOW(), NOW());
INSERT INTO `sys_menu` (`key`, `path`, `name`, `icon`, `sort`, `parent_key`, `layout`, `hidden_in_menu`, `redirect`, `component`, `created_at`, `updated_at`) VALUES('im', '/im', 'IM系统', 'message', 3, '', 0, 0, '', '', NOW(), NOW());
INSERT INTO `sys_menu` (`key`, `path`, `name`, `icon`, `sort`, `parent_key`, `layout`, `hidden_in_menu`, `redirect`, `component`, `created_at`, `updated_at`) VALUES('im-user', '/im/user', '用户管理', '', 1, 'im', 0, 0, '', './im/user', NOW(), NOW());
INSERT INTO `sys_menu` (`key`, `path`, `name`, `icon`, `sort`, `parent_key`, `layout`, `hidden_in_menu`, `redirect`, `component`, `created_at`, `updated_at`) VALUES('im-group', '/im/group', '群组管理', '', 2, 'im', 0, 0, '', './im/group', NOW(), NOW());
INSERT INTO `sys_menu` (`key`, `path`, `name`, `icon`, `sort`, `parent_key`, `layout`, `hidden_in_menu`, `redirect`, `component`, `created_at`, `updated_at`) VALUES('im-message', '/im/message', '消息管理', '', 2, 'im', 0, 0, '', './im/message', NOW(), NOW());
INSERT INTO `sys_menu` (`key`, `path`, `name`, `icon`, `sort`, `parent_key`, `layout`, `hidden_in_menu`, `redirect`, `component`, `created_at`, `updated_at`) VALUES('chat-module', '/chat/module', '应用模块管理', '', 2, 'chat', 0, 0, '', './chat/module', NOW(), NOW());