package ExpertDatabase

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// ExpertAcademicPosition 学术兼职明细：学术团体任职、专业委员会任职、政府决策咨询机构兼职等
type ExpertAcademicPosition struct {
	global.GVA_MODEL
	ExpertId         uint       `json:"expertId" form:"expertId" gorm:"column:expert_id;index;comment:专家ID;" binding:"required"`                     //专家ID
	PositionType     string     `json:"positionType" form:"positionType" gorm:"column:position_type;comment:兼职类型;"`                                  //学术团体任职/专业委员会任职/政府决策咨询机构兼职等
	OrganizationName string     `json:"organizationName" form:"organizationName" gorm:"column:organization_name;comment:任职机构名称;" binding:"required"` //任职机构名称
	PositionTitle    string     `json:"positionTitle" form:"positionTitle" gorm:"column:position_title;comment:担任职务;"`                               //担任职务
	StartDate        *time.Time `json:"startDate" form:"startDate" gorm:"column:start_date;comment:任职起始时间;"`                                         //任职起始时间
	EndDate          *time.Time `json:"endDate" form:"endDate" gorm:"column:end_date;comment:任职结束时间;"`                                               //任职结束时间（在任填空）
	Remark           string     `json:"remark" form:"remark" gorm:"column:remark;comment:备注;"`                                                       //备注

	CreatedBy uint `gorm:"column:created_by;comment:创建者"`
	UpdatedBy uint `gorm:"column:updated_by;comment:更新者"`
	DeletedBy uint `gorm:"column:deleted_by;comment:删除者"`
}

// TableName 专家学术兼职 ExpertAcademicPosition自定义表名 expert_academic_position
func (ExpertAcademicPosition) TableName() string {
	return "expert_academic_position"
}
