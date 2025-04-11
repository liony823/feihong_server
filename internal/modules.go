package modules

// 引入模块
import (
	_ "github.com/TangSengDaoDao/TangSengDaoDaoServer/modules/applet"     // 小程序模块
	_ "github.com/TangSengDaoDao/TangSengDaoDaoServer/modules/base"       // 基础模块
	_ "github.com/TangSengDaoDao/TangSengDaoDaoServer/modules/channel"    // 渠道模块
	_ "github.com/TangSengDaoDao/TangSengDaoDaoServer/modules/common"     // 公共模块
	_ "github.com/TangSengDaoDao/TangSengDaoDaoServer/modules/file"       // 文件模块
	_ "github.com/TangSengDaoDao/TangSengDaoDaoServer/modules/group"      // 群组模块
	_ "github.com/TangSengDaoDao/TangSengDaoDaoServer/modules/message"    // 消息模块
	_ "github.com/TangSengDaoDao/TangSengDaoDaoServer/modules/openapi"    // 开放平台模块
	_ "github.com/TangSengDaoDao/TangSengDaoDaoServer/modules/qrcode"     // 二维码模块
	_ "github.com/TangSengDaoDao/TangSengDaoDaoServer/modules/report"     // 报表模块
	_ "github.com/TangSengDaoDao/TangSengDaoDaoServer/modules/robot"      // 机器人模块
	_ "github.com/TangSengDaoDao/TangSengDaoDaoServer/modules/statistics" // 统计模块
	_ "github.com/TangSengDaoDao/TangSengDaoDaoServer/modules/user"       // 用户模块
	_ "github.com/TangSengDaoDao/TangSengDaoDaoServer/modules/webhook"    // 钩子模块
	_ "github.com/TangSengDaoDao/TangSengDaoDaoServer/modules/workplace"  // 工作台模块
)
