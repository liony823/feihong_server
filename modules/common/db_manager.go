package common

import (
	"strings"

	"github.com/TangSengDaoDao/TangSengDaoDaoServerLib/config"
	dbs "github.com/TangSengDaoDao/TangSengDaoDaoServerLib/pkg/db"
	"github.com/TangSengDaoDao/TangSengDaoDaoServerLib/pkg/util"
	"github.com/gocraft/dbr/v2"
)

type managerDB struct {
	session *dbr.Session
	ctx     *config.Context
}

func newManagerDB(ctx *config.Context) *managerDB {
	return &managerDB{
		session: ctx.DB(),
		ctx:     ctx,
	}
}

// 添加版本升级
func (d *managerDB) insertAppVersion(m *appVersionModel) (int64, error) {
	result, err := d.session.InsertInto("app_version").Columns(util.AttrToUnderscore(m)...).Record(m).Exec()
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return id, err
}

// 查询某个系统的最新版本
func (d *managerDB) queryNewVersion(os string) (*appVersionModel, error) {
	var model *appVersionModel
	_, err := d.session.Select("*").From("app_version").Where("os=?", os).OrderDir("created_at", false).Limit(1).Load(&model)
	return model, err
}

// 查询版本升级列表
func (d *managerDB) queryAppVersionListWithPage(pageSize, page uint64) ([]*appVersionModel, error) {
	var models []*appVersionModel
	_, err := d.session.Select("*").From("app_version").Offset((page-1)*pageSize).Limit(pageSize).OrderDir("updated_at", false).Load(&models)
	return models, err
}

// 模糊查询用户数量
func (d *managerDB) queryCount() (int64, error) {
	var count int64
	_, err := d.session.Select("count(*)").From("app_version").Load(&count)
	return count, err
}

// 查询所有背景图片
func (d *managerDB) queryChatBgs() ([]*chatBgModel, error) {
	var models []*chatBgModel
	_, err := d.session.Select("*").From("chat_bg").Load(&models)
	return models, err
}

// 查询app模块
func (d *managerDB) queryAppModule() ([]*appModuleModel, error) {
	var list []*appModuleModel
	_, err := d.session.Select("*").From("app_module").OrderDir("created_at", true).Load(&list)
	return list, err
}

// 查询某个app模块
func (d *managerDB) queryAppModuleWithSid(sid string) (*appModuleModel, error) {
	var m *appModuleModel
	_, err := d.session.Select("*").From("app_module").Where("sid=?", sid).Load(&m)
	return m, err
}

// 新增app模块
func (d *managerDB) insertAppModule(m *appModuleModel) (int64, error) {
	result, err := d.session.InsertInto("app_module").Columns(util.AttrToUnderscore(m)...).Record(m).Exec()
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return id, err
}

// 修改app模块
func (d *managerDB) updateAppModule(m *appModuleModel) error {
	_, err := d.session.Update("app_module").SetMap(map[string]interface{}{
		"name":   m.Name,
		"desc":   m.Desc,
		"status": m.Status,
	}).Where("id=?", m.Id).Exec()
	return err
}

// 删除模块
func (d *managerDB) deleteAppModule(sid string) error {
	_, err := d.session.DeleteFrom("app_module").Where("sid=?", sid).Exec()
	return err
}

// 查询菜单列表
func (d *managerDB) querySysMenuList(keyword string) ([]*sysMenuModel, error) {
	var models []*sysMenuModel
	_, err := d.session.Select("*").From("sys_menu").Where("name like ?", "%"+keyword+"%").Load(&models)
	if err != nil {
		if err == dbr.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	return models, err
}

// 新增菜单
func (d *managerDB) insertSysMenu(m *sysMenuModel) (int64, error) {
	result, err := d.session.InsertInto("sys_menu").Columns(util.AttrToUnderscore(m)...).Record(m).Exec()
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return id, err
}

// 修改菜单
func (d *managerDB) updateSysMenu(m *sysMenuModel) error {
	_, err := d.session.Update("sys_menu").SetMap(map[string]interface{}{
		"name":           m.Name,
		"desc":           m.Desc,
		"status":         m.Status,
		"icon":           m.Icon,
		"sort":           m.Sort,
		"parent_key":     m.ParentKey,
		"layout":         m.Layout,
		"hidden_in_menu": m.HiddenInMenu,
		"redirect":       m.Redirect,
		"component":      m.Component,
	}).Where("key=?", m.Key).Exec()
	return err
}

// 删除菜单
func (d *managerDB) deleteSysMenu(key string) error {
	_, err := d.session.DeleteFrom("sys_menu").Where("key=?", key).Exec()
	return err
}

// 查询菜单用户列表
func (d *managerDB) getSysMenuUserListByUID(uid string) ([]string, error) {
	var models *sysMenuUserModel
	err := d.session.Select("*").From("sys_menu_user").Where("uid=?", uid).LoadOne(&models)
	if err != nil {
		if err == dbr.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	return strings.Split(models.Menus, ","), nil
}

// 新增菜单用户
func (d *managerDB) insertSysMenuUser(m *sysMenuUserModel) (int64, error) {
	result, err := d.session.InsertInto("sys_menu_user").Columns(util.AttrToUnderscore(m)...).Record(m).Exec()
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return id, err
}

// 分配菜单
func (d *managerDB) assignMenu(uid string, menus []string) error {
	_, err := d.session.Update("sys_menu_user").Set("menus", strings.Join(menus, ",")).Where("uid=?", uid).Exec()
	return err
}

type chatBgModel struct {
	Cover string // 封面
	Url   string // 图片地址
	IsSvg int    // 1 svg图片 0 普通图片
	dbs.BaseModel
}

type appVersionModel struct {
	AppVersion  string // app版本
	OS          string // android | ios
	IsForce     int    // 是否强制更新 1:是
	UpdateDesc  string // 更新说明
	DownloadURL string // 下载地址
	Signature   string // 安装包签名
	dbs.BaseModel
}

type appModuleModel struct {
	SID    string // 模块ID
	Name   string // 模块名称
	Desc   string // 介绍
	Status int    // 状态
	dbs.BaseModel
}

type sysMenuModel struct {
	Key          string // 菜单key
	Path         string // 菜单路径
	Name         string // 菜单名称
	Desc         string // 菜单描述
	Status       int    // 菜单状态
	Icon         string // 菜单图标
	Sort         int    // 菜单排序
	ParentKey    string // 父菜单key
	Layout       bool   // 是否是布局
	HiddenInMenu bool   // 是否在菜单中隐藏
	Redirect     string // 重定向路径
	Component    string // 组件路径
	dbs.BaseModel
}

type sysMenuUserModel struct {
	Menus string // 菜单ID列表
	UID   string // 用户ID
	dbs.BaseModel
}
