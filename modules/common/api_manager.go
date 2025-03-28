package common

import (
	"errors"
	"net/http"
	"slices"
	"strings"

	"github.com/TangSengDaoDao/TangSengDaoDaoServerLib/config"
	"github.com/TangSengDaoDao/TangSengDaoDaoServerLib/pkg/log"
	"github.com/TangSengDaoDao/TangSengDaoDaoServerLib/pkg/wkhttp"
	"go.uber.org/zap"
)

// Manager 通用后台管理api
type Manager struct {
	ctx *config.Context
	log.Log
	db          *managerDB
	appconfigDB *appConfigDB
}

// NewManager NewManager
func NewManager(ctx *config.Context) *Manager {
	return &Manager{
		ctx:         ctx,
		Log:         log.NewTLog("commonManager"),
		db:          newManagerDB(ctx),
		appconfigDB: newAppConfigDB(ctx),
	}
}

// Route 配置路由规则
func (m *Manager) Route(r *wkhttp.WKHttp) {
	auth := r.Group("/v1/manager", m.ctx.BasicAuthMiddleware(r), m.ctx.AuthMiddleware(r))
	{
		auth.GET("/common/appconfig", m.appconfig)               // 获取app配置
		auth.POST("/common/appconfig", m.updateConfig)           // 修改app配置
		auth.GET("/common/appmodule", m.getAppModule)            // 获取app模块
		auth.PUT("/common/appmodule", m.updateAppModule)         // 修改app模块
		auth.POST("/common/appmodule", m.addAppModule)           // 新增app模块
		auth.DELETE("/common/:sid/appmodule", m.deleteAppModule) // 删除app模块
		auth.GET("/common/menu/current", m.getCurrentUserMenu)   // 获取当前用户菜单
		auth.GET("/common/menu", m.getMenu)                      // 获取菜单
		auth.POST("/common/menu", m.addMenu)                     // 新增菜单
		auth.PUT("/common/menu/:key", m.updateMenu)              // 修改菜单
		auth.DELETE("/common/menu/:key", m.deleteMenu)           // 删除菜单
		auth.GET("/common/menu/user/:uid", m.getMenuUser)        // 获取菜单用户
		auth.POST("/common/menu/user/:uid", m.assignMenu)        // 新增菜单用户
	}

	r.GET("/v1/manager/health", m.ctx.BasicAuthMiddleware(r), func(c *wkhttp.Context) {
		var (
			statusMap = map[string]string{
				"status": "up",
				"db":     "up",
				"redis":  "up",
			}
			lastError error
		)

		err := m.db.session.Ping()
		if err != nil {
			m.Error("db ping error", zap.Error(err))
			lastError = err
			statusMap["db"] = "down"
		}

		_, err = m.ctx.GetRedisConn().Ping()
		if err != nil {
			m.Error("redis ping error", zap.Error(err))
			lastError = err
			statusMap["redis"] = "down"
		}

		if lastError != nil {
			statusMap["status"] = "down"
			statusMap["error"] = lastError.Error()
		}

		c.JSON(http.StatusOK, statusMap)
	})
}
func (m *Manager) deleteAppModule(c *wkhttp.Context) {
	err := c.CheckLoginRoleIsSuperAdmin()
	if err != nil {
		c.ResponseError(err)
		return
	}

	sid := c.Param("sid")
	if strings.TrimSpace(sid) == "" {
		c.ResponseError(errors.New("sid不能为空！"))
		return
	}
	module, err := m.db.queryAppModuleWithSid(sid)
	if err != nil {
		m.Error("查询app模块错误", zap.Error(err))
		c.ResponseError(errors.New("查询app模块错误"))
		return
	}
	if module == nil {
		c.ResponseError(errors.New("删除的模块不存在"))
		return
	}
	err = m.db.deleteAppModule(sid)
	if err != nil {
		m.Error("删除app模块错误", zap.Error(err))
		c.ResponseError(errors.New("删除app模块错误"))
		return
	}
	c.ResponseOK()
}

