package response

// ExpertBatchImportRowError 批量导入校验失败的单行记录
type ExpertBatchImportRowError struct {
	Sheet   string `json:"sheet"`   // 出错的 sheet 名称
	Row     int    `json:"row"`     // Excel 里的实际行号（含表头，从1开始）
	Message string `json:"message"` // 出错原因
}

// ExpertBatchImportResult 批量导入结果
type ExpertBatchImportResult struct {
	Success              bool                        `json:"success"`
	TotalProfileRows     int                         `json:"totalProfileRows"`     // 背景信息表有效数据行数
	CreatedProfiles      int                         `json:"createdProfiles"`      // 新建的专家档案数
	ReusedProfiles       int                         `json:"reusedProfiles"`       // 姓名+单位已存在、复用原档案的数量
	TotalAchievementRows int                         `json:"totalAchievementRows"` // 研究成果表有效数据行数
	CreatedAchievements  int                         `json:"createdAchievements"`  // 新建的成果记录数
	SkippedAchievements  int                         `json:"skippedAchievements"`  // 该专家名下已有同标题成果、跳过的数量
	Errors               []ExpertBatchImportRowError `json:"errors"`               // 校验失败的行，成功时为空
	UnmatchedUnits       []string                    `json:"unmatchedUnits"`       // 填写的单位名在"单位管理"里找不到对应记录，不阻断导入，但需要人工核实/补建
}
