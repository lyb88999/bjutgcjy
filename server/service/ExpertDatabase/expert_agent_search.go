package ExpertDatabase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase"
	ExpertDatabaseReq "github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase/request"
	ExpertDatabaseRes "github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
)

const agentMaxProfiles = 30000

// 复用现有检索的结构筛选、分项权重、同义词及termMatch；独立传递context，避免修改全局DB。
// 内存/CPU都有候选上限，缓存未就绪时降级为字面匹配，不在交互请求中重复全表加载向量。
func searchExpertAgentCandidates(ctx context.Context, plan ExpertDatabaseRes.ExpertAgentPlan) ([]ExpertDatabaseRes.ExpertAgentCandidate, int64, []string, error) {
	warnings := []string{}
	if global.GVA_DB == nil {
		return nil, 0, warnings, fmt.Errorf("专家库尚未连接")
	}
	request := ExpertDatabaseReq.ExpertSearchReq{Name: plan.Name, UnitName: plan.UnitName, DisciplineL1: plan.DisciplineL1, TitleTier: plan.TitleTier}
	query := applyStructuralFilters(request).WithContext(ctx)
	if len(plan.ExcludedExpertIDs) > 0 {
		query = query.Where("id NOT IN ?", plan.ExcludedExpertIDs)
	}
	var profiles []ExpertDatabase.ExpertProfile
	// 只读取用于检索和解释的字段，手机号、邮箱、住址等均不进入Agent候选或模型上下文。
	columns := "id,name,unit_name,department,tech_title,org_id,discipline_l1,discipline_l2,cross_discipline,research_directions,research_keywords,policy_fields,research_objects,focus_topics,region_expertise,method_expertise,profile_link,achievement_score,influence_score,title_score,social_score,composite_score"
	if err := query.Select(columns).Order("composite_score DESC, id DESC").Limit(agentMaxProfiles + 1).Find(&profiles).Error; err != nil {
		return nil, 0, warnings, agentSearchFailure(err)
	}
	if len(profiles) > agentMaxProfiles {
		profiles = profiles[:agentMaxProfiles]
		warnings = append(warnings, "当前结构筛选范围较大，本轮只在基础分最高的30000位专家中检索；可补充单位或学科缩小范围。")
	}
	if len(profiles) == 0 {
		return []ExpertDatabaseRes.ExpertAgentCandidate{}, 0, warnings, nil
	}
	ids := make([]uint, 0, len(profiles))
	for _, profile := range profiles {
		ids = append(ids, profile.ID)
	}
	ranking, titleWeights, synonyms, err := agentSearchDictionaries(ctx)
	if err != nil {
		return nil, 0, warnings, agentSearchFailure(err)
	}
	corpus := make(map[uint]string, len(profiles))
	for _, p := range profiles {
		corpus[p.ID] = strings.Join([]string{p.Name, p.UnitName, p.TechTitle, buildProfileChunkText(p)}, " ")
	}
	achievements := map[uint][]ExpertDatabase.ExpertAchievement{}
	tagValues := map[uint][]string{}
	if len(plan.Keywords) > 0 {
		var rows []ExpertDatabase.ExpertAchievement
		if err := global.GVA_DB.WithContext(ctx).Select("id,expert_id,title,keywords").Where("expert_id IN ?", ids).Order("id DESC").Limit(250001).Find(&rows).Error; err != nil {
			return nil, 0, warnings, agentSearchFailure(err)
		}
		if len(rows) > 250000 {
			rows = rows[:250000]
			warnings = append(warnings, "成果记录超过单轮读取上限，本轮仅使用最近录入的250000条成果。")
		}
		for _, a := range rows {
			achievements[a.ExpertId] = append(achievements[a.ExpertId], a)
			corpus[a.ExpertId] += " " + a.Title + " " + a.Keywords
		}
		var tags []struct {
			ExpertID uint
			TagValue string
		}
		if err := global.GVA_DB.WithContext(ctx).Table("expert_tag_relation r").Select("r.expert_id, t.tag_value").Joins("JOIN expert_tag t ON t.id = r.tag_id AND t.deleted_at IS NULL").Where("r.expert_id IN ? AND r.deleted_at IS NULL", ids).Order("r.id DESC").Limit(300001).Scan(&tags).Error; err != nil {
			return nil, 0, warnings, agentSearchFailure(err)
		}
		if len(tags) > 300000 {
			tags = tags[:300000]
			warnings = append(warnings, "标签关联超过单轮读取上限，本轮仅使用最近录入的300000条标签关联。")
		}
		for _, tag := range tags {
			corpus[tag.ExpertID] += " " + tag.TagValue
			tagValues[tag.ExpertID] = append(tagValues[tag.ExpertID], tag.TagValue)
		}
	}
	terms := make([]keywordTerm, 0, len(plan.Keywords))
	allTexts := []string{}
	seen := map[string]bool{}
	for _, keyword := range plan.Keywords {
		expanded := expandKeywordWithSynonyms(keyword, synonyms)
		if len(expanded) > 16 {
			expanded = expanded[:16]
		}
		terms = append(terms, keywordTerm{term: keyword, expanded: expanded})
		for _, word := range expanded {
			if !seen[word] {
				seen[word] = true
				allTexts = append(allTexts, word)
			}
		}
	}
	useSemantic := false
	var embeddings map[uint][]parsedEmbeddingItem
	if len(terms) > 0 {
		embeddings, useSemantic = agentCachedEmbeddings(ctx, ids)
		if useSemantic {
			vectors, embedErr := agentEmbedTexts(ctx, allTexts)
			if embedErr != nil {
				useSemantic = false
			} else {
				byText := map[string][]float32{}
				for i, word := range allTexts {
					byText[word] = vectors[i]
				}
				for i := range terms {
					for _, word := range terms[i].expanded {
						terms[i].vectors = append(terms[i].vectors, byText[word])
					}
				}
			}
		}
		if err := ctx.Err(); err != nil {
			return nil, 0, warnings, err
		}
		if !useSemantic {
			warnings = append(warnings, "语义向量服务或缓存当前不可用，本轮使用关键词与同义词匹配。")
		}
	}
	// 同一org_id可能有不同原始单位写法；有机构表时用统一名称执行每单位人数限制。
	orgNames := map[uint]string{}
	orgIDs, seenOrg := []uint{}, map[uint]bool{}
	for _, p := range profiles {
		if p.OrgId != nil && !seenOrg[*p.OrgId] {
			seenOrg[*p.OrgId] = true
			orgIDs = append(orgIDs, *p.OrgId)
		}
	}
	if len(orgIDs) > 0 {
		var organizations []system.SysOrganization
		if err := global.GVA_DB.WithContext(ctx).Select("id,name").Where("id IN ?", orgIDs).Find(&organizations).Error; err != nil {
			return nil, 0, warnings, agentSearchFailure(err)
		}
		for _, org := range organizations {
			if org.Name != "" {
				orgNames[org.ID] = org.Name
			}
		}
	}
	search := &ExpertSearchService{}
	candidates := make([]ExpertDatabaseRes.ExpertAgentCandidate, 0, len(profiles))
	for i, p := range profiles {
		if i%32 == 0 {
			if err := ctx.Err(); err != nil {
				return nil, 0, warnings, err
			}
		}
		relevance, matchReason := 1.0, "符合已选结构化条件，按库内基础分排序"
		matchedTerms := []keywordTerm{}
		semanticMatches := []expertAgentSemanticMatch{}
		if len(terms) > 0 {
			sum, best := 0.0, -1.0
			for _, term := range terms {
				rel, matched, reason, source := search.termMatchWithSource(corpus[p.ID], embeddings[p.ID], term, useSemantic)
				if !matched {
					continue
				}
				matchedTerms = append(matchedTerms, term)
				rel = math.Min(1, math.Max(0, rel))
				if source != nil {
					semanticMatches = append(semanticMatches, expertAgentSemanticMatch{Source: *source, Relevance: rel})
				}
				sum += rel
				if rel > best {
					best = rel
					matchReason = reason
				}
			}
			if len(matchedTerms) == 0 {
				continue
			}
			relevance = sum / float64(len(terms))
			if len(terms) > 1 {
				matchReason += fmt.Sprintf("（命中%d/%d个主题词）", len(matchedTerms), len(terms))
			}
		}
		score := p.CompositeScore
		if len(terms) > 0 {
			score = relevance * (ranking["achievement"]*p.AchievementScore + ranking["influence"]*p.InfluenceScore + ranking["title"]*matchTitleWeight(p.TechTitle, titleWeights) + ranking["social"]*p.SocialScore + ranking["keyword"])
		}
		if math.IsNaN(score) || math.IsInf(score, 0) {
			score = 0
		}
		unitName := p.UnitName
		orgID := uint(0)
		if p.OrgId != nil {
			orgID = *p.OrgId
		}
		if p.OrgId != nil && orgNames[*p.OrgId] != "" {
			unitName = orgNames[*p.OrgId]
		}
		link := agentSafeURL(p.ProfileLink)
		evidenceProfile := p
		evidenceProfile.UnitName = unitName
		evidence := buildExpertAgentEvidence(evidenceProfile, achievements[p.ID], tagValues[p.ID], matchedTerms, semanticMatches)
		reason := matchReason + "；依据已发布档案及库内记录推荐。"
		candidates = append(candidates, ExpertDatabaseRes.ExpertAgentCandidate{OrgID: orgID, TopicRelevance: relevance, ID: p.ID, Name: p.Name, UnitName: unitName, Department: p.Department, TechTitle: p.TechTitle, DisciplineL1: p.DisciplineL1, ResearchDirections: p.ResearchDirections, ProfileLink: link, MatchReason: matchReason, Score: math.Round(score*100) / 100, Evidence: evidence, Reason: reason})
	}
	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].Score != candidates[j].Score {
			return candidates[i].Score > candidates[j].Score
		}
		return candidates[i].ID < candidates[j].ID
	})
	total := int64(len(candidates))
	if len(candidates) > agentMaxCandidates {
		candidates = candidates[:agentMaxCandidates]
	}
	return candidates, total, warnings, nil
}

