package system

import (
	"errors"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
)

type OrganizationService struct{}

//@function: CreateSysOrganization
//@description: 创建单位
//@param: org system.SysOrganization
//@return: err error

func (organizationService *OrganizationService) CreateSysOrganization(org system.SysOrganization) (err error) {
	return global.GVA_DB.Create(&org).Error
}

//@function: DeleteSysOrganization
//@description: 删除单位，存在子单位或已被用户关联时禁止删除
//@param: id uint
//@return: err error

func (organizationService *OrganizationService) DeleteSysOrganization(id uint) (err error) {
	var childCount int64
	if err = global.GVA_DB.Model(&system.SysOrganization{}).Where("parent_id = ?", id).Count(&childCount).Error; err != nil {
		return err
	}
	if childCount > 0 {
		return errors.New("存在下级单位，请先删除下级单位")
	}
	var userCount int64
	if err = global.GVA_DB.Model(&system.SysUser{}).Where("org_id = ?", id).Count(&userCount).Error; err != nil {
		return err
	}
	if userCount > 0 {
		return errors.New("该单位下存在关联用户，请先解除关联")
	}
	return global.GVA_DB.Delete(&system.SysOrganization{}, "id = ?", id).Error
}

//@function: UpdateSysOrganization
//@description: 更新单位
//@param: org *system.SysOrganization
//@return: err error

func (organizationService *OrganizationService) UpdateSysOrganization(org *system.SysOrganization) (err error) {
	return global.GVA_DB.Model(&system.SysOrganization{}).Where("id = ?", org.ID).Updates(map[string]interface{}{
		"name":        org.Name,
		"level":       org.Level,
		"parent_id":   org.ParentId,
		"region_code": org.RegionCode,
		"status":      org.Status,
	}).Error
}

//@function: GetSysOrganization
//@description: 根据id获取单位
//@param: id uint
//@return: org system.SysOrganization, err error

func (organizationService *OrganizationService) GetSysOrganization(id uint) (org system.SysOrganization, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&org).Error
	return
}

//@function: GetSysOrganizationInfoList
//@description: 分页获取单位列表，支持按名称关键字搜索
//@param: info request.PageInfo
//@return: list []system.SysOrganization, total int64, err error

func (organizationService *OrganizationService) GetSysOrganizationInfoList(info request.PageInfo) (list []system.SysOrganization, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&system.SysOrganization{})
	if info.Keyword != "" {
		db = db.Where("name LIKE ?", "%"+info.Keyword+"%")
	}
	if err = db.Count(&total).Error; err != nil {
		return
	}
	err = db.Order("id").Limit(limit).Offset(offset).Find(&list).Error
	return
}

//@function: GetSysOrganizationTree
//@description: 获取全部单位并按 parent_id 组装成树，供审核路由/表单级联选择使用
//@return: tree []system.SysOrganization, err error

func (organizationService *OrganizationService) GetSysOrganizationTree() (tree []system.SysOrganization, err error) {
	var all []system.SysOrganization
	if err = global.GVA_DB.Order("id").Find(&all).Error; err != nil {
		return nil, err
	}
	byParent := make(map[uint][]system.SysOrganization)
	for _, o := range all {
		var parentKey uint
		if o.ParentId != nil {
			parentKey = *o.ParentId
		}
		byParent[parentKey] = append(byParent[parentKey], o)
	}
	var attachChildren func(nodes []system.SysOrganization) []system.SysOrganization
	attachChildren = func(nodes []system.SysOrganization) []system.SysOrganization {
		for i := range nodes {
			nodes[i].Children = attachChildren(byParent[nodes[i].ID])
		}
		return nodes
	}
	tree = attachChildren(byParent[0])
	return tree, nil
}
