package ExpertDatabase

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase"
)

// parsedEmbeddingItem 是 ExpertSearchEmbeddingItem 的检索期形态：向量已经从 JSON 文本解析成
// []float32，不用每次检索都重新反序列化一遍
type parsedEmbeddingItem struct {
	ItemLabel string
	Vector    []float32
}

// embeddingCache 把全部专家的语义向量条目缓存在内存里，用 expert_search_embedding_item 表的
// MAX(updated_at) 判断新鲜度：这张表只在管理员手动跑 cmd/recompute-embeddings（或数据导入后
// 批量重算）时才会整批更新，日常检索请求不会写它，所以"查一次 MAX(updated_at)，没变就用缓存"
// 这个代价很小的判断足够把"每次检索都要拉全表向量 JSON 再反序列化"这个随专家总数线性增长、
// 实测在 7700+ 专家规模下要占大半秒的开销降到几乎为零。
var embeddingCache = struct {
	mu        sync.RWMutex
	byExpert  map[uint][]parsedEmbeddingItem
	freshness time.Time
}{}

// embeddingItemsByExpert 返回 ids 对应专家的语义向量条目（已解析成 []float32），缓存过期或
// 首次调用时会全量重建一次缓存
func (s *ExpertSearchService) embeddingItemsByExpert(ids []uint) (map[uint][]parsedEmbeddingItem, error) {
	var maxUpdated *time.Time
	if err := global.GVA_DB.Model(&ExpertDatabase.ExpertSearchEmbeddingItem{}).
		Select("MAX(updated_at)").Scan(&maxUpdated).Error; err != nil {
		return nil, err
	}

	embeddingCache.mu.RLock()
	fresh := maxUpdated != nil && !embeddingCache.freshness.IsZero() && !maxUpdated.After(embeddingCache.freshness)
	if fresh {
		result := collectByIDs(embeddingCache.byExpert, ids)
		embeddingCache.mu.RUnlock()
		return result, nil
	}
	embeddingCache.mu.RUnlock()

	var rows []ExpertDatabase.ExpertSearchEmbeddingItem
	if err := global.GVA_DB.Find(&rows).Error; err != nil {
		return nil, err
	}
	byExpert := make(map[uint][]parsedEmbeddingItem, len(rows))
	for _, r := range rows {
		var v []float32
		if json.Unmarshal([]byte(r.Vector), &v) != nil {
			continue
		}
		byExpert[r.ExpertId] = append(byExpert[r.ExpertId], parsedEmbeddingItem{ItemLabel: r.ItemLabel, Vector: v})
	}

	embeddingCache.mu.Lock()
	embeddingCache.byExpert = byExpert
	if maxUpdated != nil {
		embeddingCache.freshness = *maxUpdated
	} else {
		embeddingCache.freshness = time.Now()
	}
	embeddingCache.mu.Unlock()

	return collectByIDs(byExpert, ids), nil
}

func collectByIDs(src map[uint][]parsedEmbeddingItem, ids []uint) map[uint][]parsedEmbeddingItem {
	result := make(map[uint][]parsedEmbeddingItem, len(ids))
	for _, id := range ids {
		if items, ok := src[id]; ok {
			result[id] = items
		}
	}
	return result
}
