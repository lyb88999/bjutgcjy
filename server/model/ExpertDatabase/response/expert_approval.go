package response

// ExpertBatchApprovalFailure 批量审核里失败的单条记录及原因（比如已经被别人处理过、状态不对等）
type ExpertBatchApprovalFailure struct {
	ExpertId uint   `json:"expertId"`
	Message  string `json:"message"`
}

// ExpertBatchApprovalResult 批量审核结果：成功/失败分开统计，失败的逐条给原因，
// 不因为其中几条失败就让整批操作全部回滚——批量审核本来就是要把能过的先过掉，减少手动点击次数
type ExpertBatchApprovalResult struct {
	SuccessCount int                          `json:"successCount"`
	FailCount    int                          `json:"failCount"`
	Failures     []ExpertBatchApprovalFailure `json:"failures"`
}