// 新增app模块
func (m *Manager) addAppModule(c *wkhttp.Context) {
	err := c.CheckLoginRoleIsSuperAdmin()
	if err != nil {
		c.ResponseError(err)
		return
	}
	type ReqVO struct {
		SID    string `json:"sid"`
		Name   string `json:"name"`
		Desc   string `json:"desc"`
		Status int    `json:"status"`
	}
	var req ReqVO
	if err := c.BindJSON(&req); err != nil {
		c.ResponseError(errors.New("请求数据格式有误！"))
		return
	}

	if strings.TrimSpace(req.SID) == "" || strings.TrimSpace(req.Desc) == "" || strings.TrimSpace(req.Name) == "" {
		c.ResponseError(errors.New("名称/ID/介绍不能为空！"))
		return
	}
	module, err := m.db.queryAppModuleWithSid(req.SID)
	if err != nil {
		m.Error("查询app模块错误", zap.Error(err))
		c.ResponseError(errors.New("查询app模块错误"))
		return
	}
	if module != nil && module.SID != "" {
		c.ResponseError(errors.New("该sid模块已存在"))
		return
	}
	_, err = m.db.insertAppModule(&appModuleModel{
		SID:    req.SID,
		Name:   req.Name,
		Desc:   req.Desc,
		Status: req.Status,
	})
	if err != nil {
		m.Error("新增app模块错误", zap.Error(err))
		c.ResponseError(errors.New("新增app模块错误"))
		return
	}
	c.ResponseOK()
}
func (m *Manager) updateAppModule(c *wkhttp.Context) {
	err := c.CheckLoginRoleIsSuperAdmin()
	if err != nil {
		c.ResponseError(err)
		return
	}
	type ReqVO struct {
		SID    string `json:"sid"`
		Name   string `json:"name"`
		Desc   string `json:"desc"`
		Status int    `json:"status"`
	}
	var req ReqVO
	if err := c.BindJSON(&req); err != nil {
		c.ResponseError(errors.New("请求数据格式有误！"))
		return
	}

	if strings.TrimSpace(req.SID) == "" || strings.TrimSpace(req.Desc) == "" || strings.TrimSpace(req.Name) == "" {
		c.ResponseError(errors.New("名称/ID/介绍不能为空！"))
		return
	}
	module, err := m.db.queryAppModuleWithSid(req.SID)
	if err != nil {
		m.Error("查询app模块错误", zap.Error(err))
		c.ResponseError(errors.New("查询app模块错误"))
		return
	}
	if module == nil {
		c.ResponseError(errors.New("不存在该模块"))
		return
	}
	module.Name = req.Name
	module.Desc = req.Desc
	module.Status = req.Status
	err = m.db.updateAppModule(module)
	if err != nil {
		m.Error("修改app模块错误", zap.Error(err))
		c.ResponseError(errors.New("修改app模块错误"))
		return
	}
	c.ResponseOK()
}

