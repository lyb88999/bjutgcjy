package ExpertDatabase

import (
	"strconv"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"go.uber.org/zap"
)

// 字典类型：各级别/权重都走数据字典配置，找不到字典时回退到下面的默认值，
// 这样功能开箱即用，管理员可以随时在"系统工具 -> 字典管理"里调整而不用发版。
const (
	DictTypeAchievementLevel = "expert_achievement_level" // 成果/课题级别权重
	DictTypeAdoptionLevel    = "expert_adoption_level"    // 采纳/批示单位级别权重
	DictTypeTitleLevel       = "expert_title_level"       // 职称权重
	DictTypeSocialWeight     = "expert_social_weight"     // 社会贡献权重（学术兼职、荣誉称号）
	DictTypeRankingWeight    = "expert_ranking_weight"    // 综合排序各分项的权重 w1..w5
)

var (
	defaultAchievementLevelWeights = map[string]float64{"国家级": 10, "省部级": 5, "厅局级": 3, "一般级": 1}
	defaultAdoptionLevelWeights    = map[string]float64{"国家级": 10, "中央": 10, "省部级": 5, "厅局级": 3, "区县级": 1}
	defaultTitleLevelWeights       = map[string]float64{"院士": 8, "教授": 5, "研究员": 5, "副教授": 3, "副研究员": 3, "讲师": 1, "助理研究员": 1}
	defaultSocialWeights           = map[string]float64{"academic_position": 1, "honor_title": 3}
	// keyword 是关键词匹配本身的固定加分（w5）：整个综合得分 = relevance * (前四项 + keyword)，
	// 这一项跟其他四项一起被 relevance 整体打折，保证哪怕成果分/决策影响分是 0 的专家（还没
	// 攒够成果）也能靠这一项让综合得分随关键词匹配程度变化——见 rankedCandidates 顶部注释
	defaultRankingWeights = map[string]float64{"achievement": 1, "influence": 1, "title": 1, "social": 1, "keyword": 10}
)

type ExpertScoreService struct{}

// expertScoreSvc 供本包其他 service 文件（成果/决策影响记录变更时）触发得分重算
var expertScoreSvc = &ExpertScoreService{}

// dictWeights 按字典类型读取 label -> 权重 的映射，在 fallback 默认值之上按管理员配置的明细
// 覆盖，而不是整份替换——这样管理员之前只配置过部分分项（比如早期只有 achievement/influence/
// title/social 四项）时，代码后续新增的分项（比如这里的 keyword）依然能吃到默认值，不会因为
// 字典存在就整体退化成"没配置的分项权重是 0"
func (s *ExpertScoreService) dictWeights(dictType string, fallback map[string]float64) map[string]float64 {
	weights := make(map[string]float64, len(fallback))
	for k, v := range fallback {
		weights[k] = v
	}
	var dict system.SysDictionary
	if err := global.GVA_DB.Where("type = ?", dictType).First(&dict).Error; err != nil {
		return weights
	}
	var details []system.SysDictionaryDetail
	if err := global.GVA_DB.Where("sys_dictionary_id = ?", dict.ID).Find(&details).Error; err != nil {
		return weights
	}
	for _, d := range details {
		if w, err := strconv.ParseFloat(d.Extend, 64); err == nil {
			weights[d.Label] = w
		} else {
			weights[d.Label] = float64(d.Value)
		}
	}
	return weights
}

// calcLeveledScore 级别权重 × 相应级别篇数的求和，是成果得分与决策影响得分共用的核心算法
func calcLeveledScore(levels []string, weights map[string]float64) float64 {
	counts := make(map[string]int, len(levels))
	for _, level := range levels {
		counts[level]++
	}
	var score float64
	for level, count := range counts {
		w, ok := weights[level]
		if !ok {
			w = 1 // 未在字典里配置的级别，给保底权重，避免因为漏配级别导致该成果完全不计分
		}
		score += w * float64(count)
	}
	return score
}

// titleTierKeywords 职称原文往往是自由文本（"长聘教授"、"主任医师、教授"、"院士/教授"这种
// 组合写法很常见），不是字典里"教授/副教授/讲师..."这几个干净标签能直接精确匹配上的。之前用
// 精确匹配、查不到就退回最低档 1 分，实测导致 28% 的档案（含院士、讲席教授、长聘教授、主任
// 医师等本该是中高级的头衔）被错误打到跟讲师一个分数。改成关键词匹配：
//   1. 先查有没有"副/助理/准聘"这类降级修饰词（覆盖"研究员、准聘副教授"这种同时出现正副职称
//      的组合，按更保守的副高级算，不高估）
//   2. 再查初级职称关键词
//   3. 再查院士（比普通教授更高一档，原来的字典里都没这一档）
//   4. 最后查正高级关键词（教授/研究员/主任医师/主任药师/主任技师这些临床对应职称也算进去）
// 都没命中的才退回最低档，跟以前的兜底行为一致
var titleTierKeywords = []struct {
	dictLabel string
	fallback  float64
	keywords  []string
}{
	{"副教授", 3, []string{"副教授", "副研究员", "副主任医师", "副主任药师", "副主任技师", "副主任护师", "准聘"}},
	{"讲师", 1, []string{"讲师", "助理研究员", "助理教授", "住院医师", "研究实习员"}},
	{"院士", 8, []string{"院士"}},
	{"教授", 5, []string{"教授", "研究员", "主任医师", "主任药师", "主任技师", "主任护师", "首席研究员"}},
}

