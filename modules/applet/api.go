package applet

import (
	"github.com/TangSengDaoDao/TangSengDaoDaoServerLib/config"
	"github.com/TangSengDaoDao/TangSengDaoDaoServerLib/pkg/log"
	"github.com/TangSengDaoDao/TangSengDaoDaoServerLib/pkg/wkhttp"
)

type Applet struct {
	ctx *config.Context
	db  *DB
	log.Log
}

func NewApplet(ctx *config.Context) *Applet {
	return &Applet{
		ctx: ctx,
		db:  NewDB(ctx),
		Log: log.NewTLog("applet"),
	}
}

func (a *Applet) Route(r *wkhttp.WKHttp) {
	auth := r.Group("/v1", a.ctx.AuthMiddleware(r))
	{
		auth.GET("/applet", a.getAppletList)
	}
} 

func (a *Applet) getAppletList(c *wkhttp.Context) {
	list, err := a.db.queryAppletList()
	if err != nil {
		c.ResponseError(err)
		return
	}
	result := make([]*appletModelResp, 0)
	for _, v := range list {
		result = append(result, &appletModelResp{
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
	c.Response(result)
}

type appletModelResp struct {
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
