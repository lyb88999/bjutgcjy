<template>
  <div class="gva-table-box">
    <div class="gva-search-box">
      <el-form
        :inline="true"
        :model="searchForm"
        @keyup.enter="onSearch"
      >
        <el-form-item label="姓名">
          <el-input
            v-model="searchForm.name"
            placeholder="按姓名精确查找"
            style="width: 140px"
            clearable
          />
        </el-form-item>
        <el-form-item label="所在单位">
          <el-input
            v-model="searchForm.unitName"
            placeholder="按单位查找"
            style="width: 160px"
            clearable
          />
        </el-form-item>
        <el-form-item label="检索主题">
          <el-input
            v-model="searchForm.keyword"
            placeholder="如：防汛方案，多个词用空格分隔可联合检索"
            style="width: 220px"
            clearable
          />
        </el-form-item>
        <el-form-item label="一级学科">
          <el-select
            v-model="searchForm.disciplineL1"
            clearable
            filterable
            allow-create
            default-first-option
            placeholder="选择或输入"
            style="width: 160px"
          >
            <el-option
              v-for="d in disciplineOptions"
              :key="d"
              :label="d"
              :value="d"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="区域/国别专长">
          <el-input
            v-model="searchForm.regionExpertise"
            placeholder="筛选条件"
            clearable
          />
        </el-form-item>
        <el-form-item label="专业技术职称">
          <el-select
            v-model="searchForm.techTitle"
            clearable
            filterable
            placeholder="请选择"
            style="width: 140px"
          >
            <el-option
              v-for="t in titleOptions"
              :key="t"
              :label="t"
              :value="t"
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button
            type="primary"
            icon="search"
            @click="onSearch"
          >检索</el-button>
          <el-button
            icon="refresh"
            @click="onReset"
          >重置</el-button>
          <el-button
            icon="download"
            :loading="exporting"
            @click="handleExport"
          >导出当前结果</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="result-toolbar">
      <span class="result-count">共 <b>{{ total }}</b> 位专家</span>
      <el-select
        v-if="!hasKeyword"
        v-model="searchForm.sortBy"
        size="small"
        class="sort-select"
        @change="() => { page = 1; onSearch() }"
      >
        <el-option
          label="按综合实力排序"
          value=""
        />
        <el-option
          label="按最近更新排序"
          value="updatedAt"
        />
        <el-option
          label="按姓名排序"
          value="name"
        />
      </el-select>
    </div>

    <el-table
      v-loading="loading"
      :data="tableData"
      style="width: 100%"
    >
      <el-table-column
        label="排名"
        width="64"
        align="center"
      >
        <template #default="scope">
          <span
            class="rank-badge"
            :class="rankClass(globalRank(scope.$index))"
          >{{ globalRank(scope.$index) }}</span>
        </template>
      </el-table-column>
      <el-table-column
        align="left"
        label="姓名"
        width="140"
      >
        <template #default="scope">
          <div class="expert-name">{{ scope.row.name }}</div>
          <div
            v-if="scope.row.techTitle"
            class="expert-title"
          >{{ scope.row.techTitle }}</div>
        </template>
      </el-table-column>
      <el-table-column
        align="left"
        label="所在单位"
        prop="unitName"
        min-width="160"
        show-overflow-tooltip
      />
      <el-table-column
        align="left"
        label="一级学科"
        width="110"
      >
        <template #default="scope">
          <el-tag
            v-if="scope.row.disciplineL1"
            type="info"
            effect="plain"
            size="small"
          >{{ scope.row.disciplineL1 }}</el-tag>
          <span v-else>—</span>
        </template>
      </el-table-column>
      <el-table-column
        align="left"
        label="研究方向"
        min-width="220"
      >
        <template #default="scope">
          <div
            class="research-directions"
            :title="scope.row.researchDirections"
          >{{ scope.row.researchDirections }}</div>
          <div
            v-if="scope.row.matchReason"
            class="match-reason"
            :title="scope.row.matchReason"
          >命中：{{ scope.row.matchReason }}</div>
        </template>
      </el-table-column>
      <el-table-column
        v-if="hasKeyword"
        align="left"
        label="相关性"
        width="120"
      >
        <template #default="scope">
          <div class="score-cell">
            <div class="score-track">
              <div
                class="score-fill score-fill-relevance"
                :style="{ width: scope.row.relevance * 100 + '%' }"
              />
            </div>
            <span class="score-num">{{ (scope.row.relevance * 100).toFixed(0) }}%</span>
          </div>
        </template>
      </el-table-column>
      <el-table-column
        align="left"
        label="分项明细"
        width="190"
      >
        <template #default="scope">
          <div class="sub-scores">
            <span
              v-if="scoreColumnAvailability.hasAchievement"
              class="sub-chip"
              title="成果分"
            ><b>成果</b>{{ fmt(scope.row.achievementScore) }}</span>
            <span
              v-if="scoreColumnAvailability.hasInfluence"
              class="sub-chip"
              title="决策影响分"
            ><b>影响</b>{{ fmt(scope.row.influenceScore) }}</span>
            <span
              v-if="scoreColumnAvailability.hasSocial"
              class="sub-chip"
              title="社会贡献分"
            ><b>贡献</b>{{ fmt(scope.row.socialScore) }}</span>
          </div>
        </template>
      </el-table-column>
      <el-table-column
        align="left"
        :label="compositeLabel"
        width="160"
      >
        <template #header>
          <span class="col-header-tip">
            {{ compositeLabel }}
            <el-tooltip
              effect="dark"
              placement="top"
              :content="compositeTooltip"
            >
              <el-icon class="tip-icon"><QuestionFilled /></el-icon>
            </el-tooltip>
          </span>
        </template>
        <template #default="scope">
          <div class="score-cell">
            <div class="score-track">
              <div
                class="score-fill score-fill-composite"
                :style="{ width: scoreWidth(scope.row.realtimeScore) + '%' }"
              />
            </div>
            <span class="score-num score-num-strong">{{ scope.row.realtimeScore.toFixed(2) }}</span>
          </div>
        </template>
      </el-table-column>
      <el-table-column
        align="center"
        label="操作"
        fixed="right"
        width="90"
      >
        <template #default="scope">
          <el-button
            type="primary"
            link
            @click="viewDetail(scope.row)"
          >查看详情</el-button>
        </template>
      </el-table-column>
      <template #empty>
        <el-empty
          description="没有命中的专家，换个关键词或放宽筛选条件试试"
          :image-size="80"
        />
      </template>
    </el-table>
    <div class="gva-pagination">
      <el-pagination
        :current-page="page"
        :page-size="pageSize"
        :page-sizes="[10, 20, 50]"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        @current-change="(v) => { page = v; onSearch() }"
        @size-change="(v) => { pageSize = v; onSearch() }"
      />
    </div>
  </div>
