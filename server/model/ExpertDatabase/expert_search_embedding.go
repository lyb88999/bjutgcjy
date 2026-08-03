package ExpertDatabase

import (
	"time"
)

// ExpertSearchEmbedding 专家检索语义向量缓存：语料（研究方向/关键词/成果/标签拼起来）过
// embedding 模型转出的归一化向量，检索时跟关键词向量算余弦相似度用。corpus_hash 记的是
// 生成这条向量时的语料指纹，语料没变就不用重新调 embedding 服务，重算脚本据此跳过没变化的专家
type ExpertSearchEmbedding struct {
	ExpertId   uint      `json:"expertId" gorm:"column:expert_id;primaryKey;autoIncrement:false;comment:专家ID"`
	Model      string    `json:"model" gorm:"column:model;comment:生成向量用的embedding模型名"`
	Vector     string    `json:"vector" gorm:"column:vector;type:json;comment:归一化后的向量,JSON浮点数组"`
	CorpusHash string    `json:"corpusHash" gorm:"column:corpus_hash;comment:生成向量时的语料指纹"`
	UpdatedAt  time.Time `json:"updatedAt" gorm:"column:updated_at;comment:向量更新时间"`
}

// TableName 专家检索语义向量 ExpertSearchEmbedding自定义表名 expert_search_embedding
func (ExpertSearchEmbedding) TableName() string {
	return "expert_search_embedding"
}
