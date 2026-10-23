<template>
  <div class="gva-table-box">
    <!-- 原来这里有一整块蓝底的 hero 横幅（图标+机构名+"专家检索"标题）。去掉了：机构全称
         其实已经显示在常驻的侧边栏（$GIN_VUE_ADMIN.appName），"专家检索"这个页面标题也已经
         在顶部面包屑里（route.meta.title），横幅完全是在重复展示已经存在的信息。同目录下
         其它管理页（比如 expertProfile.vue）也都是直接进内容，没有另起一块横幅——去掉之后
         跟这些页面才是同一套页面结构，不是给这一个页面单独加的"落地页感"装饰 -->
    <div class="gva-search-box">
      <!-- 知网的检索框是一整块独立的、视觉上焊死在一起的检索工具：左边是选检索条件的下拉框，
           不是又一个单独的输入框——字段选择+输入+提交按钮没有缝隙，拼成一个白色胶囊容器
           （.search-bar），右端红色按钮跟容器边界严丝合缝；高级检索/重置/导出这些次要操作
           挪到条形下面单独一行，不再挤在同一个 el-form 里 -->
      <div
        class="search-bar"
        :class="{ 'search-bar-attached': advancedPanelOpen }"
      >
        <el-select
          v-model="searchScope"
          class="search-bar-scope"
        >
          <el-option
            label="全部字段"
            value="keyword"
          />
          <el-option
            label="姓名"
            value="name"
          />
          <el-option
            label="所在单位"
            value="unitName"
          />
        </el-select>
        <span class="search-bar-divider" />
        <el-input
          v-model="searchQuery"
          class="search-bar-input search-bar-keyword"
          :disabled="searchScope === 'keyword' && advancedActive"
          :placeholder="scopeInputPlaceholder"
          clearable
          @keyup.enter="onSearch"
        />
        <button
          type="button"
          class="search-bar-submit"
          @click="onSearch"
        >
          <el-icon><Search /></el-icon>
          检索
        </button>
      </div>

      <!-- 高级检索原来是浮在检索条边上的一个 el-popover 浮层，跟检索条本身是两套视觉语言
           （圆角卡片 vs 焊死的胶囊条），显得格格不入。现在改成挂在检索条正下方、共用同一条
           边框和圆角的"抽屉"——展开时看起来像检索条本身长出来的一部分，跟知网点开"高级检索"
           之后原地展开一张检索式表格是一个思路，不是另开一个浮动面板 -->
      <Transition name="advanced-expand">
        <div
          v-if="advancedPanelOpen"
          class="advanced-panel"
        >
          <div
            v-for="(row, idx) in advancedRows"
            :key="idx"
            class="advanced-row"
          >
            <span class="advanced-row-index">{{ idx + 1 }}</span>
            <el-select
              v-if="idx > 0"
              v-model="row.connector"
              class="advanced-connector"
            >
              <el-option
                label="且"
                value="and"
              />
              <el-option
                label="或"
                value="or"
              />
              <el-option
                label="不含"
                value="not"
              />
            </el-select>
            <span
              v-else
              class="advanced-connector advanced-connector-placeholder"
            >检索词</span>
            <el-input
              v-model="row.term"
              class="advanced-term"
              placeholder="输入检索词"
              clearable
              @keyup.enter="applyAdvanced"
            />
            <button
              type="button"
              class="advanced-row-remove"
              :disabled="advancedRows.length <= 1"
              @click="removeAdvancedRow(idx)"
            >
              <el-icon><Close /></el-icon>
            </button>
          </div>
          <div class="advanced-panel-footer">
            <button
              type="button"
              class="advanced-add-row"
              @click="addAdvancedRow"
            >+ 添加检索行</button>
            <div class="advanced-panel-footer-actions">
              <button
                type="button"
                class="search-bar-action"
                @click="clearAdvanced"
              >清空</button>
              <button
                type="button"
                class="advanced-apply"
                @click="applyAdvanced"
              >检索</button>
            </div>
          </div>
        </div>
      </Transition>

      <div class="search-bar-actions">
        <button
          type="button"
          class="search-bar-action"
          :class="{ 'is-active': advancedActive }"
          @click="advancedPanelOpen = !advancedPanelOpen"
        >高级检索{{ advancedActive ? '（已启用）' : '' }}</button>
        <span class="search-bar-action-sep">|</span>
        <button
          type="button"
          class="search-bar-action"
          @click="onReset"
        >重置</button>
        <span class="search-bar-action-sep">|</span>
        <button
          type="button"
          class="search-bar-action"
          :disabled="exporting"
          @click="handleExport"
        >{{ exporting ? '导出中…' : '导出当前结果' }}</button>
      </div>
    </div>

    <AppliedFilterChips
      :chips="appliedChips"
      @remove="removeChip"
      @clear-all="clearAllFacets"
    />

    <div class="result-toolbar">
      <span class="result-count">共 <b>{{ total }}</b> 位专家</span>
      <div
        v-if="!hasKeyword"
        class="sort-tabs"
      >
        <button
          v-for="opt in sortOptions"
          :key="opt.value"
          type="button"
          class="sort-tab"
          :class="{ active: searchForm.sortBy === opt.value }"
          @click="changeSortBy(opt.value)"
        >{{ opt.label }}</button>
      </div>
    </div>

    <div class="search-layout">
      <FacetSidebar
        :facet-counts="facetCounts"
        :model-value="facetSelection"
        @update:model-value="onFacetChange"
      />
      <div
        v-loading="loading"
        class="results-col"
      >
        <TransitionGroup
          v-if="tableData.length"
          tag="ol"
          name="card-stagger"
          class="result-list"
        >
          <ExpertResultCard
            v-for="(item, index) in tableData"
            :key="item.ID"
            :style="{ '--stagger-index': index }"
            :item="item"
            :score-column-availability="scoreColumnAvailability"
            :keyword="lastKeyword"
            :score-label="compositeLabel"
            :has-keyword="hasKeyword"
            :rank="globalRank(index)"
          />
        </TransitionGroup>
        <el-empty
          v-else-if="!loading"
          description="没有命中的专家，换个关键词或放宽筛选条件试试"
          :image-size="80"
        />
        <div class="gva-pagination">
          <el-pagination
            :current-page="page"
            :page-size="pageSize"
            :page-sizes="[10, 20, 50]"
            :total="total"
            layout="total, sizes, prev, pager, next, jumper"
            @current-change="(v) => { page = v; runSearch() }"
            @size-change="(v) => { pageSize = v; runSearch() }"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { searchExpert, exportSearchResults, getScoreColumnAvailability, getFacetCounts } from '@/api/expertSearch'
