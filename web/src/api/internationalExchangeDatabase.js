import service from '@/utils/request'

// @Tags InternationalExchangeDatabase
// @Summary 创建国际交流库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.InternationalExchangeDatabase true "创建国际交流库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /IED/createInternationalExchangeDatabase [post]
export const createInternationalExchangeDatabase = (data) => {
  return service({
    url: '/IED/createInternationalExchangeDatabase',
    method: 'post',
    data
  })
}

// @Tags InternationalExchangeDatabase
// @Summary 删除国际交流库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.InternationalExchangeDatabase true "删除国际交流库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /IED/deleteInternationalExchangeDatabase [delete]
export const deleteInternationalExchangeDatabase = (params) => {
  return service({
    url: '/IED/deleteInternationalExchangeDatabase',
    method: 'delete',
    params
  })
}

// @Tags InternationalExchangeDatabase
// @Summary 批量删除国际交流库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除国际交流库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /IED/deleteInternationalExchangeDatabase [delete]
export const deleteInternationalExchangeDatabaseByIds = (params) => {
  return service({
    url: '/IED/deleteInternationalExchangeDatabaseByIds',
    method: 'delete',
    params
  })
}

// @Tags InternationalExchangeDatabase
// @Summary 更新国际交流库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.InternationalExchangeDatabase true "更新国际交流库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /IED/updateInternationalExchangeDatabase [put]
export const updateInternationalExchangeDatabase = (data) => {
  return service({
    url: '/IED/updateInternationalExchangeDatabase',
    method: 'put',
    data
  })
}

// @Tags InternationalExchangeDatabase
// @Summary 用id查询国际交流库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.InternationalExchangeDatabase true "用id查询国际交流库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /IED/findInternationalExchangeDatabase [get]
export const findInternationalExchangeDatabase = (params) => {
  return service({
    url: '/IED/findInternationalExchangeDatabase',
    method: 'get',
    params
  })
}

// @Tags InternationalExchangeDatabase
// @Summary 分页获取国际交流库列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取国际交流库列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /IED/getInternationalExchangeDatabaseList [get]
export const getInternationalExchangeDatabaseList = (params) => {
  return service({
    url: '/IED/getInternationalExchangeDatabaseList',
    method: 'get',
    params
  })
}