// matchTitleWeight 按关键词匹配职称原文对应的权重档位，找不到字典里配置的标签时退回该档位的
// 默认权重（管理员没配置这个新档位也不会导致算分出错）
func matchTitleWeight(techTitle string, weights map[string]float64) float64 {
	for _, tier := range titleTierKeywords {
		for _, kw := range tier.keywords {
			if strings.Contains(techTitle, kw) {
				if w, ok := weights[tier.dictLabel]; ok {
					return w
				}
				return tier.fallback
			}
		}
	}
	return 1
}

// RecomputeExpertScore 重新计算某个专家的成果分/决策影响分/综合排序基础分并写回缓存字段
// 综合排序基础分不含检索相关性（与关键词无关），检索时再乘以本次查询的 relevance 系数
func (s *ExpertScoreService) RecomputeExpertScore(expertID uint) error {
	var profile ExpertDatabase.ExpertProfile
	if err := global.GVA_DB.Where("id = ?", expertID).First(&profile).Error; err != nil {
		return err
	}

	var achievements []ExpertDatabase.ExpertAchievement
	if err := global.GVA_DB.Where("expert_id = ?", expertID).Find(&achievements).Error; err != nil {
		return err
	}
	achievementLevels := make([]string, 0, len(achievements))
	for _, a := range achievements {
		achievementLevels = append(achievementLevels, a.Level)
	}
	achievementScore := calcLeveledScore(achievementLevels, s.dictWeights(DictTypeAchievementLevel, defaultAchievementLevelWeights))

	var adoptions []ExpertDatabase.ExpertAdoptionRecord
	if err := global.GVA_DB.Where("expert_id = ?", expertID).Find(&adoptions).Error; err != nil {
		return err
	}
	adoptionLevels := make([]string, 0, len(adoptions))
	for _, a := range adoptions {
		adoptionLevels = append(adoptionLevels, a.AdoptingUnitLevel)
	}
	influenceScore := calcLeveledScore(adoptionLevels, s.dictWeights(DictTypeAdoptionLevel, defaultAdoptionLevelWeights))

	titleWeights := s.dictWeights(DictTypeTitleLevel, defaultTitleLevelWeights)
	titleScore := matchTitleWeight(profile.TechTitle, titleWeights)

	var positionCount int64
	if err := global.GVA_DB.Model(&ExpertDatabase.ExpertAcademicPosition{}).Where("expert_id = ?", expertID).Count(&positionCount).Error; err != nil {
		return err
	}
	socialWeights := s.dictWeights(DictTypeSocialWeight, defaultSocialWeights)
	socialScore := float64(positionCount) * socialWeights["academic_position"]
	if profile.HonorTitle != "" {
		socialScore += socialWeights["honor_title"]
	}

	rankingWeights := s.dictWeights(DictTypeRankingWeight, defaultRankingWeights)
	compositeScore := rankingWeights["achievement"]*achievementScore +
		rankingWeights["influence"]*influenceScore +
		rankingWeights["title"]*titleScore +
		rankingWeights["social"]*socialScore

	now := time.Now()
	return global.GVA_DB.Model(&ExpertDatabase.ExpertProfile{}).Where("id = ?", expertID).Updates(map[string]interface{}{
		"achievement_score": achievementScore,
		"influence_score":   influenceScore,
		"social_score":      socialScore,
		"composite_score":   compositeScore,
		"score_updated_at":  &now,
	}).Error
}

// recomputeIfPublished 仅当专家档案已发布时才重算得分缓存，避免为草稿/审核中的档案做无意义的计算
func (s *ExpertScoreService) recomputeIfPublished(expertID uint) {
	var profile ExpertDatabase.ExpertProfile
	if err := global.GVA_DB.Where("id = ?", expertID).First(&profile).Error; err != nil {
		return
	}
	if profile.Status != "published" {
		return
	}
	if err := s.RecomputeExpertScore(expertID); err != nil {
		global.GVA_LOG.Error("重算专家得分失败", zap.Error(err))
	}
}

// RecomputeAllPublishedExpertScores 兜底批量重算所有已发布专家的得分缓存，供定时任务调用
func (s *ExpertScoreService) RecomputeAllPublishedExpertScores() error {
	var ids []uint
	if err := global.GVA_DB.Model(&ExpertDatabase.ExpertProfile{}).Where("status = ?", "published").Pluck("id", &ids).Error; err != nil {
		return err
	}
	for _, id := range ids {
		if err := s.RecomputeExpertScore(id); err != nil {
			global.GVA_LOG.Error("重算专家得分失败", zap.Error(err))
		}
	}
	return nil
}
