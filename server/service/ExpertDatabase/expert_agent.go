package ExpertDatabase

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	ExpertDatabaseReq "github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase/request"
	ExpertDatabaseRes "github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase/response"
)

var (
	ErrExpertAgentInvalidRequest  = errors.New("请检查遴选请求参数")
	ErrExpertAgentNotConfigured   = errors.New("专家遴选助手尚未配置模型服务")
	ErrExpertAgentRateLimited     = errors.New("模型服务当前请求较多，请稍后重试；未切换到其他模型")
	ErrExpertAgentUpstream        = errors.New("模型服务暂时不可用或未返回有效检索条件，请稍后重试")
	ErrExpertAgentSessionNotFound = errors.New("会话不存在或已过期，请新建会话")
	ErrExpertAgentSessionBusy     = errors.New("此会话正在处理上一条请求，请完成后再发送")
	ErrExpertAgentSessionLimit    = errors.New("会话数量已达上限，请继续现有会话或稍后再试")
)

const agentMaxCandidates = 200

type ExpertAgentService struct {
	client   *http.Client
	sessions *expertAgentSessionStore
	search   func(context.Context, ExpertDatabaseRes.ExpertAgentPlan) ([]ExpertDatabaseRes.ExpertAgentCandidate, int64, error)
}

func (s *ExpertAgentService) Status() ExpertDatabaseRes.ExpertAgentStatus {
	cfg, err := expertAgentConfigFromEnv()
	return ExpertDatabaseRes.ExpertAgentStatus{Configured: err == nil, Model: cfg.model}
}

func (s *ExpertAgentService) Chat(ctx context.Context, operatorID uint, req ExpertDatabaseReq.ExpertAgentChatReq) (result ExpertDatabaseRes.ExpertAgentChatResult, err error) {
	if err = ctx.Err(); err != nil {
		return result, err
	}
	req.Message = strings.TrimSpace(req.Message)
	if operatorID == 0 || req.Message == "" || utf8.RuneCountInString(req.Message) > 2000 || len(req.SessionID) > 64 || req.Count < 0 || req.Count > 30 || req.MaxPerUnit < 0 || req.MaxPerUnit > 30 {
		return result, ErrExpertAgentInvalidRequest
	}
	cfg, err := expertAgentConfigFromEnv()
	if err != nil {
		return result, err
	}
	store := s.sessions
	if store == nil {
		store = defaultExpertAgentSessions
	}
	id, session, created, err := store.acquire(operatorID, req.SessionID)
	if err != nil {
		return result, err
	}
	committed := false
	defer func() { store.release(id, created && !committed) }()
	// busy标记让同会话状态在本轮独占；只有完成全部检索、选择后才提交新状态。
	patch, err := s.callPlanner(ctx, cfg, req.Message, session)
	if err != nil {
		return result, err
	}
	plan, warnings, err := mergeExpertAgentPlan(session, patch, req)
	if err != nil {
		return result, err
	}
	var candidates []ExpertDatabaseRes.ExpertAgentCandidate
	var total int64
	browseAll := session.browseAll
	if mentioned, allowed := agentBrowseIntent(req.Message); mentioned {
		browseAll = allowed
	}
	needsCriteria := len(plan.Keywords) == 0 && plan.Name == "" && plan.UnitName == "" && plan.DisciplineL1 == "" && plan.TitleTier == "" && !browseAll
	if needsCriteria {
		candidates = []ExpertDatabaseRes.ExpertAgentCandidate{}
	} else if s.search != nil {
		candidates, total, err = s.search(ctx, plan)
	} else {
		var searchWarnings []string
		candidates, total, searchWarnings, err = searchExpertAgentCandidates(ctx, plan)
		warnings = append(warnings, searchWarnings...)
	}
	if err != nil {
		return result, err
	}
	if err = ctx.Err(); err != nil {
		return result, err
	}
	candidates = sanitizeExpertAgentCandidates(candidates, plan.ExcludedExpertIDs)
	if len(candidates) > agentMaxCandidates {
		candidates = candidates[:agentMaxCandidates]
	}
	selected := selectExpertAgentCandidates(candidates, plan)
	if total > int64(len(candidates)) {
		warnings = append(warnings, fmt.Sprintf("本轮共有%d位匹配候选，按检索得分取前%d位参与遴选。", total, len(candidates)))
	}
	if !needsCriteria && len(selected) < plan.Count {
		warnings = append(warnings, fmt.Sprintf("满足当前条件及单位人数限制的专家不足%d位，本轮提供%d位；未放宽条件或补造名单。", plan.Count, len(selected)))
	}
	stats := expertAgentStatistics(total, selected)
	result = ExpertDatabaseRes.ExpertAgentChatResult{
		SessionID: id, Plan: cloneExpertAgentPlan(plan), Candidates: candidates, Recommendations: selected,
		Statistics: stats, Warnings: warnings,
		Steps: []ExpertDatabaseRes.ExpertAgentStep{
			{Label: "理解需求", Detail: "将本轮要求与会话中已确认的检索条件合并。"},
			{Label: "检索专家库", Detail: fmt.Sprintf("仅检索已发布档案，找到%d位匹配候选。", total)},
			{Label: "执行遴选约束", Detail: fmt.Sprintf("执行人数、每单位人数、排除及去重条件，选出%d位专家。", len(selected))},
			{Label: "核对依据", Detail: "推荐理由来自库内档案、研究方向和成果记录，未由模型编写专家履历。"},
		},
	}
	result.Reply = fmt.Sprintf("已按当前条件推荐%d位专家，来自%d个单位、覆盖%d个已标注学科。可继续调整研究主题、人数、单位限制，或要求换一批。", len(selected), stats.UnitCount, stats.DisciplineCount)
	if len(selected) == 0 {
		result.Reply = "当前条件下没有可推荐的专家。可以放宽研究主题或单位范围，或减少排除条件后继续检索。"
	}
	if needsCriteria {
		result.Reply = "请补充要研究的主题、专家所在单位或职称要求，例如“推荐5位研究养老服务的专家，每个单位最多2位”。如需全库浏览，可以明确说“浏览全部专家”。"
		result.Steps = []ExpertDatabaseRes.ExpertAgentStep{{Label: "补充检索条件", Detail: "尚未指定主题、单位、姓名或职称，本轮未读取专家候选。"}}
	}
	store.mu.Lock()
	if err = ctx.Err(); err == nil {
		session.plan = cloneExpertAgentPlan(plan)
		session.browseAll = browseAll
		session.candidateIDs = make([]uint, 0, len(candidates))
		session.hints = make([]expertAgentCandidateHint, 0, len(candidates))
		for _, c := range candidates {
			session.candidateIDs = append(session.candidateIDs, c.ID)
			session.hints = append(session.hints, expertAgentCandidateHint{ID: c.ID, Name: c.Name, UnitName: c.UnitName})
		}
		session.selectedIDs = make([]uint, 0, len(selected))
		for _, c := range selected {
			session.selectedIDs = append(session.selectedIDs, c.ID)
		}
		session.updated = time.Now()
		committed = true
	}
	store.mu.Unlock()
	if err != nil {
		return ExpertDatabaseRes.ExpertAgentChatResult{}, err
	}
	return result, nil
}

