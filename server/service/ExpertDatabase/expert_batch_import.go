package ExpertDatabase

import (
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase"
	ExpertDatabaseResp "github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

const (
	profileSheetName     = "专家背景信息"
	achievementSheetName = "研究成果"
)

var profileSheetHeaders = []string{
	"姓名*", "性别", "民族", "政治面貌", "所在单位*", "院系/部门", "行政职务", "专业技术职称",
	"办公电话", "手机号码", "电子邮箱", "最高学历", "最高学位", "毕业院校", "所学专业",
	"一级学科", "二级学科", "研究方向", "研究关键词",
}

var achievementSheetHeaders = []string{
	"专家姓名*", "成果类型*", "成果名称*", "成果级别", "发表/立项单位",
}

var validAchievementTypes = map[string]bool{
	"著作": true, "学术论文": true, "研究报告": true, "决策咨询成果": true, "获奖成果": true, "课题项目": true,
}

var validAchievementLevels = map[string]bool{
	"国家级": true, "省部级": true, "厅局级": true, "一般级": true,
}

// BuildImportTemplate 生成批量导入模板：背景信息 + 研究成果两个 sheet，成果表按姓名关联背景信息表
func (expertProfileService *ExpertProfileService) BuildImportTemplate() (*excelize.File, error) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			global.GVA_LOG.Error(err.Error())
		}
	}()

	if err := f.SetSheetName("Sheet1", profileSheetName); err != nil {
		return nil, err
	}
	for i, h := range profileSheetHeaders {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		if err := f.SetCellValue(profileSheetName, cell, h); err != nil {
			return nil, err
		}
	}
	exampleProfile := []interface{}{
		"张三", "男", "汉族", "中共党员", "北京大学", "公共管理学院", "副院长", "教授",
		"010-62751234", "13800001234", "zhangsan@example.edu.cn", "博士研究生", "博士", "北京大学", "行政管理",
		"公共管理", "行政管理", "数字政府建设", "数字政府,基层治理",
	}
	for i, v := range exampleProfile {
		cell, _ := excelize.CoordinatesToCellName(i+1, 2)
		if err := f.SetCellValue(profileSheetName, cell, v); err != nil {
			return nil, err
		}
	}

	achievementSheetIndex, err := f.NewSheet(achievementSheetName)
	if err != nil {
		return nil, err
	}
	for i, h := range achievementSheetHeaders {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		if err := f.SetCellValue(achievementSheetName, cell, h); err != nil {
			return nil, err
		}
	}
	exampleAchievement := []interface{}{"张三", "决策咨询成果", "关于推进基层数字政务建设的建议", "厅局级", "北京大学"}
	for i, v := range exampleAchievement {
		cell, _ := excelize.CoordinatesToCellName(i+1, 2)
		if err := f.SetCellValue(achievementSheetName, cell, v); err != nil {
			return nil, err
		}
	}

	// 成果类型/成果级别列加下拉校验，减少手填出错
	typeList := []string{"著作", "学术论文", "研究报告", "决策咨询成果", "获奖成果", "课题项目"}
	levelList := []string{"国家级", "省部级", "厅局级", "一般级"}
	typeDV := excelize.NewDataValidation(true)
	typeDV.SetSqref("B2:B1000")
	if err := typeDV.SetDropList(typeList); err == nil {
		_ = f.AddDataValidation(achievementSheetName, typeDV)
	}
	levelDV := excelize.NewDataValidation(true)
	levelDV.SetSqref("D2:D1000")
	if err := levelDV.SetDropList(levelList); err == nil {
		_ = f.AddDataValidation(achievementSheetName, levelDV)
	}

	f.SetActiveSheet(achievementSheetIndex)
	f.SetActiveSheet(0)
	return f, nil
}

// matchOrgId 按单位名称的关键词匹配已知单位表，匹配不上就留空（不是伪造数据，等后续人工/管理员补充）
func matchOrgId(unitName string) *uint {
	if unitName == "" {
		return nil
	}
	var orgs []system.SysOrganization
	if err := global.GVA_DB.Where("id <> 1").Find(&orgs).Error; err != nil {
		return nil
	}
	for _, o := range orgs {
		if o.Name != "" && strings.Contains(unitName, o.Name) {
			id := o.ID
			return &id
		}
	}
	return nil
}

func cellAt(row []string, idx int) string {
	if idx < 0 || idx >= len(row) {
		return ""
	}
	return strings.TrimSpace(row[idx])
}

func isBlankRow(row []string) bool {
	for _, c := range row {
		if strings.TrimSpace(c) != "" {
			return false
		}
	}
	return true
}

