package ExpertDatabase

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// ExpertAchievement 研究成果画像明细：著作/论文/研究报告/决策咨询成果/获奖成果/课题项目统一用一张表，按 AchievementType 区分
type ExpertAchievement struct {
	global.GVA_MODEL
	ExpertId        uint       `json:"expertId" form:"expertId" gorm:"column:expert_id;index;comment:专家ID;" binding:"required"`                //专家ID
	AchievementType string     `json:"achievementType" form:"achievementType" gorm:"column:achievement_type;comment:成果类型;" binding:"required"` //著作/论文/研究报告/决策咨询成果/获奖成果/课题项目
	Title           string     `json:"title" form:"title" gorm:"column:title;comment:成果名称;" binding:"required"`                                //成果名称
	Level           string     `json:"level" form:"level" gorm:"column:level;comment:成果级别;"`                                                   //国家级/省部级/厅局级/一般级，对应字典 expert_achievement_level
	LevelSource     string     `json:"levelSource" form:"levelSource" gorm:"column:level_source;comment:级别认定来源;"`                              //如"课题下达单位"、"获奖颁发单位"
	PublishOrg      string     `json:"publishOrg" form:"publishOrg" gorm:"column:publish_org;comment:发表/出版/立项单位;"`                             //发表/出版/立项单位
	PublishDate     *time.Time `json:"publishDate" form:"publishDate" gorm:"column:publish_date;comment:发表/出版/立项时间;"`                          //发表/出版/立项时间
	Keywords        string     `json:"keywords" form:"keywords" gorm:"column:keywords;comment:关键词;"`                                           //用于相关性检索的关键词/摘要分词结果
	RefCount        int        `json:"refCount" form:"refCount" gorm:"column:ref_count;default:0;comment:检索相关次数;"`                             //检索相关次数
	Remark          string     `json:"remark" form:"remark" gorm:"column:remark;comment:备注;"`                                                  //备注

	CreatedBy uint `gorm:"column:created_by;comment:创建者"`
	UpdatedBy uint `gorm:"column:updated_by;comment:更新者"`
	DeletedBy uint `gorm:"column:deleted_by;comment:删除者"`
}

// TableName 专家成果 ExpertAchievement自定义表名 expert_achievement
func (ExpertAchievement) TableName() string {
	return "expert_achievement"
}
