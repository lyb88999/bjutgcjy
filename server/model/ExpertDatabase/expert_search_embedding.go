package ExpertDatabase

import (
	"time"
)

// ExpertSearchEmbeddingItem 专家检索语义向量缓存：一个专家对应多条语料条目（研究方向算一条、
// 每篇研究成果各算一条、标签算一条），分别转成归一化向量存起来，而不是把全部语料拼成一整段文本
// 只算一个向量——检索时对某个专家取"跟关键词最相似的那一条"作为整体相关性，避免专家名下一堆不
// 相关的成果把真正对上的那条内容"稀释"掉（比如某专家 5 篇成果只有 1 篇讲低碳经济，拼一起算的话
// 那 1 篇的信号会被另外 4 篇冲淡）。ItemLabel 顺带记下这条语料是什么（成果标题/"研究方向"等），
// 检索结果里可以直接告诉用户"具体是命中了哪一条"，不只是一个相关性百分比。
// corpus_hash 是这一条语料的指纹，语料没变就跳过、不重新调用 embedding 服务
type ExpertSearchEmbeddingItem struct {
	ID         uint      `json:"id" gorm:"column:id;primaryKey;autoIncrement;comment:自增ID"`
	ExpertId   uint      `json:"expertId" gorm:"column:expert_id;index;comment:专家ID"`
	ItemType   string    `json:"itemType" gorm:"column:item_type;index;comment:语料条目类型:profile/achievement/tag"`
	ItemRefId  uint      `json:"itemRefId" gorm:"column:item_ref_id;comment:关联的原始记录ID(比如成果ID);profile/tag类型固定为0"`
	ItemLabel  string    `json:"itemLabel" gorm:"column:item_label;comment:这条语料的展示文本,用于检索结果的命中原因提示"`
	Model      string    `json:"model" gorm:"column:model;comment:生成向量用的embedding模型名"`
	Vector     string    `json:"vector" gorm:"column:vector;type:json;comment:归一化后的向量,JSON浮点数组"`
	CorpusHash string    `json:"corpusHash" gorm:"column:corpus_hash;comment:生成向量时的语料指纹"`
	UpdatedAt  time.Time `json:"updatedAt" gorm:"column:updated_at;comment:向量更新时间"`
}

// TableName 专家检索语义向量条目 ExpertSearchEmbeddingItem自定义表名 expert_search_embedding_item
func (ExpertSearchEmbeddingItem) TableName() string {
	return "expert_search_embedding_item"
}