// 获取app模块
func (m *Manager) getAppModule(c *wkhttp.Context) {
	err := c.CheckLoginRole()
	if err != nil {
		c.ResponseError(err)
		return
	}
	modules, err := m.db.queryAppModule()
	if err != nil {
		m.Error("查询app模块错误", zap.Error(err))
		c.ResponseError(errors.New("查询app模块错误"))
		return
	}
	list := make([]*managerAppModule, 0)
	if len(modules) > 0 {
		for _, module := range modules {
			list = append(list, &managerAppModule{
				SID:    module.SID,
				Name:   module.Name,
				Desc:   module.Desc,
				Status: module.Status,
			})
		}
	}
	c.Response(list)
}
func (m *Manager) updateConfig(c *wkhttp.Context) {
	err := c.CheckLoginRoleIsSuperAdmin()
	if err != nil {
		c.ResponseError(err)
		return
	}
	type reqVO struct {
		RevokeSecond                           int    `json:"revoke_second"`
		WelcomeMessage                         string `json:"welcome_message"`
		NewUserJoinSystemGroup                 int    `json:"new_user_join_system_group"`
		SearchByPhone                          int    `json:"search_by_phone"`
		RegisterInviteOn                       int    `json:"register_invite_on"`                           // 开启注册邀请机制
		SendWelcomeMessageOn                   int    `json:"send_welcome_message_on"`                      // 开启注册登录发送欢迎语
		InviteSystemAccountJoinGroupOn         int    `json:"invite_system_account_join_group_on"`          // 开启系统账号加入群聊
		RegisterUserMustCompleteInfoOn         int    `json:"register_user_must_complete_info_on"`          // 注册用户必须填写完整信息
		ChannelPinnedMessageMaxCount           int    `json:"channel_pinned_message_max_count"`             // 频道置顶消息最大数量
		CanModifyApiUrl                        int    `json:"can_modify_api_url"`                           // 是否可以修改api地址
		IpWhiteList                            string `json:"ip_white_list"`                                // ip白名单
		LoginType                              int    `json:"login_type"`                                   // app登录类型
		SensitiveWords                         string `json:"sensitive_words"`                              // 敏感词
		DisableChangeDevice                    int    `json:"disable_change_device"`                        // 是否禁止更换设备
		SignupDeviceLimit                      int    `json:"signup_device_limit"`                          // 设备限制注册限制数
		SigleIpRegisterLimitIn12hour           int    `json:"sigle_ip_register_limit_in_12hour"`            // 单IP12小时注册限制数
		AutoClearHistoryMsg                    int    `json:"auto_clear_history_msg"`                       // 自动清除几天前历史消息
		SigninAuthCodeVisible                  int    `json:"signin_auth_code_visible"`                     // 登录授权码是否可见
		FriendOnlineStatusVisible              int    `json:"friend_online_status_visible"`                 // 好友在线状态是否可见
		MobileMsgReadStatusVisible             int    `json:"mobile_msg_read_status_visible"`               // 手机消息已读状态是否可见
		WalletPayoutMin                        int    `json:"wallet_payout_min"`                            // 钱包提现最小金额
		TransferMinAmount                      int    `json:"transfer_min_amount"`                          // 转账最小金额
		MobileEditMsg                          int    `json:"mobile_edit_msg"`                              // 手机端是否可以编辑消息
		GroupMemberSeeMember                   int    `json:"group_member_see_member"`                      // 普通群成员是否可以查看其他群成员
		MsgTimeVisible                         int    `json:"msg_time_visible"`                             // 消息时间是否可见
		PinnedConversationSync                 int    `json:"pinned_conversation_sync"`                     // 置顶会话是否同步
		OnlyInternalFriendAdd                  int    `json:"only_internal_friend_add"`                     // 仅内部号可被加好友及加非内部号好友
		OnlyInternalFriendCreateGroup          int    `json:"only_internal_friend_create_group"`            // 仅内部号可建群
		OnlyInternalFriendSendGroupRedEnvelope int    `json:"only_internal_friend_send_group_red_envelope"` // 仅内部号可发群红包
		OnlyInternalFriendSendGroupCard        int    `json:"only_internal_friend_send_group_card"`         // 仅内部号可群内推送名片
		OnlyInternalFriendGroupRobotFreeMsg    int    `json:"only_internal_friend_group_robot_free_msg"`    // 仅内部号群机器人免消息
		GroupMemberLimit                       int    `json:"group_member_limit"`                           // 群人数限制
		UserAgreementContent                   string `json:"user_agreement_content"`                       // 用户协议内容
		PrivacyPolicyContent                   string `json:"privacy_policy_content"`                       // 隐私政策内容                      // 好友分享是否可见
	}
	var req reqVO
	if err := c.BindJSON(&req); err != nil {
		c.ResponseError(errors.New("请求数据格式有误！"))
		return
	}
	appConfigM, err := m.appconfigDB.query()
	if err != nil {
		m.Error("查询应用配置失败！", zap.Error(err))
		c.ResponseError(errors.New("查询应用配置失败！"))
		return
	}
	configMap := map[string]interface{}{}
	configMap["revoke_second"] = req.RevokeSecond
	configMap["welcome_message"] = req.WelcomeMessage
	configMap["new_user_join_system_group"] = req.NewUserJoinSystemGroup
	configMap["search_by_phone"] = req.SearchByPhone
	configMap["register_invite_on"] = req.RegisterInviteOn
	configMap["send_welcome_message_on"] = req.SendWelcomeMessageOn
	configMap["invite_system_account_join_group_on"] = req.InviteSystemAccountJoinGroupOn
	configMap["register_user_must_complete_info_on"] = req.RegisterUserMustCompleteInfoOn
	configMap["channel_pinned_message_max_count"] = req.ChannelPinnedMessageMaxCount
	configMap["can_modify_api_url"] = req.CanModifyApiUrl
	configMap["ip_white_list"] = req.IpWhiteList
	configMap["login_type"] = req.LoginType
	configMap["sensitive_words"] = req.SensitiveWords
	configMap["disable_change_device"] = req.DisableChangeDevice
	configMap["signup_device_limit"] = req.SignupDeviceLimit
	configMap["sigle_ip_register_limit_in_12hour"] = req.SigleIpRegisterLimitIn12hour
	configMap["auto_clear_history_msg"] = req.AutoClearHistoryMsg
	configMap["signin_auth_code_visible"] = req.SigninAuthCodeVisible
	configMap["friend_online_status_visible"] = req.FriendOnlineStatusVisible
	configMap["mobile_msg_read_status_visible"] = req.MobileMsgReadStatusVisible
	configMap["wallet_payout_min"] = req.WalletPayoutMin
	configMap["transfer_min_amount"] = req.TransferMinAmount
	configMap["mobile_edit_msg"] = req.MobileEditMsg
	configMap["group_member_see_member"] = req.GroupMemberSeeMember
	configMap["msg_time_visible"] = req.MsgTimeVisible
	configMap["pinned_conversation_sync"] = req.PinnedConversationSync
	configMap["only_internal_friend_add"] = req.OnlyInternalFriendAdd
	configMap["only_internal_friend_create_group"] = req.OnlyInternalFriendCreateGroup
	configMap["only_internal_friend_send_group_red_envelope"] = req.OnlyInternalFriendSendGroupRedEnvelope
	configMap["only_internal_friend_send_group_card"] = req.OnlyInternalFriendSendGroupCard
	configMap["only_internal_friend_group_robot_free_msg"] = req.OnlyInternalFriendGroupRobotFreeMsg
	configMap["group_member_limit"] = req.GroupMemberLimit
	configMap["user_agreement_content"] = req.UserAgreementContent
	configMap["privacy_policy_content"] = req.PrivacyPolicyContent

	err = m.appconfigDB.updateWithMap(configMap, appConfigM.Id)
	if err != nil {
		m.Error("修改app配置信息错误", zap.Error(err))
		c.ResponseError(errors.New("修改app配置信息错误"))
		return
	}
	c.ResponseOK()
}
func (m *Manager) appconfig(c *wkhttp.Context) {
	err := c.CheckLoginRole()
	if err != nil {
		c.ResponseError(err)
		return
	}
	appconfig, err := m.appconfigDB.query()
	if err != nil {
		m.Error("查询应用配置失败！", zap.Error(err))
		c.ResponseError(errors.New("查询应用配置失败！"))
		return
	}
	var revokeSecond = 0
	var newUserJoinSystemGroup = 1
	var welcomeMessage = ""
	var searchByPhone = 1
	var registerInviteOn = 0
	var sendWelcomeMessageOn = 0
	var inviteSystemAccountJoinGroupOn = 0
	var registerUserMustCompleteInfoOn = 0
	var channelPinnedMessageMaxCount = 10
	var canModifyApiUrl = 0
	var ipWhiteList = ""
	var loginType = 0
	var sensitiveWords = ""
	var disableChangeDevice = 0
	var signupDeviceLimit = 0
	var singleIpRegisterLimitIn12hour = 0
	var autoClearHistoryMsg = 0
	var signinAuthCodeVisible = 0
	var friendOnlineStatusVisible = 0
	var mobileMsgReadStatusVisible = 0
	var walletPayoutMin = 0
	var transferMinAmount = 0
	var mobileEditMsg = 0
	var groupMemberSeeMember = 0
	var msgTimeVisible = 0
	var pinnedConversationSync = 0
	var onlyInternalFriendAdd = 0
	var onlyInternalFriendCreateGroup = 0
	var onlyInternalFriendSendGroupRedEnvelope = 0
	var onlyInternalFriendSendGroupCard = 0
	var onlyInternalFriendGroupRobotFreeMsg = 0
	var groupMemberLimit = 0
	var userAgreementContent = ""
	var privacyPolicyContent = ""
	if appconfig != nil {
		revokeSecond = appconfig.RevokeSecond
		welcomeMessage = appconfig.WelcomeMessage
		newUserJoinSystemGroup = appconfig.NewUserJoinSystemGroup
		searchByPhone = appconfig.SearchByPhone
		registerInviteOn = appconfig.RegisterInviteOn
		sendWelcomeMessageOn = appconfig.SendWelcomeMessageOn
		inviteSystemAccountJoinGroupOn = appconfig.InviteSystemAccountJoinGroupOn
		registerUserMustCompleteInfoOn = appconfig.RegisterUserMustCompleteInfoOn
		channelPinnedMessageMaxCount = appconfig.ChannelPinnedMessageMaxCount
		canModifyApiUrl = appconfig.CanModifyApiUrl
		ipWhiteList = appconfig.IpWhiteList
		loginType = appconfig.LoginType
		sensitiveWords = appconfig.SensitiveWords
		disableChangeDevice = appconfig.DisableChangeDevice
		signupDeviceLimit = appconfig.SignupDeviceLimit
		singleIpRegisterLimitIn12hour = appconfig.SigleIpRegisterLimitIn12hour
		autoClearHistoryMsg = appconfig.AutoClearHistoryMsg
		signinAuthCodeVisible = appconfig.SigninAuthCodeVisible
		friendOnlineStatusVisible = appconfig.FriendOnlineStatusVisible
		mobileMsgReadStatusVisible = appconfig.MobileMsgReadStatusVisible
		walletPayoutMin = appconfig.WalletPayoutMin
		transferMinAmount = appconfig.TransferMinAmount
		mobileEditMsg = appconfig.MobileEditMsg
		groupMemberSeeMember = appconfig.GroupMemberSeeMember
		msgTimeVisible = appconfig.MsgTimeVisible
		pinnedConversationSync = appconfig.PinnedConversationSync
		onlyInternalFriendAdd = appconfig.OnlyInternalFriendAdd
		onlyInternalFriendCreateGroup = appconfig.OnlyInternalFriendCreateGroup
		onlyInternalFriendSendGroupRedEnvelope = appconfig.OnlyInternalFriendSendGroupRedEnvelope
		onlyInternalFriendSendGroupCard = appconfig.OnlyInternalFriendSendGroupCard
		onlyInternalFriendGroupRobotFreeMsg = appconfig.OnlyInternalFriendGroupRobotFreeMsg
		groupMemberLimit = appconfig.GroupMemberLimit
		userAgreementContent = appconfig.UserAgreementContent
		privacyPolicyContent = appconfig.PrivacyPolicyContent
	}
	if revokeSecond == 0 {
		revokeSecond = 120
	}
	if welcomeMessage == "" {
		welcomeMessage = m.ctx.GetConfig().WelcomeMessage
	}
	c.Response(&AppConfigResp{
		RevokeSecond:                           revokeSecond,
		WelcomeMessage:                         welcomeMessage,
		NewUserJoinSystemGroup:                 newUserJoinSystemGroup,
		SearchByPhone:                          searchByPhone,
		RegisterInviteOn:                       registerInviteOn,
		SendWelcomeMessageOn:                   sendWelcomeMessageOn,
		InviteSystemAccountJoinGroupOn:         inviteSystemAccountJoinGroupOn,
		RegisterUserMustCompleteInfoOn:         registerUserMustCompleteInfoOn,
		ChannelPinnedMessageMaxCount:           channelPinnedMessageMaxCount,
		CanModifyApiUrl:                        canModifyApiUrl,
		IpWhiteList:                            ipWhiteList,
		LoginType:                              loginType,
		SensitiveWords:                         sensitiveWords,
		DisableChangeDevice:                    disableChangeDevice,
		SignupDeviceLimit:                      signupDeviceLimit,
		SigleIpRegisterLimitIn12hour:           singleIpRegisterLimitIn12hour,
		AutoClearHistoryMsg:                    autoClearHistoryMsg,
		SigninAuthCodeVisible:                  signinAuthCodeVisible,
		FriendOnlineStatusVisible:              friendOnlineStatusVisible,
		MobileMsgReadStatusVisible:             mobileMsgReadStatusVisible,
		WalletPayoutMin:                        walletPayoutMin,
		TransferMinAmount:                      transferMinAmount,
		MobileEditMsg:                          mobileEditMsg,
		GroupMemberSeeMember:                   groupMemberSeeMember,
		MsgTimeVisible:                         msgTimeVisible,
		PinnedConversationSync:                 pinnedConversationSync,
		OnlyInternalFriendAdd:                  onlyInternalFriendAdd,
		OnlyInternalFriendCreateGroup:          onlyInternalFriendCreateGroup,
		OnlyInternalFriendSendGroupRedEnvelope: onlyInternalFriendSendGroupRedEnvelope,
		OnlyInternalFriendSendGroupCard:        onlyInternalFriendSendGroupCard,
		OnlyInternalFriendGroupRobotFreeMsg:    onlyInternalFriendGroupRobotFreeMsg,
		GroupMemberLimit:                       groupMemberLimit,
		UserAgreementContent:                   userAgreementContent,
		PrivacyPolicyContent:                   privacyPolicyContent,
	})
}

