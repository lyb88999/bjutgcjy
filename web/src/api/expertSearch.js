import service from '@/utils/request'

// @Tags ExpertSearch
// @Summary 专家综合推荐排序检索
// @Router /expertDatabase/search [get]
export const searchExpert = (params) => {
  return service({ url: '/expertDatabase/search', method: 'get', params })
}

// @Tags ExpertSearch
// @Summary 手动重算专家得分
// @Router /expertDatabase/recomputeScore [post]
export const recomputeExpertScore = (data) => {
  return service({ url: '/expertDatabase/recomputeScore', method: 'post', data })
}

// @Tags ExpertSearch
// @Summary 导出当前检索结果
// @Router /expertDatabase/exportSearchResults [get]
export const exportSearchResults = (params) => {
  return service({ url: '/expertDatabase/exportSearchResults', method: 'get', params, responseType: 'blob' })
}

// @Tags ExpertSearch
// @Summary 分项得分是否有数据（全库是否至少一人非零），用来决定要不要显示对应的列/卡片
// @Router /expertDatabase/scoreColumnAvailability [get]
export const getScoreColumnAvailability = () => {
  return service({ url: '/expertDatabase/scoreColumnAvailability', method: 'get' })
}