import { ref, reactive, computed } from 'vue'
import FacetSidebar from './components/FacetSidebar.vue'
import AppliedFilterChips from './components/AppliedFilterChips.vue'
import ExpertResultCard from './components/ExpertResultCard.vue'

defineOptions({
  name: 'ExpertSearch'
})

// 检索条只有一个输入框，靠左边的下拉框决定这个词是当"全部字段"关键词（模糊相关性检索）还是
// 当"姓名"（精确查找）用——跟知网检索条"选字段+填词"的交互一致，不是姓名再单独摆一个输入框。
// 后端这两种查询本来就是完全不同的两条路径（keyword 走相关性排序，name 是精确匹配的结构化
// 过滤条件），下拉框选哪个决定 buildFilterParams 把 searchQuery 塞进 name 还是 keyword
const searchScope = ref('keyword')
const searchQuery = ref('')

// 原来所在单位/一级学科/专业技术职称/区域国别专长四个筛选框全部撤掉，改由左侧分面栏点选
// （facetSelection），点选后自动生效，不用点检索
const searchForm = reactive({
  sortBy: ''
})

// 分面选中值，四项都是字符串："" 表示这一维度没有筛选。单位分面的值是 org_id 的字符串形式，
// 其它三项是原始文本值——跟 FacetSidebar/后端 ExpertFacetCount.value 的口径一致
const facetSelection = reactive({
  orgId: '',
  disciplineL1: '',
  titleTier: '',
  regionExpertise: ''
})

