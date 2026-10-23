<template>
  <aside class="facet-sidebar">
    <div
      v-for="group in visibleGroups"
      :key="group.key"
      class="facet-group"
    >
      <button
        type="button"
        class="facet-head"
        @click="toggleCollapse(group.key)"
      >
        <span>{{ group.title }}</span>
        <el-icon
          class="care"
          :class="{ 'is-collapsed': isCollapsed(group.key) }"
        ><ArrowDown /></el-icon>
      </button>
      <div
        v-show="!isCollapsed(group.key)"
        class="facet-body"
      >
        <el-input
          v-if="group.searchable && group.items.length > 10"
          v-model="filterText[group.key]"
          size="small"
          class="facet-search"
          :placeholder="`搜索${group.title}…`"
          clearable
        />
        <p
          v-if="group.searchable && filteredItems(group).length === 0"
          class="facet-empty"
        >没有匹配的选项</p>
        <ul
          v-else
          class="facet-list"
        >
          <li
            v-for="item in displayedItems(group)"
            :key="item.value"
            class="facet-item"
            :class="{ active: modelValue[group.key] === item.value }"
            @click="toggle(group.key, item.value)"
          >
            <span class="fname">{{ item.label }}</span>
            <span class="fcount">
              <!-- 计数是联动重算出来的，数字变了就该让人看得到"动了"——换个筛选条件，这个数字
              跟着轻轻翻一下，比数字啪地一下跳成另一个数字更能体现"这几个维度是联动的" -->
              <Transition
                name="count-flip"
                mode="out-in"
              >
                <span :key="item.count">{{ item.count }}</span>
              </Transition>
            </span>
          </li>
        </ul>
        <el-button
          v-if="hasMore(group)"
          link
          size="small"
          class="facet-more"
          @click="toggleExpand(group.key)"
        >{{ expanded[group.key] ? '收起' : `展开全部 ${filteredItems(group).length} 项` }}</el-button>
      </div>
    </div>
  </aside>
</template>

<script setup>
import { reactive, computed } from 'vue'

// 分面检索侧栏：facetCounts 是后端联动算好的各维度计数（每个维度已经排除了自己、套用了其它
// 维度当前选中的筛选），这个组件只管展开/收起、单位分面的本地搜索框这些纯 UI 状态，选中值本身
// 是受控的（modelValue + update:modelValue），不在这里维护
const props = defineProps({
  facetCounts: {
    type: Object,
    default: () => ({ org: [], disciplineL1: [], titleTier: [], regionExpertise: [], total: 0 })
  },
  // modelValue 四个筛选字段都用字符串："" 表示这一维度没有筛选，非空字符串是当前选中的
  // facet item.value（单位分面的 value 是 org_id 的字符串形式，其它维度是原始文本值）
  modelValue: {
    type: Object,
    required: true
  }
})
const emit = defineEmits(['update:modelValue'])

const groupDefs = [
  { key: 'orgId', title: '所在单位', dataKey: 'org', searchable: true },
  { key: 'disciplineL1', title: '一级学科', dataKey: 'disciplineL1', searchable: true },
  { key: 'titleTier', title: '专业技术职称', dataKey: 'titleTier', searchable: false },
  { key: 'regionExpertise', title: '区域/国别研究专长', dataKey: 'regionExpertise', searchable: false }
]

// 分面维度全库范围内没有数据时（比如 region_expertise 现在全库都还是空的）整组隐藏，不占地方——
// 等这个字段以后真被填了数据，接口会自动开始返回非空列表，这里不用改代码，跟专家详情页/检索结果
// 隐藏全零得分列是同一套逻辑
const visibleGroups = computed(() =>
  groupDefs
    .map((g) => ({ ...g, items: props.facetCounts[g.dataKey] || [] }))
    .filter((g) => g.items.length > 0)
)

const collapsed = reactive({})
const isCollapsed = (key) => !!collapsed[key]
const toggleCollapse = (key) => {
  collapsed[key] = !collapsed[key]
}

const expanded = reactive({})
const toggleExpand = (key) => {
  expanded[key] = !expanded[key]
}

const filterText = reactive({})
const VISIBLE_COUNT = 7

const filteredItems = (group) => {
  const q = (filterText[group.key] || '').trim()
  if (!q) return group.items
  return group.items.filter((it) => it.label.includes(q))
}
const displayedItems = (group) => {
  const items = filteredItems(group)
  return expanded[group.key] ? items : items.slice(0, VISIBLE_COUNT)
}
const hasMore = (group) => filteredItems(group).length > VISIBLE_COUNT

// 点已经选中的值 = 取消这个维度的筛选；点别的值 = 切到那个值（每个维度单选，不是多选，
// 跟知网分面检索的单值筛选行为一致）
const toggle = (fieldKey, value) => {
  const next = { ...props.modelValue }
  next[fieldKey] = next[fieldKey] === value ? '' : value
  emit('update:modelValue', next)
}
</script>

<style scoped>
/* 跟 ExpertResultCard.vue 用的是同一套登录页品牌色（src/view/login/index.vue 的
   $navy/$red/$gold），只在分面检索这几个新组件里局部用，不动全局 el-color-primary。
   藏青留给结构性的地方（标题条），选中态用朱红——跟设计原型的强调色分工一致，
   见 expertSearch.vue 顶部注释 */
.facet-sidebar {
  --brand-navy: #16407a;
  --brand-navy-dark: #0c2a54;
  --brand-red: #b5272d;
  --brand-gold: #d4a017;
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.facet-group {
  background: var(--el-bg-color);
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 3px;
  overflow: hidden;
}
.facet-head {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 12px;
  background: var(--brand-navy-dark);
  border: none;
  cursor: pointer;
  font-size: 13px;
  font-weight: 600;
  color: #fff;
  text-align: left;
}
.care {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.75);
  transition: transform 160ms ease;
}
.care.is-collapsed {
  transform: rotate(-90deg);
}
.facet-body {
  padding: 8px 12px 10px;
}
.facet-search {
  margin-bottom: 6px;
}
.facet-empty {
  font-size: 12.5px;
  color: var(--el-text-color-placeholder);
  padding: 4px 2px;
  margin: 0;
}
.facet-list {
  list-style: none;
  margin: 0;
  padding: 0;
  max-height: 260px;
  overflow-y: auto;
}
.facet-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  padding: 5px 6px;
  border-radius: var(--el-border-radius-base);
  cursor: pointer;
  font-size: 13px;
  color: var(--el-text-color-regular);
  transition: background-color 140ms ease, color 140ms ease;
}
.facet-item:hover {
  background: var(--el-fill-color-light);
}
.facet-item.active {
  background: rgba(181, 39, 45, 0.1);
  color: var(--brand-red);
  font-weight: 600;
}
.fname {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.fcount {
  flex-shrink: 0;
  font-size: 11.5px;
  font-variant-numeric: tabular-nums;
  color: var(--el-text-color-placeholder);
  display: inline-block;
}
.facet-item.active .fcount {
  color: var(--brand-red);
  opacity: 0.75;
}
.facet-more {
  margin-top: 4px;
  font-size: 12px;
}
.facet-more :deep(span) {
  color: var(--brand-red);
}

.count-flip-enter-active,
.count-flip-leave-active {
  transition: opacity 150ms ease, transform 150ms ease;
}
.count-flip-enter-from {
  opacity: 0;
  transform: translateY(-4px);
}
.count-flip-leave-to {
  opacity: 0;
  transform: translateY(4px);
}
@media (prefers-reduced-motion: reduce) {
  .count-flip-enter-active,
  .count-flip-leave-active {
    transition: none;
  }
}
</style>
