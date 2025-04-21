package applet

import (
	"github.com/TangSengDaoDao/TangSengDaoDaoServerLib/config"
	"github.com/gocraft/dbr/v2"
)

type DB struct {
	session *dbr.Session
	ctx     *config.Context
}

func NewDB(ctx *config.Context) *DB {
	return &DB{
		session: ctx.DB(),
		ctx:     ctx,
	}
}

// 获取小程序列表
func (d *DB) queryAppletList() ([]*appletModel, error) {
	var models []*appletModel
	_, err := d.session.Select("*").From("applet_config").OrderDir("updated_at", false).Load(&models)
	return models, err
}