type expertAgentSemanticMatch struct {
	Source    parsedEmbeddingItem
	Relevance float64
}

// buildExpertAgentEvidence 优先放入真正参与本轮语义命中的原始记录，再补字面匹配及档案背景。
// 最多展示四条；先按相关性排列语义来源，避免最佳来源被较新的字面成果或背景信息挤掉。
func buildExpertAgentEvidence(profile ExpertDatabase.ExpertProfile, achievements []ExpertDatabase.ExpertAchievement, tags []string, matchedTerms []keywordTerm, semanticMatches []expertAgentSemanticMatch) []ExpertDatabaseRes.ExpertAgentEvidence {
	const maxEvidence = 4
	evidence := make([]ExpertDatabaseRes.ExpertAgentEvidence, 0, maxEvidence)
	seen := map[string]bool{}
	add := func(entry ExpertDatabaseRes.ExpertAgentEvidence) {
		if len(evidence) >= maxEvidence || strings.TrimSpace(entry.Title) == "" {
			return
		}
		key := fmt.Sprintf("%s:%d", entry.Kind, entry.SourceID)
		if entry.SourceID == 0 {
			key += ":" + entry.Title
		}
		if !seen[key] {
			seen[key] = true
			evidence = append(evidence, entry)
		}
	}
	byID := make(map[uint]ExpertDatabase.ExpertAchievement, len(achievements))
	for _, achievement := range achievements {
		if achievement.ID > 0 && achievement.ExpertId == profile.ID {
			byID[achievement.ID] = achievement
		}
	}
	link := agentSafeURL(profile.ProfileLink)
	orderedMatches := append([]expertAgentSemanticMatch{}, semanticMatches...)
	sort.SliceStable(orderedMatches, func(i, j int) bool {
		return orderedMatches[i].Relevance > orderedMatches[j].Relevance
	})
	for _, match := range orderedMatches {
		switch match.Source.ItemType {
		case embeddingItemTypeAchievement:
			// 必须能在本专家本轮读取的真实成果中按ID找到；不以相同或截断标题猜测来源。
			if achievement, ok := byID[match.Source.ItemRefID]; ok {
				add(ExpertDatabaseRes.ExpertAgentEvidence{Kind: "achievement", Title: truncateLabel(achievement.Title, 300), SourceID: achievement.ID})
			}
		case embeddingItemTypeProfile:
			add(ExpertDatabaseRes.ExpertAgentEvidence{Kind: "profile", Title: truncateLabel(buildProfileChunkText(profile), 300), SourceID: profile.ID, URL: link})
		case embeddingItemTypeTag:
			// 标签向量是本专家全部标签的合并条目，refID为0，不伪造某个标签记录ID。
			add(ExpertDatabaseRes.ExpertAgentEvidence{Kind: "tag", Title: truncateLabel(strings.Join(tags, "、"), 300), URL: link})
		}
	}
	if profile.ResearchDirections != "" {
		add(ExpertDatabaseRes.ExpertAgentEvidence{Kind: "researchDirection", Title: truncateLabel(profile.ResearchDirections, 300), SourceID: profile.ID, URL: link})
	}
	search := &ExpertSearchService{}
	for _, achievement := range achievements {
		if len(evidence) >= maxEvidence {
			break
		}
		if achievement.ID == 0 || achievement.ExpertId != profile.ID {
			continue
		}
		for _, term := range matchedTerms {
			if search.anyLiteralHit(achievement.Title+" "+achievement.Keywords, term.expanded) {
				add(ExpertDatabaseRes.ExpertAgentEvidence{Kind: "achievement", Title: truncateLabel(achievement.Title, 300), SourceID: achievement.ID})
				break
			}
		}
	}
	if len(evidence) == 0 {
		add(ExpertDatabaseRes.ExpertAgentEvidence{Kind: "profile", Title: truncateLabel(buildProfileChunkText(profile), 300), SourceID: profile.ID, URL: link})
	}
	if len(evidence) == 0 {
		add(ExpertDatabaseRes.ExpertAgentEvidence{Kind: "profile", Title: strings.TrimSpace(profile.Name + " " + profile.UnitName + " " + profile.TechTitle), SourceID: profile.ID, URL: link})
	}
	return evidence
}

