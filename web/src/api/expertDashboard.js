import service from '@/utils/request'

// @Tags ExpertDashboard
// @Summary 专家库统计概览
// @Router /expertDatabase/dashboardStats [get]
export const getDashboardStats = () => {
  return service({
    url: '/expertDatabase/dashboardStats',
    method: 'get'
  })
}