</template>

<script setup>
import { searchExpert, exportSearchResults, getScoreColumnAvailability } from '@/api/expertSearch'
import { getExpertTagList } from '@/api/expertTag'
import { getDict } from '@/utils/dictionary'
import { ref, reactive, computed } from 'vue'
import { useRouter } from 'vue-router'

defineOptions({
  name: 'ExpertSearch'
})

const router = useRouter()
const viewDetail = (row) => {
  router.push({ name: 'expertProfileDetail', params: { id: row.ID }})
}

// 筛选选项与主档编辑表单同源：职称来自打分权重字典，一级学科来自标签库，
// 用户不用猜库里有什么值
const titleOptions = ref([])
getDict('expert_title_level').then((items) => {
  titleOptions.value = (items || []).map((i) => i.label)
})
const disciplineOptions = ref([])
getExpertTagList({ page: 1, pageSize: 500 }).then((res) => {
  if (res.code === 0) {
    disciplineOptions.value = (res.data.list || [])
      .filter((t) => t.tagType === 'discipline_l1')
      .map((t) => t.tagValue)
  }
})

const searchForm = reactive({
  name: '',
  unitName: '',
  keyword: '',
  disciplineL1: '',
  regionExpertise: '',
  techTitle: '',
  sortBy: ''
})

