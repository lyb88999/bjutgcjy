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
	"gorm.io/gorm/clause"
)

// embeddingModelName 落库时记一下用的哪个模型生成的向量，以后换模型能一眼看出哪些记录是旧模型
// 生成的、需要整体重算（换模型后新旧向量不在同一个空间里，不能混着比余弦相似度）
const embeddingModelName = "BAAI/bge-small-zh-v1.5"

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
// 重算一批几十条、每条又是拼起来的长语料，在配置较低的服务器上（2核 CPU 实测一批 32 条能到
// 11 秒+）会远超检索用的超时，两个场景的时延预算天然不同，不能共用一个数字
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

// buildExpertCorpusText 拼出一个专家参与语义匹配的全部语料：档案自身的研究方向/关键词/关注
// 议题等字段 + 这个专家的成果标题关键词 + 标签（后两者来自 buildSearchCorpusMap，跟老的子串
// 匹配共用同一份"素材"，避免同一份数据在两个地方各拼一遍、后续改字段容易漏改一处）
func buildExpertCorpusText(p ExpertDatabase.ExpertProfile, achievementAndTagCorpus string) string {
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
	sb.WriteString(achievementAndTagCorpus)
	return sb.String()
}

// RecomputeExpertEmbeddings 重算全部已发布专家的检索语义向量。语料没变化的专家会跳过（按
// corpus_hash 比对），不会每次全量重调 embedding 服务；供运维脚本 cmd/recompute-embeddings 调用，
// 用法和打分缓存的 cmd/recompute 一致——数据导入/大批量修改之后手动跑一次
func (s *ExpertSearchService) RecomputeExpertEmbeddings() error {
	var profiles []ExpertDatabase.ExpertProfile
	if err := global.GVA_DB.Where("status = ?", "published").Find(&profiles).Error; err != nil {
		return err
	}
	if len(profiles) == 0 {
		return nil
	}

	achievementTagCorpus, err := s.buildSearchCorpusMap(profiles)
	if err != nil {
		return err
	}

	ids := make([]uint, 0, len(profiles))
	for _, p := range profiles {
		ids = append(ids, p.ID)
	}
	var existing []ExpertDatabase.ExpertSearchEmbedding
	if err := global.GVA_DB.Where("expert_id IN ?", ids).Find(&existing).Error; err != nil {
		return err
	}
	existingHash := make(map[uint]string, len(existing))
	for _, e := range existing {
		existingHash[e.ExpertId] = e.CorpusHash
	}

	type pendingItem struct {
		expertID uint
		text     string
		hash     string
	}
	var todo []pendingItem
	for _, p := range profiles {
		text := buildExpertCorpusText(p, achievementTagCorpus[p.ID])
		if text == "" {
			continue
		}
		hash := corpusHash(text)
		if existingHash[p.ID] == hash {
			continue
		}
		todo = append(todo, pendingItem{expertID: p.ID, text: text, hash: hash})
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
			record := ExpertDatabase.ExpertSearchEmbedding{
				ExpertId:   b.expertID,
				Model:      embeddingModelName,
				Vector:     string(vecJSON),
				CorpusHash: b.hash,
				UpdatedAt:  time.Now(),
			}
			if err := global.GVA_DB.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "expert_id"}},
				DoUpdates: clause.AssignmentColumns([]string{"model", "vector", "corpus_hash", "updated_at"}),
			}).Create(&record).Error; err != nil {
				return err
			}
		}
	}
	return nil
}
