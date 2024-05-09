package EngineeringEducationDatabase

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/EngineeringEducationDatabase"
    EngineeringEducationDatabaseReq "github.com/flipped-aurora/gin-vue-admin/server/model/EngineeringEducationDatabase/request"
    "gorm.io/gorm"
)

type SocialServiceDatabaseService struct {
}

// CreateSocialServiceDatabase 创建社会服务库记录
// Author [piexlmax](https://github.com/piexlmax)
func (SSDService *SocialServiceDatabaseService) CreateSocialServiceDatabase(SSD *EngineeringEducationDatabase.SocialServiceDatabase) (err error) {
	err = global.GVA_DB.Create(SSD).Error
	return err
}

// DeleteSocialServiceDatabase 删除社会服务库记录
// Author [piexlmax](https://github.com/piexlmax)
func (SSDService *SocialServiceDatabaseService)DeleteSocialServiceDatabase(ID string,userID uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
	    if err := tx.Model(&EngineeringEducationDatabase.SocialServiceDatabase{}).Where("id = ?", ID).Update("deleted_by", userID).Error; err != nil {
              return err
        }
        if err = tx.Delete(&EngineeringEducationDatabase.SocialServiceDatabase{},"id = ?",ID).Error; err != nil {
              return err
        }
        return nil
	})
	return err
}

// DeleteSocialServiceDatabaseByIds 批量删除社会服务库记录
// Author [piexlmax](https://github.com/piexlmax)
func (SSDService *SocialServiceDatabaseService)DeleteSocialServiceDatabaseByIds(IDs []string,deleted_by uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
	    if err := tx.Model(&EngineeringEducationDatabase.SocialServiceDatabase{}).Where("id in ?", IDs).Update("deleted_by", deleted_by).Error; err != nil {
            return err
        }
        if err := tx.Where("id in ?", IDs).Delete(&EngineeringEducationDatabase.SocialServiceDatabase{}).Error; err != nil {
            return err
        }
        return nil
    })
	return err
}

// UpdateSocialServiceDatabase 更新社会服务库记录
// Author [piexlmax](https://github.com/piexlmax)
func (SSDService *SocialServiceDatabaseService)UpdateSocialServiceDatabase(SSD EngineeringEducationDatabase.SocialServiceDatabase) (err error) {
	err = global.GVA_DB.Save(&SSD).Error
	return err
}

// GetSocialServiceDatabase 根据ID获取社会服务库记录
// Author [piexlmax](https://github.com/piexlmax)
func (SSDService *SocialServiceDatabaseService)GetSocialServiceDatabase(ID string) (SSD EngineeringEducationDatabase.SocialServiceDatabase, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&SSD).Error
	return
}

// GetSocialServiceDatabaseInfoList 分页获取社会服务库记录
// Author [piexlmax](https://github.com/piexlmax)
func (SSDService *SocialServiceDatabaseService)GetSocialServiceDatabaseInfoList(info EngineeringEducationDatabaseReq.SocialServiceDatabaseSearch) (list []EngineeringEducationDatabase.SocialServiceDatabase, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
    // 创建db
	db := global.GVA_DB.Model(&EngineeringEducationDatabase.SocialServiceDatabase{})
    var SSDs []EngineeringEducationDatabase.SocialServiceDatabase
    // 如果有条件搜索 下方会自动创建搜索语句
    if info.StartCreatedAt !=nil && info.EndCreatedAt !=nil {
     db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
    }
    if info.UniversityName != "" {
        db = db.Where("university_name LIKE ?","%"+ info.UniversityName+"%")
    }
	err = db.Count(&total).Error
	if err!=nil {
    	return
    }
        var OrderStr string
        orderMap := make(map[string]bool)
         	orderMap["industry_innovation_center_count"] = true
         	orderMap["regional_innovation_center_count"] = true
         	orderMap["university_name"] = true
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
	
	err = db.Find(&SSDs).Error
	return  SSDs, total, err
}
