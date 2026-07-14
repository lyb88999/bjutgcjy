package ExpertDatabase

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase"
	ExpertDatabaseReq "github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase/request"
	"gorm.io/gorm"
)

type ExpertProfileService struct{}

// CreateExpertProfile 创建专家主档
func (expertProfileService *ExpertProfileService) CreateExpertProfile(profile *ExpertDatabase.ExpertProfile) (err error) {
	return global.GVA_DB.Create(profile).Error
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
func (expertProfileService *ExpertProfileService) GetExpertProfileInfoList(info ExpertDatabaseReq.ExpertProfileSearch) (list []ExpertDatabase.ExpertProfile, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&ExpertDatabase.ExpertProfile{})
	var profiles []ExpertDatabase.ExpertProfile

	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.Name != "" {
		db = db.Where("name LIKE ?", "%"+info.Name+"%")
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
