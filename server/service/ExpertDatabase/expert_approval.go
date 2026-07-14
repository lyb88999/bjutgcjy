package ExpertDatabase

import (
	"errors"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase"
	ExpertDatabaseReq "github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase/request"
	commonRequest "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"gorm.io/gorm"
)

// 专家档案审核状态机取值
const (
	StatusDraft             = "draft"
	StatusPendingOrgReview  = "pending_org_review"
	StatusOrgRejected       = "org_rejected"
	StatusPendingCityReview = "pending_city_review"
	StatusCityRejected      = "city_rejected"
	StatusPublished         = "published"
)

type ExpertApprovalService struct{}

// getUserOrgId 查询用户所属单位ID，用于单位审核环节的数据域校验
func (s *ExpertApprovalService) getUserOrgId(userID uint) (*uint, error) {
	var user system.SysUser
	if err := global.GVA_DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return user.OrgId, nil
}

// transition 校验来源状态是否允许流转、落库新状态、写审核日志，是所有流转方法的共用核心
func (s *ExpertApprovalService) transition(tx *gorm.DB, profile *ExpertDatabase.ExpertProfile, allowedFrom []string, to string, operatorID uint, opinion string) error {
	allowed := false
	for _, from := range allowedFrom {
		if profile.Status == from {
			allowed = true
			break
		}
	}
	if !allowed {
		return errors.New("当前状态不允许该操作：" + profile.Status)
	}
	fromStatus := profile.Status
	if err := tx.Model(&ExpertDatabase.ExpertProfile{}).Where("id = ?", profile.ID).Update("status", to).Error; err != nil {
		return err
	}
	log := ExpertDatabase.ExpertApprovalLog{
		ExpertId:   profile.ID,
		FromStatus: fromStatus,
		ToStatus:   to,
		OperatorId: operatorID,
		Opinion:    opinion,
	}
	return tx.Create(&log).Error
}

// Submit 提交审核：草稿/单位退回/市级退回 -> 待单位审核
func (s *ExpertApprovalService) Submit(expertID uint, operatorID uint) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var profile ExpertDatabase.ExpertProfile
		if err := tx.Where("id = ?", expertID).First(&profile).Error; err != nil {
			return err
		}
		return s.transition(tx, &profile, []string{StatusDraft, StatusOrgRejected, StatusCityRejected}, StatusPendingOrgReview, operatorID, "")
	})
}

// OrgApprove 单位审核通过：待单位审核 -> 待市级审核；只能审核本单位提交的档案
func (s *ExpertApprovalService) OrgApprove(expertID uint, operatorID uint) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var profile ExpertDatabase.ExpertProfile
		if err := tx.Where("id = ?", expertID).First(&profile).Error; err != nil {
			return err
		}
		if err := s.assertSameOrg(operatorID, profile.OrgId); err != nil {
			return err
		}
		return s.transition(tx, &profile, []string{StatusPendingOrgReview}, StatusPendingCityReview, operatorID, "")
	})
}

// OrgReject 单位审核退回：待单位审核 -> 单位已退回，需填写意见；只能审核本单位提交的档案
func (s *ExpertApprovalService) OrgReject(expertID uint, operatorID uint, opinion string) error {
	if opinion == "" {
		return errors.New("退回时必须填写审核意见")
	}
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var profile ExpertDatabase.ExpertProfile
		if err := tx.Where("id = ?", expertID).First(&profile).Error; err != nil {
			return err
		}
		if err := s.assertSameOrg(operatorID, profile.OrgId); err != nil {
			return err
		}
		return s.transition(tx, &profile, []string{StatusPendingOrgReview}, StatusOrgRejected, operatorID, opinion)
	})
}

// CityApprove 市级审核通过：待市级审核 -> 已发布，正式收录进入检索排序范围
func (s *ExpertApprovalService) CityApprove(expertID uint, operatorID uint) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var profile ExpertDatabase.ExpertProfile
		if err := tx.Where("id = ?", expertID).First(&profile).Error; err != nil {
			return err
		}
		return s.transition(tx, &profile, []string{StatusPendingCityReview}, StatusPublished, operatorID, "")
	})
}

// CityReject 市级审核退回：待市级审核 -> 市级已退回，需填写意见
func (s *ExpertApprovalService) CityReject(expertID uint, operatorID uint, opinion string) error {
	if opinion == "" {
		return errors.New("退回时必须填写审核意见")
	}
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var profile ExpertDatabase.ExpertProfile
		if err := tx.Where("id = ?", expertID).First(&profile).Error; err != nil {
			return err
		}
		return s.transition(tx, &profile, []string{StatusPendingCityReview}, StatusCityRejected, operatorID, opinion)
	})
}