const facetCounts = ref({ org: [], disciplineL1: [], titleTier: [], regionExpertise: [], total: 0 })

// 高级检索：多行且/或/不含，跟"检索主题"单关键词输入是二选一关系，不叠加。第一行没有连接词
// （起点），后面每行的 connector 取 and/or/not。advancedActive 只要有一行填了非空白词就算启用——
// 简单输入框届时禁用，避免两边都填了但只有一个生效的歧义状态
const advancedRows = ref([{ connector: '', term: '' }])
// 展开/收起状态——原来是 el-popover 自己管的 visible，现在面板是挂在检索条下面的常驻抽屉，
// 自己维护这个开关，点"高级检索"链接切换
const advancedPanelOpen = ref(false)
const advancedActive = computed(() => advancedRows.value.some((r) => r.term.trim()))

// 只有检索范围选的是"全部字段"时，高级检索才会跟这个输入框冲突（两边都是往 keyword 这条路径塞
// 词）；选"姓名"/"所在单位"时是完全独立的结构化过滤条件（后端 LIKE 匹配，不是全文相关性排序），
// 不冲突，不用禁用
const scopeInputPlaceholder = computed(() => {
  if (searchScope.value === 'name') return '输入姓名，支持部分匹配'
  if (searchScope.value === 'unitName') return '输入单位名称，支持部分匹配'
  return advancedActive.value
    ? '已启用高级检索，如需使用请先清空'
    : '输入研究方向或关键词，多个词用空格分隔可联合检索'
})

const addAdvancedRow = () => {
  advancedRows.value.push({ connector: 'and', term: '' })
}
const removeAdvancedRow = (idx) => {
  advancedRows.value.splice(idx, 1)
  if (!advancedRows.value.length) advancedRows.value.push({ connector: '', term: '' })
}
const applyAdvanced = () => {
  advancedPanelOpen.value = false
  page.value = 1
  runSearch()
}
const clearAdvanced = () => {
  advancedRows.value = [{ connector: '', term: '' }]
  page.value = 1
  runSearch()
}

// 拼接实际发给后端的筛选参数：facetSelection 里没选的维度直接不带这个 key，而不是传空字符串——
// orgId 在后端是 *uint，传空字符串会解析失败，其它几个维度的空字符串本身语义上就是"不筛"，
// 统一按不传处理更干净
const buildFilterParams = () => {
  const params = { ...searchForm, name: '', unitName: '', keyword: '' }
  // name/unitName 都是后端 LIKE 匹配的结构化过滤条件，advancedRows/keyword 是相关性排序的检索
  // 条件，两边是独立路径（applyStructuralFilters vs rankedCandidates 的排序逻辑），可以同时
  // 生效——检索范围选"姓名"/"所在单位"不影响高级检索继续起作用
  if (searchScope.value === 'name') {
    params.name = searchQuery.value.trim()
  } else if (searchScope.value === 'unitName') {
    params.unitName = searchQuery.value.trim()
  }
  if (advancedActive.value) {
    // AdvancedRows 非空时后端完全替代 Keyword，这里干脆不发 keyword，避免语义上的"到底哪个生效"
    const rows = advancedRows.value
      .filter((r) => r.term.trim())
      .map((r, i) => ({ connector: i === 0 ? '' : r.connector, term: r.term.trim() }))
    params.advancedRows = JSON.stringify(rows)
  } else if (searchScope.value === 'keyword') {
    params.keyword = searchQuery.value.trim()
  }
  if (facetSelection.orgId) params.orgId = facetSelection.orgId
  if (facetSelection.disciplineL1) params.disciplineL1 = facetSelection.disciplineL1
  if (facetSelection.titleTier) params.titleTier = facetSelection.titleTier
  if (facetSelection.regionExpertise) params.regionExpertise = facetSelection.regionExpertise
  return params
}