// 获取当前用户菜单
func (m *Manager) getCurrentUserMenu(c *wkhttp.Context) {

	var (
		menus []*sysMenuModel
	)
	list := make([]*managerMenu, 0)

	menus, err := m.db.querySysMenuList("")
	if err != nil {
		m.Error("查询菜单失败", zap.Error(err))
		c.ResponseError(errors.New("查询菜单失败"))
		return
	}
	err = c.CheckLoginRoleIsSuperAdmin()

	if err != nil {
		err := c.CheckLoginRole()
		if err != nil {
			c.ResponseError(err)
			return
		}

		menuKeys, err := m.db.getSysMenuUserListByUID(c.GetLoginUID())
		if err != nil {
			m.Error("查询菜单失败", zap.Error(err))
			c.ResponseError(errors.New("查询菜单失败"))
			return
		}

		for _, menu := range menus {
			if slices.Contains(menuKeys, menu.Key) {
				list = append(list, m.sysMenuModelToManagerMenu(menu))
			}
		}
	}

	for _, menu := range menus {
		list = append(list, m.sysMenuModelToManagerMenu(menu))
	}
	c.Response(list)
}

// 获取菜单列表
func (m *Manager) getMenu(c *wkhttp.Context) {
	err := c.CheckLoginRole()
	if err != nil {
		c.ResponseError(err)
		return
	}
	menus, err := m.db.querySysMenuList(c.Query("keyword"))
	if err != nil {
		m.Error("查询菜单列表失败", zap.Error(err))
		c.ResponseError(errors.New("查询菜单列表失败"))
		return
	}
	list := make([]*managerMenu, 0)
	if len(menus) > 0 {
		for _, menu := range menus {
			list = append(list, m.sysMenuModelToManagerMenu(menu))
		}
	}
	c.Response(list)
}

