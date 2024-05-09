import service from '@/utils/request'

// @Tags BasicInfomationDatabase
// @Summary 创建基础信息库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.BasicInfomationDatabase true "创建基础信息库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /BID/createBasicInfomationDatabase [post]
export const createBasicInfomationDatabase = (data) => {
  return service({
    url: '/BID/createBasicInfomationDatabase',
    method: 'post',
    data
  })
}

// @Tags BasicInfomationDatabase
// @Summary 删除基础信息库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.BasicInfomationDatabase true "删除基础信息库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /BID/deleteBasicInfomationDatabase [delete]
export const deleteBasicInfomationDatabase = (params) => {
  return service({
    url: '/BID/deleteBasicInfomationDatabase',
    method: 'delete',
    params
  })
}

// @Tags BasicInfomationDatabase
// @Summary 批量删除基础信息库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除基础信息库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /BID/deleteBasicInfomationDatabase [delete]
export const deleteBasicInfomationDatabaseByIds = (params) => {
  return service({
    url: '/BID/deleteBasicInfomationDatabaseByIds',
    method: 'delete',
    params
  })
}

// @Tags BasicInfomationDatabase
// @Summary 更新基础信息库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.BasicInfomationDatabase true "更新基础信息库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /BID/updateBasicInfomationDatabase [put]
export const updateBasicInfomationDatabase = (data) => {
  return service({
    url: '/BID/updateBasicInfomationDatabase',
    method: 'put',
    data
  })
}

// @Tags BasicInfomationDatabase
// @Summary 用id查询基础信息库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.BasicInfomationDatabase true "用id查询基础信息库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /BID/findBasicInfomationDatabase [get]
export const findBasicInfomationDatabase = (params) => {
  return service({
    url: '/BID/findBasicInfomationDatabase',
    method: 'get',
    params
  })
}

// @Tags BasicInfomationDatabase
// @Summary 分页获取基础信息库列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取基础信息库列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /BID/getBasicInfomationDatabaseList [get]
export const getBasicInfomationDatabaseList = (params) => {
  return service({
    url: '/BID/getBasicInfomationDatabaseList',
    method: 'get',
    params
  })
}