const tableData = ref([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const loading = ref(false)

// 决策影响分/社会贡献分这类分项，在批量导入的花名册数据还没配上决策影响/学术兼职记录之前，
// 全库对谁都是 0，摆一堆清一色的 0 没有信息量——查一次全库有没有人非零，没有就不显示这一项，
// 等真有数据了（哪怕只有一个人）就自动恢复显示，不用改代码
const scoreColumnAvailability = ref({ hasAchievement: true, hasInfluence: true, hasSocial: true })
getScoreColumnAvailability().then((res) => {
  if (res.code === 0) {
    scoreColumnAvailability.value = res.data
  }
})

// 排名跨页连续（第 2 页第 1 行是全局第 11 名）
const globalRank = (index) => (page.value - 1) * pageSize.value + index + 1

// 跟"当前输入框里打了什么"脱钩，只认"上一次真正提交检索时带没带关键词"——不然用户打字打到一半
// （还没点检索）排序控件就跟着闪现/消失，观感很跳。高级检索同理，用 lastAdvancedActive 记
// 上一次真正提交时是不是启用的，不是当前面板里填了啥
const lastKeyword = ref('')
const lastAdvancedActive = ref(false)
const hasKeyword = computed(() => !!lastKeyword.value.trim() || lastAdvancedActive.value)

const sortOptions = [
  { value: '', label: '综合实力' },
  { value: 'updatedAt', label: '最近更新' },
  { value: 'name', label: '姓名' }
]
const changeSortBy = (value) => {
  searchForm.sortBy = value
  page.value = 1
  runSearch()
}

const compositeLabel = computed(() => (hasKeyword.value ? '综合得分' : '综合实力'))

// 分面点击（含移除标签/清除全部）都是"点了立刻生效"，跟关键词"输入完点检索才生效"是两种不同的
// 交互性质——都合理，知网本身也是这样。筛选变了要把页码退回第一页，不然可能停在一个筛完之后
// 已经不存在的页码上
const onFacetChange = (next) => {
  Object.assign(facetSelection, next)
  page.value = 1
  runSearch()
}
const removeChip = (key) => {
  onFacetChange({ ...facetSelection, [key]: '' })
}
const clearAllFacets = () => {
  onFacetChange({ orgId: '', disciplineL1: '', titleTier: '', regionExpertise: '' })
}

// 已选条件标签：单位分面显示的是机构名不是 org_id，从当前拿到的 facetCounts.org 里反查一下
// label；其它三个维度的 value 本身就是要展示的文本
const orgLabel = computed(() => {
  const found = facetCounts.value.org.find((o) => o.value === facetSelection.orgId)
  return found ? found.label : facetSelection.orgId
})
const appliedChips = computed(() => {
  const chips = []
  if (facetSelection.orgId) chips.push({ key: 'orgId', label: `所在单位：${orgLabel.value}` })
  if (facetSelection.disciplineL1) chips.push({ key: 'disciplineL1', label: `一级学科：${facetSelection.disciplineL1}` })
  if (facetSelection.titleTier) chips.push({ key: 'titleTier', label: `职称档位：${facetSelection.titleTier}` })
  if (facetSelection.regionExpertise) chips.push({ key: 'regionExpertise', label: `区域/国别：${facetSelection.regionExpertise}` })
  return chips
})

// 检索结果和分面计数一起刷新，但两个请求的失败互不拖累：分面计数请求失败时侧栏保留上一次的
// 计数（不blank、不弹错误），不能因为分面统计接口抽风就连主检索结果都出不来
const runSearch = async() => {
  loading.value = true
  const params = buildFilterParams()
  const searchPromise = searchExpert({ ...params, page: page.value, pageSize: pageSize.value })
    .then((res) => {
      if (res.code === 0) {
        tableData.value = res.data.list
        total.value = res.data.total
        // 高级检索时摘要高亮沿用"不含"以外那几行的词拼起来（ExpertResultCard 按空白/标点分词），
        // 排除词本身不该被当成命中高亮；检索范围选的是"姓名"时是精确匹配，没有可高亮的模糊关键词
        lastKeyword.value = advancedActive.value
          ? advancedRows.value.filter((r) => r.term.trim() && r.connector !== 'not').map((r) => r.term.trim()).join(' ')
          : (searchScope.value === 'keyword' ? searchQuery.value : '')
        lastAdvancedActive.value = advancedActive.value
      }
    })
    .finally(() => {
      loading.value = false
    })
  const facetPromise = getFacetCounts(params)
    .then((res) => {
      if (res.code === 0) {
        // 无匹配记录时后端切片可能为 null，统一为空数组，已选单位标签也能安全查找名称。
        facetCounts.value = {
          org: res.data?.org ?? [],
          disciplineL1: res.data?.disciplineL1 ?? [],
          titleTier: res.data?.titleTier ?? [],
          regionExpertise: res.data?.regionExpertise ?? [],
          total: res.data?.total ?? 0
        }
      }
    })
    .catch(() => {})
  await Promise.all([searchPromise, facetPromise])
}

const onSearch = () => {
  page.value = 1
  runSearch()
}

const exporting = ref(false)
const handleExport = async() => {
  exporting.value = true
  try {
    // 导出跟当前检索条件完全相同的结果（含分面筛选），不带分页参数，后端按相关性/实时得分排好序
    // 返回全部命中记录
    const res = await exportSearchResults(buildFilterParams())
    const blob = new Blob([res], {
      type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'
    })
    const url = window.URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = '专家检索结果导出.xlsx'
    a.click()
    window.URL.revokeObjectURL(url)
  } finally {
    exporting.value = false
  }
}

const onReset = () => {
  searchScope.value = 'keyword'
  searchQuery.value = ''
  searchForm.sortBy = ''
  advancedRows.value = [{ connector: '', term: '' }]
  advancedPanelOpen.value = false
  Object.assign(facetSelection, { orgId: '', disciplineL1: '', titleTier: '', regionExpertise: '' })
  page.value = 1
  runSearch()
}

runSearch()
</script>

<style scoped>
/* 挂在检索条正下方的抽屉，跟检索条共用同一条左右边框、同一个圆角半径，中间 margin-top 用负值
   把两者的边框叠在一起，看起来是检索条本身长出来的一块，不是另一张浮动卡片 */
.advanced-panel {
  margin-top: -1px;
  background: var(--el-bg-color);
  border: 1px solid var(--parchment-border);
  border-top: 1px dashed var(--parchment-border);
  border-radius: 0 0 3px 3px;
  box-shadow: 0 4px 16px rgba(15, 32, 54, 0.06);
  padding: 4px 16px 12px;
}
.advanced-row {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 7px 0;
  border-bottom: 1px dashed var(--parchment-border);
}
.advanced-row:last-of-type {
  border-bottom: none;
}
.advanced-row-index {
  flex-shrink: 0;
  width: 18px;
  height: 18px;
  border-radius: 50%;
  display: grid;
  place-items: center;
  font-size: 10.5px;
  font-variant-numeric: tabular-nums;
  color: var(--el-text-color-placeholder);
  background: var(--parchment);
}
.advanced-connector {
  width: 84px;
  flex-shrink: 0;
}
.advanced-connector :deep(.el-select__wrapper) {
  box-shadow: none !important;
  background: var(--parchment);
  border-radius: 2px;
  min-height: 30px;
}
.advanced-connector-placeholder {
  width: 84px;
  flex-shrink: 0;
  font-size: 12px;
  color: var(--el-text-color-placeholder);
}
.advanced-term {
  flex: 1;
  min-width: 0;
}
.advanced-term :deep(.el-input__wrapper) {
  box-shadow: none !important;
  background: transparent;
  padding-left: 0;
}
.advanced-row-remove {
  flex-shrink: 0;
  width: 22px;
  height: 22px;
  display: grid;
  place-items: center;
  border: none;
  background: none;
  color: var(--el-text-color-placeholder);
  cursor: pointer;
}
.advanced-row-remove:hover:not(:disabled) {
  color: var(--brand-red);
}
.advanced-row-remove:disabled {
  opacity: 0.4;
  cursor: default;
}
.advanced-panel-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  margin-top: 8px;
  padding-top: 10px;
  border-top: 1px solid var(--parchment-border);
}
.advanced-add-row {
  border: none;
  background: none;
  padding: 0;
  font-size: 12.5px;
  color: var(--brand-red);
  cursor: pointer;
}
.advanced-add-row:hover {
  text-decoration: underline;
}
.advanced-panel-footer-actions {
  display: flex;
  align-items: center;
  gap: 14px;
}
/* "检索"按钮沿用主检索条提交按钮的样式语言（同一个红色实心矩形），让这个抽屉读起来是
   同一套检索工具的延伸，不是另一个独立弹窗里随便摆的按钮 */
.advanced-apply {
  border: none;
  background: var(--brand-red);
  color: #fff;
  font-size: 12.5px;
  font-weight: 600;
  padding: 6px 18px;
  border-radius: 2px;
  cursor: pointer;
  transition: background-color 120ms ease;
}
.advanced-apply:hover {
  background: #942026;
}

.advanced-expand-enter-active,
.advanced-expand-leave-active {
  transition: max-height 220ms ease, opacity 200ms ease;
  overflow: hidden;
}
.advanced-expand-enter-from,
.advanced-expand-leave-to {
  max-height: 0;
  opacity: 0;
}
.advanced-expand-enter-to,
.advanced-expand-leave-from {
  max-height: 400px;
  opacity: 1;
}
@media (prefers-reduced-motion: reduce) {
  .advanced-expand-enter-active,
  .advanced-expand-leave-active {
    transition: none;
  }
}

/* --brand-* 是这次会话里跟检索页/结果卡片/详情页共用的品牌色。--parchment/--parchment-border
   是内容区背景，跟登录页（src/view/login/index.vue）真正在用的冷淡青蓝配色（$ice: #eaf4fb）
   同一路，border-top 挪用的是登录页"顶部一条红线"这个具体的品牌签名元素——没有另起一块横幅，
   这条线本身足够让人认出"这是同一个系统"，不用靠一整块装饰性色块 */
.gva-table-box {
  --brand-navy: #16407a;
  --brand-navy-dark: #0c2a54;
  --brand-red: #b5272d;
  --brand-gold: #d4a017;
  --parchment: #f2f6fa;
  --parchment-border: #dbe3ee;
  background: var(--parchment);
  border-top: 3px solid var(--brand-red);
}

.result-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin: 8px 0 14px;
}
.result-count {
  font-size: 13px;
  color: var(--el-text-color-secondary);
}
.result-count b {
  color: var(--el-text-color-primary);
  font-weight: 600;
}

