package ExpertDatabase

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase"
	ExpertDatabaseReq "github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase/request"
	"gorm.io/gorm"
)

type ExpertAdoptionRecordService struct{}

// CreateExpertAdoptionRecord 创建专家决策影响记录
func (expertAdoptionRecordService *ExpertAdoptionRecordService) CreateExpertAdoptionRecord(record *ExpertDatabase.ExpertAdoptionRecord) (err error) {
	return global.GVA_DB.Create(record).Error
}

// DeleteExpertAdoptionRecord 删除专家决策影响记录
func (expertAdoptionRecordService *ExpertAdoptionRecordService) DeleteExpertAdoptionRecord(ID string, userID uint) (err error) {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&ExpertDatabase.ExpertAdoptionRecord{}).Where("id = ?", ID).Update("deleted_by", userID).Error; err != nil {
			return err
		}
		return tx.Delete(&ExpertDatabase.ExpertAdoptionRecord{}, "id = ?", ID).Error
	})
}

// DeleteExpertAdoptionRecordByIds 批量删除专家决策影响记录
func (expertAdoptionRecordService *ExpertAdoptionRecordService) DeleteExpertAdoptionRecordByIds(IDs []string, deletedBy uint) (err error) {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&ExpertDatabase.ExpertAdoptionRecord{}).Where("id in ?", IDs).Update("deleted_by", deletedBy).Error; err != nil {
			return err
		}
		return tx.Where("id in ?", IDs).Delete(&ExpertDatabase.ExpertAdoptionRecord{}).Error
	})
}

// UpdateExpertAdoptionRecord 更新专家决策影响记录
func (expertAdoptionRecordService *ExpertAdoptionRecordService) UpdateExpertAdoptionRecord(record ExpertDatabase.ExpertAdoptionRecord) (err error) {
	return global.GVA_DB.Save(&record).Error
}

// GetExpertAdoptionRecord 根据ID获取专家决策影响记录
func (expertAdoptionRecordService *ExpertAdoptionRecordService) GetExpertAdoptionRecord(ID string) (record ExpertDatabase.ExpertAdoptionRecord, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&record).Error
	return
}

// GetExpertAdoptionRecordInfoList 分页获取专家决策影响记录列表
func (expertAdoptionRecordService *ExpertAdoptionRecordService) GetExpertAdoptionRecordInfoList(info ExpertDatabaseReq.ExpertAdoptionRecordSearch) (list []ExpertDatabase.ExpertAdoptionRecord, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&ExpertDatabase.ExpertAdoptionRecord{})
	var records []ExpertDatabase.ExpertAdoptionRecord

	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.ExpertId != nil {
		db = db.Where("expert_id = ?", info.ExpertId)
	}
	if info.AdoptionType != "" {
		db = db.Where("adoption_type = ?", info.AdoptionType)
	}
	if info.AdoptingUnitLevel != "" {
		db = db.Where("adopting_unit_level = ?", info.AdoptingUnitLevel)
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	orderMap := map[string]bool{"created_at": true, "adoption_date": true}
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

	err = db.Find(&records).Error
	return records, total, err
}

// GetExpertAdoptionRecordsByExpertId 获取某专家的全部决策影响记录（不分页），供算分逻辑复用
func (expertAdoptionRecordService *ExpertAdoptionRecordService) GetExpertAdoptionRecordsByExpertId(expertId uint) (list []ExpertDatabase.ExpertAdoptionRecord, err error) {
	err = global.GVA_DB.Where("expert_id = ?", expertId).Find(&list).Error
	return
}
