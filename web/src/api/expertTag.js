import service from '@/utils/request'

// @Tags ExpertTag
// @Summary 创建标签
// @Router /expertTag/createExpertTag [post]
export const createExpertTag = (data) => {
  return service({
    url: '/expertTag/createExpertTag',
    method: 'post',
    data
  })
}

// @Tags ExpertTag
// @Summary 删除标签
// @Router /expertTag/deleteExpertTag [delete]
export const deleteExpertTag = (params) => {
  return service({
    url: '/expertTag/deleteExpertTag',
    method: 'delete',
    params
  })
}

// @Tags ExpertTag
// @Summary 更新标签
// @Router /expertTag/updateExpertTag [put]
export const updateExpertTag = (data) => {
  return service({
    url: '/expertTag/updateExpertTag',
    method: 'put',
    data
  })
}

// @Tags ExpertTag
// @Summary 用id查询标签
// @Router /expertTag/findExpertTag [get]
export const findExpertTag = (params) => {
  return service({
    url: '/expertTag/findExpertTag',
    method: 'get',
    params
  })
}

// @Tags ExpertTag
// @Summary 分页获取标签列表
// @Router /expertTag/getExpertTagList [get]
export const getExpertTagList = (params) => {
  return service({
    url: '/expertTag/getExpertTagList',
    method: 'get',
    params
  })
}

// @Tags ExpertTag
// @Summary 设置专家标签关联
// @Router /expertTag/setExpertTagRelations [post]
export const setExpertTagRelations = (data) => {
  return service({
    url: '/expertTag/setExpertTagRelations',
    method: 'post',
    data
  })
}

// @Tags ExpertTag
// @Summary 获取某专家关联的标签
// @Router /expertTag/getExpertTagsByExpertId [get]
export const getExpertTagsByExpertId = (params) => {
  return service({
    url: '/expertTag/getExpertTagsByExpertId',
    method: 'get',
    params
  })
}
