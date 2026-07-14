package request

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

type ExpertTagSearch struct {
	StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`

	TagType  string `json:"tagType" form:"tagType"`
	TagValue string `json:"tagValue" form:"tagValue"`

	request.PageInfo
	Sort  string `json:"sort" form:"sort"`
	Order string `json:"order" form:"order"`
}

type ExpertTagRelationSearch struct {
	ExpertId *uint `json:"expertId" form:"expertId"`
	TagId    *uint `json:"tagId" form:"tagId"`

	request.PageInfo
}

// ExpertTagRelationSetReq 覆盖式设置某专家标签关联的入参
type ExpertTagRelationSetReq struct {
	ExpertId uint   `json:"expertId" binding:"required"`
	TagIds   []uint `json:"tagIds"`
}
