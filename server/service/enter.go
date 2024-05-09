package service

import (
	"github.com/flipped-aurora/gin-vue-admin/server/service/EngineeringEducationDatabase"
	"github.com/flipped-aurora/gin-vue-admin/server/service/example"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/flipped-aurora/gin-vue-admin/server/service/util"
)

type ServiceGroup struct {
	SystemServiceGroup                       system.ServiceGroup
	ExampleServiceGroup                      example.ServiceGroup
	EngineeringEducationDatabaseServiceGroup EngineeringEducationDatabase.ServiceGroup
	UtilServiceGroup                         util.ServiceGroup
}

var ServiceGroupApp = new(ServiceGroup)
