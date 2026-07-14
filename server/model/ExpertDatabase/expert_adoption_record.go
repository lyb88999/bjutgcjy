package ExpertDatabase

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// ExpertAdoptionRecord 决策影响画像明细：成果被内参/专报采用、部门采纳、领导批示、进入政策文件、参与政策起草/咨询论证情况
type ExpertAdoptionRecord struct {
	global.GVA_MODEL
	ExpertId          uint       `json:"expertId" form:"expertId" gorm:"column:expert_id;index;comment:专家ID;" binding:"required"`         //专家ID
	AchievementId     *uint      `json:"achievementId" form:"achievementId" gorm:"column:achievement_id;index;comment:关联成果ID;"`           //可关联到具体成果，也可独立记录
	AdoptionType      string     `json:"adoptionType" form:"adoptionType" gorm:"column:adoption_type;comment:采纳类型;" binding:"required"`   //内参/专报采用、部门采纳、领导批示、进入政策文件、参与政策起草/咨询论证
	AdoptingUnitLevel string     `json:"adoptingUnitLevel" form:"adoptingUnitLevel" gorm:"column:adopting_unit_level;comment:采纳/批示单位级别;"` //国家级/省部级/厅局级/区县级，对应字典 expert_adoption_level
	AdoptingUnitName  string     `json:"adoptingUnitName" form:"adoptingUnitName" gorm:"column:adopting_unit_name;comment:采纳/批示单位名称;"`    //具体单位名称
	AdoptionDate      *time.Time `json:"adoptionDate" form:"adoptionDate" gorm:"column:adoption_date;comment:采纳/批示时间;"`                   //采纳/批示时间
	EvidenceDesc      string     `json:"evidenceDesc" form:"evidenceDesc" gorm:"column:evidence_desc;comment:佐证材料描述;"`                    //佐证材料描述/文号
	AttachmentUrl     string     `json:"attachmentUrl" form:"attachmentUrl" gorm:"column:attachment_url;comment:佐证材料附件;"`                 //佐证材料附件

	CreatedBy uint `gorm:"column:created_by;comment:创建者"`
	UpdatedBy uint `gorm:"column:updated_by;comment:更新者"`
	DeletedBy uint `gorm:"column:deleted_by;comment:删除者"`
}

// TableName 专家决策影响记录 ExpertAdoptionRecord自定义表名 expert_adoption_record
func (ExpertAdoptionRecord) TableName() string {
	return "expert_adoption_record"
}
