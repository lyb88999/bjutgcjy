package ExpertDatabase

import (
	"sort"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase"
	ExpertDatabaseReq "github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase/request"
	ExpertDatabaseRes "github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase/response"
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

// Search 专家综合推荐排序检索：先按筛选条件圈定已发布候选集，再按
// realtimeScore = w1*achievement_score*relevance + w2*influence_score*relevance + w3*title_score + w4*social_score 排序
func (s *ExpertSearchService) Search(req ExpertDatabaseReq.ExpertSearchReq) (list []ExpertDatabaseRes.ExpertSearchItem, total int64, err error) {
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
	total = int64(len(candidates))

	rankingWeights := expertScoreSvc.dictWeights(DictTypeRankingWeight, defaultRankingWeights)
	titleWeights := expertScoreSvc.dictWeights(DictTypeTitleLevel, defaultTitleLevelWeights)

	items := make([]ExpertDatabaseRes.ExpertSearchItem, 0, len(candidates))
	for _, c := range candidates {
		relevance := s.relevance(s.buildSearchCorpus(c.ID), req.Keyword)
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
