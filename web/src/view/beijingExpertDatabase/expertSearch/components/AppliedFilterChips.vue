<template>
  <div
    v-if="chips.length"
    class="applied-chips"
  >
    <span class="applied-label">已选条件：</span>
    <TransitionGroup
      tag="span"
      name="chip-pop"
      style="display: contents"
    >
      <el-tag
        v-for="chip in chips"
        :key="chip.key"
        closable
        class="applied-chip"
        @close="$emit('remove', chip.key)"
      >{{ chip.label }}</el-tag>
    </TransitionGroup>
    <el-button
      link
      size="small"
      class="clear-all"
      @click="$emit('clear-all')"
    >清除全部</el-button>
  </div>
</template>

<script setup>
// 已选条件标签行：纯展示组件，选中值和"哪个标签对应哪个筛选维度"都由父组件（expertSearch.vue）
// 算好传进来，这里只管渲染和把用户的点击 emit 回去
defineProps({
  // chips: Array<{ key: string, label: string }>，key 跟 FacetSidebar 的 modelValue 字段名对齐
  // （orgId/disciplineL1/titleTier/regionExpertise），移除某个标签时父组件按 key 清掉对应字段
  chips: {
    type: Array,
    default: () => []
  }
})
defineEmits(['remove', 'clear-all'])
</script>

<style scoped>
.applied-chips {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px;
  padding: 10px 0;
  margin-bottom: 4px;
  border-bottom: 1px solid var(--el-border-color-lighter);
}
.applied-label {
  font-size: 12.5px;
  color: var(--el-text-color-secondary);
  flex-shrink: 0;
}
.applied-chip {
  /* 已选条件是"当前生效的筛选"，用朱红这个强调色而不是结构色藏青——跟设计原型
     （expert-search-cnki-prototype.html）里 .chip 用 brick 强调、facet-head 用 navy 结构
     的分工一致，藏青留给品牌/结构性的地方 */
  --el-tag-bg-color: rgba(181, 39, 45, 0.08);
  --el-tag-border-color: rgba(181, 39, 45, 0.28);
  --el-tag-text-color: #b5272d;
}
.clear-all {
  font-size: 12.5px;
  margin-left: 4px;
}
.clear-all :deep(span) {
  color: #b5272d;
}

.chip-pop-enter-active,
.chip-pop-leave-active {
  transition: opacity 160ms ease, transform 160ms ease;
}
.chip-pop-enter-from,
.chip-pop-leave-to {
  opacity: 0;
  transform: scale(0.85);
}
@media (prefers-reduced-motion: reduce) {
  .chip-pop-enter-active,
  .chip-pop-leave-active {
    transition: none;
  }
}
</style>
