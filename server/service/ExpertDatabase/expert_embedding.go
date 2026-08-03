package ExpertDatabase

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase"
)

// embeddingModelName 落库时记一下用的哪个模型生成的向量，以后换模型能一眼看出哪些记录是旧模型
// 生成的、需要整体重算（换模型后新旧向量不在同一个空间里，不能混着比余弦相似度）
const embeddingModelName = "BAAI/bge-small-zh-v1.5"

// 语料条目类型：一个专家拆成好几条独立语料分别算向量，而不是拼成一整段文本只算一个向量，
// 检索时取"跟关键词最相似的那一条"作为整体相关性，避免不相关内容把真正对上的那条稀释掉
const (
	embeddingItemTypeProfile     = "profile"
	embeddingItemTypeAchievement = "achievement"
	embeddingItemTypeTag         = "tag"
)

type embedRequestBody struct {
	Texts []string `json:"texts"`
}

type embedResponseBody struct {
	Embeddings [][]float32 `json:"embeddings"`
}

// embedTexts 调用本地自建的轻量语义向量服务（仓库根目录 embedding-service/），把文本转成
// 归一化后的向量；服务地址在 config.yaml 的 expert-embedding.url 配置，默认 127.0.0.1:8901。
// 这一层出错时调用方一律要能优雅退回关键词子串匹配，不能让语义检索的故障拖垮整个检索功能。
// timeout 由调用方指定而不是固定读配置——实时检索只编码一个关键词，一次请求要快；离线批量
// 重算一批几十条语料，在配置较低的服务器上（2核 CPU 实测一批 32 条能到 11 秒+）会远超检索用
// 的超时，两个场景的时延预算天然不同，不能共用一个数字
func embedTexts(texts []string, timeout time.Duration) ([][]float32, error) {
	if len(texts) == 0 {
		return nil, nil
	}
	cfg := global.GVA_CONFIG.ExpertEmbedding
	if cfg.Url == "" {
		return nil, fmt.Errorf("expert-embedding.url 未配置")
	}
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	body, err := json.Marshal(embedRequestBody{Texts: texts})
	if err != nil {
		return nil, err
	}
	client := &http.Client{Timeout: timeout}
	resp, err := client.Post(cfg.Url+"/embed", "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("embedding 服务返回状态码 %d", resp.StatusCode)
	}
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var out embedResponseBody
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, err
	}
	if len(out.Embeddings) != len(texts) {
		return nil, fmt.Errorf("embedding 服务返回向量数(%d)与请求文本数(%d)不符", len(out.Embeddings), len(texts))
	}
	return out.Embeddings, nil
}

// corpusHash 语料指纹，只用来判断语料有没有变化（没变就不用重新调 embedding 服务），不用于安全用途
func corpusHash(text string) string {
	sum := sha256.Sum256([]byte(text))
	return hex.EncodeToString(sum[:16])
}

// cosineSim 两个向量的余弦相似度。embedding 服务已经把向量归一化到单位长度了，直接点积即可，
// 不用再各自除以模长
func cosineSim(a, b []float32) float64 {
	if len(a) == 0 || len(b) == 0 || len(a) != len(b) {
		return 0
	}
	var dot float64
	for i := range a {
		dot += float64(a[i]) * float64(b[i])
	}
	return dot
}

// buildProfileChunkText 专家档案自己填的研究方向/关键词等字段拼成一条独立语料（不含成果和标签，
// 那两块各自拆成单独的条目）
func buildProfileChunkText(p ExpertDatabase.ExpertProfile) string {
	var sb strings.Builder
	for _, f := range []string{
		p.ResearchDirections, p.ResearchKeywords, p.FocusTopics, p.ResearchObjects,
		p.MethodExpertise, p.RegionExpertise, p.DisciplineL1, p.DisciplineL2,
		p.CrossDiscipline, p.PolicyFields,
	} {
		if f != "" {
			sb.WriteString(f)
			sb.WriteString(" ")
		}
	}
	return strings.TrimSpace(sb.String())
}