// 将sysMenuModel转成managerMenu
func (m *Manager) sysMenuModelToManagerMenu(menu *sysMenuModel) *managerMenu {
	return &managerMenu{
		Key:          menu.Key,
		Path:         menu.Path,
		Name:         menu.Name,
		Desc:         menu.Desc,
		Status:       menu.Status,
		Icon:         menu.Icon,
		Sort:         menu.Sort,
		ParentKey:    menu.ParentKey,
		Layout:       menu.Layout,
		HiddenInMenu: menu.HiddenInMenu,
		Redirect:     menu.Redirect,
		Component:    menu.Component,
		CreatedAt:    menu.CreatedAt.String(),
		UpdatedAt:    menu.UpdatedAt.String(),
	}
}

// 新增菜单
func (m *Manager) addMenu(c *wkhttp.Context) {
	err := c.CheckLoginRole()
	if err != nil {
		c.ResponseError(err)
		return
	}
	type managerMenuReqVO struct {
		Key          string `json:"key"`
		Path         string `json:"path"`
		Name         string `json:"name"`
		Desc         string `json:"desc"`
		Status       int    `json:"status"`
		Icon         string `json:"icon"`
		Sort         int    `json:"sort"`
		ParentKey    string `json:"parent_key"`
		Layout       bool   `json:"layout"`
		HiddenInMenu bool   `json:"hidden_in_menu"`
		Redirect     string `json:"redirect"`
		Component    string `json:"component"`
	}
	var req managerMenuReqVO
	if err := c.BindJSON(&req); err != nil {
		c.ResponseError(errors.New("请求数据格式有误！"))
		return
	}
	id, err := m.db.insertSysMenu(&sysMenuModel{
		Key:    req.Key,
		Path:   req.Path,
		Name:   req.Name,
		Desc:   req.Desc,
		Status: req.Status,
	})
	if err != nil {
		m.Error("新增菜单失败", zap.Error(err))
		c.ResponseError(errors.New("新增菜单失败"))
		return
	}
	c.Response(map[string]interface{}{
		"id": id,
	})
}

