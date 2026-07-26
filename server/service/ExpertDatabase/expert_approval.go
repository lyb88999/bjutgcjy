package ExpertDatabase

import (
	"errors"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase"
	ExpertDatabaseReq "github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase/request"
	ExpertDatabaseResp "github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase/response"
	commonRequest "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"gorm.io/gorm"
)

// 专家档案审核状态机取值
const (
	StatusDraft             = "draft"
	StatusPendingOrgReview  = "pending_org_review"
	StatusOrgRejected       = "org_rejected"
	StatusPendingCityReview = "pending_city_review"
	StatusCityRejected      = "city_rejected"
	StatusPublished         = "published"
)

// DictTypeReviewSwitch 专家库审核开关字典类型：管理员在"系统工具 -> 字典管理"里改这个字典的
// 明细值，就能控制批量导入/手动新增/提交审核这三个入口是否要走三级审核流程，不用改代码、不用重启
const DictTypeReviewSwitch = "expert_review_switch"

// reviewRequired 读取审核开关；字典或明细缺失时默认按"需要审核"处理（安全默认值，不会因为
// 管理员没配置字典就意外整体跳过审核流程），字典存在且明细 value 为 0 时表示关闭审核、直接发布
func reviewRequired() bool {
	var dict system.SysDictionary
	if err := global.GVA_DB.Where("type = ?", DictTypeReviewSwitch).First(&dict).Error; err != nil {
		return true
	}
	var detail system.SysDictionaryDetail
	if err := global.GVA_DB.Where("sys_dictionary_id = ?", dict.ID).Order("sort").First(&detail).Error; err != nil {
		return true
	}
	return detail.Value != 0
}

type ExpertApprovalService struct{}

// transition 校验来源状态是否允许流转、落库新状态、写审核日志，是所有流转方法的共用核心
func (s *ExpertApprovalService) transition(tx *gorm.DB, profile *ExpertDatabase.ExpertProfile, allowedFrom []string, to string, operatorID uint, opinion string) error {
	allowed := false
	for _, from := range allowedFrom {
		if profile.Status == from {
			allowed = true
			break
		}
	}
	if !allowed {
		return errors.New("当前状态不允许该操作：" + profile.Status)
	}
	fromStatus := profile.Status
	if err := tx.Model(&ExpertDatabase.ExpertProfile{}).Where("id = ?", profile.ID).Update("status", to).Error; err != nil {
		return err
	}
	log := ExpertDatabase.ExpertApprovalLog{
		ExpertId:   profile.ID,
		FromStatus: fromStatus,
		ToStatus:   to,
		OperatorId: operatorID,
		Opinion:    opinion,
	}
	return tx.Create(&log).Error
}

// Submit 提交审核：草稿/单位退回/市级退回 -> 待单位审核；审核开关关闭时直接跳到已发布
func (s *ExpertApprovalService) Submit(expertID uint, operatorID uint) error {
	if !reviewRequired() {
		err := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
			var profile ExpertDatabase.ExpertProfile
			if err := tx.Where("id = ?", expertID).First(&profile).Error; err != nil {
				return err
			}
			if err := s.transition(tx, &profile, []string{StatusDraft, StatusOrgRejected, StatusCityRejected}, StatusPublished, operatorID, ""); err != nil {
				return err
			}
			return tx.Model(&ExpertDatabase.ExpertProfile{}).Where("id = ?", expertID).Update("review_bypassed", true).Error
		})
		if err != nil {
			return err
		}
		expertScoreSvc.recomputeIfPublished(expertID)
		return nil
	}
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var profile ExpertDatabase.ExpertProfile
		if err := tx.Where("id = ?", expertID).First(&profile).Error; err != nil {
			return err
		}
		// 待单位审核列表按 org_id 精确匹配审核员所属单位，没有 org_id 的档案如果按常规流程进
		// "待单位审核"，会静默卡死在这一步且无人能审。这类档案没有单位可以归口，跳过单位审核，
		// 直接进入市级审核——市级审核不按单位限定范围，市级审核员和管理员都能处理，流程不会被卡住。
		target := StatusPendingOrgReview
		if profile.OrgId == nil {
			target = StatusPendingCityReview
		}
		return s.transition(tx, &profile, []string{StatusDraft, StatusOrgRejected, StatusCityRejected}, target, operatorID, "")
	})
}

// OrgApprove 单位审核通过：待单位审核 -> 待市级审核；只能审核本单位提交的档案
func (s *ExpertApprovalService) OrgApprove(expertID uint, operatorID uint) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var profile ExpertDatabase.ExpertProfile
		if err := tx.Where("id = ?", expertID).First(&profile).Error; err != nil {
			return err
		}
		if err := s.assertSameOrg(operatorID, profile.OrgId); err != nil {
			return err
		}
		return s.transition(tx, &profile, []string{StatusPendingOrgReview}, StatusPendingCityReview, operatorID, "")
	})
}

