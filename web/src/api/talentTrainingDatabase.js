import service from '@/utils/request'

// @Tags TalentTrainingDatabase
// @Summary 创建人才培养库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.TalentTrainingDatabase true "创建人才培养库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /TTD/createTalentTrainingDatabase [post]
export const createTalentTrainingDatabase = (data) => {
  return service({
    url: '/TTD/createTalentTrainingDatabase',
    method: 'post',
    data
  })
}

// @Tags TalentTrainingDatabase
// @Summary 删除人才培养库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.TalentTrainingDatabase true "删除人才培养库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /TTD/deleteTalentTrainingDatabase [delete]
export const deleteTalentTrainingDatabase = (params) => {
  return service({
    url: '/TTD/deleteTalentTrainingDatabase',
    method: 'delete',
    params
  })
}

// @Tags TalentTrainingDatabase
// @Summary 批量删除人才培养库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除人才培养库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /TTD/deleteTalentTrainingDatabase [delete]
export const deleteTalentTrainingDatabaseByIds = (params) => {
  return service({
    url: '/TTD/deleteTalentTrainingDatabaseByIds',
    method: 'delete',
    params
  })
}

// @Tags TalentTrainingDatabase
// @Summary 更新人才培养库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.TalentTrainingDatabase true "更新人才培养库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /TTD/updateTalentTrainingDatabase [put]
export const updateTalentTrainingDatabase = (data) => {
  return service({
    url: '/TTD/updateTalentTrainingDatabase',
    method: 'put',
    data
  })
}

// @Tags TalentTrainingDatabase
// @Summary 用id查询人才培养库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.TalentTrainingDatabase true "用id查询人才培养库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /TTD/findTalentTrainingDatabase [get]
export const findTalentTrainingDatabase = (params) => {
  return service({
    url: '/TTD/findTalentTrainingDatabase',
    method: 'get',
    params
  })
}

// @Tags TalentTrainingDatabase
// @Summary 分页获取人才培养库列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取人才培养库列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /TTD/getTalentTrainingDatabaseList [get]
export const getTalentTrainingDatabaseList = (params) => {
  return service({
    url: '/TTD/getTalentTrainingDatabaseList',
    method: 'get',
    params
  })
}
