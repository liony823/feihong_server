package applet

import (
	"embed"

	"github.com/TangSengDaoDao/TangSengDaoDaoServerLib/config"
	"github.com/TangSengDaoDao/TangSengDaoDaoServerLib/pkg/register"
)

//go:embed sql
var sqlFS embed.FS

//go:embed swagger/api.yaml
var swaggerContent string

func init() {
	register.AddModule(func(ctx interface{}) register.Module {
		return register.Module{
			Name: "applet",
			SetupAPI: func() register.APIRouter {
				return NewApplet(ctx.(*config.Context))
			},
			Swagger: swaggerContent,
			SQLDir:  register.NewSQLFS(sqlFS),
		}
	})
	register.AddModule(func(ctx interface{}) register.Module {
		return register.Module{
			Name: "applet_manager",
			SetupAPI: func() register.APIRouter {
				return NewManager(ctx.(*config.Context))
			},
		}
	})
}
