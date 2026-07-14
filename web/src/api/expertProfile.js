import service from '@/utils/request'

// @Tags ExpertProfile
// @Summary 创建专家主档
// @Router /expertProfile/createExpertProfile [post]
export const createExpertProfile = (data) => {
  return service({
    url: '/expertProfile/createExpertProfile',
    method: 'post',
    data
  })
}

// @Tags ExpertProfile
// @Summary 删除专家主档
// @Router /expertProfile/deleteExpertProfile [delete]
export const deleteExpertProfile = (params) => {
  return service({
    url: '/expertProfile/deleteExpertProfile',
    method: 'delete',
    params
  })
}

// @Tags ExpertProfile
// @Summary 批量删除专家主档
// @Router /expertProfile/deleteExpertProfileByIds [delete]
export const deleteExpertProfileByIds = (params) => {
  return service({
    url: '/expertProfile/deleteExpertProfileByIds',
    method: 'delete',
    params
  })
}

// @Tags ExpertProfile
// @Summary 更新专家主档
// @Router /expertProfile/updateExpertProfile [put]
export const updateExpertProfile = (data) => {
  return service({
    url: '/expertProfile/updateExpertProfile',
    method: 'put',
    data
  })
}

// @Tags ExpertProfile
// @Summary 用id查询专家主档
// @Router /expertProfile/findExpertProfile [get]
export const findExpertProfile = (params) => {
  return service({
    url: '/expertProfile/findExpertProfile',
    method: 'get',
    params
  })
}

// @Tags ExpertProfile
// @Summary 分页获取专家主档列表
// @Router /expertProfile/getExpertProfileList [get]
export const getExpertProfileList = (params) => {
  return service({
    url: '/expertProfile/getExpertProfileList',
    method: 'get',
    params
  })
}
