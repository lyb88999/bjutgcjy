package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/router/EngineeringEducationDatabase"
	"github.com/flipped-aurora/gin-vue-admin/server/router/example"
	"github.com/flipped-aurora/gin-vue-admin/server/router/system"
	"github.com/flipped-aurora/gin-vue-admin/server/router/util"
)

type RouterGroup struct {
	System                       system.RouterGroup
	Example                      example.RouterGroup
	EngineeringEducationDatabase EngineeringEducationDatabase.RouterGroup
	Util                         util.RouterGroup
}

var RouterGroupApp = new(RouterGroup)