func agentSearchDictionaries(ctx context.Context) (map[string]float64, map[string]float64, [][]string, error) {
	ranking, title := map[string]float64{}, map[string]float64{}
	for key, value := range defaultRankingWeights {
		ranking[key] = value
	}
	for key, value := range defaultTitleLevelWeights {
		title[key] = value
	}
	groups := [][]string{}
	var dictionaries []system.SysDictionary
	if err := global.GVA_DB.WithContext(ctx).Where("type IN ?", []string{DictTypeRankingWeight, DictTypeTitleLevel, DictTypeSearchSynonym}).Find(&dictionaries).Error; err != nil {
		return nil, nil, nil, err
	}
	types, ids := map[uint]string{}, []uint{}
	for _, dictionary := range dictionaries {
		types[dictionary.ID] = dictionary.Type
		ids = append(ids, dictionary.ID)
	}
	if len(ids) == 0 {
		return ranking, title, groups, nil
	}
	var details []system.SysDictionaryDetail
	if err := global.GVA_DB.WithContext(ctx).Where("sys_dictionary_id IN ?", ids).Find(&details).Error; err != nil {
		return nil, nil, nil, err
	}
	for _, detail := range details {
		kind := types[uint(detail.SysDictionaryID)]
		if kind == DictTypeSearchSynonym {
			if detail.Status == nil || *detail.Status {
				if terms := splitKeyword(detail.Extend); len(terms) > 1 {
					groups = append(groups, terms)
				}
			}
			continue
		}
		weight, err := strconv.ParseFloat(detail.Extend, 64)
		if err != nil {
			weight = float64(detail.Value)
		}
		if math.IsNaN(weight) || math.IsInf(weight, 0) {
			continue
		}
		if kind == DictTypeRankingWeight {
			ranking[detail.Label] = weight
		} else {
			title[detail.Label] = weight
		}
	}
	return ranking, title, groups, nil
}

