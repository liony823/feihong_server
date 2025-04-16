package user

import (
	"github.com/TangSengDaoDao/TangSengDaoDaoServerLib/common"
	"github.com/TangSengDaoDao/TangSengDaoDaoServerLib/pkg/wkhttp"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// 获取密保问题
func (u *User) getSecurityQuestion(c *wkhttp.Context) {
	uid := c.GetLoginUID()
	question, err := u.db.queryUserSecurity(uid)

	if err != nil {
		u.Error("查询密保问题失败", zap.Error(err))
		c.ResponseError(errors.New(common.ErrSecurityQueryFailed))
		return
	}

	c.Response(map[string]interface{}{
		"question": question.Question,
	})
}

// 用户添加密保问题
func (u *User) addSecurityQuestion(c *wkhttp.Context) {
	uid := c.GetLoginUID()
	var req setSecurityQuestionReq
	if err := c.BindJSON(&req); err != nil {
		c.ResponseError(errors.New(common.ErrRequestDataError))
		return
	}

	if req.Question == "" {
		c.ResponseError(errors.New(common.ErrSecurityQuestionEmpty))
		return
	}

	if req.Answer == "" {
		c.ResponseError(errors.New(common.ErrSecurityAnswerEmpty))
		return
	}

	err := u.db.insertUserSecurity(&UserSecurityModel{
		UID:      uid,
		Question: req.Question,
		Answer:   req.Answer,
	})
	if err != nil {
		u.Error("添加密保问题失败", zap.Error(err))
		c.ResponseError(errors.New(common.ErrSecurityAddFailed))
		return
	}

	c.ResponseOK()
}

// 用户修改密保问题
func (u *User) updateSecurityQuestion(c *wkhttp.Context) {
	uid := c.GetLoginUID()
	var req setSecurityQuestionReq
	if err := c.BindJSON(&req); err != nil {
		c.ResponseError(errors.New(common.ErrRequestDataError))
		return
	}

	if req.Question == "" {
		c.ResponseError(errors.New(common.ErrSecurityQuestionEmpty))
		return
	}

	if req.Answer == "" {
		c.ResponseError(errors.New(common.ErrSecurityAnswerEmpty))
		return
	}

	userSecurity, err := u.db.queryUserSecurity(uid)
	if err != nil {
		u.Error("查询密保问题失败", zap.Error(err))
		c.ResponseError(errors.New(common.ErrSecurityQueryFailed))
		return
	}

	if userSecurity == nil {
		c.ResponseError(errors.New(common.ErrSecurityQuestionNotExist))
		return
	}

	err = u.db.updateUserSecurity(&UserSecurityModel{
		UID:      uid,
		Question: req.Question,
		Answer:   req.Answer,
	})
	if err != nil {
		u.Error("修改密保问题失败", zap.Error(err))
		c.ResponseError(errors.New(common.ErrSecurityUpdateFailed))
		return
	}

	c.ResponseOK()
}

type setSecurityQuestionReq struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}
