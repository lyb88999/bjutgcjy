package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// SysOrganization 单位/机构，用于专家库三级审核（个人申报-单位审核-市级审核）的数据归属与审核路由
type SysOrganization struct {
	global.GVA_MODEL
	Name       string            `json:"name" form:"name" gorm:"column:name;comment:单位名称"`                      // 单位名称
	Level      string            `json:"level" form:"level" gorm:"column:level;comment:单位级别(市级/区级/单位级等)"`       // 单位级别
	ParentId   *uint             `json:"parentId" form:"parentId" gorm:"column:parent_id;comment:上级单位ID"`       // 上级单位ID
	RegionCode string            `json:"regionCode" form:"regionCode" gorm:"column:region_code;comment:行政区划代码"` // 行政区划代码
	Status     *bool             `json:"status" form:"status" gorm:"column:status;default:true;comment:启用状态"`   // 启用状态
	Children   []SysOrganization `json:"children" gorm:"-"`                                                     // 子单位，仅用于树形展示，不落库
}

func (SysOrganization) TableName() string {
	return "sys_organizations"
}
