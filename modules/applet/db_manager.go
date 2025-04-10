package applet

import (
	"github.com/TangSengDaoDao/TangSengDaoDaoServer/pkg/util"
	"github.com/TangSengDaoDao/TangSengDaoDaoServerLib/config"
	dbs "github.com/TangSengDaoDao/TangSengDaoDaoServerLib/pkg/db"
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

// 添加小程序
func (d *managerDB) insertApplet(m *appletModel) (int64, error) {
	result, err := d.session.InsertInto("applet_config").Columns(util.AttrToUnderscore(m)...).Record(m).Exec()
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return id, err
}

// 查询小程序列表
func (d *managerDB) queryAppletListWithPage(pageSize, page uint64, keyword string) ([]*appletModel, error) {
	var models []*appletModel
	_, err := d.session.Select("*").From("applet_config").Where("display_name like ?", "%"+keyword+"%").Offset((page-1)*pageSize).Limit(pageSize).OrderDir("updated_at", false).Load(&models)
	return models, err
}

// 查询小程序数量
func (d *managerDB) queryAppletCount(keyword string) (int64, error) {
	var count int64
	_, err := d.session.Select("count(*)").From("applet_config").Where("display_name like ?", "%"+keyword+"%").Load(&count)
	return count, err
}

// 查询默认小程序
func (d *managerDB) queryDefaultApplet() ([]*appletModel, error) {
	var model []*appletModel
	_, err := d.session.Select("*").From("applet_config").Where("is_default = 1").OrderDir("updated_at", false).Load(&model)
	return model, err
}

// 查询小程序
func (d *managerDB) queryAppletOnDefault() (*appletModel, error) {
	var model *appletModel
	_, err := d.session.Select("*").From("applet_config").Where("is_default = 1").Load(&model)
	return model, err
}

// 更新小程序
func (d *managerDB) updateApplet(m *appletModel) error {
	_, err := d.session.Update("applet_config").SetMap(map[string]interface{}{
		"display_name": m.DisplayName,
		"icon":         m.Icon,
		"link":         m.Link,
		"priority":     m.Priority,
		"status":       m.Status,
	}).Where("key = ?", m.Key).Exec()
	return err
}

// updateAppletTx 更新小程序事物
func (d *managerDB) setDefaultAppletTx(m *appletModel, tx *dbr.Tx) error {
	_, err := tx.Update("applet_config").SetMap(map[string]interface{}{
		"is_default": m.IsDefault,
	}).Where("key = ?", m.Key).Exec()
	return err
}

// 删除小程序
func (d *managerDB) deleteApplet(key string) error {
	_, err := d.session.DeleteFrom("applet_config").Where("key = ?", key).Exec()
	return err
}

type appletModel struct {
	DisplayName string
	Key         string
	Icon        string
	Link        string
	Priority    int64
	Status      int
	IsDefault   int
	dbs.BaseModel
}
