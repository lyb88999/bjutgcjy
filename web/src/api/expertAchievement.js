import service from '@/utils/request'

// @Tags ExpertAchievement
// @Summary 创建专家成果
// @Router /expertAchievement/createExpertAchievement [post]
export const createExpertAchievement = (data) => {
  return service({
    url: '/expertAchievement/createExpertAchievement',
    method: 'post',
    data
  })
}

// @Tags ExpertAchievement
// @Summary 删除专家成果
// @Router /expertAchievement/deleteExpertAchievement [delete]
export const deleteExpertAchievement = (params) => {
  return service({
    url: '/expertAchievement/deleteExpertAchievement',
    method: 'delete',
    params
  })
}

// @Tags ExpertAchievement
// @Summary 批量删除专家成果
// @Router /expertAchievement/deleteExpertAchievementByIds [delete]
export const deleteExpertAchievementByIds = (params) => {
  return service({
    url: '/expertAchievement/deleteExpertAchievementByIds',
    method: 'delete',
    params
  })
}

// @Tags ExpertAchievement
// @Summary 更新专家成果
// @Router /expertAchievement/updateExpertAchievement [put]
export const updateExpertAchievement = (data) => {
  return service({
    url: '/expertAchievement/updateExpertAchievement',
    method: 'put',
    data
  })
}

// @Tags ExpertAchievement
// @Summary 用id查询专家成果
// @Router /expertAchievement/findExpertAchievement [get]
export const findExpertAchievement = (params) => {
  return service({
    url: '/expertAchievement/findExpertAchievement',
    method: 'get',
    params
  })
}

// @Tags ExpertAchievement
// @Summary 分页获取专家成果列表
// @Router /expertAchievement/getExpertAchievementList [get]
export const getExpertAchievementList = (params) => {
  return service({
    url: '/expertAchievement/getExpertAchievementList',
    method: 'get',
    params
  })
}
