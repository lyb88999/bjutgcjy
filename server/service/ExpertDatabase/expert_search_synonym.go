package ExpertDatabase

import (
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
)

// DictTypeSearchSynonym 检索同义词字典类型："系统工具-字典管理"里新建/编辑这个类型的字典明细，
// 每条明细的"扩展值"填一组用顿号或逗号分隔的相关词（比如"碳中和、低碳、碳达峰、双碳"），
// 检索词命中其中一个，整组词都会参与匹配。这是给 embedding 模型判断不出来的近义关系打的
// 一个确定性补丁——"碳中和"和"低碳"这两个词 bge-small-zh 判断得不够像（实测过，见
// rankedCandidates 里 minSemanticRelevance 的注释），调模型/调阈值都解决不了，只能靠人工
// 配置兜底；社科口的政策热词更新很快，字典可以随时加，不用发版
const DictTypeSearchSynonym = "expert_search_synonym"

// searchSynonymGroups 读取全部启用状态的同义词组，每组是一份去重后的词列表
func (s *ExpertSearchService) searchSynonymGroups() [][]string {
	var dict system.SysDictionary
	if err := global.GVA_DB.Where("type = ?", DictTypeSearchSynonym).First(&dict).Error; err != nil {
		return nil
	}
	var details []system.SysDictionaryDetail
	if err := global.GVA_DB.Where("sys_dictionary_id = ?", dict.ID).Find(&details).Error; err != nil {
		return nil
	}
	groups := make([][]string, 0, len(details))
	for _, d := range details {
		if d.Status != nil && !*d.Status {
			continue
		}
		terms := splitKeyword(d.Extend)
		if len(terms) > 1 {
			groups = append(groups, terms)
		}
	}
	return groups
}

// expandKeywordWithSynonyms 检索词命中某个同义词组时，把整组词并进来一起参与匹配；
// 检索词本身始终排在第一位，没命中任何一组时原样返回只有一个元素的切片
func expandKeywordWithSynonyms(keyword string, groups [][]string) []string {
	expanded := []string{keyword}
	seen := map[string]bool{keyword: true}
	for _, group := range groups {
		matched := false
		for _, t := range group {
			if t == keyword {
				matched = true
				break
			}
		}
		if !matched {
			continue
		}
		for _, t := range group {
			if !seen[t] {
				seen[t] = true
				expanded = append(expanded, t)
			}
		}
	}
	return expanded
}

// anyLiteralHit 展开后的任意一个词只要原样出现在语料里就算命中——用于子串匹配兜底路径和
// 语义匹配的字面命中保底，两处都要按展开后的整组词判断，不能只看用户原始输入的那一个词
func (s *ExpertSearchService) anyLiteralHit(corpus string, keywords []string) bool {
	for _, kw := range keywords {
		if strings.Contains(corpus, kw) {
			return true
		}
	}
	return false
}
