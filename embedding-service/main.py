import os

from fastapi import FastAPI
from pydantic import BaseModel
from sentence_transformers import SentenceTransformer

MODEL_NAME = os.environ.get("EMBEDDING_MODEL", "BAAI/bge-small-zh-v1.5")

app = FastAPI()
model = SentenceTransformer(MODEL_NAME)


class EmbedRequest(BaseModel):
    texts: list[str]


class EmbedResponse(BaseModel):
    embeddings: list[list[float]]
    dim: int


@app.get("/health")
def health():
    return {"status": "ok", "model": MODEL_NAME, "dim": model.get_sentence_embedding_dimension()}


@app.post("/embed", response_model=EmbedResponse)
def embed(req: EmbedRequest):
    # normalize_embeddings=True 把向量归一化到单位长度，Go 那边算相似度直接点积就行，不用再算模长
    vectors = model.encode(req.texts, normalize_embeddings=True, convert_to_numpy=True)
    return EmbedResponse(embeddings=vectors.tolist(), dim=model.get_sentence_embedding_dimension())
