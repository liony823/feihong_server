package common

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/TangSengDaoDao/TangSengDaoDaoServerLib/config"
	"go.uber.org/zap"
)

var onceSerce sync.Once

// IService IService
type IService interface {
	GetAppConfig() (*AppConfigResp, error)
	// 获取短编号
	GetShortno() (string, error)
	SetShortnoUsed(shortno string, business string) error
}

// NewService NewService
func NewService(ctx *config.Context) IService {
	return newService(ctx)
}

type service struct {
	ctx         *config.Context
	appConfigDB *appConfigDB
	shortnoDB   *shortnoDB
	shortnoLock sync.RWMutex
}

func newService(ctx *config.Context) *service {
	// if ctx.GetConfig().ShortNo.NumOn {
	onceSerce.Do(func() {
		go runGenShortnoTask(ctx)
	})
	// }

	return &service{
		ctx:         ctx,
		appConfigDB: newAppConfigDB(ctx),
		shortnoDB:   newShortnoDB(ctx),
	}
}

// GetAppConfig GetAppConfig
func (s *service) GetAppConfig() (*AppConfigResp, error) {
	appConfigM, err := s.appConfigDB.query()
	if err != nil {
		return nil, err
	}

	return &AppConfigResp{
		RSAPublicKey:                   appConfigM.RSAPublicKey,
		Version:                        appConfigM.Version,
		SuperToken:                     appConfigM.SuperToken,
		SuperTokenOn:                   appConfigM.SuperTokenOn,
		WelcomeMessage:                 appConfigM.WelcomeMessage,
		NewUserJoinSystemGroup:         appConfigM.NewUserJoinSystemGroup,
		SearchByPhone:                  appConfigM.SearchByPhone,
		RegisterInviteOn:               appConfigM.RegisterInviteOn,
		SendWelcomeMessageOn:           appConfigM.SendWelcomeMessageOn,
		InviteSystemAccountJoinGroupOn: appConfigM.InviteSystemAccountJoinGroupOn,
		RegisterUserMustCompleteInfoOn: appConfigM.RegisterUserMustCompleteInfoOn,
		ChannelPinnedMessageMaxCount:   appConfigM.ChannelPinnedMessageMaxCount,
		CanModifyApiUrl:                appConfigM.CanModifyApiUrl,

		// 新增字段
		IpWhiteList:                            appConfigM.IpWhiteList,
		LoginType:                              appConfigM.LoginType,
		SensitiveWords:                         appConfigM.SensitiveWords,
		DisableChangeDevice:                    appConfigM.DisableChangeDevice,
		SignupDeviceLimit:                      appConfigM.SignupDeviceLimit,
		SigleIpRegisterLimitIn12hour:           appConfigM.SigleIpRegisterLimitIn12hour,
		AutoClearHistoryMsg:                    appConfigM.AutoClearHistoryMsg,
		SigninAuthCodeVisible:                  appConfigM.SigninAuthCodeVisible,
		FriendOnlineStatusVisible:              appConfigM.FriendOnlineStatusVisible,
		MobileMsgReadStatusVisible:             appConfigM.MobileMsgReadStatusVisible,
		WalletPayoutMin:                        appConfigM.WalletPayoutMin,
		TransferMinAmount:                      appConfigM.TransferMinAmount,
		MobileEditMsg:                          appConfigM.MobileEditMsg,
		GroupMemberSeeMember:                   appConfigM.GroupMemberSeeMember,
		MsgTimeVisible:                         appConfigM.MsgTimeVisible,
		PinnedConversationSync:                 appConfigM.PinnedConversationSync,
		OnlyInternalFriendAdd:                  appConfigM.OnlyInternalFriendAdd,
		OnlyInternalFriendCreateGroup:          appConfigM.OnlyInternalFriendCreateGroup,
		OnlyInternalFriendSendGroupRedEnvelope: appConfigM.OnlyInternalFriendSendGroupRedEnvelope,
		OnlyInternalFriendSendGroupCard:        appConfigM.OnlyInternalFriendSendGroupCard,
		OnlyInternalFriendGroupRobotFreeMsg:    appConfigM.OnlyInternalFriendGroupRobotFreeMsg,
		GroupMemberLimit:                       appConfigM.GroupMemberLimit,
		UserAgreementContent:                   appConfigM.UserAgreementContent,
		PrivacyPolicyContent:                   appConfigM.PrivacyPolicyContent,
	}, nil
}

func (s *service) GetShortno() (string, error) {

	s.shortnoLock.Lock() // 这里需要加锁 要不然多线程下会出现shortNo重复的问题
	defer s.shortnoLock.Unlock()

	shortnoM, err := s.shortnoDB.queryVail()
	if err != nil {
		return "", err
	}
	if shortnoM == nil {
		return "", errors.New("没有短编号可分配")
	}
	err = s.shortnoDB.updateLock(shortnoM.Shortno, 1)
	if err != nil {
		return "", err
	}
	return shortnoM.Shortno, nil
}

func (s *service) SetShortnoUsed(shortno string, business string) error {
	return s.shortnoDB.updateUsed(shortno, 1, business)
}

// 开启生成短编号任务
func runGenShortnoTask(ctx *config.Context) {
	shortnoDB := newShortnoDB(ctx)
	errorSleep := time.Second * 2
	for {
		count, err := shortnoDB.queryVailCount()
		if err != nil {
			time.Sleep(errorSleep)
			continue
		}
		if count < 10000 {
			shortnos := generateNums(ctx.GetConfig().ShortNo.NumLen, 100)
			if len(shortnos) > 0 {
				err = shortnoDB.inserts(shortnos)
				if err != nil {
					ctx.Error("添加短编号失败！", zap.Error(err))
				}
			}
		}
		time.Sleep(time.Second * 30)
	}
}