.sort-tabs {
  display: flex;
  gap: 2px;
  background: var(--el-bg-color);
  border: 1px solid var(--parchment-border);
  border-radius: 3px;
  padding: 3px;
}
.sort-tab {
  border: none;
  background: none;
  padding: 5px 12px;
  font-size: 12.5px;
  color: var(--el-text-color-secondary);
  border-radius: 2px;
  cursor: pointer;
  transition: background-color 160ms ease, color 160ms ease, box-shadow 160ms ease;
}
.sort-tab.active {
  /* 选中态用强调红，不是结构色藏青——跟原型的强调色分工保持一致 */
  background: var(--parchment);
  color: var(--brand-red);
  font-weight: 600;
}

.gva-search-box {
  border: 1px solid var(--parchment-border);
  border-radius: 3px;
}

/* 独立的一整块检索条：白底、圆角、投影，姓名/检索主题两个输入框和检索按钮拼成一个没有缝隙的
   胶囊容器——靠 el-input 自己的边框/投影全部去掉（只留容器整体的边框+投影）来做到"焊死在一起"
   的观感，跟 Element 默认松散排开的 el-form-item 完全是两种视觉逻辑 */
.search-bar {
  display: flex;
  align-items: stretch;
  height: 46px;
  background: var(--el-bg-color);
  border: 1px solid var(--parchment-border);
  border-radius: 3px;
  box-shadow: 0 1px 2px rgba(15, 32, 54, 0.05), 0 4px 16px rgba(15, 32, 54, 0.06);
  overflow: hidden;
}
/* 高级检索抽屉展开时，检索条下沿跟抽屉上沿拼在一起，检索条自己就不需要下圆角/下边框了，
   不然会跟抽屉的方角衔接处露出一小段多余的圆弧 */
