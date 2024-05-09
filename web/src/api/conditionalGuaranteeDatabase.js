import service from '@/utils/request'

// @Tags ConditionalGuaranteeDatabase
// @Summary 创建条件保障库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.ConditionalGuaranteeDatabase true "创建条件保障库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /CGD/createConditionalGuaranteeDatabase [post]
export const createConditionalGuaranteeDatabase = (data) => {
  return service({
    url: '/CGD/createConditionalGuaranteeDatabase',
    method: 'post',
    data
  })
}

// @Tags ConditionalGuaranteeDatabase
// @Summary 删除条件保障库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.ConditionalGuaranteeDatabase true "删除条件保障库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /CGD/deleteConditionalGuaranteeDatabase [delete]
export const deleteConditionalGuaranteeDatabase = (params) => {
  return service({
    url: '/CGD/deleteConditionalGuaranteeDatabase',
    method: 'delete',
    params
  })
}

// @Tags ConditionalGuaranteeDatabase
// @Summary 批量删除条件保障库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除条件保障库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /CGD/deleteConditionalGuaranteeDatabase [delete]
export const deleteConditionalGuaranteeDatabaseByIds = (params) => {
  return service({
    url: '/CGD/deleteConditionalGuaranteeDatabaseByIds',
    method: 'delete',
    params
  })
}

// @Tags ConditionalGuaranteeDatabase
// @Summary 更新条件保障库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.ConditionalGuaranteeDatabase true "更新条件保障库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /CGD/updateConditionalGuaranteeDatabase [put]
export const updateConditionalGuaranteeDatabase = (data) => {
  return service({
    url: '/CGD/updateConditionalGuaranteeDatabase',
    method: 'put',
    data
  })
}

// @Tags ConditionalGuaranteeDatabase
// @Summary 用id查询条件保障库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.ConditionalGuaranteeDatabase true "用id查询条件保障库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /CGD/findConditionalGuaranteeDatabase [get]
export const findConditionalGuaranteeDatabase = (params) => {
  return service({
    url: '/CGD/findConditionalGuaranteeDatabase',
    method: 'get',
    params
  })
}

// @Tags ConditionalGuaranteeDatabase
// @Summary 分页获取条件保障库列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取条件保障库列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /CGD/getConditionalGuaranteeDatabaseList [get]
export const getConditionalGuaranteeDatabaseList = (params) => {
  return service({
    url: '/CGD/getConditionalGuaranteeDatabaseList',
    method: 'get',
    params
  })
}
