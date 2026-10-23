# 本地语义向量服务

服务使用 `BAAI/bge-small-zh-v1.5`，输出 512 维归一化向量，与当前专家库中已存的向量配套。模型文件和 Python 虚拟环境均被 Git 忽略。

## 首次准备

从本目录执行以下命令。项目已有可用 `venv` 时，可以直接使用它。

```sh
python3 -m venv venv
venv/bin/pip install -r requirements.txt
venv/bin/python - <<'PY'
from pathlib import Path
from huggingface_hub import snapshot_download

snapshot_download(
    repo_id="BAAI/bge-small-zh-v1.5",
    revision="7999e1d3359715c523056ef9478215996d62a620",
    cache_dir=str(Path(".cache/hub").resolve()),
    token=False,
    allow_patterns=[
        "config.json", "model.safetensors", "modules.json", "1_Pooling/config.json",
        "sentence_bert_config.json", "config_sentence_transformers.json",
        "tokenizer.json", "tokenizer_config.json", "special_tokens_map.json", "vocab.txt",
    ],
)
PY
```

模型文件约 103 MiB。此步骤只下载模型，不访问或重算专家数据。

## 启动与核验

在仓库根目录执行：

```sh
sh embedding-service/start-local.sh
curl http://127.0.0.1:8901/health
```

启动脚本使用固定模型快照，默认离线加载，监听本机 `8901`。健康响应应包含 `model: BAAI/bge-small-zh-v1.5` 和 `dim: 512`。`POST /embed` 接收 `{"texts":["社区养老"]}`，返回 `embeddings`、`dim` 和 `model`。

服务启动后，后端还需要预热数据库中已有的向量缓存。正常启动 `server/main.go` 会自动在后台预热；预热完成后助手即可使用语义检索。测试运行器若绕过了 `main.go`，也需要调用 `ExpertSearchService.WarmEmbeddingCache()`。

若服务故障，或响应中的模型、维度、归一化校验不通过，助手会使用关键词和同义词匹配并在结果中提示。模型文件更新不等于专家向量更新，本次启用过程不执行 `recompute-embeddings`。