// 修改菜单
func (m *Manager) updateMenu(c *wkhttp.Context) {
	err := c.CheckLoginRole()
	if err != nil {
		c.ResponseError(err)
		return
	}
	type managerMenuReqVO struct {
		Key          string `json:"key"`
		Path         string `json:"path"`
		Name         string `json:"name"`
		Desc         string `json:"desc"`
		Status       int    `json:"status"`
		Icon         string `json:"icon"`
		Sort         int    `json:"sort"`
		ParentKey    string `json:"parent_key"`
		Layout       bool   `json:"layout"`
		HiddenInMenu bool   `json:"hidden_in_menu"`
		Redirect     string `json:"redirect"`
		Component    string `json:"component"`
	}
	var req managerMenuReqVO
	if err := c.BindJSON(&req); err != nil {
		c.ResponseError(errors.New("请求数据格式有误！"))
		return
	}

	err = m.db.updateSysMenu(&sysMenuModel{
		Key:          req.Key,
		Path:         req.Path,
		Name:         req.Name,
		Desc:         req.Desc,
		Status:       req.Status,
		Icon:         req.Icon,
		Sort:         req.Sort,
		ParentKey:    req.ParentKey,
		Layout:       req.Layout,
		HiddenInMenu: req.HiddenInMenu,
		Redirect:     req.Redirect,
		Component:    req.Component,
	})
	if err != nil {
		m.Error("修改菜单失败", zap.Error(err))
		c.ResponseError(errors.New("修改菜单失败"))
		return
	}
	c.ResponseOK()
}

