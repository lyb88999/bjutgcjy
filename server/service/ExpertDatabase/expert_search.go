package ExpertDatabase

import (
	"sort"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase"
	ExpertDatabaseReq "github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase/request"
	ExpertDatabaseRes "github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase/response"
	"github.com/xuri/excelize/v2"
)

type ExpertSearchService struct{}

// buildSearchCorpus 拼接某专家用于相关性计算的检索语料：成果标题/关键词 + 关联标签
func (s *ExpertSearchService) buildSearchCorpus(expertID uint) string {
	var achievements []ExpertDatabase.ExpertAchievement
	global.GVA_DB.Where("expert_id = ?", expertID).Find(&achievements)
	var sb strings.Builder
	for _, a := range achievements {
		sb.WriteString(a.Title)
		sb.WriteString(" ")
		sb.WriteString(a.Keywords)
		sb.WriteString(" ")
	}
	var tags []ExpertDatabase.ExpertTag
	global.GVA_DB.
		Joins("JOIN expert_tag_relation ON expert_tag_relation.tag_id = expert_tag.id").
		Where("expert_tag_relation.expert_id = ?", expertID).
		Find(&tags)
	for _, t := range tags {
		sb.WriteString(t.TagValue)
		sb.WriteString(" ")
	}
	return sb.String()
}

// splitKeyword 检索词切分 V1：按常见分隔符拆词，够用即可；后续量大再引入分词库/ES（见技术方案 2.2）
func splitKeyword(keyword string) []string {
	replacer := strings.NewReplacer(",", " ", "，", " ", "、", " ", "/", " ", ";", " ", "；", " ")
	return strings.Fields(replacer.Replace(keyword))
}

// relevance 相关性系数 V1：检索词命中语料的比例，未输入关键词时按 1（不做主题区分）
func (s *ExpertSearchService) relevance(corpus string, keyword string) float64 {
	terms := splitKeyword(keyword)
	if len(terms) == 0 {
		return 1
	}
	hit := 0
	for _, t := range terms {
		if strings.Contains(corpus, t) {
			hit++
		}
	}
	return float64(hit) / float64(len(terms))
}

// rankedCandidates 按筛选条件圈定已发布候选集，并按
// realtimeScore = w1*achievement_score*relevance + w2*influence_score*relevance + w3*title_score + w4*social_score
// 算好相关性/实时得分、排好序，返回不分页的完整结果——Search()（列表分页）和 ExportSearchResults()
// （导出全部）共用同一份排序逻辑，避免排序算法在两个地方各写一遍、后续改权重容易漏改一处
func (s *ExpertSearchService) rankedCandidates(req ExpertDatabaseReq.ExpertSearchReq) (items []ExpertDatabaseRes.ExpertSearchItem, err error) {
	db := global.GVA_DB.Model(&ExpertDatabase.ExpertProfile{}).Where("status = ?", "published")
	if req.DisciplineL1 != "" {
		db = db.Where("discipline_l1 = ?", req.DisciplineL1)
	}
	if req.RegionExpertise != "" {
		db = db.Where("region_expertise LIKE ?", "%"+req.RegionExpertise+"%")
	}
	if req.TechTitle != "" {
		db = db.Where("tech_title = ?", req.TechTitle)
	}

	var candidates []ExpertDatabase.ExpertProfile
	if err = db.Find(&candidates).Error; err != nil {
		return
	}

	rankingWeights := expertScoreSvc.dictWeights(DictTypeRankingWeight, defaultRankingWeights)
	titleWeights := expertScoreSvc.dictWeights(DictTypeTitleLevel, defaultTitleLevelWeights)

	items = make([]ExpertDatabaseRes.ExpertSearchItem, 0, len(candidates))
	for _, c := range candidates {
		relevance := s.relevance(s.buildSearchCorpus(c.ID), req.Keyword)
		// 填了关键词却一个词都没命中的候选人直接跳过，不进结果集——不然职称权重、社会贡献分
		// 这些不受相关性影响的分项会把一堆跟关键词毫不沾边的人顶到排名前面，"检索"就退化成了
		// "不管搜什么都是把全库按职称排一遍"，relevance 字段形同虚设。关键词为空时 relevance()
		// 总是返回 1，这里的判断天然不影响不带关键词的默认浏览场景。
		if req.Keyword != "" && relevance == 0 {
			continue
		}
		titleScore, ok := titleWeights[c.TechTitle]
		if !ok {
			titleScore = 1
		}
		realtimeScore := rankingWeights["achievement"]*c.AchievementScore*relevance +
			rankingWeights["influence"]*c.InfluenceScore*relevance +
			rankingWeights["title"]*titleScore +
			rankingWeights["social"]*c.SocialScore
		items = append(items, ExpertDatabaseRes.ExpertSearchItem{
			ExpertProfile: c,
			Relevance:     relevance,
			RealtimeScore: realtimeScore,
		})
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].RealtimeScore > items[j].RealtimeScore
	})
	return items, nil
}

// Search 专家综合推荐排序检索：圈定候选集、排好序之后再分页返回
func (s *ExpertSearchService) Search(req ExpertDatabaseReq.ExpertSearchReq) (list []ExpertDatabaseRes.ExpertSearchItem, total int64, err error) {
	items, err := s.rankedCandidates(req)
	if err != nil {
		return
	}
	total = int64(len(items))

	page := req.Page
	pageSize := req.PageSize
	if pageSize == 0 {
		pageSize = 10
	}
	if page == 0 {
		page = 1
	}
	start := (page - 1) * pageSize
	if start >= len(items) {
		return []ExpertDatabaseRes.ExpertSearchItem{}, total, nil
	}
	end := start + pageSize
	if end > len(items) {
		end = len(items)
	}
	return items[start:end], total, nil
}

var searchExportHeaders = []string{
	"姓名", "所在单位", "专业技术职称", "一级学科", "研究关键词",
	"相关性", "实时综合得分", "成果分", "决策影响分", "社会贡献分",
}

// ExportSearchResults 导出当前检索条件下命中的全部结果（按相关性/实时得分排好序，不分页）
func (s *ExpertSearchService) ExportSearchResults(req ExpertDatabaseReq.ExpertSearchReq) (*excelize.File, error) {
	items, err := s.rankedCandidates(req)
	if err != nil {
		return nil, err
	}

	f := excelize.NewFile()
	sheetName := "检索结果"
	if err := f.SetSheetName("Sheet1", sheetName); err != nil {
		return nil, err
	}
	for i, h := range searchExportHeaders {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		if err := f.SetCellValue(sheetName, cell, h); err != nil {
			return nil, err
		}
	}
	for rowIdx, item := range items {
		row := []interface{}{
			item.Name, item.UnitName, item.TechTitle, item.DisciplineL1, item.ResearchKeywords,
			item.Relevance, item.RealtimeScore, item.AchievementScore, item.InfluenceScore, item.SocialScore,
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