.search-bar-attached {
  border-radius: 3px 3px 0 0;
  border-bottom-color: transparent;
}
.search-bar-input {
  flex-shrink: 0;
}
.search-bar-input :deep(.el-input__wrapper) {
  box-shadow: none !important;
  background: transparent;
  border-radius: 0;
  padding: 0 14px;
  height: 100%;
}
/* 左边的检索范围下拉框，去掉 el-select 自己的边框/投影，用暖色底跟右边的输入区分出
   "这是选择器不是输入框"，跟知网检索条左侧的字段选择器是同一个视觉逻辑 */
.search-bar-scope {
  flex-shrink: 0;
  width: 110px;
}
.search-bar-scope :deep(.el-select__wrapper) {
  box-shadow: none !important;
  background: var(--parchment);
  border-radius: 0;
  height: 100%;
  min-height: 0;
  padding: 0 10px;
}
.search-bar-keyword {
  flex: 1;
  min-width: 0;
}
.search-bar-divider {
  width: 1px;
  flex-shrink: 0;
  background: var(--parchment-border);
  margin: 9px 0;
}
.search-bar-submit {
  flex-shrink: 0;
  border: none;
  background: var(--brand-red);
  color: #fff;
  padding: 0 26px;
  font-size: 14px;
  font-weight: 600;
  letter-spacing: 0.02em;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 6px;
  transition: background-color 120ms ease;
}
.search-bar-submit:hover {
  background: #942026;
}