// 删除菜单
func (m *Manager) deleteMenu(c *wkhttp.Context) {
	err := c.CheckLoginRole()
	if err != nil {
		c.ResponseError(err)
		return
	}
	err = m.db.deleteSysMenu(c.Param("key"))
	if err != nil {
		m.Error("删除菜单失败", zap.Error(err))
		c.ResponseError(errors.New("删除菜单失败"))
		return
	}
	c.ResponseOK()
}

// 获取菜单用户
func (m *Manager) getMenuUser(c *wkhttp.Context) {
	err := c.CheckLoginRole()
	if err != nil {
		c.ResponseError(err)
		return
	}
	uid := c.Param("uid")
	menus, err := m.db.getSysMenuUserListByUID(uid)
	if err != nil {
		m.Error("查询菜单用户失败", zap.Error(err))
		c.ResponseError(errors.New("查询菜单用户失败"))
		return
	}
	c.Response(menus)
}

// 分配菜单用户
func (m *Manager) assignMenu(c *wkhttp.Context) {
	err := c.CheckLoginRole()
	if err != nil {
		c.ResponseError(err)
		return
	}
	uid := c.Param("uid")
	type managerMenuUserReqVO struct {
		Menus []string `json:"menus"`
	}
	var req managerMenuUserReqVO
	if err := c.BindJSON(&req); err != nil {
		c.ResponseError(errors.New("请求数据格式有误！"))
		return
	}
	_, err = m.db.getSysMenuUserListByUID(uid)
	if err != nil {
		_, err = m.db.insertSysMenuUser(&sysMenuUserModel{
			UID:   uid,
			Menus: strings.Join(req.Menus, ","),
		})
		if err != nil {
			m.Error("新增菜单用户失败", zap.Error(err))
			c.ResponseError(errors.New("新增菜单用户失败"))
			return
		}
		c.ResponseOK()
		return
	}

	err = m.db.assignMenu(uid, req.Menus)
	if err != nil {
		m.Error("分配菜单失败", zap.Error(err))
		c.ResponseError(errors.New("分配菜单失败"))
		return
	}
	c.ResponseOK()
}

type managerMenu struct {
	Key          string `json:"key"`
	Path         string `json:"path"`
	Name         string `json:"name"`
	Desc         string `json:"desc"`
	Status       int    `json:"status"`
	Icon         string `json:"icon"`
	Sort         int    `json:"sort"`
	ParentKey    string `json:"parent_key"`
	Layout       bool   `json:"layout"`
	HiddenInMenu bool   `json:"hidden_in_menu"`
	Redirect     string `json:"redirect"`
	Component    string `json:"component"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type managerAppModule struct {
	SID    string `json:"sid"`
	Name   string `json:"name"`
	Desc   string `json:"desc"`
	Status int    `json:"status"` // 模块状态 1.可选 0.不可选 2.选中不可编辑
}
