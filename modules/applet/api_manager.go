package applet

import (
	"errors"
	"strings"

	"github.com/TangSengDaoDao/TangSengDaoDaoServerLib/config"
	"github.com/TangSengDaoDao/TangSengDaoDaoServerLib/pkg/log"
	"github.com/TangSengDaoDao/TangSengDaoDaoServerLib/pkg/wkhttp"
	"go.uber.org/zap"
)

type Manager struct {
	ctx *config.Context
	db  *managerDB
	log.Log
}

func NewManager(ctx *config.Context) *Manager {
	return &Manager{
		ctx: ctx,
		db:  newManagerDB(ctx),
		Log: log.NewTLog("appletManager"),
	}
}

func (m *Manager) Route(r *wkhttp.WKHttp) {
	auth := r.Group("/v1/manager", m.ctx.BasicAuthMiddleware(r), m.ctx.AuthMiddleware(r), m.ctx.AdminOperateRecordMiddleware(r))
	{
		auth.GET("/applet", m.getAppletList)
		auth.POST("/applet", m.addApplet)
		auth.PUT("/applet", m.updateApplet)
		auth.DELETE("/applet/:key", m.deleteApplet)
		auth.PUT("/applet/default", m.setDefaultApplet)
	}
}

func (m *Manager) getAppletList(c *wkhttp.Context) {
	err := c.CheckLoginRole()
	if err != nil {
		c.ResponseError(err)
		return
	}
	keyword := c.Query("keyword")
	pageIndex, pageSize := c.GetPage()
	var list []*appletModel
	var count int64

	list, err = m.db.queryAppletListWithPage(uint64(pageSize), uint64(pageIndex), keyword)
	if err != nil {
		m.Error("获取小程序列表失败", zap.Error(err))
		c.ResponseError(err)
		return
	}
	count, err = m.db.queryAppletCount(keyword)
	if err != nil {
		m.Error("获取小程序数量失败", zap.Error(err))
		c.ResponseError(err)
		return
	}
	result := make([]*managerAppletModel, 0)
	for _, v := range list {
		result = append(result, &managerAppletModel{
			ID:          v.Id,
			DisplayName: v.DisplayName,
			Key:         v.Key,
			Icon:        v.Icon,
			Link:        v.Link,
			Priority:    v.Priority,
			Status:      v.Status,
			IsDefault:   v.IsDefault,
			CreatedAt:   v.CreatedAt.String(),
			UpdatedAt:   v.UpdatedAt.String(),
		})
	}
	c.Response(map[string]interface{}{
		"list":  result,
		"count": count,
	})
}

// 添加小程序
func (m *Manager) addApplet(c *wkhttp.Context) {
	err := c.CheckLoginRole()
	if err != nil {
		c.ResponseError(err)
		return
	}

	var req addAppletReq
	err = c.BindJSON(&req)
	if err != nil {
		c.ResponseError(errors.New("请求数据格式有误！"))
		return
	}
	err = req.check()
	if err != nil {
		c.ResponseError(err)
		return
	}

	_, err = m.db.insertApplet(&appletModel{
		DisplayName: req.DisplayName,
		Key:         req.Key,
		Icon:        req.Icon,
		Link:        req.Link,
		Priority:    req.Priority,
	})
	if err != nil {
		m.Error("添加小程序失败", zap.Error(err))
		c.ResponseError(errors.New("添加小程序失败"))
		return
	}
	c.ResponseOK()
}

// 更新小程序
func (m *Manager) updateApplet(c *wkhttp.Context) {
	err := c.CheckLoginRole()
	if err != nil {
		c.ResponseError(err)
		return
	}
	var req addAppletReq
	err = c.BindJSON(&req)
	if err != nil {
		c.ResponseError(errors.New("请求数据格式有误！"))
		return
	}
	err = req.check()
	if err != nil {
		c.ResponseError(err)
		return
	}

	err = m.db.updateApplet(&appletModel{
		DisplayName: req.DisplayName,
		Key:         req.Key,
		Icon:        req.Icon,
		Link:        req.Link,
		Priority:    req.Priority,
		Status:      req.Status,
	})
	if err != nil {
		c.ResponseError(err)
		return
	}
	c.ResponseOK()
}

