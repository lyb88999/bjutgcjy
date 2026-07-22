import service from '@/utils/request'

// @Tags ExpertOrgUser
// @Summary 单位审核员为本单位新建个人申报人账号
// @Router /expertOrgUser/createOrgApplicant [post]
export const createOrgApplicant = (data) => {
  return service({
    url: '/expertOrgUser/createOrgApplicant',
    method: 'post',
    data
  })
}

// @Tags ExpertOrgUser
// @Summary 分页获取本单位个人申报人账号列表
// @Router /expertOrgUser/getOrgUserList [get]
export const getOrgUserList = (params) => {
  return service({
    url: '/expertOrgUser/getOrgUserList',
    method: 'get',
    params
  })
}

// @Tags ExpertOrgUser
// @Summary 启用/冻结本单位个人申报人账号
// @Router /expertOrgUser/toggleOrgUserEnable [post]
export const toggleOrgUserEnable = (data) => {
  return service({
    url: '/expertOrgUser/toggleOrgUserEnable',
    method: 'post',
    data
  })
}
