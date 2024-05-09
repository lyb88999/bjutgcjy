package EngineeringEducationDatabase

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/EngineeringEducationDatabase"
    EngineeringEducationDatabaseReq "github.com/flipped-aurora/gin-vue-admin/server/model/EngineeringEducationDatabase/request"
    "gorm.io/gorm"
)

type BasicInfomationDatabaseService struct {
}

// CreateBasicInfomationDatabase 创建基础信息库记录
// Author [piexlmax](https://github.com/piexlmax)
func (BIDService *BasicInfomationDatabaseService) CreateBasicInfomationDatabase(BID *EngineeringEducationDatabase.BasicInfomationDatabase) (err error) {
	err = global.GVA_DB.Create(BID).Error
	return err
}

// DeleteBasicInfomationDatabase 删除基础信息库记录
// Author [piexlmax](https://github.com/piexlmax)
func (BIDService *BasicInfomationDatabaseService)DeleteBasicInfomationDatabase(ID string,userID uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
	    if err := tx.Model(&EngineeringEducationDatabase.BasicInfomationDatabase{}).Where("id = ?", ID).Update("deleted_by", userID).Error; err != nil {
              return err
        }
        if err = tx.Delete(&EngineeringEducationDatabase.BasicInfomationDatabase{},"id = ?",ID).Error; err != nil {
              return err
        }
        return nil
	})
	return err
}

// DeleteBasicInfomationDatabaseByIds 批量删除基础信息库记录
// Author [piexlmax](https://github.com/piexlmax)
func (BIDService *BasicInfomationDatabaseService)DeleteBasicInfomationDatabaseByIds(IDs []string,deleted_by uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
	    if err := tx.Model(&EngineeringEducationDatabase.BasicInfomationDatabase{}).Where("id in ?", IDs).Update("deleted_by", deleted_by).Error; err != nil {
            return err
        }
        if err := tx.Where("id in ?", IDs).Delete(&EngineeringEducationDatabase.BasicInfomationDatabase{}).Error; err != nil {
            return err
        }
        return nil
    })
	return err
}

// UpdateBasicInfomationDatabase 更新基础信息库记录
// Author [piexlmax](https://github.com/piexlmax)
func (BIDService *BasicInfomationDatabaseService)UpdateBasicInfomationDatabase(BID EngineeringEducationDatabase.BasicInfomationDatabase) (err error) {
	err = global.GVA_DB.Save(&BID).Error
	return err
}

// GetBasicInfomationDatabase 根据ID获取基础信息库记录
// Author [piexlmax](https://github.com/piexlmax)
func (BIDService *BasicInfomationDatabaseService)GetBasicInfomationDatabase(ID string) (BID EngineeringEducationDatabase.BasicInfomationDatabase, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&BID).Error
	return
}

// GetBasicInfomationDatabaseInfoList 分页获取基础信息库记录
// Author [piexlmax](https://github.com/piexlmax)
func (BIDService *BasicInfomationDatabaseService)GetBasicInfomationDatabaseInfoList(info EngineeringEducationDatabaseReq.BasicInfomationDatabaseSearch) (list []EngineeringEducationDatabase.BasicInfomationDatabase, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
    // 创建db
	db := global.GVA_DB.Model(&EngineeringEducationDatabase.BasicInfomationDatabase{})
    var BIDs []EngineeringEducationDatabase.BasicInfomationDatabase
    // 如果有条件搜索 下方会自动创建搜索语句
    if info.StartCreatedAt !=nil && info.EndCreatedAt !=nil {
     db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
    }
	err = db.Count(&total).Error
	if err!=nil {
    	return
    }
        var OrderStr string
        orderMap := make(map[string]bool)
         	orderMap["school_code"] = true
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
	
	err = db.Find(&BIDs).Error
	return  BIDs, total, err
}
