// 自动生成模板SocialServiceDatabase
package EngineeringEducationDatabase

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"time"
)

// 社会服务库 结构体  SocialServiceDatabase
type SocialServiceDatabase struct {
	global.GVA_MODEL
	UniversityName                string    `json:"universityName" form:"universityName" gorm:"primarykey;column:university_name;comment:学校名称;" binding:"required"`                                                  //学校名称
	IndustryInnovationCenterCount *int      `json:"industryInnovationCenterCount" form:"industryInnovationCenterCount" gorm:"column:industry_innovation_center_count;comment:行业协同创新中心数或产学研基地数量;" binding:"required"` //行业协同创新中心数或产学研基地数
	RegionalInnovationCenterCount *int      `json:"regionalInnovationCenterCount" form:"regionalInnovationCenterCount" gorm:"column:regional_innovation_center_count;comment:区域协同创新中心数或产学研基地数量;" binding:"required"` //区域协同创新中心数或产学研基地数
	CreatedBy                     uint      `gorm:"column:created_by;comment:创建者"`
	UpdatedBy                     uint      `gorm:"column:updated_by;comment:更新者"`
	DeletedBy                     uint      `gorm:"column:deleted_by;comment:删除者"`
	AcquisitionTime               time.Time `json:"acquisitionTime" form:"acquisitionTime" gorm:"column:acquisition_time;comment:数据采集时间;"`
}

// TableName 社会服务库 SocialServiceDatabase自定义表名 social_service_database
func (SocialServiceDatabase) TableName() string {
	return "social_service_database"
}
