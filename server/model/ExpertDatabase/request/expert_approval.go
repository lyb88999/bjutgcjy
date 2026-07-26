package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

// ExpertApprovalAction 审核流转通用入参：提交/单位通过不需要意见，退回类操作需要填写 opinion
type ExpertApprovalAction struct {
	ExpertId uint   `json:"expertId" binding:"required"`
	Opinion  string `json:"opinion"`
}

// ExpertAdminSetStatus 管理员直接改写状态（跳过流程，但仍记录日志）
type ExpertAdminSetStatus struct {
	ExpertId uint   `json:"expertId" binding:"required"`
	Status   string `json:"status" binding:"required"`
	Opinion  string `json:"opinion"`
}

// ExpertBatchApprovalAction 批量审核通过入参：一次性对多条记录做同一个审核动作
type ExpertBatchApprovalAction struct {
	ExpertIds []uint `json:"expertIds" binding:"required"`
}

// ExpertApprovalLogSearch 审核日志查询
type ExpertApprovalLogSearch struct {
	ExpertId uint `json:"expertId" form:"expertId" binding:"required"`

	request.PageInfo
}