// ImportBatch 解析并导入批量上传的 Excel：先把两个 sheet 全部校验一遍，
// 有任何一行不合格就整体拒绝、不落库，避免"导入了一半"这种含糊状态；
// 全部通过才在一个事务里真正写入
func (expertProfileService *ExpertProfileService) ImportBatch(src io.Reader, operatorID uint) (ExpertDatabaseResp.ExpertBatchImportResult, error) {
	result := ExpertDatabaseResp.ExpertBatchImportResult{}

	f, err := excelize.OpenReader(src)
	if err != nil {
		return result, errors.New("无法解析上传的文件，请确认是不是 .xlsx 格式：" + err.Error())
	}
	defer f.Close()

	profileRows, err := f.GetRows(profileSheetName)
	if err != nil {
		return result, fmt.Errorf("找不到「%s」这个 sheet，请使用模板文件填写，不要重命名 sheet", profileSheetName)
	}
	achievementRows, err := f.GetRows(achievementSheetName)
	if err != nil {
		return result, fmt.Errorf("找不到「%s」这个 sheet，请使用模板文件填写，不要重命名 sheet", achievementSheetName)
	}

	var errs []ExpertDatabaseResp.ExpertBatchImportRowError

	type profileRow struct {
		rowNum int
		name   string
		fields map[string]string
	}
	nameToRow := map[string]*profileRow{}
	var profileList []*profileRow

	for i, row := range profileRows {
		rowNum := i + 1
		if rowNum == 1 {
			continue // 表头
		}
		if isBlankRow(row) {
			continue
		}
		name := cellAt(row, 0)
		unitName := cellAt(row, 4)
		if name == "" {
			errs = append(errs, ExpertDatabaseResp.ExpertBatchImportRowError{Sheet: profileSheetName, Row: rowNum, Message: "姓名不能为空"})
			continue
		}
		if unitName == "" {
			errs = append(errs, ExpertDatabaseResp.ExpertBatchImportRowError{Sheet: profileSheetName, Row: rowNum, Message: "所在单位不能为空"})
			continue
		}
		if _, dup := nameToRow[name]; dup {
			errs = append(errs, ExpertDatabaseResp.ExpertBatchImportRowError{Sheet: profileSheetName, Row: rowNum, Message: "姓名「" + name + "」在背景信息表里重复出现，成果表无法唯一关联，请改成可区分的姓名或分开提交"})
			continue
		}
		pr := &profileRow{
			rowNum: rowNum,
			name:   name,
			fields: map[string]string{
				"gender": cellAt(row, 1), "ethnicity": cellAt(row, 2), "politicalStatus": cellAt(row, 3),
				"unitName": unitName, "department": cellAt(row, 5), "adminTitle": cellAt(row, 6), "techTitle": cellAt(row, 7),
				"phone": cellAt(row, 8), "mobile": cellAt(row, 9), "email": cellAt(row, 10),
				"highestEducation": cellAt(row, 11), "highestDegree": cellAt(row, 12), "graduateSchool": cellAt(row, 13), "major": cellAt(row, 14),
				"disciplineL1": cellAt(row, 15), "disciplineL2": cellAt(row, 16), "researchDirections": cellAt(row, 17), "researchKeywords": cellAt(row, 18),
			},
		}
		nameToRow[name] = pr
		profileList = append(profileList, pr)
	}
	result.TotalProfileRows = len(profileList)

	type achievementRowData struct {
		rowNum          int
		expertName      string
		achievementType string
		title           string
		level           string
		publishOrg      string
	}
	var achievementList []*achievementRowData

	for i, row := range achievementRows {
		rowNum := i + 1
		if rowNum == 1 {
			continue
		}
		if isBlankRow(row) {
			continue
		}
		expertName := cellAt(row, 0)
		achievementType := cellAt(row, 1)
		title := cellAt(row, 2)
		level := cellAt(row, 3)
		publishOrg := cellAt(row, 4)

		if expertName == "" {
			errs = append(errs, ExpertDatabaseResp.ExpertBatchImportRowError{Sheet: achievementSheetName, Row: rowNum, Message: "专家姓名不能为空"})
			continue
		}
		if _, ok := nameToRow[expertName]; !ok {
			errs = append(errs, ExpertDatabaseResp.ExpertBatchImportRowError{Sheet: achievementSheetName, Row: rowNum, Message: "专家姓名「" + expertName + "」在「" + profileSheetName + "」表里找不到，请先在背景信息表里加上这个人"})
			continue
		}
		if achievementType == "" {
			errs = append(errs, ExpertDatabaseResp.ExpertBatchImportRowError{Sheet: achievementSheetName, Row: rowNum, Message: "成果类型不能为空"})
			continue
		}
		if !validAchievementTypes[achievementType] {
			errs = append(errs, ExpertDatabaseResp.ExpertBatchImportRowError{Sheet: achievementSheetName, Row: rowNum, Message: "成果类型「" + achievementType + "」不合法，只能是：著作/学术论文/研究报告/决策咨询成果/获奖成果/课题项目"})
			continue
		}
		if title == "" {
			errs = append(errs, ExpertDatabaseResp.ExpertBatchImportRowError{Sheet: achievementSheetName, Row: rowNum, Message: "成果名称不能为空"})
			continue
		}
		if level != "" && !validAchievementLevels[level] {
			errs = append(errs, ExpertDatabaseResp.ExpertBatchImportRowError{Sheet: achievementSheetName, Row: rowNum, Message: "成果级别「" + level + "」不合法，只能是：国家级/省部级/厅局级/一般级，留空表示待定"})
			continue
		}
		achievementList = append(achievementList, &achievementRowData{
			rowNum: rowNum, expertName: expertName, achievementType: achievementType, title: title, level: level, publishOrg: publishOrg,
		})
	}
	result.TotalAchievementRows = len(achievementList)

	if len(errs) > 0 {
		result.Success = false
		result.Errors = errs
		return result, nil
	}

	// 审核开关关闭时，批量导入的新档案直接落库为已发布，不进草稿/审核流程
	skipReview := !reviewRequired()
	initialStatus := StatusDraft
	if skipReview {
		initialStatus = StatusPublished
	}

	nameToExpertID := map[string]uint{}
	unmatchedUnits := map[string]bool{}
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		for _, pr := range profileList {
			var existing ExpertDatabase.ExpertProfile
			findErr := tx.Where("name = ? AND unit_name = ?", pr.name, pr.fields["unitName"]).First(&existing).Error
			if findErr == nil {
				nameToExpertID[pr.name] = existing.ID
				result.ReusedProfiles++
				continue
			}
			if !errors.Is(findErr, gorm.ErrRecordNotFound) {
				return findErr
			}

			// 单位名在"单位管理"里按关键词匹配不上，不阻断导入——只是先留空 org_id（提交审核时
			// 会兜底交给市级审核），同时把这个单位名记下来，导入完成后提示管理员去核实/补建，
			// 而不是自动建一个新单位（自由文本打法五花八门，自动建容易堆出一堆重复/近似的单位记录）
			orgId := matchOrgId(pr.fields["unitName"])
			if orgId == nil {
				unmatchedUnits[pr.fields["unitName"]] = true
			}

			profile := ExpertDatabase.ExpertProfile{
				Name: pr.name, Gender: pr.fields["gender"], Ethnicity: pr.fields["ethnicity"], PoliticalStatus: pr.fields["politicalStatus"],
				UnitName: pr.fields["unitName"], Department: pr.fields["department"], AdminTitle: pr.fields["adminTitle"], TechTitle: pr.fields["techTitle"],
				Phone: pr.fields["phone"], Mobile: pr.fields["mobile"], Email: pr.fields["email"],
				HighestEducation: pr.fields["highestEducation"], HighestDegree: pr.fields["highestDegree"], GraduateSchool: pr.fields["graduateSchool"], Major: pr.fields["major"],
				DisciplineL1: pr.fields["disciplineL1"], DisciplineL2: pr.fields["disciplineL2"], ResearchDirections: pr.fields["researchDirections"], ResearchKeywords: pr.fields["researchKeywords"],
				Status:         initialStatus,
				ReviewBypassed: skipReview,
				OrgId:          orgId,
				SubmittedBy:    &operatorID,
				CreatedBy:      operatorID,
				UpdatedBy:      operatorID,
			}
			if err := tx.Create(&profile).Error; err != nil {
				return err
			}
			nameToExpertID[pr.name] = profile.ID
			result.CreatedProfiles++
		}

		for _, ar := range achievementList {
			expertID := nameToExpertID[ar.expertName]
			// 去重策略：同一专家名下已经有一条标题完全相同的成果，就跳过而不是重复插入——
			// 标题是一份成果最直观的身份标识（论文/报告/课题基本不会真的重名），重复上传同一份
			// 模板不会越导越多；级别/发表单位如果这次填的不一样，不覆盖旧记录，避免静默改掉已有数据
			var dup int64
			if err := tx.Model(&ExpertDatabase.ExpertAchievement{}).
				Where("expert_id = ? AND title = ?", expertID, ar.title).
				Count(&dup).Error; err != nil {
				return err
			}
			if dup > 0 {
				result.SkippedAchievements++
				continue
			}
			achievement := ExpertDatabase.ExpertAchievement{
				ExpertId: expertID, AchievementType: ar.achievementType, Title: ar.title, Level: ar.level, PublishOrg: ar.publishOrg,
				CreatedBy: operatorID, UpdatedBy: operatorID,
			}
			if err := tx.Create(&achievement).Error; err != nil {
				return err
			}
			result.CreatedAchievements++
		}
		return nil
	})
	if err != nil {
		return ExpertDatabaseResp.ExpertBatchImportResult{}, errors.New("导入失败：" + err.Error())
	}

	if skipReview {
		// 关闭审核时新档案已经是已发布状态；复用的既有档案如果本来就是已发布，
		// 这次新增的成果也要计入得分——recomputeIfPublished 对非已发布档案是空操作，不会误改草稿/审核中的记录
		for _, expertID := range nameToExpertID {
			expertScoreSvc.recomputeIfPublished(expertID)
		}
	}

	for unitName := range unmatchedUnits {
		result.UnmatchedUnits = append(result.UnmatchedUnits, unitName)
	}
	sort.Strings(result.UnmatchedUnits)

	result.Success = true
	return result, nil
}
