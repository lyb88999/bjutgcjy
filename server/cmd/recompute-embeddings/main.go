package main

import (
	"fmt"
	"os"

	"go.uber.org/zap"

	"github.com/flipped-aurora/gin-vue-admin/server/core"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/initialize"
	ExpertDatabase "github.com/flipped-aurora/gin-vue-admin/server/service/ExpertDatabase"
)

// 一次性维护命令：重算全部已发布专家的检索语义向量（语料没变化的专家会跳过，不会每次全量
// 重调 embedding 服务）。需要先把 embedding-service 起起来（默认 127.0.0.1:8901），否则会报错。
// 用法：go run ./cmd/recompute-embeddings
func main() {
	global.GVA_VP = core.Viper()
	global.GVA_LOG = core.Zap()
	zap.ReplaceGlobals(global.GVA_LOG)
	global.GVA_DB = initialize.Gorm()
	if global.GVA_DB == nil {
		fmt.Println("数据库连接失败")
		os.Exit(1)
	}

	svc := ExpertDatabase.ExpertSearchService{}
	if err := svc.RecomputeExpertEmbeddings(); err != nil {
		fmt.Println("重算失败:", err)
		os.Exit(1)
	}
	fmt.Println("全部已发布专家的检索语义向量已重算")
}
