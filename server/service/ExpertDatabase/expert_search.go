package ExpertDatabase

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

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
	itemsByExpert := map[uint][]ExpertDatabase.ExpertSearchEmbeddingItem{}
	var keywordVectors [][]float32
	expandedKeywords := []string{req.Keyword}
	useSemanticMatch := false
	if req.Keyword != "" && len(candidates) > 0 {
		if corpus, err = s.buildSearchCorpusMap(candidates); err != nil {
			return nil, err
		}

		// 先按同义词字典把检索词展开成一组词（没配同义词的话就还是只有原词一个）。"碳中和"和
		// "低碳"这类近义关系 embedding 模型自己判断不出来（实测过），这一步是确定性的人工兜底，
		// 不依赖模型判不判断得准——见 expert_search_synonym.go 顶部注释
		expandedKeywords = expandKeywordWithSynonyms(req.Keyword, s.searchSynonymGroups())

		// 语义向量匹配优先于关键词子串匹配：子串匹配要求检索词整体原样出现在语料里，"人工智能"
		// 搜不出"具身智能"这类相关但不同字面的内容。每个专家名下拆成好几条独立语料条目分别存好
		// 向量（研究方向一条、每篇成果各一条、标签一条，见 cmd/recompute-embeddings），这里把
		// 展开后的每个词都算一次向量，再挨个专家取"跟任意一个展开词、任意一条语料最相似的那一对"。
		// embedding 服务调不通时（网络问题/服务没起来）整体退回子串匹配，不能让语义检索的故障
		// 拖垮基本检索能力
		keywordEmbedTimeout := time.Duration(global.GVA_CONFIG.ExpertEmbedding.TimeoutSec) * time.Second
		if vecs, embedErr := embedTexts(expandedKeywords, keywordEmbedTimeout); embedErr == nil && len(vecs) == len(expandedKeywords) {
			keywordVectors = vecs
			ids := make([]uint, 0, len(candidates))
			for _, c := range candidates {
				ids = append(ids, c.ID)
			}
			var embRows []ExpertDatabase.ExpertSearchEmbeddingItem
			if dbErr := global.GVA_DB.Where("expert_id IN ?", ids).Find(&embRows).Error; dbErr == nil {
				for _, row := range embRows {
					itemsByExpert[row.ExpertId] = append(itemsByExpert[row.ExpertId], row)
				}
				useSemanticMatch = true
			}
		}
	}

	// minSemanticRelevance 语义相似度低于这个阈值就当作没命中。改成按条目（研究方向/每篇成果/
	// 标签各自独立）取最大相似度之后，同一个阈值不能沿用旧的 0.5——按条目匹配相当于给每个专家
	// 多了好几次"够到"关键词的机会，同样的阈值下噪音会变多（实测"京津冀协同发展"从旧方案的 177
	// 条冲到 227 条）。校到 0.6 之后"人工智能""养老服务"命中的都是干净结果，"京津冀协同发展"
	// 压到 37 条（比旧方案的 177 还准）。已知局限依旧在："碳中和"搜不出"低碳"相关的专家——查过
	// 陆小成这批人按条目算的最高分也就 0.46~0.48，本质是 bge-small-zh 对这两个词的语义关联判断
	// 得不够高，不是语料拼接方式的问题，调阈值解决不了，只能靠更大的模型或者人工同义词表兜底
	const minSemanticRelevance = 0.6
	// literalHitRelevanceFloor 关键词原样出现在语料里时的相关性保底值——语义相似度是模糊估计，
	// 偶尔会打偏低分，但字面命中是无可辩驳的强信号，不该让语义分把这种结果排没了
	const literalHitRelevanceFloor = 0.6

	items = make([]ExpertDatabaseRes.ExpertSearchItem, 0, len(candidates))
	for _, c := range candidates {
		var relevance float64
		var matchReason string
		matched := true
		switch {
		case req.Keyword == "":
			relevance = 1
		case useSemanticMatch:
			if embItems, ok := itemsByExpert[c.ID]; ok && len(embItems) > 0 {
				bestSim := -1.0
				for _, it := range embItems {
					var v []float32
					if json.Unmarshal([]byte(it.Vector), &v) != nil {
						continue
					}
					for _, kwVec := range keywordVectors {
						if sim := cosineSim(kwVec, v); sim > bestSim {
							bestSim = sim
							matchReason = it.ItemLabel
						}
					}
				}
				if bestSim < 0 {
					bestSim = 0
				}
				relevance = bestSim
				matched = bestSim >= minSemanticRelevance
				if s.anyLiteralHit(corpus[c.ID], expandedKeywords) {
					matched = true
					if relevance < literalHitRelevanceFloor {
						relevance = literalHitRelevanceFloor
					}
				}
			} else {
				// 这个专家还没生成向量条目（刚发布/刚导入，还没跑过 cmd/recompute-embeddings），
				// 不能因为向量缺失就把人整个漏掉，退回子串匹配
				matched = s.anyLiteralHit(corpus[c.ID], expandedKeywords)
				if matched {
					relevance = 1
					matchReason = "包含关键词「" + req.Keyword + "」"
				}
			}
		default:
			matched = s.anyLiteralHit(corpus[c.ID], expandedKeywords)
			if matched {
				relevance = 1
				matchReason = "包含关键词「" + req.Keyword + "」"
			}
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
			MatchReason:   matchReason,
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
	"相关性", "命中依据", "实时综合得分", "成果分", "决策影响分", "社会贡献分",
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
			item.Relevance, item.MatchReason, item.RealtimeScore, item.AchievementScore, item.InfluenceScore, item.SocialScore,
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
