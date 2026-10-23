package ExpertDatabase

import (
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

// buildSearchCorpusMap 一次性预取一批候选专家的检索语料。语料来源四块：① 姓名/所在单位/专业
// 技术职称——检索条"全部字段"这个选项名字上就该覆盖这几个最基础的身份字段，之前漏了，导致
// "全部字段"模式下按姓名/单位搜人反而搜不到（除非名字碰巧出现在研究方向文本里），跟"姓名"/
// "所在单位"两个独立检索范围选项（前端 el-select，走 name/unitName 精确+模糊结构化过滤）
// 是两条互补的路——那两个是精确的结构化过滤，这里是"全部字段"关键词检索也能顺带命中同一批
// 信息，两者不冲突；② 专家档案自己填的研究方向/关键词/关注议题等字段；③ 成果标题/关键词；
// ④ 关联标签。③④整批只发两条 SQL——替代原先逐个专家各查两次的写法，候选集几百上千条时检索
// 耗时不再随人数线性膨胀。标签 JOIN 只认未解除的关联（relation 是软删除，解除过的标签不该再
// 计入语料）
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
			p.Name, p.UnitName, p.TechTitle,
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

// tieBreakKey 给排序平局打散用的确定性伪随机数：同一个 id 永远算出同一个值（同一次查询内
// 多次排序、翻页顺序不会跳），但跟 id 本身的大小无关联，不会让平局回到"按 id/导入批次聚堆"
func tieBreakKey(id uint) uint32 {
	x := uint32(id)
	x = ((x >> 16) ^ x) * 0x45d9f3b
	x = ((x >> 16) ^ x) * 0x45d9f3b
	x = (x >> 16) ^ x
	return x
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
func (s *ExpertSearchService) termMatch(corpus string, embItems []parsedEmbeddingItem, kt keywordTerm, useSemanticMatch bool) (relevance float64, matched bool, reason string) {
	relevance, matched, reason, _ = s.termMatchWithSource(corpus, embItems, kt, useSemanticMatch)
	return
}

// termMatchWithSource 额外保留纯语义命中的原始条目身份。展示标签可能被截断或重名，不能用它反查成果ID。
// source 是条目元数据的独立拷贝，不带向量切片；字面命中或未达到语义阈值时不返回语义来源。
func (s *ExpertSearchService) termMatchWithSource(corpus string, embItems []parsedEmbeddingItem, kt keywordTerm, useSemanticMatch bool) (relevance float64, matched bool, reason string, source *parsedEmbeddingItem) {
	literalHit := s.anyLiteralHit(corpus, kt.expanded)
	if literalHit {
		if strings.Contains(corpus, kt.term) {
			reason = "包含关键词「" + kt.term + "」"
		} else {
			for _, alias := range kt.expanded {
				if strings.Contains(corpus, alias) {
					reason = "命中关键词「" + kt.term + "」的同义词「" + alias + "」"
					break
				}
			}
		}
	}
	if !useSemanticMatch || len(embItems) == 0 {
		if literalHit {
			relevance = 1
			matched = true
		}
		return
	}
	bestSim := -1.0
	var bestSource parsedEmbeddingItem
	for _, it := range embItems {
		for _, kwVec := range kt.vectors {
			if sim := cosineSim(kwVec, it.Vector); sim > bestSim {
				bestSim = sim
				if !literalHit {
					reason = it.ItemLabel
					bestSource = parsedEmbeddingItem{ItemType: it.ItemType, ItemRefID: it.ItemRefID, ItemLabel: it.ItemLabel}
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
		source = &bestSource
	}
	return
}

// rankedCandidates 按筛选条件圈定已发布候选集，并按
// realtimeScore = relevance * (w1*achievement_score + w2*influence_score + w3*title_score + w4*social_score + w5)
// 算好相关性/实时得分、排好序，返回不分页的完整结果——Search()（列表分页）和 ExportSearchResults()
// （导出全部）共用同一份排序逻辑，避免排序算法在两个地方各写一遍、后续改权重容易漏改一处。
//
// relevance 是乘在整个括号上的，不是只乘成果分/决策影响分那两项——早期版本只乘了这两项，
// 职称权重、社会贡献分完全不受关键词匹配程度影响，导致：① 这两项本来就是 0 的专家（还没攒够
// 成果），综合得分对这批人形同摆设、看起来"怎么搜都不变"；② 更普遍的问题是，就算成果分不是 0，
// 职称/社会贡献这两块通常量级不小（职称权重 1~5，社会贡献分不封顶），relevance 只能通过很小
// 一部分分项体现出来，实测下来"关键词匹配得好不好"对最终排名的影响弱到几乎看不出来。改成整体
// 相乘之后，relevance 从 1.0 掉到 0.6（语义匹配的下限）时，不管这个专家职称多高、社会贡献分
// 多厚，综合得分都会跟着打对应的折扣，"匹配程度"才真正统筹进了这一个分数里。w5 是关键词匹配
// 本身的固定加分（默认 10，权重 key 是 keyword），同样被 relevance 乘着，纯粹没匹配上关键词
// 的候选人在筛选阶段已经被跳过，不会拿到这一项
func (s *ExpertSearchService) rankedCandidates(req ExpertDatabaseReq.ExpertSearchReq) (items []ExpertDatabaseRes.ExpertSearchItem, err error) {
	var candidates []ExpertDatabase.ExpertProfile
	if err = applyStructuralFilters(req).Find(&candidates).Error; err != nil {
		return
	}

	rankingWeights := expertScoreSvc.dictWeights(DictTypeRankingWeight, defaultRankingWeights)
	titleWeights := expertScoreSvc.dictWeights(DictTypeTitleLevel, defaultTitleLevelWeights)

	// 高级检索（多行且/或/不含，见 expert_search_advanced.go）优先于简单关键词：AdvancedRows
	// 非空且过滤掉空词后还剩至少一行时，terms 从这些行的原文取；否则落回简单关键词的空格分词。
	// 两条路径共享下面同一套语料构建/同义词展开/单次批量 embedding 调用，只在下方判定循环里
	// 按 useAdvanced 分叉——不会因为走了高级检索就多打好几次 embedding 服务请求
	var terms []string
	var advancedRows []ExpertDatabaseReq.AdvancedSearchRow
	useAdvanced := false
	if req.AdvancedRows != "" {
		rows, parseErr := parseAdvancedRows(req.AdvancedRows)
		if parseErr != nil {
			return nil, parseErr
		}
		if len(rows) > 0 {
			advancedRows = rows
			useAdvanced = true
			terms = make([]string, len(rows))
			for i, r := range rows {
				terms[i] = r.Term
			}
		}
	}
	if !useAdvanced {
		// 检索词按常见分隔符拆成多个独立词（比如"人工智能 京津冀"拆成两个），下面按词分别判断
		// 命中——命中任意一个词就收录（OR），相关性按"命中词的平均得分"算，命中的词越多平均分
		// 越高，天然做到"命中越多排名越靠前"，不要求整句原样命中（几乎不可能，之前多关键词检索
		// 基本搜不出东西的根源就在这——旧实现把整个输入串当一个词去比对）
		terms = splitKeyword(req.Keyword)
	}
	corpus := map[uint]string{}
	itemsByExpert := map[uint][]parsedEmbeddingItem{}
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
			if byExpert, embErr := s.embeddingItemsByExpert(ids); embErr == nil {
				itemsByExpert = byExpert
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
			if useAdvanced {
				// 高级检索：按连接词把每一行的判定结果从左到右折叠（见 combineTermResults），
				// 不是简单的"命中任意一个就算"。matchReason 只拼命中的/排除生效的那几行，
				// 不列进没起作用的判断，读起来接近自然语言
				var acc termResult
				var reasonParts []string
				for i, kt := range keywordTerms {
					termRel, termMatched, reason := s.termMatch(corpus[c.ID], embItems, kt, useSemanticMatch)
					row := termResult{matched: termMatched, relevance: termRel}
					if i == 0 {
						acc = row
						if termMatched {
							reasonParts = append(reasonParts, reason)
						}
						continue
					}
					connector := advancedRows[i].Connector
					acc = combineTermResults(acc, connector, row)
					switch connector {
					case "not":
						if !termMatched {
							reasonParts = append(reasonParts, "不含「"+advancedRows[i].Term+"」")
						}
					case "and":
						if termMatched {
							reasonParts = append(reasonParts, "且"+reason)
						}
					default:
						if termMatched {
							reasonParts = append(reasonParts, "或"+reason)
						}
					}
				}
				matchedAny = acc.matched
				relevance = acc.relevance
				matchReason = strings.Join(reasonParts, "")
			} else {
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
		titleScore := matchTitleWeight(c.TechTitle, titleWeights)
		realtimeScore := relevance * (rankingWeights["achievement"]*c.AchievementScore +
			rankingWeights["influence"]*c.InfluenceScore +
			rankingWeights["title"]*titleScore +
			rankingWeights["social"]*c.SocialScore)
		// keyword 这一项是关键词匹配本身的固定加分，只在真的填了关键词时才加——默认浏览场景
		// len(terms)==0 时 relevance 恒为 1，如果不加这个判断，浏览列表里显示的"综合实力"会比
		// 专家详情页缓存的 compositeScore 平白多出这一整项，两个页面的数字对不上
		if len(terms) > 0 {
			realtimeScore += relevance * rankingWeights["keyword"]
		}
		items = append(items, ExpertDatabaseRes.ExpertSearchItem{
			ExpertProfile: c,
			Relevance:     relevance,
			RealtimeScore: realtimeScore,
			MatchReason:   matchReason,
		})
	}

	// 默认浏览（没填关键词）时排序方式可以按 SortBy 换成"最近更新"或"姓名"——这两种场景下
	// 综合实力分不是用户真正想看的顺序。带关键词检索必须按相关性/综合得分排序，忽略 SortBy，
	// 不然搜索结果的顺序跟检索意图脱节
	if len(terms) == 0 {
		switch req.SortBy {
		case "updatedAt":
			sort.Slice(items, func(i, j int) bool {
				return items[i].UpdatedAt.After(items[j].UpdatedAt)
			})
			return items, nil
		case "name":
			sort.Slice(items, func(i, j int) bool {
				return items[i].Name < items[j].Name
			})
			return items, nil
		}
	}
	// 批量导入的档案还没攒够成果/决策采纳记录时，综合得分基本只剩职称分这一项，会有大量专家分数
	// 完全相同（实测一批几千人导入后，69% 的人扎堆在同一个分数）。这种平局不该按 id 顺序排——
	// 那等于按"哪个学校哪批导入的"分组，同一学校的人会连成一大片，看着像是漏了别的学校。
	// 用哈希打散一下顺序：同一批查询里顺序仍然稳定（不是每次刷新都跳），只是不再跟着 id/学校聚堆。
	sort.Slice(items, func(i, j int) bool {
		if items[i].RealtimeScore != items[j].RealtimeScore {
			return items[i].RealtimeScore > items[j].RealtimeScore
		}
		return tieBreakKey(items[i].ID) < tieBreakKey(items[j].ID)
	})
	return items, nil
}

// browseModeQuery 没填关键词时的浏览路径：筛选和排序全部下推到数据库，配一条不带 LIMIT 的
// COUNT 拿总数——不像 rankedCandidates 那样把整个候选集读进 Go 内存再排序切片。relevance 恒为
// 1（浏览模式无所谓相关性），realtimeScore 直接用持久化的 composite_score：浏览模式下
// realtimeScore 的算法跟 RecomputeExpertScore 写 composite_score 时用的是同一个公式（relevance=1
// 时两者代数上相等），不用为了排序重新拉字典权重、重新跑一遍职称匹配——前提是得分缓存是新的，
// 如果字典权重改过还没跑 cmd/recompute，浏览模式的顺序会短暂跟专家详情页的缓存值一样"滞后"，
// 这跟这个项目里其它读 composite_score 缓存的地方（比如专家详情页）是同一种、已经存在的行为，
// 不是这次改造新引入的不一致。
//
// 排序上，同分的情况用 id 做稳定 tiebreak（ORDER BY composite_score DESC, id DESC），替代原来
// rankedCandidates 里 Go 算的哈希打散——这是一个小的、有意接受的行为变化：完全同分且没填关键词
// 时，结果相对顺序会比原来更贴近 id 顺序，不再是打散的。只在这一种情况下才会体现，影响面很小，
// 换来的是真正的数据库分页（不用整表读进内存）。
//
// forExport=true 时不分页，返回全部命中记录，供 ExportSearchResults 的浏览模式分支使用
func (s *ExpertSearchService) browseModeQuery(req ExpertDatabaseReq.ExpertSearchReq, forExport bool) (items []ExpertDatabaseRes.ExpertSearchItem, total int64, err error) {
	if err = applyStructuralFilters(req).Count(&total).Error; err != nil {
		return
	}

	q := applyStructuralFilters(req)
	switch req.SortBy {
	case "updatedAt":
		q = q.Order("updated_at DESC")
	case "name":
		q = q.Order("name ASC")
	default:
		q = q.Order("composite_score DESC, id DESC")
	}
	if !forExport {
		pageSize := req.PageSize
		if pageSize == 0 {
			pageSize = 10
		}
		page := req.Page
		if page == 0 {
			page = 1
		}
		q = q.Limit(pageSize).Offset((page - 1) * pageSize)
	}

	var profiles []ExpertDatabase.ExpertProfile
	if err = q.Find(&profiles).Error; err != nil {
		return
	}
	items = make([]ExpertDatabaseRes.ExpertSearchItem, 0, len(profiles))
	for _, p := range profiles {
		items = append(items, ExpertDatabaseRes.ExpertSearchItem{
			ExpertProfile: p,
			Relevance:     1,
			RealtimeScore: p.CompositeScore,
		})
	}
	return
}

// Search 专家综合推荐排序检索：有关键词走 rankedCandidates（相关性排序要对候选集整体做语义匹配，
// 排完序才能分页，没法在数据库层直接 LIMIT），没填关键词走 browseModeQuery（结构化筛选+排序+
// 分页全部下推到数据库，不用把整个候选集读进内存）
func (s *ExpertSearchService) Search(req ExpertDatabaseReq.ExpertSearchReq) (list []ExpertDatabaseRes.ExpertSearchItem, total int64, err error) {
	// AdvancedRows 非空时哪怕 Keyword 是空的，也得走 rankedCandidates 做关键词匹配——
	// 不能因为 Keyword 恰好没填就误判成浏览模式，跳过高级检索的过滤条件
	if req.Keyword == "" && req.AdvancedRows == "" {
		return s.browseModeQuery(req, false)
	}

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

// GetScoreColumnAvailability 全库已发布专家里，成果分/决策影响分/社会贡献分是不是至少有一个人
// 非零——花名册类批量导入的专家往往还没有成果/决策影响/学术兼职记录，某一项对全库所有人都是 0
// 时，界面上继续显示一整列 0 没有意义，前端拿这个结果决定隐藏对应的列/卡片
func (s *ExpertSearchService) GetScoreColumnAvailability() (ExpertDatabaseRes.ExpertScoreColumnAvailability, error) {
	var result ExpertDatabaseRes.ExpertScoreColumnAvailability
	var probe []uint
	if err := global.GVA_DB.Model(&ExpertDatabase.ExpertProfile{}).Where("status = ? AND achievement_score > 0", "published").Limit(1).Pluck("id", &probe).Error; err != nil {
		return result, err
	}
	result.HasAchievement = len(probe) > 0

	probe = probe[:0]
	if err := global.GVA_DB.Model(&ExpertDatabase.ExpertProfile{}).Where("status = ? AND influence_score > 0", "published").Limit(1).Pluck("id", &probe).Error; err != nil {
		return result, err
	}
	result.HasInfluence = len(probe) > 0

	probe = probe[:0]
	if err := global.GVA_DB.Model(&ExpertDatabase.ExpertProfile{}).Where("status = ? AND social_score > 0", "published").Limit(1).Pluck("id", &probe).Error; err != nil {
		return result, err
	}
	result.HasSocial = len(probe) > 0

	return result, nil
}

var searchExportHeaders = []string{
	"姓名", "所在单位", "专业技术职称", "一级学科", "研究关键词",
	"相关性", "命中依据", "实时综合得分", "成果分", "决策影响分", "职称分", "社会贡献分",
}

// ExportSearchResults 导出当前检索条件下命中的全部结果（按相关性/实时得分排好序，不分页）。
// 跟 Search 一样按有无关键词分支：浏览模式走 browseModeQuery 的不分页路径（导出顺序变成
// ORDER BY composite_score DESC, id DESC，是浏览模式排序调整的自然延伸，不是单独的行为差异）
func (s *ExpertSearchService) ExportSearchResults(req ExpertDatabaseReq.ExpertSearchReq) (*excelize.File, error) {
	var items []ExpertDatabaseRes.ExpertSearchItem
	var err error
	if req.Keyword == "" && req.AdvancedRows == "" {
		items, _, err = s.browseModeQuery(req, true)
	} else {
		items, err = s.rankedCandidates(req)
	}
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
			item.Relevance, item.MatchReason, item.RealtimeScore, item.AchievementScore, item.InfluenceScore, item.TitleScore, item.SocialScore,
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
