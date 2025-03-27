package common

import (
	"github.com/TangSengDaoDao/TangSengDaoDaoServerLib/config"
	ldb "github.com/TangSengDaoDao/TangSengDaoDaoServerLib/pkg/db"
	"github.com/TangSengDaoDao/TangSengDaoDaoServerLib/pkg/util"
	"github.com/gocraft/dbr/v2"
)

type appConfigDB struct {
	session *dbr.Session
	ctx     *config.Context
}

func newAppConfigDB(ctx *config.Context) *appConfigDB {

	return &appConfigDB{
		session: ctx.DB(),
		ctx:     ctx,
	}
}

func (a *appConfigDB) query() (*appConfigModel, error) {
	var m *appConfigModel
	_, err := a.session.Select("*").From("app_config").OrderDesc("created_at").Load(&m)
	return m, err
}

func (a *appConfigDB) insert(m *appConfigModel) error {
	_, err := a.session.InsertInto("app_config").Columns(util.AttrToUnderscore(m)...).Record(m).Exec()
	return err
}
func (a *appConfigDB) updateWithMap(configMap map[string]interface{}, id int64) error {
	_, err := a.session.Update("app_config").SetMap(configMap).Where("id=?", id).Exec()
	return err
}

type appConfigModel struct {
	RSAPrivateKey                  string
	RSAPublicKey                   string
	Version                        int
	SuperToken                     string
	SuperTokenOn                   int
	RevokeSecond                   int    // 消息可撤回时长
	WelcomeMessage                 string // 登录欢迎语
	NewUserJoinSystemGroup         int    // 新用户是否加入系统群聊
	SearchByPhone                  int    // 是否可通过手机号搜索
	RegisterInviteOn               int    // 开启注册邀请机制
	SendWelcomeMessageOn           int    // 开启注册登录发送欢迎语
	InviteSystemAccountJoinGroupOn int    // 开启系统账号加入群聊
	RegisterUserMustCompleteInfoOn int    // 注册用户是否必须完善个人信息
	ChannelPinnedMessageMaxCount   int    // 频道置顶消息最大数量
	CanModifyApiUrl                int    // 是否可以修改API地址

	IpWhiteList                            string // 后台IP白名单
	LoginType                              int    // app登录类型
	SensitiveWords                         string // 敏感词 (多个敏感词用英文的 | 符号分割)
	DisableChangeDevice                    int    // 是否禁止更换设备: 1 禁止、0 不禁止
	SignupDeviceLimit                      int    // 设备限制注册限制数, 0为不限制
	SigleIpRegisterLimitIn12hour           int    // 单IP12小时注册限制数, 0为不限制
	AutoClearHistoryMsg                    int    // 自动清除几天前历史消息, 0 不自动清除
	SigninAuthCodeVisible                  int    // 登录授权码是否可见
	FriendOnlineStatusVisible              int    // 好友在线状态是否可见
	MobileMsgReadStatusVisible             int    // 手机消息已读状态是否可见
	WalletPayoutMin                        int    // 钱包提现最小金额，单位为分，0为不限制
	TransferMinAmount                      int    // 转账最小金额，单位为分，0为不限制
	MobileEditMsg                          int    // 手机端是否可以编辑消息
	GroupMemberSeeMember                   int    // 普通群成员是否可以查看其他群成员
	MsgTimeVisible                         int    // 消息时间是否可见
	PinnedConversationSync                 int    // 置顶会话是否同步
	OnlyInternalFriendAdd                  int    // 仅内部号可被加好友及加非内部号好友
	OnlyInternalFriendCreateGroup          int    // 仅内部号可建群
	OnlyInternalFriendSendGroupRedEnvelope int    // 仅内部号可发群红包
	OnlyInternalFriendSendGroupCard        int    // 仅内部号可群内推送名片
	OnlyInternalFriendGroupRobotFreeMsg    int    // 群机器人免消息
	GroupMemberLimit                       int    // 群人数限制: 0 不限制
	UserAgreementContent                   string // 用户协议内容
	PrivacyPolicyContent                   string // 隐私政策内容

	ldb.BaseModel
}
