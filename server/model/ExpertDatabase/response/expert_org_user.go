package response

import "time"

// OrgUserItem 本单位账号列表的单条展示
type OrgUserItem struct {
	ID          uint      `json:"ID"`
	CreatedAt   time.Time `json:"CreatedAt"`
	Username    string    `json:"username"`
	NickName    string    `json:"nickName"`
	Phone       string    `json:"phone"`
	Email       string    `json:"email"`
	Enable      int       `json:"enable"`
	AuthorityId uint      `json:"authorityId"`
}
