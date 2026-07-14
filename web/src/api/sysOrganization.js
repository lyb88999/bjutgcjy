import service from '@/utils/request'

// @Tags SysOrganization
// @Summary 创建单位
// @Router /sysOrganization/createSysOrganization [post]
export const createSysOrganization = (data) => {
  return service({
    url: '/sysOrganization/createSysOrganization',
    method: 'post',
    data
  })
}

// @Tags SysOrganization
// @Summary 删除单位
// @Router /sysOrganization/deleteSysOrganization [delete]
export const deleteSysOrganization = (params) => {
  return service({
    url: '/sysOrganization/deleteSysOrganization',
    method: 'delete',
    params
  })
}

// @Tags SysOrganization
// @Summary 更新单位
// @Router /sysOrganization/updateSysOrganization [put]
export const updateSysOrganization = (data) => {
  return service({
    url: '/sysOrganization/updateSysOrganization',
    method: 'put',
    data
  })
}

// @Tags SysOrganization
// @Summary 用id查询单位
// @Router /sysOrganization/findSysOrganization [get]
export const findSysOrganization = (params) => {
  return service({
    url: '/sysOrganization/findSysOrganization',
    method: 'get',
    params
  })
}

// @Tags SysOrganization
// @Summary 分页获取单位列表
// @Router /sysOrganization/getSysOrganizationList [get]
export const getSysOrganizationList = (params) => {
  return service({
    url: '/sysOrganization/getSysOrganizationList',
    method: 'get',
    params
  })
}

// @Tags SysOrganization
// @Summary 获取单位树
// @Router /sysOrganization/getSysOrganizationTree [get]
export const getSysOrganizationTree = (params) => {
  return service({
    url: '/sysOrganization/getSysOrganizationTree',
    method: 'get',
    params
  })
}
