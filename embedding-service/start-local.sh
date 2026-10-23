#!/bin/sh
set -eu

service_dir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
model_revision=7999e1d3359715c523056ef9478215996d62a620
model_path="$service_dir/.cache/hub/models--BAAI--bge-small-zh-v1.5/snapshots/$model_revision"

if [ ! -x "$service_dir/venv/bin/python" ] || [ ! -f "$model_path/model.safetensors" ]; then
    echo "本地运行环境或模型尚未准备，请先按 embedding-service/README.md 安装。" >&2
    exit 1
fi

export EMBEDDING_MODEL=BAAI/bge-small-zh-v1.5
export EMBEDDING_MODEL_PATH="$model_path"
export HF_HUB_OFFLINE=1
export TRANSFORMERS_OFFLINE=1
export HF_HUB_DISABLE_TELEMETRY=1

exec "$service_dir/venv/bin/python" -m uvicorn main:app --app-dir "$service_dir" --host 127.0.0.1 --port 8901 "$@"