func mergeExpertAgentPlan(previous *expertAgentSession, patch expertAgentPlanPatch, req ExpertDatabaseReq.ExpertAgentChatReq) (ExpertDatabaseRes.ExpertAgentPlan, []string, error) {
	plan := cloneExpertAgentPlan(previous.plan)
	warnings := []string{}
	if patch.Keywords != nil {
		plan.Keywords = append([]string{}, (*patch.Keywords)...)
	}
	for _, field := range []struct {
		source *string
		target *string
	}{
		{patch.Name, &plan.Name}, {patch.UnitName, &plan.UnitName}, {patch.DisciplineL1, &plan.DisciplineL1}, {patch.TitleTier, &plan.TitleTier},
	} {
		if field.source != nil {
			*field.target = strings.TrimSpace(*field.source)
		}
	}
	if patch.Count != nil {
		plan.Count = *patch.Count
	}
	if patch.MaxPerUnit != nil {
		plan.MaxPerUnit = *patch.MaxPerUnit
	}
	if patch.Diversify != nil {
		plan.Diversify = *patch.Diversify
	}
	if req.Count > 0 {
		plan.Count = req.Count
	}
	if req.MaxPerUnit > 0 {
		plan.MaxPerUnit = req.MaxPerUnit
	}
	if plan.Count < 1 || plan.Count > 30 || plan.MaxPerUnit < 0 || plan.MaxPerUnit > 30 || len(plan.Keywords) > 8 {
		return plan, warnings, ErrExpertAgentUpstream
	}
	for _, value := range []string{plan.Name, plan.UnitName, plan.DisciplineL1} {
		if utf8.RuneCountInString(value) > 120 {
			return plan, warnings, ErrExpertAgentUpstream
		}
	}
	switch plan.TitleTier {
	case "", "院士", "教授", "副教授", "讲师", "未分类":
	default:
		return plan, warnings, ErrExpertAgentUpstream
	}
	keywords := []string{}
	seenWords := map[string]bool{}
	for _, word := range plan.Keywords {
		word = strings.TrimSpace(word)
		if utf8.RuneCountInString(word) > 100 {
			return plan, warnings, ErrExpertAgentUpstream
		}
		if word != "" && !seenWords[word] {
			seenWords[word] = true
			keywords = append(keywords, word)
		}
	}
	plan.Keywords = keywords
	allowed := map[uint]bool{}
	for _, id := range previous.candidateIDs {
		allowed[id] = true
	}
	for _, id := range previous.plan.ExcludedExpertIDs {
		allowed[id] = true
	}
	if patch.ExcludedExpertIDs != nil && len(*patch.ExcludedExpertIDs) > 200 {
		return plan, warnings, ErrExpertAgentUpstream
	}
	batchMentioned, batchAllowed := agentAnotherBatchIntent(req.Message)
	changeBatch := patch.ExcludePrevious != nil && *patch.ExcludePrevious
	if batchMentioned {
		changeBatch = batchAllowed
	}
	if agentClearExclusionsIntent(req.Message) {
		// 明确取消排除比模型残留的换批标记优先，不能刚清空又把上轮名单加回来。
		plan.ExcludedExpertIDs = []uint{}
	} else if changeBatch || batchMentioned {
		// “换一批”针对已推荐名单，不是全部候选。模型偶尔会把整份候选ID都塞进排除数组，
		// 因此历史排除与上轮推荐由服务决定；额外排除须由用户本轮点名或明确ID支持。
		plan.ExcludedExpertIDs = append([]uint{}, previous.plan.ExcludedExpertIDs...)
		if changeBatch {
			plan.ExcludedExpertIDs = append(plan.ExcludedExpertIDs, previous.selectedIDs...)
		}
		justified := map[uint]bool{}
		for _, id := range plan.ExcludedExpertIDs {
			justified[id] = true
		}
		nameCounts := map[string]int{}
		for _, hint := range previous.hints {
			nameCounts[strings.TrimSpace(hint.Name)]++
		}
		for _, hint := range previous.hints {
			name := strings.TrimSpace(hint.Name)
			if nameCounts[name] != 1 {
				name = ""
			} // 同名候选须使用明确ID，避免一并误排。
			if allowed[hint.ID] && agentExplicitlyExcludesCandidate(req.Message, hint.ID, name) {
				justified[hint.ID] = true
				plan.ExcludedExpertIDs = append(plan.ExcludedExpertIDs, hint.ID)
			}
		}
		ignored := false
		if patch.ExcludedExpertIDs != nil {
			for _, id := range *patch.ExcludedExpertIDs {
				if !justified[id] {
					ignored = true
				}
			}
		}
		if ignored {
			warnings = append(warnings, "已忽略本轮未明确要求的额外排除；换批仅排除历史已排除及上轮推荐的专家。")
		}
	} else if patch.ExcludedExpertIDs != nil {
		plan.ExcludedExpertIDs = []uint{}
		ignored := false
		for _, id := range *patch.ExcludedExpertIDs {
			if id > 0 && allowed[id] {
				plan.ExcludedExpertIDs = append(plan.ExcludedExpertIDs, id)
			} else {
				ignored = true
			}
		}
		if ignored {
			warnings = append(warnings, "已忽略不属于本会话真实候选的排除编号。")
		}
	}
	seenIDs := map[uint]bool{}
	excluded := []uint{}
	for _, id := range plan.ExcludedExpertIDs {
		if id > 0 && !seenIDs[id] {
			seenIDs[id] = true
			excluded = append(excluded, id)
		}
	}
	if len(excluded) > 1000 {
		return plan, warnings, ErrExpertAgentInvalidRequest
	}
	plan.ExcludedExpertIDs = excluded
	return plan, warnings, nil
}

