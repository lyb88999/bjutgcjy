import service from '@/utils/request'

// @Tags ExpertAdoptionRecord
// @Summary 创建专家决策影响记录
// @Router /expertAdoptionRecord/createExpertAdoptionRecord [post]
export const createExpertAdoptionRecord = (data) => {
  return service({
    url: '/expertAdoptionRecord/createExpertAdoptionRecord',
    method: 'post',
    data
  })
}

// @Tags ExpertAdoptionRecord
// @Summary 删除专家决策影响记录
// @Router /expertAdoptionRecord/deleteExpertAdoptionRecord [delete]
export const deleteExpertAdoptionRecord = (params) => {
  return service({
    url: '/expertAdoptionRecord/deleteExpertAdoptionRecord',
    method: 'delete',
    params
  })
}

// @Tags ExpertAdoptionRecord
// @Summary 批量删除专家决策影响记录
// @Router /expertAdoptionRecord/deleteExpertAdoptionRecordByIds [delete]
export const deleteExpertAdoptionRecordByIds = (params) => {
  return service({
    url: '/expertAdoptionRecord/deleteExpertAdoptionRecordByIds',
    method: 'delete',
    params
  })
}

// @Tags ExpertAdoptionRecord
// @Summary 更新专家决策影响记录
// @Router /expertAdoptionRecord/updateExpertAdoptionRecord [put]
export const updateExpertAdoptionRecord = (data) => {
  return service({
    url: '/expertAdoptionRecord/updateExpertAdoptionRecord',
    method: 'put',
    data
  })
}

// @Tags ExpertAdoptionRecord
// @Summary 用id查询专家决策影响记录
// @Router /expertAdoptionRecord/findExpertAdoptionRecord [get]
export const findExpertAdoptionRecord = (params) => {
  return service({
    url: '/expertAdoptionRecord/findExpertAdoptionRecord',
    method: 'get',
    params
  })
}

// @Tags ExpertAdoptionRecord
// @Summary 分页获取专家决策影响记录列表
// @Router /expertAdoptionRecord/getExpertAdoptionRecordList [get]
export const getExpertAdoptionRecordList = (params) => {
  return service({
    url: '/expertAdoptionRecord/getExpertAdoptionRecordList',
    method: 'get',
    params
  })
}