// BatchOrgApprove 批量单位审核通过：逐条复用 OrgApprove 的校验和流转逻辑，某一条失败不影响其余记录，
// 失败原因逐条收集返回，方便审核员知道哪几条没过、为什么（比如混进了一条不是本单位提交的记录）
func (s *ExpertApprovalService) BatchOrgApprove(expertIDs []uint, operatorID uint) ExpertDatabaseResp.ExpertBatchApprovalResult {
	result := ExpertDatabaseResp.ExpertBatchApprovalResult{}
	for _, id := range expertIDs {
		if err := s.OrgApprove(id, operatorID); err != nil {
			result.FailCount++
			result.Failures = append(result.Failures, ExpertDatabaseResp.ExpertBatchApprovalFailure{ExpertId: id, Message: err.Error()})
			continue
		}
		result.SuccessCount++
	}
	return result
}

// OrgReject 单位审核退回：待单位审核 -> 单位已退回，需填写意见；只能审核本单位提交的档案
func (s *ExpertApprovalService) OrgReject(expertID uint, operatorID uint, opinion string) error {
	if opinion == "" {
		return errors.New("退回时必须填写审核意见")
	}
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var profile ExpertDatabase.ExpertProfile
		if err := tx.Where("id = ?", expertID).First(&profile).Error; err != nil {
			return err
		}
		if err := s.assertSameOrg(operatorID, profile.OrgId); err != nil {
			return err
		}
		return s.transition(tx, &profile, []string{StatusPendingOrgReview}, StatusOrgRejected, operatorID, opinion)
	})
}

// CityApprove 市级审核通过：待市级审核 -> 已发布，正式收录进入检索排序范围；
// 走到这里说明是真的过完了三级审核，顺带把"免审核发布"标记清掉（哪怕它之前曾经被审核开关绕过一次）
func (s *ExpertApprovalService) CityApprove(expertID uint, operatorID uint) error {
	err := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var profile ExpertDatabase.ExpertProfile
		if err := tx.Where("id = ?", expertID).First(&profile).Error; err != nil {
			return err
		}
		if err := s.transition(tx, &profile, []string{StatusPendingCityReview}, StatusPublished, operatorID, ""); err != nil {
			return err
		}
		return tx.Model(&ExpertDatabase.ExpertProfile{}).Where("id = ?", expertID).Update("review_bypassed", false).Error
	})
	if err != nil {
		return err
	}
	// 首次发布，得分缓存此前从未算过，这里立即算一次，不用等每日兜底任务
	expertScoreSvc.recomputeIfPublished(expertID)
	return nil
}

// BatchCityApprove 批量市级审核通过：逐条复用 CityApprove（含发布 + 得分重算），某一条失败不影响其余记录
func (s *ExpertApprovalService) BatchCityApprove(expertIDs []uint, operatorID uint) ExpertDatabaseResp.ExpertBatchApprovalResult {
	result := ExpertDatabaseResp.ExpertBatchApprovalResult{}
	for _, id := range expertIDs {
		if err := s.CityApprove(id, operatorID); err != nil {
			result.FailCount++
			result.Failures = append(result.Failures, ExpertDatabaseResp.ExpertBatchApprovalFailure{ExpertId: id, Message: err.Error()})
			continue
		}
		result.SuccessCount++
	}
	return result
}

// CityReject 市级审核退回：待市级审核 -> 市级已退回，需填写意见
func (s *ExpertApprovalService) CityReject(expertID uint, operatorID uint, opinion string) error {
	if opinion == "" {
		return errors.New("退回时必须填写审核意见")
	}
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var profile ExpertDatabase.ExpertProfile
		if err := tx.Where("id = ?", expertID).First(&profile).Error; err != nil {
			return err
		}
		return s.transition(tx, &profile, []string{StatusPendingCityReview}, StatusCityRejected, operatorID, opinion)
	})
}

// AdminSetStatus 管理员直接改写状态，跳过流程限制，但仍记录审核日志以便追溯
func (s *ExpertApprovalService) AdminSetStatus(expertID uint, status string, operatorID uint, opinion string) error {
	validStatuses := map[string]bool{
		StatusDraft: true, StatusPendingOrgReview: true, StatusOrgRejected: true,
		StatusPendingCityReview: true, StatusCityRejected: true, StatusPublished: true,
	}
	if !validStatuses[status] {
		return errors.New("非法的审核状态：" + status)
	}
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var profile ExpertDatabase.ExpertProfile
		if err := tx.Where("id = ?", expertID).First(&profile).Error; err != nil {
			return err
		}
		fromStatus := profile.Status
		if err := tx.Model(&ExpertDatabase.ExpertProfile{}).Where("id = ?", expertID).Update("status", status).Error; err != nil {
			return err
		}
		log := ExpertDatabase.ExpertApprovalLog{
			ExpertId:   expertID,
			FromStatus: fromStatus,
			ToStatus:   status,
			OperatorId: operatorID,
			Opinion:    opinion,
		}
		return tx.Create(&log).Error
	})
}

