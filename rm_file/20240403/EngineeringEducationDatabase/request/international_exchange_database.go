package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
	
)

type InternationalExchangeDatabaseSearch struct{
    
        StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
        EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
    
                      UniversityName  string `json:"universityName" form:"universityName" `
                      InternationalConferencesHosted  *int `json:"internationalConferencesHosted" form:"internationalConferencesHosted" `
                      ForeignFacultyCount  *int `json:"foreignFacultyCount" form:"foreignFacultyCount" `
                      InternationalStudentsCount  *int `json:"internationalStudentsCount" form:"internationalStudentsCount" `
                      StudentsAbroadExchangeCount  *int `json:"studentsAbroadExchangeCount" form:"studentsAbroadExchangeCount" `
    request.PageInfo
    Sort  string `json:"sort" form:"sort"`
    Order string `json:"order" form:"order"`
}