func agentInstructionClauses(message string) []string {
	return strings.FieldsFunc(message, func(r rune) bool { return strings.ContainsRune("，,。；;！？!?\n", r) })
}

func agentIntentNegated(prefix string) bool {
	prefix = strings.TrimSpace(prefix)
	// 只检查紧邻指令的否定，不把“不改人数，只换一批”里的其他否定当成拒绝换批。
	for {
		trimmed := false
		for _, filler := range []string{"重新", "继续", "给我", "帮我", "再"} {
			if strings.HasSuffix(prefix, filler) {
				prefix = strings.TrimSpace(strings.TrimSuffix(prefix, filler))
				trimmed = true
				break
			}
		}
		if !trimmed {
			break
		}
	}
	for _, negative := range []string{"不要", "不用", "不想", "不需要", "无需", "不必", "别", "取消", "停止", "不再", "不"} {
		if strings.HasSuffix(prefix, negative) {
			return true
		}
	}
	return false
}

func agentPhraseIntent(message string, phrases []string) (mentioned, allowed bool) {
	for _, clause := range agentInstructionClauses(message) {
		last := -1
		for _, phrase := range phrases {
			if index := strings.LastIndex(clause, phrase); index >= 0 && index >= last {
				last = index
				mentioned, allowed = true, !agentIntentNegated(clause[:index])
			}
		}
	}
	return mentioned, allowed
}