// superAdminAuthorityId 系统自带的超级管理员角色ID，管理员账号本身不挂靠任何单位，
// 审核流程里的"本单位"校验对它天然无法满足——管理员理应有全权限，不能被单位归口卡住
const superAdminAuthorityId = 1

// assertSameOrg 校验操作人所属单位与专家档案的申报单位一致，防止跨单位审核；
// 超级管理员账号跳过这层校验，可以审核任意单位提交的档案
func (s *ExpertApprovalService) assertSameOrg(operatorID uint, profileOrgId *uint) error {
	var operator system.SysUser
	if err := global.GVA_DB.Where("id = ?", operatorID).First(&operator).Error; err != nil {
		return err
	}
	if operator.AuthorityId == superAdminAuthorityId {
		return nil
	}
	if operator.OrgId == nil || profileOrgId == nil || *operator.OrgId != *profileOrgId {
		return errors.New("无权审核其他单位提交的专家档案")
	}
	return nil
}

// listByFilter 分页查询专家档案的通用底座，供下面三个审核台列表接口复用
func (s *ExpertApprovalService) listByFilter(filter func(*gorm.DB) *gorm.DB, page commonRequest.PageInfo) (list []ExpertDatabase.ExpertProfile, total int64, err error) {
	limit := page.PageSize
	offset := page.PageSize * (page.Page - 1)
	db := filter(global.GVA_DB.Model(&ExpertDatabase.ExpertProfile{}))
	if err = db.Count(&total).Error; err != nil {
		return
	}
	// 按更新时间倒序，不是按创建时间——申报人编辑完草稿后，那条记录应该排到"我发起的"最前面，
	// 方便马上找到并提交，而不是淹没在一长串按创建顺序排列的旧记录里
	db = db.Order("updated_at desc")
	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}
	err = db.Find(&list).Error
	return
}

// GetMyDrafts 获取当前用户自己申报的专家档案列表（对应前端"我发起的"）
func (s *ExpertApprovalService) GetMyDrafts(userID uint, page commonRequest.PageInfo) ([]ExpertDatabase.ExpertProfile, int64, error) {
	return s.listByFilter(func(db *gorm.DB) *gorm.DB {
		return db.Where("submitted_by = ?", userID)
	}, page)
}

// GetPendingOrgReview 获取待本单位审核的专家档案列表（对应前端"待我审核的"，单位审核员视角）
// 审核台是所有角色共用的同一个页面，市级审核员/个人申报人等没有关联单位的账号也会触发这个查询，
// 未关联单位时返回空列表而不是报错，避免这类角色一打开"待我审核的"tab 就看到误导性的错误提示。
// 超级管理员是例外：它本身不挂靠任何单位，但理应能看到所有单位的待审记录（呼应 assertSameOrg
// 里对超级管理员的放行），否则就算能调用 orgApprove 接口，页面上也永远看不到能审的东西
func (s *ExpertApprovalService) GetPendingOrgReview(userID uint, page commonRequest.PageInfo) ([]ExpertDatabase.ExpertProfile, int64, error) {
	var operator system.SysUser
	if err := global.GVA_DB.Where("id = ?", userID).First(&operator).Error; err != nil {
		return nil, 0, err
	}
	if operator.AuthorityId == superAdminAuthorityId {
		return s.listByFilter(func(db *gorm.DB) *gorm.DB {
			return db.Where("status = ?", StatusPendingOrgReview)
		}, page)
	}
	if operator.OrgId == nil {
		return []ExpertDatabase.ExpertProfile{}, 0, nil
	}
	return s.listByFilter(func(db *gorm.DB) *gorm.DB {
		return db.Where("org_id = ? AND status = ?", operator.OrgId, StatusPendingOrgReview)
	}, page)
}

// GetPendingCityReview 获取待市级审核的专家档案列表（对应前端"待我审核的"，市级审核员视角，不限单位）
func (s *ExpertApprovalService) GetPendingCityReview(page commonRequest.PageInfo) ([]ExpertDatabase.ExpertProfile, int64, error) {
	return s.listByFilter(func(db *gorm.DB) *gorm.DB {
		return db.Where("status = ?", StatusPendingCityReview)
	}, page)
}

// GetExpertApprovalLogList 获取某专家档案的审核流转历史，按时间倒序（对应前端"历史记录"）
func (s *ExpertApprovalService) GetExpertApprovalLogList(info ExpertDatabaseReq.ExpertApprovalLogSearch) (list []ExpertDatabase.ExpertApprovalLog, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&ExpertDatabase.ExpertApprovalLog{}).Where("expert_id = ?", info.ExpertId)
	if err = db.Count(&total).Error; err != nil {
		return
	}
	db = db.Order("id desc")
	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}
	err = db.Find(&list).Error
	return
}
