// 自动生成模板InviteCode
package util

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	
	
)

// 邀请码 结构体  InviteCode
type InviteCode struct {
 global.GVA_MODEL 
      Code  string `json:"code" form:"code" gorm:"column:code;comment:生成的邀请码;size:255;"`  //code字段 
      Status  *bool `json:"status" form:"status" gorm:"column:status;comment:是否已被使用;"`  //status字段 
      Uuid  string `json:"uuid" form:"uuid" gorm:"column:uuid;comment:生成该邀请码的用户ID;"`  //uuid字段 
      CreatedBy  uint   `gorm:"column:created_by;comment:创建者"`
      UpdatedBy  uint   `gorm:"column:updated_by;comment:更新者"`
      DeletedBy  uint   `gorm:"column:deleted_by;comment:删除者"`
}


// TableName 邀请码 InviteCode自定义表名 invite_code
func (InviteCode) TableName() string {
  return "invite_code"
}

