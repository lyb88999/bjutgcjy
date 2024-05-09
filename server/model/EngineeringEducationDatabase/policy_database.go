// 自动生成模板PolicyDatabase
package EngineeringEducationDatabase

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"time"
)

// 政策数据库 结构体  PolicyDatabase
type PolicyDatabase struct {
	global.GVA_MODEL
	PolicyName            string    `json:"policyName" form:"policyName" gorm:"column:policy_name;comment:政策名称;" binding:"required"`                                       //政策名称
	PolicyReleaseDate     string    `json:"policyReleaseDate" form:"policyReleaseDate" gorm:"column:policy_release_date;comment:政策发布时间;" binding:"required"`               //政策发布时间
	PolicyCountryOrRegion string    `json:"policyCountryOrRegion" form:"policyCountryOrRegion" gorm:"column:policy_country_or_region;comment:政策国家或地区;" binding:"required"` //政策国家或地区
	DataSource            string    `json:"dataSource" form:"dataSource" gorm:"column:data_source;comment:数据来源;" binding:"required"`                                       //数据来源
	PolicyContent         string    `json:"policyContent" form:"policyContent" gorm:"column:policy_content;comment:政策内容;type:text;" binding:"required"`                    //政策内容
	AdditionalRemarks     string    `json:"additionalRemarks" form:"additionalRemarks" gorm:"column:additional_remarks;comment:其他备注;"`                                     //其他备注
	CreatedBy             uint      `gorm:"column:created_by;comment:创建者"`
	UpdatedBy             uint      `gorm:"column:updated_by;comment:更新者"`
	DeletedBy             uint      `gorm:"column:deleted_by;comment:删除者"`
	AcquisitionTime       time.Time `json:"acquisitionTime" form:"acquisitionTime" gorm:"column:acquisition_time;comment:数据采集时间;"`
}

// TableName 政策数据库 PolicyDatabase自定义表名 policy_database
func (PolicyDatabase) TableName() string {
	return "policy_database"
}
