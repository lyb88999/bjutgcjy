package ExpertDatabase

import (
	"bytes"
	"strings"
	"testing"

	"github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/xuri/excelize/v2"
)

// buildImportFile 按批量导入模板的格式在内存里拼一个 xlsx：两个 sheet、第一行表头、后面数据行
func buildImportFile(t *testing.T, profileRows [][]interface{}, achievementRows [][]interface{}) *bytes.Reader {
	t.Helper()
	f := excelize.NewFile()
	if err := f.SetSheetName("Sheet1", profileSheetName); err != nil {
		t.Fatalf("建 sheet 失败: %v", err)
	}
	if _, err := f.NewSheet(achievementSheetName); err != nil {
		t.Fatalf("建 sheet 失败: %v", err)
	}
	writeRows := func(sheet string, headers []string, rows [][]interface{}) {
		for i, h := range headers {
			cell, _ := excelize.CoordinatesToCellName(i+1, 1)
			if err := f.SetCellValue(sheet, cell, h); err != nil {
				t.Fatalf("写表头失败: %v", err)
			}
		}
		for r, row := range rows {
			for i, v := range row {
				cell, _ := excelize.CoordinatesToCellName(i+1, r+2)
				if err := f.SetCellValue(sheet, cell, v); err != nil {
					t.Fatalf("写数据失败: %v", err)
				}
			}
		}
	}
	writeRows(profileSheetName, profileSheetHeaders, profileRows)
	writeRows(achievementSheetName, achievementSheetHeaders, achievementRows)
	buf, err := f.WriteToBuffer()
	if err != nil {
		t.Fatalf("生成 xlsx 失败: %v", err)
	}
	return bytes.NewReader(buf.Bytes())
}

// profileRowFor 生成一行只填必填字段（姓名、所在单位）的背景信息行
func profileRowFor(name, unitName string) []interface{} {
	return []interface{}{name, "", "", "", unitName}
}

func TestImportBatch_ValidFileCreatesProfilesAndMatchesOrg(t *testing.T) {
	db := setupTestDB(t)
	// matchOrgId 会排除 id=1（根节点），所以先占位一条，再建真正参与匹配的单位
	if err := db.Create(&system.SysOrganization{Name: "根节点"}).Error; err != nil {
		t.Fatalf("建占位单位失败: %v", err)
	}
	pku := system.SysOrganization{Name: "北京大学"}
	if err := db.Create(&pku).Error; err != nil {
		t.Fatalf("建单位失败: %v", err)
	}

	src := buildImportFile(t,
		[][]interface{}{
			profileRowFor("张三", "北京大学公共管理学院"),
			profileRowFor("李四", "不存在的野单位"),
		},
		[][]interface{}{
			{"张三", "决策咨询成果", "关于基层治理的建议", "厅局级", "北京大学"},
			{"张三", "学术论文", "数字政府研究", "", ""},
		},
	)

	svc := &ExpertProfileService{}
	result, err := svc.ImportBatch(src, 42)
	if err != nil {
		t.Fatalf("导入不应报错: %v", err)
	}
	if !result.Success {
		t.Fatalf("导入应成功，错误：%+v", result.Errors)
	}
	if result.CreatedProfiles != 2 || result.CreatedAchievements != 2 {
		t.Fatalf("应新建 2 档案 2 成果，实际 %d/%d", result.CreatedProfiles, result.CreatedAchievements)
	}

	var zhangsan ExpertDatabase.ExpertProfile
	if err := db.Where("name = ?", "张三").First(&zhangsan).Error; err != nil {
		t.Fatalf("查张三失败: %v", err)
	}
	if zhangsan.OrgId == nil || *zhangsan.OrgId != pku.ID {
		t.Fatalf("单位名含「北京大学」应匹配到已有单位，实际 orgId=%v", zhangsan.OrgId)
	}
	if zhangsan.Status != StatusDraft {
		t.Fatalf("审核开关默认开启，导入档案应是草稿，实际 %s", zhangsan.Status)
	}
	var lisi ExpertDatabase.ExpertProfile
	if err := db.Where("name = ?", "李四").First(&lisi).Error; err != nil {
		t.Fatalf("查李四失败: %v", err)
	}
	if lisi.OrgId != nil {
		t.Fatal("匹配不上的单位应留空 org_id，而不是乱关联")
	}
	if len(result.UnmatchedUnits) != 1 || result.UnmatchedUnits[0] != "不存在的野单位" {
		t.Fatalf("匹配不上的单位名应回报给调用方，实际 %v", result.UnmatchedUnits)
	}
}

func TestImportBatch_AnyInvalidRowRejectsWholeFile(t *testing.T) {
	db := setupTestDB(t)
	src := buildImportFile(t,
		[][]interface{}{
			profileRowFor("张三", "北京大学"),
			profileRowFor("", "缺姓名单位"), // 非法：姓名为空
		},
		[][]interface{}{
			{"张三", "野类型", "标题", "", ""},  // 非法：成果类型不在枚举里
			{"王五", "学术论文", "标题2", "", ""}, // 非法：背景信息表里没有王五
		},
	)

	svc := &ExpertProfileService{}
	result, err := svc.ImportBatch(src, 42)
	if err != nil {
		t.Fatalf("校验不通过应通过 result 报告，不应返回 err: %v", err)
	}
	if result.Success {
		t.Fatal("有非法行时应整体拒绝")
	}
	if len(result.Errors) != 3 {
		t.Fatalf("应报出 3 条行级错误，实际 %d 条：%+v", len(result.Errors), result.Errors)
	}
	var count int64
	db.Model(&ExpertDatabase.ExpertProfile{}).Count(&count)
	if count != 0 {
		t.Fatalf("整体拒绝时不应有任何落库，实际写入 %d 条", count)
	}
}

func TestImportBatch_ReimportReusesProfilesAndSkipsDuplicateAchievements(t *testing.T) {
	setupTestDB(t)
	build := func() *bytes.Reader {
		return buildImportFile(t,
			[][]interface{}{profileRowFor("张三", "北京大学")},
			[][]interface{}{{"张三", "学术论文", "同一篇论文", "省部级", ""}},
		)
	}

	svc := &ExpertProfileService{}
	if _, err := svc.ImportBatch(build(), 42); err != nil {
		t.Fatalf("首次导入不应报错: %v", err)
	}
	result, err := svc.ImportBatch(build(), 42)
	if err != nil {
		t.Fatalf("重复导入不应报错: %v", err)
	}
	if result.CreatedProfiles != 0 || result.ReusedProfiles != 1 {
		t.Fatalf("同名同单位应复用既有档案，实际新建 %d 复用 %d", result.CreatedProfiles, result.ReusedProfiles)
	}
	if result.CreatedAchievements != 0 || result.SkippedAchievements != 1 {
		t.Fatalf("同名成果应跳过而不是重复插入，实际新建 %d 跳过 %d", result.CreatedAchievements, result.SkippedAchievements)
	}
}

func TestImportBatch_MissingSheetGivesActionableError(t *testing.T) {
	setupTestDB(t)
	f := excelize.NewFile()
	buf, err := f.WriteToBuffer()
	if err != nil {
		t.Fatalf("生成 xlsx 失败: %v", err)
	}

	svc := &ExpertProfileService{}
	_, err = svc.ImportBatch(bytes.NewReader(buf.Bytes()), 42)
	if err == nil || !strings.Contains(err.Error(), profileSheetName) {
		t.Fatalf("缺 sheet 时应报出具体 sheet 名，实际 err=%v", err)
	}
}
