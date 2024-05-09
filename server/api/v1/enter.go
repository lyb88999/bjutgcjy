package v1

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/EngineeringEducationDatabase"
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/example"
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/system"
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/util"
)

type ApiGroup struct {
	SystemApiGroup                       system.ApiGroup
	ExampleApiGroup                      example.ApiGroup
	EngineeringEducationDatabaseApiGroup EngineeringEducationDatabase.ApiGroup
	UtilApiGroup                         util.ApiGroup
}

var ApiGroupApp = new(ApiGroup)
