<template>
  <div class="gva-table-box">
    <div class="gva-search-box">
      <el-form
        :inline="true"
        :model="searchForm"
        @keyup.enter="onSearch"
      >
        <el-form-item label="检索主题">
          <el-input
            v-model="searchForm.keyword"
            placeholder="如：防汛方案"
            style="width: 220px"
            clearable
          />
        </el-form-item>
        <el-form-item label="一级学科">
          <el-input
            v-model="searchForm.disciplineL1"
            placeholder="筛选条件"
            clearable
          />
        </el-form-item>
        <el-form-item label="区域/国别专长">
          <el-input
            v-model="searchForm.regionExpertise"
            placeholder="筛选条件"
            clearable
          />
        </el-form-item>
        <el-form-item label="专业技术职称">
          <el-input
            v-model="searchForm.techTitle"
            placeholder="筛选条件"
            clearable
          />
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
        prop="researchDirections"
        min-width="180"
        show-overflow-tooltip
      />
      <el-table-column
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
        align="right"
        label="成果分"
        width="80"
      >
        <template #default="scope">
          <span class="sub-score">{{ fmt(scope.row.achievementScore) }}</span>
        </template>
      </el-table-column>
      <el-table-column
        align="right"
        label="决策影响分"
        width="96"
      >
        <template #default="scope">
          <span class="sub-score">{{ fmt(scope.row.influenceScore) }}</span>
        </template>
      </el-table-column>
      <el-table-column
        align="right"
        label="社会贡献分"
        width="96"
      >
        <template #default="scope">
          <span class="sub-score">{{ fmt(scope.row.socialScore) }}</span>
        </template>
      </el-table-column>
      <el-table-column
        align="left"
        label="综合得分"
        width="150"
      >
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
    <p class="search-hint">
      综合得分 = w1×成果得分×相关性 + w2×决策影响得分×相关性 + w3×职称权重 +
      w4×社会贡献得分，权重可在"系统工具-字典管理"里调整（字典类型 expert_ranking_weight）。
    </p>
  </div>
</template>

<script setup>
import { searchExpert, exportSearchResults } from '@/api/expertSearch'
import { ref, reactive, computed } from 'vue'

defineOptions({
  name: 'ExpertSearch'
})

const searchForm = reactive({
  keyword: '',
  disciplineL1: '',
  regionExpertise: '',
  techTitle: ''
})

const tableData = ref([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const loading = ref(false)

const fmt = (v) => Math.round((v || 0) * 10) / 10

// 排名跨页连续（第 2 页第 1 行是全局第 11 名），金银铜徽章只属于全局前三
const globalRank = (index) => (page.value - 1) * pageSize.value + index + 1
const rankClass = (rank) =>
  rank === 1 ? 'rank-gold' : rank === 2 ? 'rank-silver' : rank === 3 ? 'rank-bronze' : 'rank-plain'

// 综合得分条按当前结果集里的最高分归一化，第一名满条，其余按比例，扫一眼就能看出分差有多大
const maxScore = computed(() =>
  tableData.value.reduce((m, r) => Math.max(m, r.realtimeScore || 0), 0)
)
const scoreWidth = (score) => (maxScore.value ? ((score || 0) / maxScore.value) * 100 : 0)

const onSearch = async() => {
  loading.value = true
  try {
    const res = await searchExpert({ ...searchForm, page: page.value, pageSize: pageSize.value })
    if (res.code === 0) {
      tableData.value = res.data.list
      total.value = res.data.total
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
  searchForm.keyword = ''
  searchForm.disciplineL1 = ''
  searchForm.regionExpertise = ''
  searchForm.techTitle = ''
  page.value = 1
  onSearch()
}

onSearch()
</script>

<style scoped>
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

.score-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}
.score-track {
  flex: 1;
  height: 6px;
  border-radius: 3px;
  background: var(--el-fill-color-light);
  overflow: hidden;
}
.score-fill {
  height: 100%;
  border-radius: 3px;
}
.score-fill-relevance { background: #69b1ff; }
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
.sub-score {
  font-variant-numeric: tabular-nums;
  color: var(--el-text-color-regular);
}

.search-hint {
  margin-top: 12px;
  font-size: 12px;
  color: var(--el-text-color-secondary);
}
</style>
