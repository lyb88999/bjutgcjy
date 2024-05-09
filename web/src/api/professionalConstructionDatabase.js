import service from '@/utils/request'

// @Tags ProfessionalConstructionDatabase
// @Summary 创建专业建设库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.ProfessionalConstructionDatabase true "创建专业建设库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /PCD/createProfessionalConstructionDatabase [post]
export const createProfessionalConstructionDatabase = (data) => {
  return service({
    url: '/PCD/createProfessionalConstructionDatabase',
    method: 'post',
    data
  })
}

// @Tags ProfessionalConstructionDatabase
// @Summary 删除专业建设库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.ProfessionalConstructionDatabase true "删除专业建设库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /PCD/deleteProfessionalConstructionDatabase [delete]
export const deleteProfessionalConstructionDatabase = (params) => {
  return service({
    url: '/PCD/deleteProfessionalConstructionDatabase',
    method: 'delete',
    params
  })
}

// @Tags ProfessionalConstructionDatabase
// @Summary 批量删除专业建设库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除专业建设库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /PCD/deleteProfessionalConstructionDatabase [delete]
export const deleteProfessionalConstructionDatabaseByIds = (params) => {
  return service({
    url: '/PCD/deleteProfessionalConstructionDatabaseByIds',
    method: 'delete',
    params
  })
}

// @Tags ProfessionalConstructionDatabase
// @Summary 更新专业建设库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.ProfessionalConstructionDatabase true "更新专业建设库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /PCD/updateProfessionalConstructionDatabase [put]
export const updateProfessionalConstructionDatabase = (data) => {
  return service({
    url: '/PCD/updateProfessionalConstructionDatabase',
    method: 'put',
    data
  })
}

// @Tags ProfessionalConstructionDatabase
// @Summary 用id查询专业建设库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.ProfessionalConstructionDatabase true "用id查询专业建设库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /PCD/findProfessionalConstructionDatabase [get]
export const findProfessionalConstructionDatabase = (params) => {
  return service({
    url: '/PCD/findProfessionalConstructionDatabase',
    method: 'get',
    params
  })
}

// @Tags ProfessionalConstructionDatabase
// @Summary 分页获取专业建设库列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取专业建设库列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /PCD/getProfessionalConstructionDatabaseList [get]
export const getProfessionalConstructionDatabaseList = (params) => {
  return service({
    url: '/PCD/getProfessionalConstructionDatabaseList',
    method: 'get',
    params
  })
}
