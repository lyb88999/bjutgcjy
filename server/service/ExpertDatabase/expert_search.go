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

// keywordTerm 检索词拆分后的一个独立词，各自展开同义词、各自算好语义向量。多词检索时每个词
// 独立判断命中，互不影响，最后再汇总（见 rankedCandidates 里的合并逻辑）
type keywordTerm struct {
	term     string
	expanded []string
	vectors  [][]float32
}

// termMatch 单个检索词对某位专家的命中判定，语义向量匹配优先于字面子串匹配（子串匹配要求
// 词整体原样出现在语料里，"人工智能"搜不出"具身智能"这类相关但不同字面的内容）。
// embedding 服务调不通、或这个专家还没生成向量条目（刚发布/刚导入，没跑过
// cmd/recompute-embeddings）时整体退回字面子串匹配，不能让语义检索的故障/数据缺失拖垮
// 基本检索能力
func (s *ExpertSearchService) termMatch(corpus string, embItems []ExpertDatabase.ExpertSearchEmbeddingItem, kt keywordTerm, useSemanticMatch bool) (relevance float64, matched bool, reason string) {
	literalHit := s.anyLiteralHit(corpus, kt.expanded)
	if literalHit {
		reason = "包含关键词「" + kt.term + "」"
	}
	if !useSemanticMatch || len(embItems) == 0 {
		if literalHit {
			relevance = 1
			matched = true
		}
		return
	}
	bestSim := -1.0
	for _, it := range embItems {
		var v []float32
		if json.Unmarshal([]byte(it.Vector), &v) != nil {
			continue
		}
		for _, kwVec := range kt.vectors {
			if sim := cosineSim(kwVec, v); sim > bestSim {
				bestSim = sim
				if !literalHit {
					reason = it.ItemLabel
				}
			}
		}
	}
	switch {
	case literalHit:
		relevance = literalHitRelevanceFloor
		if bestSim > relevance {
			relevance = bestSim
		}
		matched = true
	case bestSim >= minSemanticRelevance:
		relevance = bestSim
		matched = true
	}
	return
}

// rankedCandidates 按筛选条件圈定已发布候选集，并按
// realtimeScore = w1*achievement_score*relevance + w2*influence_score*relevance + w3*title_score
//                + w4*social_score + w5*relevance
// 算好相关性/实时得分、排好序，返回不分页的完整结果——Search()（列表分页）和 ExportSearchResults()
// （导出全部）共用同一份排序逻辑，避免排序算法在两个地方各写一遍、后续改权重容易漏改一处。
// w5*relevance 这一项是单独给"关键词匹配程度本身"的得分：achievement/influence 两项是
// 乘relevance，对成果分/决策影响分本来就是 0 的专家（比如刚入库、还没攒够成果的专家）不管
// 关键词匹配得多好，这两项乘出来还是 0，综合得分对这批人形同摆设、看起来"怎么搜都不变"；
// 加一项不依赖其他分项、只看 relevance 本身的得分，保证关键词匹配程度总能实实在在体现到
// 综合得分里
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

	// 检索词按常见分隔符拆成多个独立词（比如"人工智能 京津冀"拆成两个），下面按词分别判断命中——
	// 命中任意一个词就收录（OR），相关性按"命中词的平均得分"算，命中的词越多平均分越高，天然做到
	// "命中越多排名越靠前"，不要求整句原样命中（几乎不可能，之前多关键词检索基本搜不出东西的
	// 根源就在这——旧实现把整个输入串当一个词去比对）
	terms := splitKeyword(req.Keyword)
	corpus := map[uint]string{}
	itemsByExpert := map[uint][]ExpertDatabase.ExpertSearchEmbeddingItem{}
	keywordTerms := make([]keywordTerm, len(terms))
	useSemanticMatch := false
	if len(terms) > 0 && len(candidates) > 0 {
		if corpus, err = s.buildSearchCorpusMap(candidates); err != nil {
			return nil, err
		}

		// 先按同义词字典把每个词各自展开成一组词（没配同义词的话就还是只有原词一个）。"碳中和"和
		// "低碳"这类近义关系 embedding 模型自己判断不出来（实测过），这一步是确定性的人工兜底，
		// 不依赖模型判不判断得准——见 expert_search_synonym.go 顶部注释
		synonymGroups := s.searchSynonymGroups()
		allTextsSet := make(map[string]bool)
		for i, t := range terms {
			expanded := expandKeywordWithSynonyms(t, synonymGroups)
			keywordTerms[i] = keywordTerm{term: t, expanded: expanded}
			for _, e := range expanded {
				allTextsSet[e] = true
			}
		}

		// 所有词（含同义词展开）合并去重后只调一次 embedding 服务，再按词把向量分发回去——不管
		// 检索词有几个，实时检索这一步始终只发一次请求。embedding 服务调不通时（网络问题/服务
		// 没起来）整体退回子串匹配，不能让语义检索的故障拖垮基本检索能力
		allTexts := make([]string, 0, len(allTextsSet))
		for txt := range allTextsSet {
			allTexts = append(allTexts, txt)
		}
		keywordEmbedTimeout := time.Duration(global.GVA_CONFIG.ExpertEmbedding.TimeoutSec) * time.Second
		if vecs, embedErr := embedTexts(allTexts, keywordEmbedTimeout); embedErr == nil && len(vecs) == len(allTexts) {
			vecByText := make(map[string][]float32, len(allTexts))
			for i, txt := range allTexts {
				vecByText[txt] = vecs[i]
			}
			for i := range keywordTerms {
				vs := make([][]float32, 0, len(keywordTerms[i].expanded))
				for _, e := range keywordTerms[i].expanded {
					vs = append(vs, vecByText[e])
				}
				keywordTerms[i].vectors = vs
			}

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

	items = make([]ExpertDatabaseRes.ExpertSearchItem, 0, len(candidates))
	for _, c := range candidates {
		var relevance float64
		var matchReason string
		matchedAny := len(terms) == 0
		if !matchedAny {
			embItems := itemsByExpert[c.ID]
			bestTermRelevance := -1.0
			matchedCount := 0
			var relevanceSum float64
			for _, kt := range keywordTerms {
				termRel, termMatched, reason := s.termMatch(corpus[c.ID], embItems, kt, useSemanticMatch)
				relevanceSum += termRel
				if termMatched {
					matchedAny = true
					matchedCount++
				}
				if termRel > bestTermRelevance {
					bestTermRelevance = termRel
					matchReason = reason
				}
			}
			relevance = relevanceSum / float64(len(terms))
			if len(terms) > 1 && matchedAny {
				matchReason = fmt.Sprintf("%s（命中 %d/%d 个关键词）", matchReason, matchedCount, len(terms))
			}
		} else {
			relevance = 1
		}
		// 填了关键词却一个词都没命中的候选人直接跳过，不进结果集——不然职称权重、社会贡献分这些
		// 不受相关性影响的分项会把一堆跟关键词毫不沾边的人顶到排名前面，"检索"就退化成了"不管搜
		// 什么都是把全库按职称排一遍"。没填关键词时 matchedAny 恒为 true，这里的判断天然不影响
		// 默认浏览场景。
		if !matchedAny {
			continue
		}
		titleScore, ok := titleWeights[c.TechTitle]
		if !ok {
			titleScore = 1
		}
		realtimeScore := rankingWeights["achievement"]*c.AchievementScore*relevance +
			rankingWeights["influence"]*c.InfluenceScore*relevance +
			rankingWeights["title"]*titleScore +
			rankingWeights["social"]*c.SocialScore +
			rankingWeights["keyword"]*relevance
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
