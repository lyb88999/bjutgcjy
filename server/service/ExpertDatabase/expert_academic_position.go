package ExpertDatabase

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase"
	ExpertDatabaseReq "github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase/request"
	"gorm.io/gorm"
)

type ExpertAcademicPositionService struct{}

// CreateExpertAcademicPosition 创建专家学术兼职
func (expertAcademicPositionService *ExpertAcademicPositionService) CreateExpertAcademicPosition(position *ExpertDatabase.ExpertAcademicPosition) (err error) {
	return global.GVA_DB.Create(position).Error
}

// DeleteExpertAcademicPosition 删除专家学术兼职
func (expertAcademicPositionService *ExpertAcademicPositionService) DeleteExpertAcademicPosition(ID string, userID uint) (err error) {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&ExpertDatabase.ExpertAcademicPosition{}).Where("id = ?", ID).Update("deleted_by", userID).Error; err != nil {
			return err
		}
		return tx.Delete(&ExpertDatabase.ExpertAcademicPosition{}, "id = ?", ID).Error
	})
}

// DeleteExpertAcademicPositionByIds 批量删除专家学术兼职
func (expertAcademicPositionService *ExpertAcademicPositionService) DeleteExpertAcademicPositionByIds(IDs []string, deletedBy uint) (err error) {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&ExpertDatabase.ExpertAcademicPosition{}).Where("id in ?", IDs).Update("deleted_by", deletedBy).Error; err != nil {
			return err
		}
		return tx.Where("id in ?", IDs).Delete(&ExpertDatabase.ExpertAcademicPosition{}).Error
	})
}

// UpdateExpertAcademicPosition 更新专家学术兼职
func (expertAcademicPositionService *ExpertAcademicPositionService) UpdateExpertAcademicPosition(position ExpertDatabase.ExpertAcademicPosition) (err error) {
	return global.GVA_DB.Save(&position).Error
}

// GetExpertAcademicPosition 根据ID获取专家学术兼职
func (expertAcademicPositionService *ExpertAcademicPositionService) GetExpertAcademicPosition(ID string) (position ExpertDatabase.ExpertAcademicPosition, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&position).Error
	return
}

// GetExpertAcademicPositionInfoList 分页获取专家学术兼职列表
func (expertAcademicPositionService *ExpertAcademicPositionService) GetExpertAcademicPositionInfoList(info ExpertDatabaseReq.ExpertAcademicPositionSearch) (list []ExpertDatabase.ExpertAcademicPosition, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&ExpertDatabase.ExpertAcademicPosition{})
	var positions []ExpertDatabase.ExpertAcademicPosition

	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.ExpertId != nil {
		db = db.Where("expert_id = ?", info.ExpertId)
	}
	if info.PositionType != "" {
		db = db.Where("position_type = ?", info.PositionType)
	}
	if info.OrganizationName != "" {
		db = db.Where("organization_name LIKE ?", "%"+info.OrganizationName+"%")
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	orderMap := map[string]bool{"created_at": true, "start_date": true}
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

	err = db.Find(&positions).Error
	return positions, total, err
}
