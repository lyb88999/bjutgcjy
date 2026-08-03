package ExpertDatabase

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase"
	ExpertDatabaseReq "github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase/request"
	ExpertDatabaseRes "github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase/response"
	"github.com/xuri/excelize/v2"
)

type ExpertSearchService struct{}

// buildSearchCorpusMap 一次性预取一批候选专家的检索语料。语料来源三块：① 专家档案自己填的
// 研究方向/关键词/关注议题等字段——这块以前漏掉了，导致专家档案里明明白白写着的主题词都搜不到，
// 是最基础的召回缺口；② 成果标题/关键词；③ 关联标签。②③整批只发两条 SQL——替代原先逐个专家
// 各查两次的写法，候选集几百上千条时检索耗时不再随人数线性膨胀。标签 JOIN 只认未解除的关联
// （relation 是软删除，解除过的标签不该再计入语料）
func (s *ExpertSearchService) buildSearchCorpusMap(candidates []ExpertDatabase.ExpertProfile) (map[uint]string, error) {
	builders := make(map[uint]*strings.Builder, len(candidates))
	appendTerm := func(expertID uint, term string) {
		if term == "" {
			return
		}
		sb, ok := builders[expertID]
		if !ok {
			sb = &strings.Builder{}
			builders[expertID] = sb
		}
		sb.WriteString(term)
		sb.WriteString(" ")
	}

	expertIDs := make([]uint, 0, len(candidates))
	for _, p := range candidates {
		expertIDs = append(expertIDs, p.ID)
		for _, f := range []string{
			p.ResearchDirections, p.ResearchKeywords, p.FocusTopics, p.ResearchObjects,
			p.MethodExpertise, p.RegionExpertise, p.DisciplineL1, p.DisciplineL2,
			p.CrossDiscipline, p.PolicyFields,
		} {
			appendTerm(p.ID, f)
		}
	}

	var achievements []ExpertDatabase.ExpertAchievement
	if err := global.GVA_DB.Where("expert_id IN ?", expertIDs).Find(&achievements).Error; err != nil {
		return nil, err
	}
	for _, a := range achievements {
		appendTerm(a.ExpertId, a.Title)
		appendTerm(a.ExpertId, a.Keywords)
	}

	type expertTagRow struct {
		ExpertId uint
		TagValue string
	}
	var tagRows []expertTagRow
	if err := global.GVA_DB.Model(&ExpertDatabase.ExpertTag{}).
		Select("expert_tag_relation.expert_id AS expert_id, expert_tag.tag_value AS tag_value").
		Joins("JOIN expert_tag_relation ON expert_tag_relation.tag_id = expert_tag.id AND expert_tag_relation.deleted_at IS NULL").
		Where("expert_tag_relation.expert_id IN ?", expertIDs).
		Scan(&tagRows).Error; err != nil {
		return nil, err
	}
	for _, row := range tagRows {
		appendTerm(row.ExpertId, row.TagValue)
	}

	corpus := make(map[uint]string, len(builders))
	for id, sb := range builders {
		corpus[id] = sb.String()
	}
	return corpus, nil
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
	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.UnitName != "" {
		db = db.Where("unit_name LIKE ?", "%"+req.UnitName+"%")
	}
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

	// 只有带关键词检索才需要算相关性；默认浏览（不带关键词）时 relevance 恒为 1，下面这些查询整个跳过
	corpus := map[uint]string{}
	candidateVectors := map[uint][]float32{}
	var keywordVector []float32
	useSemanticMatch := false
	if req.Keyword != "" && len(candidates) > 0 {
		if corpus, err = s.buildSearchCorpusMap(candidates); err != nil {
			return nil, err
		}

		// 语义向量匹配优先于关键词子串匹配：子串匹配要求检索词整体原样出现在语料里，"人工智能"
		// 搜不出"具身智能"这类相关但不同字面的内容；向量是离线预算好存在 expert_search_embedding
		// 里的（见 cmd/recompute-embeddings），这里只需要现算一次关键词的向量。embedding 服务
		// 调不通时（网络问题/服务没起来）整体退回子串匹配，不能让语义检索的故障拖垮基本检索能力
		if vecs, embedErr := embedTexts([]string{req.Keyword}); embedErr == nil && len(vecs) == 1 {
			keywordVector = vecs[0]
			ids := make([]uint, 0, len(candidates))
			for _, c := range candidates {
				ids = append(ids, c.ID)
			}
			var embRows []ExpertDatabase.ExpertSearchEmbedding
			if dbErr := global.GVA_DB.Where("expert_id IN ?", ids).Find(&embRows).Error; dbErr == nil {
				for _, row := range embRows {
					var v []float32
					if json.Unmarshal([]byte(row.Vector), &v) == nil {
						candidateVectors[row.ExpertId] = v
					}
				}
				useSemanticMatch = true
			}
		}
	}

	// minSemanticRelevance 语义相似度低于这个阈值就当作没命中——余弦相似度是连续值，几乎不会
	// 正好等于 0，不做下限截断的话每次检索都会把全库按相似度排一遍，"不相关"这个概念就没有了。
	// 0.5 是多组关键词实测折中出来的经验值。这里有个已知的真实局限，记录下来别以后重新踩一遍：
	// 调低到 0.45 能多召回一点"碳中和→低碳"这类近义词场景，但会把"京津冀协同发展"这类较长复合
	// 短语的命中数从 177 冲到 300（超过全库 70%），显然是噪音而不是真的相关——bge-small-zh 这个
	// 量级的模型对短语越长、越具体，区分度反而越模糊，全局阈值没法同时让两类查询都满意。0.5 对
	// 多数场景更稳，代价是偶尔漏掉个别近义词（碳中和/低碳这种），比多数查询混进一堆噪音更可接受
	const minSemanticRelevance = 0.5
	// literalHitRelevanceFloor 关键词原样出现在语料里时的相关性保底值——语义相似度是模糊估计，
	// 可能因为句子整体风格不像而给出偏低的分数（范明志"数据法治、数字司法与人工智能法律治理"
	// 实测只有 0.49，卡在阈值下面），但字面命中是无可辩驳的强信号，不该让语义分把这种结果排没了
	const literalHitRelevanceFloor = 0.6

	items = make([]ExpertDatabaseRes.ExpertSearchItem, 0, len(candidates))
	for _, c := range candidates {
		var relevance float64
		matched := true
		switch {
		case req.Keyword == "":
			relevance = 1
		case useSemanticMatch:
			if vec, ok := candidateVectors[c.ID]; ok {
				sim := cosineSim(keywordVector, vec)
				if sim < 0 {
					sim = 0
				}
				relevance = sim
				matched = sim >= minSemanticRelevance
				if literalHit := s.relevance(corpus[c.ID], req.Keyword) > 0; literalHit {
					matched = true
					if relevance < literalHitRelevanceFloor {
						relevance = literalHitRelevanceFloor
					}
				}
			} else {
				// 这个专家还没生成向量（刚发布/刚导入，还没跑过 cmd/recompute-embeddings），
				// 不能因为向量缺失就把人整个漏掉，退回子串匹配
				relevance = s.relevance(corpus[c.ID], req.Keyword)
				matched = relevance > 0
			}
		default:
			relevance = s.relevance(corpus[c.ID], req.Keyword)
			matched = relevance > 0
		}
		// 填了关键词却没命中的候选人直接跳过，不进结果集——不然职称权重、社会贡献分这些不受
		// 相关性影响的分项会把一堆跟关键词毫不沾边的人顶到排名前面，"检索"就退化成了"不管搜
		// 什么都是把全库按职称排一遍"。关键词为空时 relevance 恒为 1 且 matched 恒为 true，
		// 这里的判断天然不影响不带关键词的默认浏览场景。
		if req.Keyword != "" && !matched {
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
	if len(items) > exportRowLimit {
		return nil, fmt.Errorf("命中 %d 条记录，超过单次导出上限 %d 条，请缩小检索范围后再导出", len(items), exportRowLimit)
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