func agentCachedEmbeddings(ctx context.Context, ids []uint) (map[uint][]parsedEmbeddingItem, bool) {
	if global.GVA_CONFIG.ExpertEmbedding.Url == "" {
		return nil, false
	}
	var maxUpdated *time.Time
	if err := global.GVA_DB.WithContext(ctx).Model(&ExpertDatabase.ExpertSearchEmbeddingItem{}).Select("MAX(updated_at)").Scan(&maxUpdated).Error; err != nil || maxUpdated == nil {
		return nil, false
	}
	embeddingCache.mu.RLock()
	defer embeddingCache.mu.RUnlock()
	if len(embeddingCache.byExpert) == 0 || !embeddingCache.freshness.Equal(*maxUpdated) {
		return nil, false
	}
	return collectByIDs(embeddingCache.byExpert, ids), true
}

func agentEmbedTexts(ctx context.Context, texts []string) ([][]float32, error) {
	config := global.GVA_CONFIG.ExpertEmbedding
	timeout := time.Duration(config.TimeoutSec) * time.Second
	if timeout <= 0 || timeout > 10*time.Second {
		timeout = 5 * time.Second
	}
	embedCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	body, err := json.Marshal(embedRequestBody{Texts: texts})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(embedCtx, http.MethodPost, strings.TrimRight(config.Url, "/")+"/embed", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: timeout, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("语义服务不可用")
	}
	raw, err := io.ReadAll(io.LimitReader(response.Body, 4*1024*1024))
	if err != nil {
		return nil, err
	}
	var decoded embedResponseBody
	if json.Unmarshal(raw, &decoded) != nil || len(decoded.Embeddings) != len(texts) {
		return nil, fmt.Errorf("语义向量格式无效")
	}
	if err := validateAgentEmbeddingResponse(decoded, len(texts)); err != nil {
		return nil, err
	}
	return decoded.Embeddings, nil
}

func agentSafeURL(raw string) string {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil || u.Host == "" || u.User != nil || (u.Scheme != "https" && u.Scheme != "http") {
		return ""
	}
	return u.String()
}
