package ExpertDatabase

import (
	"errors"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase"
	ExpertDatabaseReq "github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

// ErrProfileAccessDenied 数据域校验不通过时统一返回这个错误，API 层据此把权限问题原样提示给
// 用户，而不是笼统的"查询失败/更新失败"
var ErrProfileAccessDenied = errors.New("无权访问该专家档案")

// 专家库三个角色ID：个人申报人/单位审核员/市级审核员。专家主档列表对审核员角色做数据域收敛
// （见 GetExpertProfileInfoList），个人申报人角色ID被本单位账号管理功能用来固定新建账号的角色
const (
	authorityIndividualApplicant = 9001
	authorityOrgReviewer         = 9002
	authorityCityReviewer        = 9003
)

type ExpertProfileService struct{}

// CreateExpertProfile 创建专家主档；审核开关关闭时跳过草稿状态，直接落库为已发布并重算得分
func (expertProfileService *ExpertProfileService) CreateExpertProfile(profile *ExpertDatabase.ExpertProfile) (err error) {
	skipReview := !reviewRequired()
	if skipReview {
		profile.Status = StatusPublished
		profile.ReviewBypassed = true
	}
	if err = global.GVA_DB.Create(profile).Error; err != nil {
		return err
	}
	if skipReview {
		expertScoreSvc.recomputeIfPublished(profile.ID)
	}
	return nil
}

// canViewProfile 单条详情的可见性规则，必须跟列表页 GetExpertProfileInfoList 的数据域收敛完全
// 一致——列表上看不到的记录，拿着 ID 直接调详情接口也不能看到，否则收敛就只是"藏了个入口"
func canViewProfile(operator system.SysUser, profile ExpertDatabase.ExpertProfile) bool {
	switch operator.AuthorityId {
	case authorityOrgReviewer:
		if profile.Status == StatusPublished {
			return true
		}
		return operator.OrgId != nil && profile.OrgId != nil &&
			*operator.OrgId == *profile.OrgId && profile.Status != StatusDraft
	case authorityCityReviewer:
		return profile.Status == StatusPublished || profile.Status == StatusPendingCityReview
	default:
		return true
	}
}

// canModifyProfile 更新/删除比查看更严：个人申报人只能动自己创建的档案；单位审核员只能动本单位
// 的档案（他人草稿除外）；市级审核员只能动其可见范围内的档案；管理员等其余角色不受限
func canModifyProfile(operator system.SysUser, profile ExpertDatabase.ExpertProfile) bool {
	switch operator.AuthorityId {
	case authorityIndividualApplicant:
		return profile.CreatedBy == operator.ID
	case authorityOrgReviewer:
		if operator.OrgId == nil || profile.OrgId == nil || *operator.OrgId != *profile.OrgId {
			return false
		}
		return profile.Status != StatusDraft || profile.CreatedBy == operator.ID
	case authorityCityReviewer:
		return profile.Status == StatusPublished || profile.Status == StatusPendingCityReview
	default:
		return true
	}
}

func loadOperator(operatorID uint) (operator system.SysUser, err error) {
	err = global.GVA_DB.Where("id = ?", operatorID).First(&operator).Error
	return
}

// DeleteExpertProfile 删除专家主档；操作人对这条档案没有修改权时整体拒绝
func (expertProfileService *ExpertProfileService) DeleteExpertProfile(ID string, operatorID uint) (err error) {
	return expertProfileService.DeleteExpertProfileByIds([]string{ID}, operatorID)
}

// DeleteExpertProfileByIds 批量删除专家主档；任何一条越出操作人的数据域就整体拒绝——删除是
// 管理操作，全成或全败比"部分删掉了但不知道哪几条没删"更容易核对
func (expertProfileService *ExpertProfileService) DeleteExpertProfileByIds(IDs []string, operatorID uint) (err error) {
	operator, err := loadOperator(operatorID)
	if err != nil {
		return err
	}
	var profiles []ExpertDatabase.ExpertProfile
	if err = global.GVA_DB.Where("id in ?", IDs).Find(&profiles).Error; err != nil {
		return err
	}
	for _, p := range profiles {
		if !canModifyProfile(operator, p) {
			return ErrProfileAccessDenied
		}
	}
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&ExpertDatabase.ExpertProfile{}).Where("id in ?", IDs).Update("deleted_by", operatorID).Error; err != nil {
			return err
		}
		return tx.Where("id in ?", IDs).Delete(&ExpertDatabase.ExpertProfile{}).Error
	})
}

