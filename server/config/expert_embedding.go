package config

// ExpertEmbedding 专家检索语义向量服务：本地自建的轻量 embedding 微服务地址，
// 用于把关键词/专家语料转成向量做余弦相似度匹配（见 service/ExpertDatabase/expert_embedding.go）
type ExpertEmbedding struct {
	Url        string `mapstructure:"url" json:"url" yaml:"url"`
	TimeoutSec int    `mapstructure:"timeout-sec" json:"timeoutSec" yaml:"timeout-sec"`
}