// AdminSetStatus 管理员直接改写状态，跳过流程限制，但仍记录审核日志以便追溯
func (s *ExpertApprovalService) AdminSetStatus(expertID uint, status string, operatorID uint, opinion string) error {
	validStatuses := map[string]bool{
		StatusDraft: true, StatusPendingOrgReview: true, StatusOrgRejected: true,
		StatusPendingCityReview: true, StatusCityRejected: true, StatusPublished: true,
	}
	if !validStatuses[status] {
		return errors.New("非法的审核状态：" + status)
	}
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var profile ExpertDatabase.ExpertProfile
		if err := tx.Where("id = ?", expertID).First(&profile).Error; err != nil {
			return err
		}
		fromStatus := profile.Status
		if err := tx.Model(&ExpertDatabase.ExpertProfile{}).Where("id = ?", expertID).Update("status", status).Error; err != nil {
			return err
		}
		log := ExpertDatabase.ExpertApprovalLog{
			ExpertId:   expertID,
			FromStatus: fromStatus,
			ToStatus:   status,
			OperatorId: operatorID,
			Opinion:    opinion,
		}
		return tx.Create(&log).Error
	})
}

// assertSameOrg 校验操作人所属单位与专家档案的申报单位一致，防止跨单位审核
func (s *ExpertApprovalService) assertSameOrg(operatorID uint, profileOrgId *uint) error {
	operatorOrgId, err := s.getUserOrgId(operatorID)
	if err != nil {
		return err
	}
	if operatorOrgId == nil || profileOrgId == nil || *operatorOrgId != *profileOrgId {
		return errors.New("无权审核其他单位提交的专家档案")
	}
	return nil
}

// listByFilter 分页查询专家档案的通用底座，供下面三个审核台列表接口复用
func (s *ExpertApprovalService) listByFilter(filter func(*gorm.DB) *gorm.DB, page commonRequest.PageInfo) (list []ExpertDatabase.ExpertProfile, total int64, err error) {
	limit := page.PageSize
	offset := page.PageSize * (page.Page - 1)
	db := filter(global.GVA_DB.Model(&ExpertDatabase.ExpertProfile{}))
	if err = db.Count(&total).Error; err != nil {
		return
	}
	db = db.Order("id desc")
	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}
	err = db.Find(&list).Error
	return
}

// GetMyDrafts 获取当前用户自己申报的专家档案列表（对应前端"我发起的"）
func (s *ExpertApprovalService) GetMyDrafts(userID uint, page commonRequest.PageInfo) ([]ExpertDatabase.ExpertProfile, int64, error) {
	return s.listByFilter(func(db *gorm.DB) *gorm.DB {
		return db.Where("submitted_by = ?", userID)
	}, page)
}

// GetPendingOrgReview 获取待本单位审核的专家档案列表（对应前端"待我审核的"，单位审核员视角）
func (s *ExpertApprovalService) GetPendingOrgReview(userID uint, page commonRequest.PageInfo) ([]ExpertDatabase.ExpertProfile, int64, error) {
	orgId, err := s.getUserOrgId(userID)
	if err != nil {
		return nil, 0, err
	}
	if orgId == nil {
		return nil, 0, errors.New("当前用户未关联单位，无法查看待审核列表")
	}
	return s.listByFilter(func(db *gorm.DB) *gorm.DB {
		return db.Where("org_id = ? AND status = ?", orgId, StatusPendingOrgReview)
	}, page)
}

// GetPendingCityReview 获取待市级审核的专家档案列表（对应前端"待我审核的"，市级审核员视角，不限单位）
func (s *ExpertApprovalService) GetPendingCityReview(page commonRequest.PageInfo) ([]ExpertDatabase.ExpertProfile, int64, error) {
	return s.listByFilter(func(db *gorm.DB) *gorm.DB {
		return db.Where("status = ?", StatusPendingCityReview)
	}, page)
}

// GetExpertApprovalLogList 获取某专家档案的审核流转历史，按时间倒序（对应前端"历史记录"）
func (s *ExpertApprovalService) GetExpertApprovalLogList(info ExpertDatabaseReq.ExpertApprovalLogSearch) (list []ExpertDatabase.ExpertApprovalLog, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&ExpertDatabase.ExpertApprovalLog{}).Where("expert_id = ?", info.ExpertId)
	if err = db.Count(&total).Error; err != nil {
		return
	}
	db = db.Order("id desc")
	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}
	err = db.Find(&list).Error
	return
}