func generateNums(len int, count int) []string {
	var nums = make([]string, 0, count)
	rd := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := count; i > 0; i-- {
		var num = rd.Int63n(1e16)
		nums = append(nums, fmt.Sprintf("%016d", num)[0:len])
	}
	return nums

}

type AppConfigResp struct {
	RSAPublicKey                   string `json:"rsa_public_key"`                      // 公钥
	Version                        int    `json:"version"`                             // 版本
	SuperToken                     string `json:"super_token"`                         // 超级token
	SuperTokenOn                   int    `json:"super_token_on"`                      // 超级token是否开启
	WelcomeMessage                 string `json:"welcome_message"`                     // 登录欢迎语
	NewUserJoinSystemGroup         int    `json:"new_user_join_system_group"`          // 新用户是否加入系统群聊
	SearchByPhone                  int    `json:"search_by_phone"`                     // 是否可通过手机号搜索
	RegisterInviteOn               int    `json:"register_invite_on"`                  // 是否开启注册邀请
	SendWelcomeMessageOn           int    `json:"send_welcome_message_on"`             // 是否发送登录欢迎语
	InviteSystemAccountJoinGroupOn int    `json:"invite_system_account_join_group_on"` // 是否允许邀请系统账号进入群聊
	RegisterUserMustCompleteInfoOn int    `json:"register_user_must_complete_info_on"` // 是否要求注册用户必须填写完整信息
	ChannelPinnedMessageMaxCount   int    `json:"channel_pinned_message_max_count"`    // 频道置顶消息最大数量
	CanModifyApiUrl                int    `json:"can_modify_api_url"`                  // 是否可以修改API地址

	ShortnoEditOff int `json:"shortno_edit_off"` // 是否关闭短编号编辑
	RevokeSecond   int `json:"revoke_second"`    // 消息可撤回时长

	// 新增字段
	IpWhiteList                            string `json:"ip_white_list"`                                // 后台IP白名单
	LoginType                              int    `json:"login_type"`                                   // app登录类型
	SensitiveWords                         string `json:"sensitive_words"`                              // 敏感词
	DisableChangeDevice                    int    `json:"disable_change_device"`                        // 是否禁止更换设备
	SignupDeviceLimit                      int    `json:"signup_device_limit"`                          // 设备限制注册限制数
	SigleIpRegisterLimitIn12hour           int    `json:"sigle_ip_register_limit_in12hour"`             // 单IP12小时注册限制数
	AutoClearHistoryMsg                    int    `json:"auto_clear_history_msg"`                       // 自动清除几天前历史消息
	MiniProgramVisible                     int    `json:"mini_program_visible"`                         // 小程序页是否可见
	DiscoveryVisible                       int    `json:"discovery_visible"`                            // 发现页是否可见
	ChargeAndPayoutVisible                 int    `json:"charge_and_payout_visible"`                    // 充值和提现是否可见
	VoiceCallVisible                       int    `json:"voice_call_visible"`                           // 语音通话是否可见
	VideoCallVisible                       int    `json:"video_call_visible"`                           // 视频通话是否可见
	SignupInviteCodeVisible                int    `json:"signup_invite_code_visible"`                   // 注册邀请码是否可见
	SigninAuthCodeVisible                  int    `json:"signin_auth_code_visible"`                     // 登录授权码是否可见
	FriendOnlineStatusVisible              int    `json:"friend_online_status_visible"`                 // 好友在线状态是否可见
	MobileMsgReadStatusVisible             int    `json:"mobile_msg_read_status_visible"`               // 手机消息已读状态是否可见
	SignRedEnvelopeVisible                 int    `json:"sign_red_envelope_visible"`                    // 签到红包模块是否开启
	MineWalletVisible                      int    `json:"mine_wallet_visible"`                          // 我的钱包是否开启
	WalletPayoutMin                        int    `json:"wallet_payout_min"`                            // 钱包提现最小金额
	RedEnvelopeVisible                     int    `json:"red_envelope_visible"`                         // 红包模块是否开启
	TransferVisible                        int    `json:"transfer_visible"`                             // 转账模块是否开启
	TransferMinAmount                      int    `json:"transfer_min_amount"`                          // 转账最小金额
	MobileEditMsg                          int    `json:"mobile_edit_msg"`                              // 手机端是否可以编辑消息
	GroupMemberSeeMember                   int    `json:"group_member_see_member"`                      // 普通群成员是否可以查看其他群成员
	MsgTimeVisible                         int    `json:"msg_time_visible"`                             // 消息时间是否可见
	PinnedConversationSync                 int    `json:"pinned_conversation_sync"`                     // 置顶会话是否同步
	OnlyInternalFriendAdd                  int    `json:"only_internal_friend_add"`                     // 仅内部号可被加好友及加非内部号好友
	OnlyInternalFriendCreateGroup          int    `json:"only_internal_friend_create_group"`            // 仅内部号可建群
	OnlyInternalFriendSendGroupRedEnvelope int    `json:"only_internal_friend_send_group_red_envelope"` // 仅内部号可发群红包
	OnlyInternalFriendSendGroupCard        int    `json:"only_internal_friend_send_group_card"`         // 仅内部号可群内推送名片
	OnlyInternalFriendGroupRobotFreeMsg    int    `json:"only_internal_friend_group_robot_free_msg"`    // 群机器人免消息
	GroupMemberLimit                       int    `json:"group_member_limit"`                           // 群人数限制
	UserAgreementContent                   string `json:"user_agreement_content"`                       // 用户协议内容
	PrivacyPolicyContent                   string `json:"privacy_policy_content"`                       // 隐私政策内容
	MomentsVisible                         int    `json:"moments_visible"`                              // 好友分享是否可见
}