func agentAnotherBatchIntent(message string) (bool, bool) {
	return agentPhraseIntent(message, []string{"换一批", "换一组", "换一位", "换一人", "换另一批", "再来一批"})
}

func agentClearExclusionsIntent(message string) bool {
	_, allowed := agentPhraseIntent(message, []string{"取消排除", "取消所有排除", "取消之前的排除", "清空排除", "清空所有排除", "清除排除", "重置排除", "允许重复", "可以重复", "恢复被排除的专家"})
	return allowed
}

var agentExplicitIDPattern = regexp.MustCompile(`(?i)(?:\bid|编号|#|＃)\s*(?:[:：=]|为)?\s*([0-9]+)\b`)

func agentExplicitlyExcludesCandidate(message string, id uint, name string) bool {
	for _, clause := range agentInstructionClauses(message) {
		var mentions [][2]int
		if name != "" {
			if index := strings.Index(clause, name); index >= 0 {
				mentions = append(mentions, [2]int{index, index + len(name)})
			}
		}
		for _, match := range agentExplicitIDPattern.FindAllStringSubmatchIndex(clause, -1) {
			value, err := strconv.ParseUint(clause[match[2]:match[3]], 10, 64)
			if err == nil && value == uint64(id) {
				mentions = append(mentions, [2]int{match[0], match[1]})
			}
		}
		for _, mention := range mentions {
			prefix, suffix := clause[:mention[0]], clause[mention[1]:]
			position, exclude := -1, false
			for _, action := range []struct {
				word    string
				exclude bool
			}{
				{"排除", true}, {"不要", true}, {"别选", true}, {"不选", true}, {"不考虑", true}, {"不推荐", true}, {"剔除", true}, {"去掉", true}, {"避开", true},
				{"保留", false}, {"包括", false}, {"看看", false}, {"介绍", false}, {"比较", false},
			} {
				if index := strings.LastIndex(prefix, action.word); index >= 0 && index >= position {
					position = index
					exclude = action.exclude && !agentIntentNegated(prefix[:index])
				}
			}
			if position >= 0 {
				if exclude {
					return true
				}
				continue
			}
			// 后置表达如“张三也不要”。仅接受紧随该候选的指令，不把后面另一人的排除归给他。
			suffix = strings.TrimSpace(suffix)
			for _, filler := range []string{"也", "先", "暂时", "这次"} {
				suffix = strings.TrimPrefix(suffix, filler)
			}
			for _, action := range []string{"不要", "排除", "不考虑", "不选", "不推荐", "剔除", "去掉"} {
				if strings.HasPrefix(suffix, action) {
					return true
				}
			}
		}
	}
	return false
}

func sanitizeExpertAgentCandidates(in []ExpertDatabaseRes.ExpertAgentCandidate, excluded []uint) []ExpertDatabaseRes.ExpertAgentCandidate {
	seen := map[uint]bool{}
	for _, id := range excluded {
		seen[id] = true
	}
	out := make([]ExpertDatabaseRes.ExpertAgentCandidate, 0, len(in))
	for _, c := range in {
		if c.ID == 0 || seen[c.ID] {
			continue
		}
		seen[c.ID] = true
		if c.Evidence == nil {
			c.Evidence = []ExpertDatabaseRes.ExpertAgentEvidence{}
		}
		out = append(out, c)
	}
	return out
}

func agentUnitKey(name string) string { return strings.ToLower(strings.Join(strings.Fields(name), "")) }

