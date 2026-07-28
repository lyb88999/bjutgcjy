import service from '@/utils/request'

// @Tags ExpertApproval
// @Summary 提交专家档案审核
// @Router /expertApproval/submit [post]
export const submitExpertProfile = (data) => {
  return service({ url: '/expertApproval/submit', method: 'post', data })
}

// @Tags ExpertApproval
// @Summary 单位审核通过
// @Router /expertApproval/orgApprove [post]
export const orgApproveExpertProfile = (data) => {
  return service({ url: '/expertApproval/orgApprove', method: 'post', data })
}

// @Tags ExpertApproval
// @Summary 批量单位审核通过
// @Router /expertApproval/batchOrgApprove [post]
export const batchOrgApproveExpertProfile = (data) => {
  return service({ url: '/expertApproval/batchOrgApprove', method: 'post', data })
}

// @Tags ExpertApproval
// @Summary 单位审核退回
// @Router /expertApproval/orgReject [post]
export const orgRejectExpertProfile = (data) => {
  return service({ url: '/expertApproval/orgReject', method: 'post', data })
}

// @Tags ExpertApproval
// @Summary 市级审核通过
// @Router /expertApproval/cityApprove [post]
export const cityApproveExpertProfile = (data) => {
  return service({ url: '/expertApproval/cityApprove', method: 'post', data })
}

// @Tags ExpertApproval
// @Summary 批量市级审核通过
// @Router /expertApproval/batchCityApprove [post]
export const batchCityApproveExpertProfile = (data) => {
  return service({ url: '/expertApproval/batchCityApprove', method: 'post', data })
}

// @Tags ExpertApproval
// @Summary 市级审核退回
// @Router /expertApproval/cityReject [post]
export const cityRejectExpertProfile = (data) => {
  return service({ url: '/expertApproval/cityReject', method: 'post', data })
}

// @Tags ExpertApproval
// @Summary 管理员直接改写审核状态
// @Router /expertApproval/adminSetStatus [post]
export const adminSetExpertProfileStatus = (data) => {
  return service({ url: '/expertApproval/adminSetStatus', method: 'post', data })
}

// @Tags ExpertApproval
// @Summary 我发起的专家档案列表
// @Router /expertApproval/myDrafts [get]
export const getMyDrafts = (params) => {
  return service({ url: '/expertApproval/myDrafts', method: 'get', params })
}

// @Tags ExpertApproval
// @Summary 待本单位审核的专家档案列表
// @Router /expertApproval/pendingOrgReview [get]
export const getPendingOrgReview = (params) => {
  return service({ url: '/expertApproval/pendingOrgReview', method: 'get', params })
}

// @Tags ExpertApproval
// @Summary 待市级审核的专家档案列表
// @Router /expertApproval/pendingCityReview [get]
export const getPendingCityReview = (params) => {
  return service({ url: '/expertApproval/pendingCityReview', method: 'get', params })
}

// @Tags ExpertApproval
// @Summary 专家档案审核历史
// @Router /expertApproval/getApprovalLogList [get]
export const getExpertApprovalLogList = (params) => {
  return service({ url: '/expertApproval/getApprovalLogList', method: 'get', params })
}
