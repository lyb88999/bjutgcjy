import service from '@/utils/request'

// @Tags ExpertAcademicPosition
// @Summary 创建专家学术兼职
// @Router /expertAcademicPosition/createExpertAcademicPosition [post]
export const createExpertAcademicPosition = (data) => {
  return service({
    url: '/expertAcademicPosition/createExpertAcademicPosition',
    method: 'post',
    data
  })
}

// @Tags ExpertAcademicPosition
// @Summary 删除专家学术兼职
// @Router /expertAcademicPosition/deleteExpertAcademicPosition [delete]
export const deleteExpertAcademicPosition = (params) => {
  return service({
    url: '/expertAcademicPosition/deleteExpertAcademicPosition',
    method: 'delete',
    params
  })
}

// @Tags ExpertAcademicPosition
// @Summary 批量删除专家学术兼职
// @Router /expertAcademicPosition/deleteExpertAcademicPositionByIds [delete]
export const deleteExpertAcademicPositionByIds = (params) => {
  return service({
    url: '/expertAcademicPosition/deleteExpertAcademicPositionByIds',
    method: 'delete',
    params
  })
}

// @Tags ExpertAcademicPosition
// @Summary 更新专家学术兼职
// @Router /expertAcademicPosition/updateExpertAcademicPosition [put]
export const updateExpertAcademicPosition = (data) => {
  return service({
    url: '/expertAcademicPosition/updateExpertAcademicPosition',
    method: 'put',
    data
  })
}

// @Tags ExpertAcademicPosition
// @Summary 用id查询专家学术兼职
// @Router /expertAcademicPosition/findExpertAcademicPosition [get]
export const findExpertAcademicPosition = (params) => {
  return service({
    url: '/expertAcademicPosition/findExpertAcademicPosition',
    method: 'get',
    params
  })
}

// @Tags ExpertAcademicPosition
// @Summary 分页获取专家学术兼职列表
// @Router /expertAcademicPosition/getExpertAcademicPositionList [get]
export const getExpertAcademicPositionList = (params) => {
  return service({
    url: '/expertAcademicPosition/getExpertAcademicPositionList',
    method: 'get',
    params
  })
}