// UpdateExpertProfile 更新专家主档。这个接口只负责画像/背景字段的编辑：审核状态、得分缓存、
// 归属人这些字段一律以库里现值为准，请求体里就算夹带了也会被覆盖回去——改状态必须走审核
// 流转接口或管理员的 AdminSetStatus，改得分只能靠重算，否则任何能编辑档案的人都能把自己
// 的草稿直接"改"成已发布
func (expertProfileService *ExpertProfileService) UpdateExpertProfile(profile ExpertDatabase.ExpertProfile, operatorID uint) (err error) {
	operator, err := loadOperator(operatorID)
	if err != nil {
		return err
	}
	var existing ExpertDatabase.ExpertProfile
	if err = global.GVA_DB.Where("id = ?", profile.ID).First(&existing).Error; err != nil {
		return err
	}
	if !canModifyProfile(operator, existing) {
		return ErrProfileAccessDenied
	}
	profile.CreatedAt = existing.CreatedAt
	profile.Status = existing.Status
	profile.CurrentStep = existing.CurrentStep
	profile.SubmittedBy = existing.SubmittedBy
	profile.ReviewBypassed = existing.ReviewBypassed
	profile.AchievementScore = existing.AchievementScore
	profile.InfluenceScore = existing.InfluenceScore
	profile.SocialScore = existing.SocialScore
	profile.CompositeScore = existing.CompositeScore
	profile.ScoreUpdatedAt = existing.ScoreUpdatedAt
	profile.CreatedBy = existing.CreatedBy
	profile.DeletedBy = existing.DeletedBy
	return global.GVA_DB.Save(&profile).Error
}

// GetExpertProfile 根据ID获取专家主档，按操作人角色做数据域校验（规则见 canViewProfile）
func (expertProfileService *ExpertProfileService) GetExpertProfile(ID string, operatorID uint) (profile ExpertDatabase.ExpertProfile, err error) {
	operator, err := loadOperator(operatorID)
	if err != nil {
		return
	}
	if err = global.GVA_DB.Where("id = ?", ID).First(&profile).Error; err != nil {
		return
	}
	if !canViewProfile(operator, profile) {
		return ExpertDatabase.ExpertProfile{}, ErrProfileAccessDenied
	}
	return
}

