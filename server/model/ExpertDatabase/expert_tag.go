package ExpertDatabase

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// ExpertTag 标签库：学科/关键词/政策领域/区域/研究方法等多值标签的统一存储，配合 ExpertTagRelation 做多对多关联
type ExpertTag struct {
	global.GVA_MODEL
	TagType  string `json:"tagType" form:"tagType" gorm:"column:tag_type;index;comment:标签类型;" binding:"required"` //discipline_l1/keyword/policy_field/region/method...
	TagValue string `json:"tagValue" form:"tagValue" gorm:"column:tag_value;comment:标签值;" binding:"required"`     //标签值

	CreatedBy uint `gorm:"column:created_by;comment:创建者"`
	UpdatedBy uint `gorm:"column:updated_by;comment:更新者"`
	DeletedBy uint `gorm:"column:deleted_by;comment:删除者"`
}

// TableName 专家标签 ExpertTag自定义表名 expert_tag
func (ExpertTag) TableName() string {
	return "expert_tag"
}

// ExpertTagRelation 专家-标签关联表
type ExpertTagRelation struct {
	global.GVA_MODEL
	ExpertId uint `json:"expertId" form:"expertId" gorm:"column:expert_id;index;comment:专家ID;" binding:"required"`
	TagId    uint `json:"tagId" form:"tagId" gorm:"column:tag_id;index;comment:标签ID;" binding:"required"`
}

// TableName 专家标签关联 ExpertTagRelation自定义表名 expert_tag_relation
func (ExpertTagRelation) TableName() string {
	return "expert_tag_relation"
}