// pendingEmbeddingItem 一条待生成/更新向量的语料条目
type pendingEmbeddingItem struct {
	expertID  uint
	itemType  string
	itemRefID uint
	itemLabel string
	text      string
	hash      string
}

// embeddingItemKey 用 (专家ID, 条目类型, 关联记录ID) 三元组标识一条语料条目，同一个专家的
// 研究方向和每一篇成果各自是独立的条目
func embeddingItemKey(expertID uint, itemType string, itemRefID uint) string {
	return fmt.Sprintf("%d|%s|%d", expertID, itemType, itemRefID)
}

// truncateLabel 按 rune 截断展示文本，避免中文单字节截断出现乱码，也避免检索结果里的
// "命中原因"提示因为一条很长的研究方向/成果标题而把界面撑爆
func truncateLabel(s string, maxRunes int) string {
	r := []rune(s)
	if len(r) <= maxRunes {
		return s
	}
	return string(r[:maxRunes]) + "…"
}

// RecomputeExpertEmbeddings 重算全部已发布专家的检索语义向量条目：研究方向算一条、每篇研究
// 成果各算一条、标签算一条。语料没变化的条目会跳过（按 corpus_hash 比对），成果被删掉或改没了
// 的旧条目也会一并清理，不会留着占检索时的比对开销。供运维脚本 cmd/recompute-embeddings 调用，
// 用法和打分缓存的 cmd/recompute 一致——数据导入/大批量修改之后手动跑一次
func (s *ExpertSearchService) RecomputeExpertEmbeddings() error {
	var profiles []ExpertDatabase.ExpertProfile
	if err := global.GVA_DB.Where("status = ?", "published").Find(&profiles).Error; err != nil {
		return err
	}
	if len(profiles) == 0 {
		return nil
	}
	ids := make([]uint, 0, len(profiles))
	for _, p := range profiles {
		ids = append(ids, p.ID)
	}

	var achievements []ExpertDatabase.ExpertAchievement
	if err := global.GVA_DB.Where("expert_id IN ?", ids).Find(&achievements).Error; err != nil {
		return err
	}
	achByExpert := make(map[uint][]ExpertDatabase.ExpertAchievement, len(profiles))
	for _, a := range achievements {
		achByExpert[a.ExpertId] = append(achByExpert[a.ExpertId], a)
	}

	type expertTagRow struct {
		ExpertId uint
		TagValue string
	}
	var tagRows []expertTagRow
	if err := global.GVA_DB.Model(&ExpertDatabase.ExpertTag{}).
		Select("expert_tag_relation.expert_id AS expert_id, expert_tag.tag_value AS tag_value").
		Joins("JOIN expert_tag_relation ON expert_tag_relation.tag_id = expert_tag.id AND expert_tag_relation.deleted_at IS NULL").
		Where("expert_tag_relation.expert_id IN ?", ids).
		Scan(&tagRows).Error; err != nil {
		return err
	}
	tagsByExpert := make(map[uint][]string, len(profiles))
	for _, row := range tagRows {
		tagsByExpert[row.ExpertId] = append(tagsByExpert[row.ExpertId], row.TagValue)
	}

	var existing []ExpertDatabase.ExpertSearchEmbeddingItem
	if err := global.GVA_DB.Where("expert_id IN ?", ids).Find(&existing).Error; err != nil {
		return err
	}
	existingByKey := make(map[string]ExpertDatabase.ExpertSearchEmbeddingItem, len(existing))
	for _, e := range existing {
		existingByKey[embeddingItemKey(e.ExpertId, e.ItemType, e.ItemRefId)] = e
	}

	desiredKeys := make(map[string]bool, len(existing))
	var todo []pendingEmbeddingItem
	addItem := func(expertID uint, itemType string, itemRefID uint, itemLabel, text string) {
		text = strings.TrimSpace(text)
		if text == "" {
			return
		}
		key := embeddingItemKey(expertID, itemType, itemRefID)
		desiredKeys[key] = true
		hash := corpusHash(text)
		if prev, ok := existingByKey[key]; ok && prev.CorpusHash == hash {
			return // 语料没变，不用重新调 embedding
		}
		todo = append(todo, pendingEmbeddingItem{
			expertID: expertID, itemType: itemType, itemRefID: itemRefID,
			itemLabel: itemLabel, text: text, hash: hash,
		})
	}

	for _, p := range profiles {
		if profileText := buildProfileChunkText(p); profileText != "" {
			addItem(p.ID, embeddingItemTypeProfile, 0, "研究方向："+truncateLabel(profileText, 40), profileText)
		}
		for _, a := range achByExpert[p.ID] {
			title := strings.Join(strings.Fields(a.Title), " ") // 有些成果标题里带换行，展示时拉平成一行
			addItem(p.ID, embeddingItemTypeAchievement, a.ID, "成果：《"+truncateLabel(title, 30)+"》", a.Title+" "+a.Keywords)
		}
		if tags := tagsByExpert[p.ID]; len(tags) > 0 {
			tagText := strings.Join(tags, " ")
			addItem(p.ID, embeddingItemTypeTag, 0, "标签："+truncateLabel(strings.Join(tags, "、"), 40), tagText)
		}
	}

	// 清理不再需要的旧条目：成果被删了/改空了，对应的向量条目也要跟着删，不然检索时还会拿
	// 一条早就不存在的成果去比相似度
	var staleIDs []uint
	for key, e := range existingByKey {
		if !desiredKeys[key] {
			staleIDs = append(staleIDs, e.ID)
		}
	}
	if len(staleIDs) > 0 {
		if err := global.GVA_DB.Delete(&ExpertDatabase.ExpertSearchEmbeddingItem{}, staleIDs).Error; err != nil {
			return err
		}
	}

	if len(todo) == 0 {
		return nil
	}

	// 分批调用，避免一次性把几百条语料丢给 embedding 服务导致单次请求体过大/超时。这是离线
	// 批处理，不影响用户实时检索体验，用一个远比检索场景宽松的超时（实测低配 CPU 上一批 32 条
	// 长语料能跑到 11 秒+，检索用的 5 秒对这里完全不够）
	const batchSize = 16
	const batchTimeout = 60 * time.Second
	for start := 0; start < len(todo); start += batchSize {
		end := start + batchSize
		if end > len(todo) {
			end = len(todo)
		}
		batch := todo[start:end]
		texts := make([]string, len(batch))
		for i, b := range batch {
			texts[i] = b.text
		}
		vectors, err := embedTexts(texts, batchTimeout)
		if err != nil {
			return fmt.Errorf("embedding 服务调用失败(第 %d/%d 批): %w", start/batchSize+1, (len(todo)+batchSize-1)/batchSize, err)
		}
		for i, b := range batch {
			vecJSON, err := json.Marshal(vectors[i])
			if err != nil {
				return err
			}
			key := embeddingItemKey(b.expertID, b.itemType, b.itemRefID)
			if prev, ok := existingByKey[key]; ok {
				if err := global.GVA_DB.Model(&ExpertDatabase.ExpertSearchEmbeddingItem{}).Where("id = ?", prev.ID).
					Updates(map[string]interface{}{
						"item_label":  b.itemLabel,
						"model":       embeddingModelName,
						"vector":      string(vecJSON),
						"corpus_hash": b.hash,
						"updated_at":  time.Now(),
					}).Error; err != nil {
					return err
				}
				continue
			}
			record := ExpertDatabase.ExpertSearchEmbeddingItem{
				ExpertId:   b.expertID,
				ItemType:   b.itemType,
				ItemRefId:  b.itemRefID,
				ItemLabel:  b.itemLabel,
				Model:      embeddingModelName,
				Vector:     string(vecJSON),
				CorpusHash: b.hash,
				UpdatedAt:  time.Now(),
			}
			if err := global.GVA_DB.Create(&record).Error; err != nil {
				return err
			}
		}
	}
	return nil
}
