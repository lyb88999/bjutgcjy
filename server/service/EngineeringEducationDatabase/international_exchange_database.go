package EngineeringEducationDatabase

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/EngineeringEducationDatabase"
    EngineeringEducationDatabaseReq "github.com/flipped-aurora/gin-vue-admin/server/model/EngineeringEducationDatabase/request"
    "gorm.io/gorm"
)

type InternationalExchangeDatabaseService struct {
}

// CreateInternationalExchangeDatabase 创建国际交流库记录
// Author [piexlmax](https://github.com/piexlmax)
func (IEDService *InternationalExchangeDatabaseService) CreateInternationalExchangeDatabase(IED *EngineeringEducationDatabase.InternationalExchangeDatabase) (err error) {
	err = global.GVA_DB.Create(IED).Error
	return err
}

// DeleteInternationalExchangeDatabase 删除国际交流库记录
// Author [piexlmax](https://github.com/piexlmax)
func (IEDService *InternationalExchangeDatabaseService)DeleteInternationalExchangeDatabase(ID string,userID uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
	    if err := tx.Model(&EngineeringEducationDatabase.InternationalExchangeDatabase{}).Where("id = ?", ID).Update("deleted_by", userID).Error; err != nil {
              return err
        }
        if err = tx.Delete(&EngineeringEducationDatabase.InternationalExchangeDatabase{},"id = ?",ID).Error; err != nil {
              return err
        }
        return nil
	})
	return err
}

// DeleteInternationalExchangeDatabaseByIds 批量删除国际交流库记录
// Author [piexlmax](https://github.com/piexlmax)
func (IEDService *InternationalExchangeDatabaseService)DeleteInternationalExchangeDatabaseByIds(IDs []string,deleted_by uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
	    if err := tx.Model(&EngineeringEducationDatabase.InternationalExchangeDatabase{}).Where("id in ?", IDs).Update("deleted_by", deleted_by).Error; err != nil {
            return err
        }
        if err := tx.Where("id in ?", IDs).Delete(&EngineeringEducationDatabase.InternationalExchangeDatabase{}).Error; err != nil {
            return err
        }
        return nil
    })
	return err
}

// UpdateInternationalExchangeDatabase 更新国际交流库记录
// Author [piexlmax](https://github.com/piexlmax)
func (IEDService *InternationalExchangeDatabaseService)UpdateInternationalExchangeDatabase(IED EngineeringEducationDatabase.InternationalExchangeDatabase) (err error) {
	err = global.GVA_DB.Save(&IED).Error
	return err
}

// GetInternationalExchangeDatabase 根据ID获取国际交流库记录
// Author [piexlmax](https://github.com/piexlmax)
func (IEDService *InternationalExchangeDatabaseService)GetInternationalExchangeDatabase(ID string) (IED EngineeringEducationDatabase.InternationalExchangeDatabase, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&IED).Error
	return
}

// GetInternationalExchangeDatabaseInfoList 分页获取国际交流库记录
// Author [piexlmax](https://github.com/piexlmax)
func (IEDService *InternationalExchangeDatabaseService)GetInternationalExchangeDatabaseInfoList(info EngineeringEducationDatabaseReq.InternationalExchangeDatabaseSearch) (list []EngineeringEducationDatabase.InternationalExchangeDatabase, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
    // 创建db
	db := global.GVA_DB.Model(&EngineeringEducationDatabase.InternationalExchangeDatabase{})
    var IEDs []EngineeringEducationDatabase.InternationalExchangeDatabase
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
         	orderMap["international_conferences_hosted"] = true
         	orderMap["foreign_faculty_count"] = true
         	orderMap["international_students_count"] = true
         	orderMap["students_abroad_exchange_count"] = true
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
	
	err = db.Find(&IEDs).Error
	return  IEDs, total, err
}
