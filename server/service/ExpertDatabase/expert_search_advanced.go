package ExpertDatabase

import (
	"encoding/json"
	"strings"

	ExpertDatabaseReq "github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase/request"
)

// termResult 高级检索里单独一行（或者到目前为止累计折叠出的结果）的判定状态：命中没命中、
// 相关性多少。combineTermResults 就是拿这个结构互相折叠
type termResult struct {
	matched   bool
	relevance float64
}

// combineTermResults 把累计结果 acc 和这一行的结果 row 按连接词左折叠，三种连接词对应模糊逻辑
// 的且/或/非：
//   - 且：两边都命中才算命中，相关性取两边较小值——不能让一边的高分掩盖另一边勉强及格的匹配，
//     "且"查询的可信度应该由最弱的那个词决定
//   - 或：任意一边命中就算命中，相关性取两边较大值——采信证据更强的那一条
//   - 不含：累计结果命中、且这一行没命中才算命中，相关性沿用累计结果不变——纯排除，不该给
//     "成功排除了某个词"这件事本身加分
//
// 未知/空连接词按"或"兜底（正常情况下前端总会传三者之一，这里只是防御性兜底，不应该真的触发）
func combineTermResults(acc termResult, connector string, row termResult) termResult {
	switch connector {
	case "and":
		if !acc.matched || !row.matched {
			return termResult{matched: false, relevance: 0}
		}
		rel := acc.relevance
		if row.relevance < rel {
			rel = row.relevance
		}
		return termResult{matched: true, relevance: rel}
	case "not":
		if !acc.matched || row.matched {
			return termResult{matched: false, relevance: 0}
		}
		return termResult{matched: true, relevance: acc.relevance}
	default: // "or" 以及未知值兜底
		if !acc.matched && !row.matched {
			return termResult{matched: false, relevance: 0}
		}
		rel := acc.relevance
		if row.relevance > rel {
			rel = row.relevance
		}
		return termResult{matched: true, relevance: rel}
	}
}

// parseAdvancedRows 解析高级检索的 JSON 数组，过滤掉词是空白的行（防御性 trim，前端正常不会
// 提交空行，但导出/直接调接口的场景不能假设输入干净）。JSON 本身格式不对时直接返回 error，
// 让调用方（rankedCandidates）原样往上抛，走现有的 "检索失败" 错误处理，不需要额外的错误处理
func parseAdvancedRows(raw string) ([]ExpertDatabaseReq.AdvancedSearchRow, error) {
	var rows []ExpertDatabaseReq.AdvancedSearchRow
	if err := json.Unmarshal([]byte(raw), &rows); err != nil {
		return nil, err
	}
	filtered := make([]ExpertDatabaseReq.AdvancedSearchRow, 0, len(rows))
	for _, r := range rows {
		term := strings.TrimSpace(r.Term)
		if term == "" {
			continue
		}
		filtered = append(filtered, ExpertDatabaseReq.AdvancedSearchRow{Connector: r.Connector, Term: term})
	}
	return filtered, nil
}