.search-bar-actions {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 2px 2px;
  font-size: 12.5px;
}
.search-bar-action {
  border: none;
  background: none;
  padding: 0;
  font-size: 12.5px;
  color: var(--el-text-color-secondary);
  cursor: pointer;
}
.search-bar-action:hover {
  color: var(--brand-red);
}
.search-bar-action.is-active {
  color: var(--brand-red);
  font-weight: 600;
}
.search-bar-action:disabled {
  cursor: default;
  opacity: 0.6;
}
.search-bar-action-sep {
  color: var(--parchment-border);
  font-size: 12px;
}

.search-layout {
  display: grid;
  grid-template-columns: 240px 1fr;
  gap: 20px;
  align-items: start;
}

/* 分页器原来完全是 Element 默认蓝（.gva-pagination .is-active 全局规则读 var(--el-color-primary)），
   跟页面其它地方已经换成的藏青/朱红完全脱节，最容易被一眼看出"没换肤"。这里局部覆盖掉当前页
   高亮的颜色，不动全局 el-color-primary（别的页面的分页器还是默认蓝） */
.results-col :deep(.gva-pagination .is-active) {
  background: var(--brand-navy) !important;
}
@media (max-width: 860px) {
  .search-layout {
    grid-template-columns: 1fr;
  }
}

.results-col {
  min-width: 0;
}
.result-list {
  list-style: none;
  margin: 0;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

/* 结果卡片入场带一点点错位延迟——不是纯装饰，是让"筛完/搜完之后结果整批换新了"这件事被看见。
   延迟按 --stagger-index（ExpertResultCard.vue 传下来的行内 style）算，封顶到第 10 张卡片，
   避免一页 50 条时最后几张卡拖太久才出现 */
.card-stagger-enter-active {
  transition: opacity 260ms ease, transform 260ms ease;
  transition-delay: calc(min(var(--stagger-index, 0), 10) * 28ms);
}
.card-stagger-enter-from {
  opacity: 0;
  transform: translateY(8px);
}
.card-stagger-move {
  transition: transform 220ms ease;
}
@media (prefers-reduced-motion: reduce) {
  .card-stagger-enter-active,
  .card-stagger-move {
    transition: none !important;
  }
}
</style>
