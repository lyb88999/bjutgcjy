import service from '@/utils/request'

// @Tags PolicyDatabase
// @Summary 创建政策数据库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.PolicyDatabase true "创建政策数据库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /PD/createPolicyDatabase [post]
export const createPolicyDatabase = (data) => {
  return service({
    url: '/PD/createPolicyDatabase',
    method: 'post',
    data
  })
}

// @Tags PolicyDatabase
// @Summary 删除政策数据库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.PolicyDatabase true "删除政策数据库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /PD/deletePolicyDatabase [delete]
export const deletePolicyDatabase = (params) => {
  return service({
    url: '/PD/deletePolicyDatabase',
    method: 'delete',
    params
  })
}

// @Tags PolicyDatabase
// @Summary 批量删除政策数据库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除政策数据库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /PD/deletePolicyDatabase [delete]
export const deletePolicyDatabaseByIds = (params) => {
  return service({
    url: '/PD/deletePolicyDatabaseByIds',
    method: 'delete',
    params
  })
}

// @Tags PolicyDatabase
// @Summary 更新政策数据库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.PolicyDatabase true "更新政策数据库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /PD/updatePolicyDatabase [put]
export const updatePolicyDatabase = (data) => {
  return service({
    url: '/PD/updatePolicyDatabase',
    method: 'put',
    data
  })
}

// @Tags PolicyDatabase
// @Summary 用id查询政策数据库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.PolicyDatabase true "用id查询政策数据库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /PD/findPolicyDatabase [get]
export const findPolicyDatabase = (params) => {
  return service({
    url: '/PD/findPolicyDatabase',
    method: 'get',
    params
  })
}

// @Tags PolicyDatabase
// @Summary 分页获取政策数据库列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取政策数据库列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /PD/getPolicyDatabaseList [get]
export const getPolicyDatabaseList = (params) => {
  return service({
    url: '/PD/getPolicyDatabaseList',
    method: 'get',
    params
  })
}
