package statistics

import (
	"errors"

	"github.com/TangSengDaoDao/TangSengDaoDaoServer/modules/group"
	"github.com/TangSengDaoDao/TangSengDaoDaoServer/modules/message"
	"github.com/TangSengDaoDao/TangSengDaoDaoServer/modules/user"
	"github.com/TangSengDaoDao/TangSengDaoDaoServerLib/config"
	"github.com/TangSengDaoDao/TangSengDaoDaoServerLib/pkg/log"
	"github.com/TangSengDaoDao/TangSengDaoDaoServerLib/pkg/wkhttp"
	"go.uber.org/zap"
)

type Manager struct {
	ctx *config.Context
	log.Log

	messageService message.IService
	userService    user.IService
	groupService   group.IService
}

func NewManager(ctx *config.Context) *Manager {
	return &Manager{
		ctx:            ctx,
		Log:            log.NewTLog("statisticsManager"),
		messageService: message.NewService(ctx),
		userService:    user.NewService(ctx),
		groupService:   group.NewService(ctx),
	}
}

func (m *Manager) Route(r *wkhttp.WKHttp) {
	auth := r.Group("/v1/manager", m.ctx.BasicAuthMiddleware(r), m.ctx.AuthMiddleware(r))
	{
		auth.GET("/statistics/countnum", m.countNum)                                                // 统计数量
		auth.GET("/statistics/registeruser/:start_date/:end_date", m.registerUserListWithDateSpace) // 某个时间区间的注册统计数据
		auth.GET("/statistics/createdgroup/:start_date/:end_date", m.createGroupWithDateSpace)      // 某个时间段的建群数据
	}
}

// 某个时间区间的注册数据
func (m *Manager) registerUserListWithDateSpace(c *wkhttp.Context) {
	err := c.CheckLoginRole()
	if err != nil {
		c.ResponseError(err)
		return
	}
	startDate := c.Param("start_date")
	endDate := c.Param("end_date")
	if startDate == "" || endDate == "" {
		c.ResponseError(errors.New("查询日期不能为空"))
		return
	}
	list, err := m.userService.GetRegisterCountWithDateSpace(startDate, endDate)
	if err != nil {
		c.ResponseError(errors.New("查询注册用户数量错误"))
		return
	}
	c.Response(list)
}

// 获取某个时间段的建群数量
func (m *Manager) createGroupWithDateSpace(c *wkhttp.Context) {
	err := c.CheckLoginRole()
	if err != nil {
		c.ResponseError(err)
		return
	}
	startDate := c.Param("start_date")
	endDate := c.Param("end_date")
	if startDate == "" || endDate == "" {
		c.ResponseError(errors.New("查询日期不能为空"))
		return
	}
	list, err := m.groupService.GetGroupWithDateSpace(startDate, endDate)
	if err != nil {
		c.ResponseError(errors.New("查询注册用户数量错误"))
		return
	}
	c.Response(list)
}

// 统计数量
func (m *Manager) countNum(c *wkhttp.Context) {
	err := c.CheckLoginRole()
	if err != nil {
		c.ResponseError(err)
		return
	}
	date := c.Query("date")
	// 获取总用户数
	totalUserCount, err := m.userService.GetAllUserCount()
	if err != nil {
		m.Error("查询用户数量错误", zap.Error(err))
		c.ResponseError(errors.New("查询用户数量错误"))
		return
	}
	// 查询某天注册量
	registerCount, err := m.userService.GetRegisterWithDate(date)
	if err != nil {
		m.Error("查询某天用户注册量错误", zap.Error(err))
		c.ResponseError(errors.New("查询某天用户注册量错误"))
		return
	}
	// 查询总群数
	totalGroupCount, err := m.groupService.GetAllGroupCount()
	if err != nil {
		m.Error("查询总群数量错误", zap.Error(err))
		c.ResponseError(errors.New("查询总群数量错误"))
		return
	}
	// 查询某天的新建群数量
	groupCreatedCount, err := m.groupService.GetCreatedCountWithDate(date)
	if err != nil {
		m.Error("查询某天群新建数量错误", zap.Error(err))
		c.ResponseError(errors.New("查询某天群新建数量错误"))
		return
	}
	// 查询总在线数量
	onlineCount, err := m.userService.GetOnlineCount()
	if err != nil {
		m.Error("查询总在线用户数量错误", zap.Error(err))
		c.ResponseError(errors.New("查询总在线用户数量错误"))
		return
	}
	messageCount, err := m.messageService.GetMsgCount()
	if err != nil {
		m.Error("查询总消息数错误", zap.Error(err))
		c.ResponseError(errors.New("查询总消息数量错误"))
		return
	}
	c.Response(&countNum{
		UserTotalCount:   totalUserCount,
		RegisterCount:    registerCount,
		GroupTotalCount:  totalGroupCount,
		GroupCreateCount: groupCreatedCount,
		OnlineTotalCount: onlineCount,
		MsgTotalCount:    messageCount,
	})
}

type countNum struct {
	UserTotalCount   int64 `json:"user_total_count"`   // 用户总数
	RegisterCount    int64 `json:"register_count"`     // 注册数量
	GroupTotalCount  int64 `json:"group_total_count"`  // 群总数
	GroupCreateCount int64 `json:"group_create_count"` // 群创建数量
	OnlineTotalCount int64 `json:"online_total_count"` // 总在线数量
	MsgTotalCount    int64 `json:"msg_total_count"`    // 总消息数量
}
