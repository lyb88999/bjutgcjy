package ExpertDatabase

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase"
	ExpertDatabaseReq "github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase/request"
	"gorm.io/gorm"
)

type ExpertAchievementService struct{}

// CreateExpertAchievement 创建专家成果
func (expertAchievementService *ExpertAchievementService) CreateExpertAchievement(achievement *ExpertDatabase.ExpertAchievement) (err error) {
	return global.GVA_DB.Create(achievement).Error
}

// DeleteExpertAchievement 删除专家成果
func (expertAchievementService *ExpertAchievementService) DeleteExpertAchievement(ID string, userID uint) (err error) {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&ExpertDatabase.ExpertAchievement{}).Where("id = ?", ID).Update("deleted_by", userID).Error; err != nil {
			return err
		}
		return tx.Delete(&ExpertDatabase.ExpertAchievement{}, "id = ?", ID).Error
	})
}

// DeleteExpertAchievementByIds 批量删除专家成果
func (expertAchievementService *ExpertAchievementService) DeleteExpertAchievementByIds(IDs []string, deletedBy uint) (err error) {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&ExpertDatabase.ExpertAchievement{}).Where("id in ?", IDs).Update("deleted_by", deletedBy).Error; err != nil {
			return err
		}
		return tx.Where("id in ?", IDs).Delete(&ExpertDatabase.ExpertAchievement{}).Error
	})
}

// UpdateExpertAchievement 更新专家成果
func (expertAchievementService *ExpertAchievementService) UpdateExpertAchievement(achievement ExpertDatabase.ExpertAchievement) (err error) {
	return global.GVA_DB.Save(&achievement).Error
}

// GetExpertAchievement 根据ID获取专家成果
func (expertAchievementService *ExpertAchievementService) GetExpertAchievement(ID string) (achievement ExpertDatabase.ExpertAchievement, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&achievement).Error
	return
}

// GetExpertAchievementInfoList 分页获取专家成果列表
func (expertAchievementService *ExpertAchievementService) GetExpertAchievementInfoList(info ExpertDatabaseReq.ExpertAchievementSearch) (list []ExpertDatabase.ExpertAchievement, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&ExpertDatabase.ExpertAchievement{})
	var achievements []ExpertDatabase.ExpertAchievement

	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.ExpertId != nil {
		db = db.Where("expert_id = ?", info.ExpertId)
	}
	if info.AchievementType != "" {
		db = db.Where("achievement_type = ?", info.AchievementType)
	}
	if info.Level != "" {
		db = db.Where("level = ?", info.Level)
	}
	if info.Title != "" {
		db = db.Where("title LIKE ?", "%"+info.Title+"%")
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	orderMap := map[string]bool{"created_at": true, "publish_date": true, "ref_count": true}
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

	err = db.Find(&achievements).Error
	return achievements, total, err
}

// GetExpertAchievementsByExpertId 获取某专家的全部成果（不分页），供算分逻辑复用
func (expertAchievementService *ExpertAchievementService) GetExpertAchievementsByExpertId(expertId uint) (list []ExpertDatabase.ExpertAchievement, err error) {
	err = global.GVA_DB.Where("expert_id = ?", expertId).Find(&list).Error
	return
}
