package main

import (
	"fmt"
	"os"
	"strconv"

	"go.uber.org/zap"

	"github.com/flipped-aurora/gin-vue-admin/server/core"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/initialize"
	ExpertDatabase "github.com/flipped-aurora/gin-vue-admin/server/service/ExpertDatabase"
)

// 一次性维护命令：重算指定专家（或全部已发布专家）的得分缓存
// 用法：go run ./cmd/recompute [expertID]
func main() {
	global.GVA_VP = core.Viper()
	global.GVA_LOG = core.Zap()
	zap.ReplaceGlobals(global.GVA_LOG)
	global.GVA_DB = initialize.Gorm()
	if global.GVA_DB == nil {
		fmt.Println("数据库连接失败")
		os.Exit(1)
	}

	svc := ExpertDatabase.ExpertScoreService{}
	if len(os.Args) > 1 {
		id, err := strconv.ParseUint(os.Args[1], 10, 64)
		if err != nil {
			fmt.Println("非法的专家ID:", os.Args[1])
			os.Exit(1)
		}
		if err := svc.RecomputeExpertScore(uint(id)); err != nil {
			fmt.Println("重算失败:", err)
			os.Exit(1)
		}
		fmt.Printf("专家 %d 得分已重算\n", id)
		return
	}
	if err := svc.RecomputeAllPublishedExpertScores(); err != nil {
		fmt.Println("重算失败:", err)
		os.Exit(1)
	}
	fmt.Println("全部已发布专家得分已重算")
}
