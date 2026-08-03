#!/bin/bash
# 健康检查看门狗：定期探活，服务真的挂了 systemd 的 Restart=always 已经会自动拉起来，
# 这个脚本专门兜"进程还在但卡死不响应"这种 Restart 机制本身抓不到的情况
set -euo pipefail

if ! curl -fsS --max-time 5 http://127.0.0.1:8901/health > /dev/null 2>&1; then
    logger -t embedding-service-watchdog "健康检查失败，重启 embedding-service"
    systemctl restart embedding-service
fi
