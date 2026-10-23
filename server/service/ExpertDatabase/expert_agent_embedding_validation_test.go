package ExpertDatabase

import (
	"math"
	"testing"
)

func TestAgentEmbeddingCompatibilityRejectsMismatchedSpaces(t *testing.T) {
	vector := make([]float32, expertAgentEmbeddingDimension)
	vector[0] = 1
	for _, test := range []struct {
		name     string
		response embedResponseBody
		valid    bool
	}{
		{"matching model", embedResponseBody{Embeddings: [][]float32{vector}, Dim: 512, Model: embeddingModelName}, true},
		{"legacy service without metadata", embedResponseBody{Embeddings: [][]float32{vector}}, true},
		{"different model", embedResponseBody{Embeddings: [][]float32{vector}, Model: "different-model"}, false},
		{"incorrect declared dimension", embedResponseBody{Embeddings: [][]float32{vector}, Dim: 768}, false},
		{"truncated vector", embedResponseBody{Embeddings: [][]float32{{1, 0}}}, false},
		{"missing result", embedResponseBody{Embeddings: [][]float32{}}, false},
		{"zero vector", embedResponseBody{Embeddings: [][]float32{make([]float32, 512)}}, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			if err := validateAgentEmbeddingResponse(test.response, 1); (err == nil) != test.valid {
				t.Fatalf("embedding compatibility accepted=%v, expected=%v: %v", err == nil, test.valid, err)
			}
		})
	}
	for _, number := range []float32{float32(math.Inf(1)), float32(math.NaN())} {
		invalid := append([]float32{}, vector...)
		invalid[1] = number
		if err := validateAgentEmbeddingResponse(embedResponseBody{Embeddings: [][]float32{invalid}}, 1); err == nil {
			t.Fatal("nonfinite embedding values must not enter ranking")
		}
	}
}
