package request

import "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"

// CreateOrgApplicantRequest 单位审核员为本单位新建个人申报人账号——
// 特意不接受 orgId/authorityId 参数，这两个值一律由后端按操作人自己的所属单位/固定角色
// 推导，客户端传什么都不会被采纳，避免审核员越权建别的单位或别的角色的账号
type CreateOrgApplicantRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	NickName string `json:"nickName"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
}

// OrgUserSearch 查询本单位账号列表
type OrgUserSearch struct {
	request.PageInfo
}

// ToggleOrgUserEnableRequest 启用/冻结本单位账号
type ToggleOrgUserEnableRequest struct {
	ID     uint `json:"ID" binding:"required"`
	Enable int  `json:"enable" binding:"required"` // 1正常 2冻结
}
