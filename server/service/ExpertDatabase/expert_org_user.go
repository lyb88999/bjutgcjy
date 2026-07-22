package ExpertDatabase

import (
	"errors"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	ExpertDatabaseReq "github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase/request"
	ExpertDatabaseResp "github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gofrs/uuid/v5"
	"gorm.io/gorm"
)

type ExpertOrgUserService struct{}

// getOperatorOrgId 查询操作人自己的所属单位ID，本单位账号管理功能的所有操作都以此为唯一数据域边界
func (s *ExpertOrgUserService) getOperatorOrgId(operatorID uint) (*uint, error) {
	var operator system.SysUser
	if err := global.GVA_DB.Where("id = ?", operatorID).First(&operator).Error; err != nil {
		return nil, err
	}
	return operator.OrgId, nil
}

// CreateOrgApplicant 单位审核员为本单位新建个人申报人账号。
// 所属单位固定取操作人自己的 org_id，角色固定是"个人申报人"——请求体里不接受这两个字段，
// 从根本上避免审核员越权建别的单位或别的角色（比如另一个审核员、管理员）的账号
func (s *ExpertOrgUserService) CreateOrgApplicant(operatorID uint, req ExpertDatabaseReq.CreateOrgApplicantRequest) error {
	orgId, err := s.getOperatorOrgId(operatorID)
	if err != nil {
		return err
	}
	if orgId == nil {
		return errors.New("你的账号还没有关联所属单位，请先联系管理员配置好所属单位，才能给本单位新建账号")
	}

	var existing system.SysUser
	findErr := global.GVA_DB.Where("username = ?", req.Username).First(&existing).Error
	if findErr == nil {
		return errors.New("用户名已被占用")
	}
	if !errors.Is(findErr, gorm.ErrRecordNotFound) {
		return findErr
	}

	nickName := req.NickName
	if nickName == "" {
		nickName = req.Username
	}

	user := system.SysUser{
		UUID:        uuid.Must(uuid.NewV4()),
		Username:    req.Username,
		Password:    utils.BcryptHash(req.Password),
		NickName:    nickName,
		Phone:       req.Phone,
		Email:       req.Email,
		AuthorityId: authorityIndividualApplicant,
		Authorities: []system.SysAuthority{{AuthorityId: authorityIndividualApplicant}},
		Enable:      1,
		OrgId:       orgId,
	}
	return global.GVA_DB.Create(&user).Error
}

// GetOrgUserList 分页获取本单位的个人申报人账号列表——只看得到自己单位、"个人申报人"角色的账号，
// 看不到别的单位的账号，也看不到本单位其他审核员/管理员账号
func (s *ExpertOrgUserService) GetOrgUserList(operatorID uint, page ExpertDatabaseReq.OrgUserSearch) (list []ExpertDatabaseResp.OrgUserItem, total int64, err error) {
	orgId, err := s.getOperatorOrgId(operatorID)
	if err != nil {
		return nil, 0, err
	}
	if orgId == nil {
		return []ExpertDatabaseResp.OrgUserItem{}, 0, nil
	}

	limit := page.PageSize
	offset := page.PageSize * (page.Page - 1)
	db := global.GVA_DB.Model(&system.SysUser{}).Where("org_id = ? AND authority_id = ?", orgId, authorityIndividualApplicant)
	if err = db.Count(&total).Error; err != nil {
		return
	}
	var users []system.SysUser
	if err = db.Order("id desc").Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return
	}
	for _, u := range users {
		list = append(list, ExpertDatabaseResp.OrgUserItem{
			ID: u.ID, CreatedAt: u.CreatedAt, Username: u.Username, NickName: u.NickName,
			Phone: u.Phone, Email: u.Email, Enable: u.Enable, AuthorityId: u.AuthorityId,
		})
	}
	return
}

// ToggleOrgUserEnable 启用/冻结本单位的个人申报人账号——操作前二次校验目标账号确实属于
// 操作人自己的单位、确实是"个人申报人"角色，防止拿别的单位或别的角色的用户ID越权操作
func (s *ExpertOrgUserService) ToggleOrgUserEnable(operatorID uint, targetUserID uint, enable int) error {
	orgId, err := s.getOperatorOrgId(operatorID)
	if err != nil {
		return err
	}
	if orgId == nil {
		return errors.New("你的账号还没有关联所属单位")
	}
	var target system.SysUser
	if err := global.GVA_DB.Where("id = ?", targetUserID).First(&target).Error; err != nil {
		return err
	}
	if target.OrgId == nil || *target.OrgId != *orgId || target.AuthorityId != authorityIndividualApplicant {
		return errors.New("无权操作其他单位或其他角色的账号")
	}
	return global.GVA_DB.Model(&system.SysUser{}).Where("id = ?", targetUserID).Update("enable", enable).Error
}