const tableData = ref([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const loading = ref(false)

const fmt = (v) => Math.round((v || 0) * 10) / 10

// 决策影响分/社会贡献分这类分项，在批量导入的花名册数据还没配上决策影响/学术兼职记录之前，
// 全库对谁都是 0，摆一列清一色的 0 没有信息量——查一次全库有没有人非零，没有就不显示这一项，
// 等真有数据了（哪怕只有一个人）就自动恢复显示，不用改代码
const scoreColumnAvailability = ref({ hasAchievement: true, hasInfluence: true, hasSocial: true })
getScoreColumnAvailability().then((res) => {
  if (res.code === 0) {
    scoreColumnAvailability.value = res.data
  }
})

// 排名跨页连续（第 2 页第 1 行是全局第 11 名），金银铜徽章只属于全局前三
const globalRank = (index) => (page.value - 1) * pageSize.value + index + 1
const rankClass = (rank) =>
  rank === 1 ? 'rank-gold' : rank === 2 ? 'rank-silver' : rank === 3 ? 'rank-bronze' : 'rank-plain'

// 综合得分条按当前结果集里的最高分归一化，第一名满条，其余按比例，扫一眼就能看出分差有多大
const maxScore = computed(() =>
  tableData.value.reduce((m, r) => Math.max(m, r.realtimeScore || 0), 0)
)
const scoreWidth = (score) => (maxScore.value ? ((score || 0) / maxScore.value) * 100 : 0)

// 跟"当前输入框里打了什么"脱钩，只认"上一次真正提交检索时带没带关键词"——不然用户打字打到一半
// （还没点检索）相关性列就跟着闪现/消失，观感很跳
const lastKeyword = ref('')
const hasKeyword = computed(() => !!lastKeyword.value.trim())

const compositeLabel = computed(() => (hasKeyword.value ? '综合得分' : '综合实力'))
const compositeTooltip = computed(() =>
  hasKeyword.value
    ? '综合得分 = 相关性 ×（成果分×w1 + 决策影响分×w2 + 职称权重×w3 + 社会贡献分×w4 + w5），反映的是与本次检索词的匹配程度，权重可在"系统工具-字典管理"（expert_ranking_weight）调整'
    : '未输入检索词时相关性恒为 1，这里显示的是专家整体实力（成果、决策影响、职称、社会贡献的加权和），不代表和某个主题相关，权重可在"系统工具-字典管理"（expert_ranking_weight）调整'
)

const onSearch = async() => {
  loading.value = true
  try {
    const res = await searchExpert({ ...searchForm, page: page.value, pageSize: pageSize.value })
    if (res.code === 0) {
      tableData.value = res.data.list
      total.value = res.data.total
      lastKeyword.value = searchForm.keyword
    }
  } finally {
    loading.value = false
  }
}

const exporting = ref(false)
const handleExport = async() => {
  exporting.value = true
  try {
    // 导出跟当前检索条件完全相同的结果，不带分页参数，后端按相关性/实时得分排好序返回全部命中记录
    const res = await exportSearchResults(searchForm)
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
  searchForm.name = ''
  searchForm.unitName = ''
  searchForm.keyword = ''
  searchForm.disciplineL1 = ''
  searchForm.regionExpertise = ''
  searchForm.techTitle = ''
  searchForm.sortBy = ''
  page.value = 1
  onSearch()
}

onSearch()
</script>

<style scoped>
.result-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}
.result-count {
  font-size: 13px;
  color: var(--el-text-color-secondary);
}
.result-count b {
  color: var(--el-text-color-primary);
  font-weight: 600;
}
.sort-select {
  width: 150px;
}

.rank-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  border-radius: 50%;
  font-size: 12px;
  font-weight: 700;
  color: #fff;
}
.rank-gold { background: linear-gradient(135deg, #ffd666, #fa8c16); }
.rank-silver { background: linear-gradient(135deg, #d9d9d9, #8c8c8c); }
.rank-bronze { background: linear-gradient(135deg, #ffc069, #ad6800); }
.rank-plain { background: var(--el-fill-color); color: var(--el-text-color-secondary); }

.expert-name {
  font-weight: 600;
  color: var(--el-text-color-primary);
}
.expert-title {
  font-size: 12px;
  color: var(--el-text-color-secondary);
  margin-top: 2px;
}

.research-directions {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.match-reason {
  margin-top: 2px;
  font-size: 12px;
  color: var(--el-color-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.score-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}
.score-track {
  flex: 1;
  height: 7px;
  border-radius: 4px;
  background: var(--el-fill-color-light);
  overflow: hidden;
}
.score-fill {
  height: 100%;
  border-radius: 4px;
}
.score-fill-relevance { background: linear-gradient(90deg, #95c2ff, #69b1ff); }
.score-fill-composite { background: linear-gradient(90deg, #1677ff, #69b1ff); }
.score-num {
  min-width: 36px;
  text-align: right;
  font-size: 12px;
  color: var(--el-text-color-regular);
  font-variant-numeric: tabular-nums;
}
.score-num-strong {
  font-size: 13px;
  font-weight: 700;
  color: var(--el-color-primary);
  min-width: 44px;
}

.sub-scores {
  display: flex;
  flex-wrap: wrap;
  gap: 4px 6px;
}
.sub-chip {
  display: inline-flex;
  align-items: baseline;
  gap: 3px;
  padding: 1px 7px;
  border-radius: 10px;
  background: var(--el-fill-color-light);
  font-size: 12px;
  color: var(--el-text-color-regular);
  font-variant-numeric: tabular-nums;
  white-space: nowrap;
}
.sub-chip b {
  font-size: 11px;
  font-weight: 400;
  color: var(--el-text-color-secondary);
}

.col-header-tip {
  display: inline-flex;
  align-items: center;
  gap: 4px;
}
.tip-icon {
  font-size: 14px;
  color: var(--el-text-color-secondary);
  cursor: help;
}
</style>
