package ExpertDatabase

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase"
	ExpertDatabaseReq "github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase/request"
	"gorm.io/gorm"
)

type ExpertAchievementService struct{}

// CreateExpertAchievement 创建专家成果；能不能给这个专家挂成果，跟能不能改这个专家档案本身
// 是同一套数据域规则（见 checkCanModifyExpertSubResource）
func (expertAchievementService *ExpertAchievementService) CreateExpertAchievement(achievement *ExpertDatabase.ExpertAchievement, operatorID uint) (err error) {
	if err = checkCanModifyExpertSubResource(operatorID, achievement.ExpertId); err != nil {
		return err
	}
	if err = global.GVA_DB.Create(achievement).Error; err != nil {
		return err
	}
	expertScoreSvc.recomputeIfPublished(achievement.ExpertId)
	return nil
}

// DeleteExpertAchievement 删除专家成果
func (expertAchievementService *ExpertAchievementService) DeleteExpertAchievement(ID string, userID uint) (err error) {
	var achievement ExpertDatabase.ExpertAchievement
	if err = global.GVA_DB.Where("id = ?", ID).First(&achievement).Error; err != nil {
		return err
	}
	if err = checkCanModifyExpertSubResource(userID, achievement.ExpertId); err != nil {
		return err
	}
	if err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&ExpertDatabase.ExpertAchievement{}).Where("id = ?", ID).Update("deleted_by", userID).Error; err != nil {
			return err
		}
		return tx.Delete(&ExpertDatabase.ExpertAchievement{}, "id = ?", ID).Error
	}); err != nil {
		return err
	}
	expertScoreSvc.recomputeIfPublished(achievement.ExpertId)
	return nil
}

// DeleteExpertAchievementByIds 批量删除专家成果；跟批量删档案一样，只要有一条越出数据域就整体拒绝
func (expertAchievementService *ExpertAchievementService) DeleteExpertAchievementByIds(IDs []string, deletedBy uint) (err error) {
	var achievements []ExpertDatabase.ExpertAchievement
	if err = global.GVA_DB.Where("id in ?", IDs).Find(&achievements).Error; err != nil {
		return err
	}
	for _, a := range achievements {
		if err = checkCanModifyExpertSubResource(deletedBy, a.ExpertId); err != nil {
			return err
		}
	}
	if err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&ExpertDatabase.ExpertAchievement{}).Where("id in ?", IDs).Update("deleted_by", deletedBy).Error; err != nil {
			return err
		}
		return tx.Where("id in ?", IDs).Delete(&ExpertDatabase.ExpertAchievement{}).Error
	}); err != nil {
		return err
	}
	for _, a := range achievements {
		expertScoreSvc.recomputeIfPublished(a.ExpertId)
	}
	return nil
}

// UpdateExpertAchievement 更新专家成果；expertId 以库里现值为准，不允许通过更新请求把一条成果
// 挪去挂在别的专家名下（既没有这个使用场景，放开了也会绕过数据域校验）
func (expertAchievementService *ExpertAchievementService) UpdateExpertAchievement(achievement ExpertDatabase.ExpertAchievement, operatorID uint) (err error) {
	var existing ExpertDatabase.ExpertAchievement
	if err = global.GVA_DB.Where("id = ?", achievement.ID).First(&existing).Error; err != nil {
		return err
	}
	if err = checkCanModifyExpertSubResource(operatorID, existing.ExpertId); err != nil {
		return err
	}
	achievement.ExpertId = existing.ExpertId
	achievement.CreatedAt = existing.CreatedAt
	achievement.CreatedBy = existing.CreatedBy
	if err = global.GVA_DB.Save(&achievement).Error; err != nil {
		return err
	}
	expertScoreSvc.recomputeIfPublished(achievement.ExpertId)
	return nil
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
