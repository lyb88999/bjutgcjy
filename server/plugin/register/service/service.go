package service

import (
	"errors"
	uuid2 "github.com/gofrs/uuid/v5"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	systemReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	plugGlobal "github.com/flipped-aurora/gin-vue-admin/server/plugin/register/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/register/model"
	userService "github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/mojocn/base64Captcha"
	uuid "github.com/satori/go.uuid"
)

type RegisterService struct{}

func (e *RegisterService) PlugService(req model.Request) (res *system.SysUser, err error) {
	if err := utils.Verify(req, utils.LoginVerify); err != nil {
		return res, err
	}
	var (
		store = base64Captcha.DefaultMemStore
		user  system.SysUser
		us    *userService.UserService
	)
	// 验证码验证
	if !store.Verify(req.CaptchaId, req.Captcha, true) {
		return res, errors.New("验证码错误")
	}
	// 推荐人验证
	var referrer system.SysUser
	if err = global.GVA_DB.Where("username = ?", req.ReferrerUsername).First(&referrer).Error; err != nil {
		return res, errors.New("推荐人不存在")
	}
	// 邀请码验证
	var invitation system.InviteCode
	if err = global.GVA_DB.Where("code = ?", req.InviteCode).First(&invitation).Error; err != nil {
		return res, errors.New("邀请码不存在")
	}
	if invitation.Status == 1 {
		return res, errors.New("邀请码已被使用")
	}
	if invitation.UUID != referrer.UUID {
		return res, errors.New("邀请码与推荐人不匹配")
	}
	u := &system.SysUser{Username: req.Username, Password: req.Password}
	err = global.GVA_DB.Where("username = ?", u.Username).Preload("Authorities").Preload("Authority").First(&user).Error
	if err == nil {
		return res, errors.New("用户名已注册")
	}
	if user.Username != "" && user.Password != "" {
		return res, errors.New("用户名已注册")
	}

	var sysAuthority systemReq.Register
	sysAuthority.Username = u.Username
	sysAuthority.NickName = u.NickName
	sysAuthority.Password = u.Password
	// TODO 角色ID
	sysAuthority.AuthorityId = plugGlobal.GlobalConfig.AuthorityId
	sysAuthority.AuthorityIds = append(sysAuthority.AuthorityIds, plugGlobal.GlobalConfig.AuthorityId)
	// 因为上面定义过，且得到了数据库默认的值，所以直接使用
	user.Password = u.Password
	user.UUID = uuid2.UUID(uuid.NewV4())
	user.Username = u.Username
	user.NickName = u.Username
	// 将邀请码与用户关联
	user.InviteCode = invitation.Code
	user.AuthorityId = plugGlobal.GlobalConfig.AuthorityId

	for _, v := range sysAuthority.AuthorityIds {
		user.Authorities = append(user.Authorities, system.SysAuthority{
			AuthorityId: v,
			//DefaultRouter: "dashboard", //TODO 疑问，系统注册的时候有这个参数
		})
	}

	if rest, err := us.Register(user); err != nil {
		return &rest, errors.New("注册失败!")
	}
	// 标记邀请码已被使用，并关联新用户
	invitation.Status = 1
	if err = global.GVA_DB.Save(&invitation).Error; err != nil {
		return res, errors.New("邀请码更新失败")
	}
	if res, err = us.Login(&user); err != nil {
		return res, errors.New("登陆失败!")
	}
	return res, nil
	// 前面的代码 拿不到正确的 user，所以需要再次查询一次
}