// Explicit refusal also clears a previous full-database browsing request.
func agentBrowseIntent(message string) (mentioned, allowed bool) {
	clauses := strings.FieldsFunc(message, func(r rune) bool {
		return strings.ContainsRune("，,。；;！？!?\n", r)
	})
	for _, clause := range clauses {
		for _, phrase := range []string{"浏览全部", "全部专家", "全库", "不限方向", "不限主题", "不限领域", "任意专家", "随机专家"} {
			index := strings.Index(clause, phrase)
			if index < 0 {
				continue
			}
			mentioned, allowed = true, true
			prefix := []rune(clause[:index])
			if len(prefix) > 12 {
				prefix = prefix[len(prefix)-12:]
			}
			for _, negative := range []string{"不要", "不想", "不需要", "无需", "不必", "不浏览", "不查", "不看", "不从", "不再", "别", "取消", "停止", "拒绝"} {
				if strings.Contains(string(prefix), negative) {
					allowed = false
					break
				}
			}
		}
	}
	return mentioned, allowed
}

func agentCandidateUnitKey(candidate ExpertDatabaseRes.ExpertAgentCandidate) string {
	if candidate.OrgID > 0 {
		return fmt.Sprintf("org:%d", candidate.OrgID)
	}
	return "name:" + agentUnitKey(candidate.UnitName)
}

func selectExpertAgentCandidates(candidates []ExpertDatabaseRes.ExpertAgentCandidate, plan ExpertDatabaseRes.ExpertAgentPlan) []ExpertDatabaseRes.ExpertAgentCandidate {
	selected := make([]ExpertDatabaseRes.ExpertAgentCandidate, 0, plan.Count)
	units, disciplines, used := map[string]int{}, map[string]int{}, map[uint]bool{}
	for len(selected) < plan.Count {
		best, bestNovelty := -1, -1
		for i, c := range candidates {
			unit := agentCandidateUnitKey(c)
			if used[c.ID] || (plan.MaxPerUnit > 0 && units[unit] >= plan.MaxPerUnit) {
				continue
			}
			novelty := 0
			if plan.Diversify {
				if c.DisciplineL1 != "" && disciplines[c.DisciplineL1] == 0 {
					novelty += 2
				}
				if (c.OrgID > 0 || strings.TrimSpace(c.UnitName) != "") && units[unit] == 0 {
					novelty++
				}
			}
			if best == -1 || novelty > bestNovelty || (novelty == bestNovelty && c.Score > candidates[best].Score) {
				best, bestNovelty = i, novelty
			}
		}
		if best < 0 {
			break
		}
		c := candidates[best]
		used[c.ID] = true
		units[agentCandidateUnitKey(c)]++
		disciplines[c.DisciplineL1]++
		selected = append(selected, c)
	}
	return selected
}

func expertAgentStatistics(total int64, selected []ExpertDatabaseRes.ExpertAgentCandidate) ExpertDatabaseRes.ExpertAgentStatistics {
	stats := ExpertDatabaseRes.ExpertAgentStatistics{CandidateCount: int(total), SelectedCount: len(selected), Units: []ExpertDatabaseRes.ExpertAgentDistribution{}, Disciplines: []ExpertDatabaseRes.ExpertAgentDistribution{}}
	units, disciplines := map[string]int{}, map[string]int{}
	unitLabels := map[string]string{}
	for _, c := range selected {
		unit := strings.TrimSpace(c.UnitName)
		if unit == "" {
			unit = "未标注单位"
		}
		key := agentCandidateUnitKey(c)
		units[key]++
		if unitLabels[key] == "" {
			unitLabels[key] = unit
		}
		discipline := strings.TrimSpace(c.DisciplineL1)
		if discipline == "" {
			discipline = "未标注学科"
		}
		disciplines[discipline]++
	}
	for key, count := range units {
		label := unitLabels[key]
		stats.Units = append(stats.Units, ExpertDatabaseRes.ExpertAgentDistribution{Label: label, Count: count})
		if label != "未标注单位" {
			stats.UnitCount++
		}
	}
	for label, count := range disciplines {
		stats.Disciplines = append(stats.Disciplines, ExpertDatabaseRes.ExpertAgentDistribution{Label: label, Count: count})
		if label != "未标注学科" {
			stats.DisciplineCount++
		}
	}
	for _, distribution := range [][]ExpertDatabaseRes.ExpertAgentDistribution{stats.Units, stats.Disciplines} {
		sort.Slice(distribution, func(i, j int) bool {
			if distribution[i].Count != distribution[j].Count {
				return distribution[i].Count > distribution[j].Count
			}
			return distribution[i].Label < distribution[j].Label
		})
	}
	return stats
}
