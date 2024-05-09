import service from '@/utils/request'

// @Tags SubjectConstructionDatabase
// @Summary 创建学科建设库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SubjectConstructionDatabase true "创建学科建设库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /SCD/createSubjectConstructionDatabase [post]
export const createSubjectConstructionDatabase = (data) => {
  return service({
    url: '/SCD/createSubjectConstructionDatabase',
    method: 'post',
    data
  })
}

// @Tags SubjectConstructionDatabase
// @Summary 删除学科建设库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SubjectConstructionDatabase true "删除学科建设库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /SCD/deleteSubjectConstructionDatabase [delete]
export const deleteSubjectConstructionDatabase = (params) => {
  return service({
    url: '/SCD/deleteSubjectConstructionDatabase',
    method: 'delete',
    params
  })
}

// @Tags SubjectConstructionDatabase
// @Summary 批量删除学科建设库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除学科建设库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /SCD/deleteSubjectConstructionDatabase [delete]
export const deleteSubjectConstructionDatabaseByIds = (params) => {
  return service({
    url: '/SCD/deleteSubjectConstructionDatabaseByIds',
    method: 'delete',
    params
  })
}

// @Tags SubjectConstructionDatabase
// @Summary 更新学科建设库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SubjectConstructionDatabase true "更新学科建设库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /SCD/updateSubjectConstructionDatabase [put]
export const updateSubjectConstructionDatabase = (data) => {
  return service({
    url: '/SCD/updateSubjectConstructionDatabase',
    method: 'put',
    data
  })
}

// @Tags SubjectConstructionDatabase
// @Summary 用id查询学科建设库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.SubjectConstructionDatabase true "用id查询学科建设库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /SCD/findSubjectConstructionDatabase [get]
export const findSubjectConstructionDatabase = (params) => {
  return service({
    url: '/SCD/findSubjectConstructionDatabase',
    method: 'get',
    params
  })
}

// @Tags SubjectConstructionDatabase
// @Summary 分页获取学科建设库列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取学科建设库列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /SCD/getSubjectConstructionDatabaseList [get]
export const getSubjectConstructionDatabaseList = (params) => {
  return service({
    url: '/SCD/getSubjectConstructionDatabaseList',
    method: 'get',
    params
  })
}
