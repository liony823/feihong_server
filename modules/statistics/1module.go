package statistics

import (
	_ "embed"

	"github.com/TangSengDaoDao/TangSengDaoDaoServerLib/config"
	"github.com/TangSengDaoDao/TangSengDaoDaoServerLib/pkg/register"
)

//go:embed swagger/api.yaml
var swaggerContent string

func init() {
	register.AddModule(func(ctx interface{}) register.Module {
		x := ctx.(*config.Context)
		return register.Module{
			Name: "statistics_manager",
			SetupAPI: func() register.APIRouter {
				return NewManager(x)
			},
			Swagger: swaggerContent,
		}
	})
}
