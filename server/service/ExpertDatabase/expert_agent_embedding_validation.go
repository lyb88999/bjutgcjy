package ExpertDatabase

import (
	"fmt"
	"math"
)

const expertAgentEmbeddingDimension = 512

// A different vector space must trigger the existing keyword fallback rather than silently miss experts.
func validateAgentEmbeddingResponse(response embedResponseBody, count int) error {
	if len(response.Embeddings) != count || count == 0 {
		return fmt.Errorf("语义向量数量无效")
	}
	if response.Model != "" && response.Model != embeddingModelName {
		return fmt.Errorf("语义模型与现有专家向量不一致")
	}
	if response.Dim != 0 && response.Dim != expertAgentEmbeddingDimension {
		return fmt.Errorf("语义向量维度与现有专家向量不一致")
	}
	for _, vector := range response.Embeddings {
		if len(vector) != expertAgentEmbeddingDimension {
			return fmt.Errorf("语义向量维度与现有专家向量不一致")
		}
		var norm float64
		for _, number := range vector {
			x := float64(number)
			if math.IsNaN(x) || math.IsInf(x, 0) {
				return fmt.Errorf("语义向量数值无效")
			}
			norm += x * x
		}
		if math.Abs(norm-1) > 0.05 {
			return fmt.Errorf("语义向量未正确归一化")
		}
	}
	return nil
}
