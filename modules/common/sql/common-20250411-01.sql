-- +migrate Up

ALTER TABLE `sys_menu`
	DROP PRIMARY KEY,
ADD PRIMARY KEY  (`id`,`key`)  ;