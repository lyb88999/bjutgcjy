package EngineeringEducationDatabase

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/EngineeringEducationDatabase"
	EngineeringEducationDatabaseReq "github.com/flipped-aurora/gin-vue-admin/server/model/EngineeringEducationDatabase/request"
	"gorm.io/gorm"
)

type PolicyDatabaseService struct {
}

// CreatePolicyDatabase 创建政策数据库记录
// Author [piexlmax](https://github.com/piexlmax)
func (PDService *PolicyDatabaseService) CreatePolicyDatabase(PD *EngineeringEducationDatabase.PolicyDatabase) (err error) {
	err = global.GVA_DB.Create(PD).Error
	return err
}

// DeletePolicyDatabase 删除政策数据库记录
// Author [piexlmax](https://github.com/piexlmax)
func (PDService *PolicyDatabaseService) DeletePolicyDatabase(ID string, userID uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&EngineeringEducationDatabase.PolicyDatabase{}).Where("id = ?", ID).Update("deleted_by", userID).Error; err != nil {
			return err
		}
		if err = tx.Delete(&EngineeringEducationDatabase.PolicyDatabase{}, "id = ?", ID).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// DeletePolicyDatabaseByIds 批量删除政策数据库记录
// Author [piexlmax](https://github.com/piexlmax)
func (PDService *PolicyDatabaseService) DeletePolicyDatabaseByIds(IDs []string, deleted_by uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&EngineeringEducationDatabase.PolicyDatabase{}).Where("id in ?", IDs).Update("deleted_by", deleted_by).Error; err != nil {
			return err
		}
		if err := tx.Where("id in ?", IDs).Delete(&EngineeringEducationDatabase.PolicyDatabase{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// UpdatePolicyDatabase 更新政策数据库记录
// Author [piexlmax](https://github.com/piexlmax)
func (PDService *PolicyDatabaseService) UpdatePolicyDatabase(PD EngineeringEducationDatabase.PolicyDatabase) (err error) {
	err = global.GVA_DB.Save(&PD).Error
	return err
}

// GetPolicyDatabase 根据ID获取政策数据库记录
// Author [piexlmax](https://github.com/piexlmax)
func (PDService *PolicyDatabaseService) GetPolicyDatabase(ID string) (PD EngineeringEducationDatabase.PolicyDatabase, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&PD).Error
	return
}

// GetPolicyDatabaseInfoList 分页获取政策数据库记录
// Author [piexlmax](https://github.com/piexlmax)
func (PDService *PolicyDatabaseService) GetPolicyDatabaseInfoList(info EngineeringEducationDatabaseReq.PolicyDatabaseSearch) (list []EngineeringEducationDatabase.PolicyDatabase, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&EngineeringEducationDatabase.PolicyDatabase{})
	var PDs []EngineeringEducationDatabase.PolicyDatabase
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	// 按照政策发布时间查询
	//if info.StartPolicyReleaseDate != nil && info.EndPolicyReleaseDate != nil {
	//	db = db.Where("policy_release_date BETWEEN ? AND ? ", info.StartPolicyReleaseDate, info.EndPolicyReleaseDate)
	//}
	if info.PolicyCountryOrRegion != "" {
		db = db.Where("policy_country_or_region = ?", info.PolicyCountryOrRegion)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	var OrderStr string
	orderMap := make(map[string]bool)
	orderMap["policy_release_date"] = true
	if orderMap[info.Sort] {
		OrderStr = info.Sort
		if info.Order == "descending" {
			OrderStr = OrderStr + " desc"
		}
		db = db.Order(OrderStr)
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&PDs).Error
	return PDs, total, err
}
