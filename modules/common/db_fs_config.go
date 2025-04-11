package common

import (
	"github.com/TangSengDaoDao/TangSengDaoDaoServerLib/config"
	dbs "github.com/TangSengDaoDao/TangSengDaoDaoServerLib/pkg/db"
	"github.com/gocraft/dbr/v2"
)

type fsConfigDB struct {
	session *dbr.Session
	ctx     *config.Context
}

func newFsConfigDB(ctx *config.Context) *fsConfigDB {
	return &fsConfigDB{
		session: ctx.DB(),
		ctx:     ctx,
	}
}

// 查询文件上传配置列表
func (d *fsConfigDB) queryFSConfigList() ([]*fsConfigModel, error) {
	var models []*fsConfigModel
	_, err := d.session.Select("*").From("fs_config").OrderDir("updated_at", false).Load(&models)
	return models, err
}

// 根据key修改文件上传配置
func (d *fsConfigDB) updateFSConfigWithKey(m *fsConfigModel) error {
	_, err := d.session.Update("fs_config").SetMap(map[string]interface{}{
		"title":   m.Title,
		"status":  m.Status,
		"options": m.Options,
	}).Where("`key`=?", m.Key).Exec()
	return err
}

type fsConfigModel struct {
	Title   string
	Key     string
	Options string
	Status  int32
	dbs.BaseModel
}