// 删除小程序
func (m *Manager) deleteApplet(c *wkhttp.Context) {
	err := c.CheckLoginRole()
	if err != nil {
		c.ResponseError(err)
		return
	}
	key := c.Param("key")
	if key == "" {
		c.ResponseError(errors.New("小程序key不能为空"))
		return
	}
	err = m.db.deleteApplet(key)
	if err != nil {
		m.Error("删除小程序失败", zap.Error(err))
		c.ResponseError(errors.New("删除小程序失败"))
		return
	}
	c.ResponseOK()
}

// 设置默认小程序
func (m *Manager) setDefaultApplet(c *wkhttp.Context) {
	err := c.CheckLoginRole()
	if err != nil {
		c.ResponseError(err)
		return
	}
	type reqVO struct {
		Key       string `json:"key"`
		IsDefault int    `json:"is_default"`
	}
	var req reqVO
	err = c.BindJSON(&req)
	if err != nil {
		c.ResponseError(errors.New("请求数据格式有误！"))
		return
	}
	if req.Key == "" {
		c.ResponseError(errors.New("小程序key不能为空"))
		return
	}
	if req.IsDefault != 0 && req.IsDefault != 1 {
		c.ResponseError(errors.New("is_default参数错误"))
		return
	}

	defaultApplets, err := m.db.queryDefaultApplet()
	if err != nil {
		m.Error("获取默认小程序失败", zap.Error(err))
		c.ResponseError(errors.New("获取默认小程序失败"))
		return
	}
	tx, err := m.db.session.Begin()
	if err != nil {
		m.Error("开启事物失败", zap.Error(err))
		c.ResponseError(errors.New("开启事物失败"))
		return
	}

	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
			panic(err)
		}
	}()

	if req.IsDefault == 1 {
		for _, applet := range defaultApplets {
			if applet.Key == req.Key {
				continue
			}
			err := m.db.setDefaultAppletTx(&appletModel{
				Key:       applet.Key,
				IsDefault: 0,
			}, tx)
			if err != nil {
				tx.Rollback()
				m.Error("设置默认小程序失败", zap.Error(err))
				c.ResponseError(errors.New("设置默认小程序失败"))
				return
			}
		}

	}

	err = m.db.setDefaultAppletTx(&appletModel{
		Key:       req.Key,
		IsDefault: req.IsDefault,
	}, tx)
	if err != nil {
		m.Error("设置默认小程序失败", zap.Error(err))
		c.ResponseError(errors.New("设置默认小程序失败"))
		return
	}

	err = tx.Commit()
	if err != nil {
		m.Error("设置默认小程序失败", zap.Error(err))
		c.ResponseError(errors.New("设置默认小程序失败"))
		return
	}
	c.ResponseOK()
}
func (r addAppletReq) check() error {
	if strings.TrimSpace(r.DisplayName) == "" {
		return errors.New("小程序名称不能为空")
	}
	if strings.TrimSpace(r.Key) == "" {
		return errors.New("小程序key不能为空")
	}
	if strings.TrimSpace(r.Icon) == "" {
		return errors.New("小程序icon不能为空")
	}
	if strings.TrimSpace(r.Link) == "" {
		return errors.New("小程序link不能为空")
	}
	return nil
}

type addAppletReq struct {
	DisplayName string `json:"display_name"`
	Key         string `json:"key"`
	Icon        string `json:"icon"`
	Link        string `json:"link"`
	Priority    int64  `json:"priority"`
	Status      int    `json:"status"`
}

type managerAppletModel struct {
	ID          int64  `json:"id"`
	DisplayName string `json:"display_name"`
	Key         string `json:"key"`
	Icon        string `json:"icon"`
	Link        string `json:"link"`
	Priority    int64  `json:"priority"`
	Status      int    `json:"status"`
	IsDefault   int    `json:"is_default"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