// GetExpertProfileInfoList 分页获取专家主档列表
//
// 单位审核员/市级审核员在审核台之外，还能通过这个列表页查到全库专家——原先没有按角色收敛，
// 相当于谁都能翻到别人、别的单位还没提交的草稿，跟"审核台只显示轮到自己审的记录"这个预期不一致。
// 这里补上：单位审核员只看得到已发布的 + 本单位已提交（非草稿）的；市级审核员只看得到已发布的 +
// 待市级审核的；管理员和个人申报人不受影响，保持原有的全量可见。
func (expertProfileService *ExpertProfileService) GetExpertProfileInfoList(info ExpertDatabaseReq.ExpertProfileSearch, operatorID uint) (list []ExpertDatabase.ExpertProfile, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&ExpertDatabase.ExpertProfile{})
	var profiles []ExpertDatabase.ExpertProfile

	var operator system.SysUser
	if err = global.GVA_DB.Where("id = ?", operatorID).First(&operator).Error; err != nil {
		return
	}
	switch operator.AuthorityId {
	case authorityOrgReviewer:
		if operator.OrgId != nil {
			db = db.Where("status = ? OR (org_id = ? AND status <> ?)", StatusPublished, *operator.OrgId, StatusDraft)
		} else {
			db = db.Where("status = ?", StatusPublished)
		}
	case authorityCityReviewer:
		db = db.Where("status = ? OR status = ?", StatusPublished, StatusPendingCityReview)
	}

	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.Name != "" {
		db = db.Where("name LIKE ?", "%"+info.Name+"%")
	}
	if info.UnitName != "" {
		db = db.Where("unit_name LIKE ?", "%"+info.UnitName+"%")
	}
	if info.TechTitle != "" {
		db = db.Where("tech_title = ?", info.TechTitle)
	}
	if info.DisciplineL1 != "" {
		db = db.Where("discipline_l1 = ?", info.DisciplineL1)
	}
	if info.DisciplineL2 != "" {
		db = db.Where("discipline_l2 = ?", info.DisciplineL2)
	}
	if info.ResearchKeywords != "" {
		db = db.Where("research_keywords LIKE ?", "%"+info.ResearchKeywords+"%")
	}
	if info.Status != "" {
		db = db.Where("status = ?", info.Status)
	}
	if info.OrgId != nil {
		db = db.Where("org_id = ?", info.OrgId)
	}
	if info.OrgUnmatched {
		db = db.Where("org_id IS NULL")
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	orderMap := map[string]bool{
		"created_at":        true,
		"achievement_score": true,
		"influence_score":   true,
		"composite_score":   true,
	}
	if orderMap[info.Sort] {
		orderStr := info.Sort
		if info.Order == "descending" {
			orderStr += " desc"
		}
		db = db.Order(orderStr)
	} else {
		db = db.Order("id desc")
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&profiles).Error
	return profiles, total, err
}

var profileExportHeaders = []string{
	"姓名", "性别", "民族", "政治面貌", "所在单位", "院系/部门", "行政职务", "专业技术职称",
	"办公电话", "手机号码", "电子邮箱", "最高学历", "最高学位", "毕业院校", "所学专业",
	"一级学科", "二级学科", "研究方向", "研究关键词", "审核状态", "综合排序得分", "创建日期",
}

// ExportExpertProfiles 导出当前筛选条件下命中的专家列表（尊重跟列表页完全相同的筛选条件和角色数据域
// 收敛，见 GetExpertProfileInfoList）；PageSize 强制清零，拿全部匹配记录而不是当前这一页
func (expertProfileService *ExpertProfileService) ExportExpertProfiles(info ExpertDatabaseReq.ExpertProfileSearch, operatorID uint) (*excelize.File, error) {
	info.Page = 1
	info.PageSize = 0
	list, _, err := expertProfileService.GetExpertProfileInfoList(info, operatorID)
	if err != nil {
		return nil, err
	}

	f := excelize.NewFile()
	sheetName := "专家列表"
	if err := f.SetSheetName("Sheet1", sheetName); err != nil {
		return nil, err
	}
	for i, h := range profileExportHeaders {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		if err := f.SetCellValue(sheetName, cell, h); err != nil {
			return nil, err
		}
	}
	for rowIdx, p := range list {
		row := []interface{}{
			p.Name, p.Gender, p.Ethnicity, p.PoliticalStatus, p.UnitName, p.Department, p.AdminTitle, p.TechTitle,
			p.Phone, p.Mobile, p.Email, p.HighestEducation, p.HighestDegree, p.GraduateSchool, p.Major,
			p.DisciplineL1, p.DisciplineL2, p.ResearchDirections, p.ResearchKeywords,
			statusLabel[p.Status], p.CompositeScore, p.CreatedAt.Format("2006-01-02"),
		}
		for i, v := range row {
			cell, _ := excelize.CoordinatesToCellName(i+1, rowIdx+2)
			if err := f.SetCellValue(sheetName, cell, v); err != nil {
				return nil, err
			}
		}
	}
	return f, nil
}
