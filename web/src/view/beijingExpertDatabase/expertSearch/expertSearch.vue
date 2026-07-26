<template>
  <div class="gva-table-box">
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchForm" @keyup.enter="onSearch">
        <el-form-item label="检索主题">
          <el-input v-model="searchForm.keyword" placeholder="如：防汛方案" style="width: 220px" clearable />
        </el-form-item>
        <el-form-item label="一级学科">
          <el-input v-model="searchForm.disciplineL1" placeholder="筛选条件" clearable />
        </el-form-item>
        <el-form-item label="区域/国别专长">
          <el-input v-model="searchForm.regionExpertise" placeholder="筛选条件" clearable />
        </el-form-item>
        <el-form-item label="专业技术职称">
          <el-input v-model="searchForm.techTitle" placeholder="筛选条件" clearable />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSearch">检索</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
          <el-button icon="download" :loading="exporting" @click="handleExport">导出当前结果</el-button>
        </el-form-item>
      </el-form>
    </div>

    <el-table :data="tableData" style="width: 100%" v-loading="loading">
      <el-table-column type="index" label="排名" width="60" />
      <el-table-column align="left" label="姓名" prop="name" width="110" />
      <el-table-column align="left" label="所在单位" prop="unitName" width="160" />
      <el-table-column align="left" label="专业技术职称" prop="techTitle" width="120" />
      <el-table-column align="left" label="一级学科" prop="disciplineL1" width="120" />
      <el-table-column align="left" label="研究方向" prop="researchDirections" min-width="180" show-overflow-tooltip />
      <el-table-column align="left" label="相关性" prop="relevance" width="90">
        <template #default="scope">{{ (scope.row.relevance * 100).toFixed(0) }}%</template>
      </el-table-column>
      <el-table-column align="left" label="成果分" prop="achievementScore" width="90" />
      <el-table-column align="left" label="决策影响分" prop="influenceScore" width="100" />
      <el-table-column align="left" label="社会贡献分" prop="socialScore" width="100" />
      <el-table-column align="left" label="综合得分" prop="realtimeScore" width="100">
        <template #default="scope">{{ scope.row.realtimeScore.toFixed(2) }}</template>
      </el-table-column>
    </el-table>
    <div class="gva-pagination">
      <el-pagination
        :current-page="page"
        :page-size="pageSize"
        :page-sizes="[10, 20, 50]"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        @current-change="v => { page = v; onSearch() }"
        @size-change="v => { pageSize = v; onSearch() }"
      />
    </div>
    <p class="search-hint">综合得分 = w1×成果得分×相关性 + w2×决策影响得分×相关性 + w3×职称权重 + w4×社会贡献得分，权重可在"系统工具-字典管理"里调整（字典类型 expert_ranking_weight）。</p>
  </div>
</template>

<script setup>
import { searchExpert, exportSearchResults } from '@/api/expertSearch'
import { ref, reactive } from 'vue'

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

const onSearch = async () => {
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
const handleExport = async () => {
  exporting.value = true
  try {
    // 导出跟当前检索条件完全相同的结果，不带分页参数，后端按相关性/实时得分排好序返回全部命中记录
    const res = await exportSearchResults(searchForm)
    const blob = new Blob([res], { type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' })
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
.search-hint {
  margin-top: 12px;
  font-size: 12px;
  color: #909399;
}
</style>
