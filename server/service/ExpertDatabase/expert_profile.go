package ExpertDatabase

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase"
	ExpertDatabaseReq "github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"gorm.io/gorm"
)

// 专家库三个角色ID：个人申报人/单位审核员/市级审核员。专家主档列表对审核员角色做数据域收敛
// （见 GetExpertProfileInfoList），个人申报人角色ID被本单位账号管理功能用来固定新建账号的角色
const (
	authorityIndividualApplicant = 9001
	authorityOrgReviewer         = 9002
	authorityCityReviewer        = 9003
)

type ExpertProfileService struct{}

// CreateExpertProfile 创建专家主档；审核开关关闭时跳过草稿状态，直接落库为已发布并重算得分
func (expertProfileService *ExpertProfileService) CreateExpertProfile(profile *ExpertDatabase.ExpertProfile) (err error) {
	skipReview := !reviewRequired()
	if skipReview {
		profile.Status = StatusPublished
		profile.ReviewBypassed = true
	}
	if err = global.GVA_DB.Create(profile).Error; err != nil {
		return err
	}
	if skipReview {
		expertScoreSvc.recomputeIfPublished(profile.ID)
	}
	return nil
}

// DeleteExpertProfile 删除专家主档
func (expertProfileService *ExpertProfileService) DeleteExpertProfile(ID string, userID uint) (err error) {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&ExpertDatabase.ExpertProfile{}).Where("id = ?", ID).Update("deleted_by", userID).Error; err != nil {
			return err
		}
		return tx.Delete(&ExpertDatabase.ExpertProfile{}, "id = ?", ID).Error
	})
}

// DeleteExpertProfileByIds 批量删除专家主档
func (expertProfileService *ExpertProfileService) DeleteExpertProfileByIds(IDs []string, deletedBy uint) (err error) {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&ExpertDatabase.ExpertProfile{}).Where("id in ?", IDs).Update("deleted_by", deletedBy).Error; err != nil {
			return err
		}
		return tx.Where("id in ?", IDs).Delete(&ExpertDatabase.ExpertProfile{}).Error
	})
}

// UpdateExpertProfile 更新专家主档
func (expertProfileService *ExpertProfileService) UpdateExpertProfile(profile ExpertDatabase.ExpertProfile) (err error) {
	return global.GVA_DB.Save(&profile).Error
}

// GetExpertProfile 根据ID获取专家主档
func (expertProfileService *ExpertProfileService) GetExpertProfile(ID string) (profile ExpertDatabase.ExpertProfile, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&profile).Error
	return
}

// GetExpertProfileInfoList 分页获取专家主档列表
//
// 单位审核员/市级审核员在审核台之外，还能通过这个列表页查到全库专家——原先没有按角色收敛，
// 相当于谁都能翻到别人、别的单位还没提交的草稿，跟"审核台只显示轮到自己审的记录"这个预期不一致。
// 这里补上：单位审核员只看得到已发布的 + 本单位已提交（非草稿）的；市级审核员只看得到已发布的 +
// 待市级审核的；管理员和个人申报人不受影响，保持原有的全量可见。
func (expertProfileService *ExpertProfileService) GetExpertProfileInfoList(info ExpertDatabaseReq.ExpertProfileSearch, operatorID uint) (list []ExpertDatabase.ExpertProfile, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&ExpertDatabase.ExpertProfile{})
	var profiles []ExpertDatabase.ExpertProfile

	var operator system.SysUser
	if err = global.GVA_DB.Where("id = ?", operatorID).First(&operator).Error; err != nil {
		return
	}
	switch operator.AuthorityId {
	case authorityOrgReviewer:
		if operator.OrgId != nil {
			db = db.Where("status = ? OR (org_id = ? AND status <> ?)", StatusPublished, *operator.OrgId, StatusDraft)
		} else {
			db = db.Where("status = ?", StatusPublished)
		}
	case authorityCityReviewer:
		db = db.Where("status = ? OR status = ?", StatusPublished, StatusPendingCityReview)
	}

	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.Name != "" {
		db = db.Where("name LIKE ?", "%"+info.Name+"%")
	}
	if info.UnitName != "" {
		db = db.Where("unit_name LIKE ?", "%"+info.UnitName+"%")
	}
	if info.TechTitle != "" {
		db = db.Where("tech_title = ?", info.TechTitle)
	}
	if info.DisciplineL1 != "" {
		db = db.Where("discipline_l1 = ?", info.DisciplineL1)
	}
	if info.DisciplineL2 != "" {
		db = db.Where("discipline_l2 = ?", info.DisciplineL2)
	}
	if info.ResearchKeywords != "" {
		db = db.Where("research_keywords LIKE ?", "%"+info.ResearchKeywords+"%")
	}
	if info.Status != "" {
		db = db.Where("status = ?", info.Status)
	}
	if info.OrgId != nil {
		db = db.Where("org_id = ?", info.OrgId)
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	orderMap := map[string]bool{
		"created_at":        true,
		"achievement_score": true,
		"influence_score":   true,
		"composite_score":   true,
	}
	if orderMap[info.Sort] {
		orderStr := info.Sort
		if info.Order == "descending" {
			orderStr += " desc"
		}
		db = db.Order(orderStr)
	} else {
		db = db.Order("id desc")
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&profiles).Error
	return profiles, total, err
}
