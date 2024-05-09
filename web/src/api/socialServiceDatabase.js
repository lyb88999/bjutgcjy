import service from '@/utils/request'

// @Tags SocialServiceDatabase
// @Summary 创建社会服务库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SocialServiceDatabase true "创建社会服务库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /SSD/createSocialServiceDatabase [post]
export const createSocialServiceDatabase = (data) => {
  return service({
    url: '/SSD/createSocialServiceDatabase',
    method: 'post',
    data
  })
}

// @Tags SocialServiceDatabase
// @Summary 删除社会服务库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SocialServiceDatabase true "删除社会服务库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /SSD/deleteSocialServiceDatabase [delete]
export const deleteSocialServiceDatabase = (params) => {
  return service({
    url: '/SSD/deleteSocialServiceDatabase',
    method: 'delete',
    params
  })
}

// @Tags SocialServiceDatabase
// @Summary 批量删除社会服务库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除社会服务库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /SSD/deleteSocialServiceDatabase [delete]
export const deleteSocialServiceDatabaseByIds = (params) => {
  return service({
    url: '/SSD/deleteSocialServiceDatabaseByIds',
    method: 'delete',
    params
  })
}

// @Tags SocialServiceDatabase
// @Summary 更新社会服务库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SocialServiceDatabase true "更新社会服务库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /SSD/updateSocialServiceDatabase [put]
export const updateSocialServiceDatabase = (data) => {
  return service({
    url: '/SSD/updateSocialServiceDatabase',
    method: 'put',
    data
  })
}

// @Tags SocialServiceDatabase
// @Summary 用id查询社会服务库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.SocialServiceDatabase true "用id查询社会服务库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /SSD/findSocialServiceDatabase [get]
export const findSocialServiceDatabase = (params) => {
  return service({
    url: '/SSD/findSocialServiceDatabase',
    method: 'get',
    params
  })
}

// @Tags SocialServiceDatabase
// @Summary 分页获取社会服务库列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取社会服务库列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /SSD/getSocialServiceDatabaseList [get]
export const getSocialServiceDatabaseList = (params) => {
  return service({
    url: '/SSD/getSocialServiceDatabaseList',
    method: 'get',
    params
  })
}
