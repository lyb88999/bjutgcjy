package ExpertDatabase

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// ExpertApprovalLog 专家档案审核流转留痕：每次提交/审核/退回都写一条，用于审计与审核时效统计
type ExpertApprovalLog struct {
	global.GVA_MODEL
	ExpertId   uint   `json:"expertId" form:"expertId" gorm:"column:expert_id;index;comment:专家ID;" binding:"required"`
	FromStatus string `json:"fromStatus" form:"fromStatus" gorm:"column:from_status;comment:流转前状态;"`
	ToStatus   string `json:"toStatus" form:"toStatus" gorm:"column:to_status;comment:流转后状态;"`
	OperatorId uint   `json:"operatorId" form:"operatorId" gorm:"column:operator_id;comment:操作人用户ID;"`
	Opinion    string `json:"opinion" form:"opinion" gorm:"column:opinion;comment:审核意见;"`
}

// TableName 专家审核日志 ExpertApprovalLog自定义表名 expert_approval_log
func (ExpertApprovalLog) TableName() string {
	return "expert_approval_log"
}
