package util

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/util"
    utilReq "github.com/flipped-aurora/gin-vue-admin/server/model/util/request"
    "gorm.io/gorm"
)

type InviteCodeService struct {
}

// CreateInviteCode 创建邀请码记录
// Author [piexlmax](https://github.com/piexlmax)
func (inviteCodeService *InviteCodeService) CreateInviteCode(inviteCode *util.InviteCode) (err error) {
	err = global.GVA_DB.Create(inviteCode).Error
	return err
}

// DeleteInviteCode 删除邀请码记录
// Author [piexlmax](https://github.com/piexlmax)
func (inviteCodeService *InviteCodeService)DeleteInviteCode(ID string,userID uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
	    if err := tx.Model(&util.InviteCode{}).Where("id = ?", ID).Update("deleted_by", userID).Error; err != nil {
              return err
        }
        if err = tx.Delete(&util.InviteCode{},"id = ?",ID).Error; err != nil {
              return err
        }
        return nil
	})
	return err
}

// DeleteInviteCodeByIds 批量删除邀请码记录
// Author [piexlmax](https://github.com/piexlmax)
func (inviteCodeService *InviteCodeService)DeleteInviteCodeByIds(IDs []string,deleted_by uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
	    if err := tx.Model(&util.InviteCode{}).Where("id in ?", IDs).Update("deleted_by", deleted_by).Error; err != nil {
            return err
        }
        if err := tx.Where("id in ?", IDs).Delete(&util.InviteCode{}).Error; err != nil {
            return err
        }
        return nil
    })
	return err
}

// UpdateInviteCode 更新邀请码记录
// Author [piexlmax](https://github.com/piexlmax)
func (inviteCodeService *InviteCodeService)UpdateInviteCode(inviteCode util.InviteCode) (err error) {
	err = global.GVA_DB.Save(&inviteCode).Error
	return err
}

// GetInviteCode 根据ID获取邀请码记录
// Author [piexlmax](https://github.com/piexlmax)
func (inviteCodeService *InviteCodeService)GetInviteCode(ID string) (inviteCode util.InviteCode, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&inviteCode).Error
	return
}

// GetInviteCodeInfoList 分页获取邀请码记录
// Author [piexlmax](https://github.com/piexlmax)
func (inviteCodeService *InviteCodeService)GetInviteCodeInfoList(info utilReq.InviteCodeSearch) (list []util.InviteCode, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
    // 创建db
	db := global.GVA_DB.Model(&util.InviteCode{})
    var inviteCodes []util.InviteCode
    // 如果有条件搜索 下方会自动创建搜索语句
    if info.StartCreatedAt !=nil && info.EndCreatedAt !=nil {
     db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
    }
    if info.Status != nil {
        db = db.Where("status = ?",info.Status)
    }
	err = db.Count(&total).Error
	if err!=nil {
    	return
    }

	if limit != 0 {
       db = db.Limit(limit).Offset(offset)
    }
	
	err = db.Find(&inviteCodes).Error
	return  inviteCodes, total, err
}
