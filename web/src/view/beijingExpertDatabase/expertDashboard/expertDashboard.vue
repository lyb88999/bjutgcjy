<template>
  <div class="expert-dashboard">
    <el-row :gutter="16" class="stat-cards">
      <el-col :span="6">
        <el-card shadow="never"><div class="stat-label">已发布专家数</div><div class="stat-value">{{ stats.totalPublished }}</div></el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="never"><div class="stat-label">待单位审核</div><div class="stat-value">{{ stats.pendingOrgReview }}</div></el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="never"><div class="stat-label">待市级审核</div><div class="stat-value">{{ stats.pendingCityReview }}</div></el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="never"><div class="stat-label">免审核发布</div><div class="stat-value">{{ stats.reviewBypassed }}</div></el-card>
      </el-col>
    </el-row>

    <el-row :gutter="16" class="chart-row">
      <el-col :span="12">
        <el-card shadow="never">
          <div class="chart-title">审核状态分布</div>
          <div ref="statusChartRef" class="chart-box" />
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card shadow="never">
          <div class="chart-title">最近30天新增趋势</div>
          <div ref="trendChartRef" class="chart-box" />
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="16" class="chart-row">
      <el-col :span="12">
        <el-card shadow="never">
          <div class="chart-title">按单位分布（已发布，Top 10）</div>
          <div ref="orgChartRef" class="chart-box" />
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card shadow="never">
          <div class="chart-title">按一级学科分布（已发布，Top 10）</div>
          <div ref="disciplineChartRef" class="chart-box" />
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { getDashboardStats } from '@/api/expertDashboard'
import * as echarts from 'echarts'
import { nextTick, onMounted, onUnmounted, ref } from 'vue'
import { useWindowResize } from '@/hooks/use-windows-resize'

defineOptions({
  name: 'ExpertDashboard'
})

const stats = ref({
  totalPublished: 0, totalDraft: 0, pendingOrgReview: 0, pendingCityReview: 0, reviewBypassed: 0,
  statusBreakdown: [], orgBreakdown: [], disciplineBreakdown: [], recentTrend: []
})

const statusChartRef = ref(null)
const trendChartRef = ref(null)
const orgChartRef = ref(null)
const disciplineChartRef = ref(null)
let statusChart, trendChart, orgChart, disciplineChart

useWindowResize(() => {
  statusChart && statusChart.resize()
  trendChart && trendChart.resize()
  orgChart && orgChart.resize()
  disciplineChart && disciplineChart.resize()
})

const renderStatusChart = () => {
  statusChart = echarts.init(statusChartRef.value)
  statusChart.setOption({
    tooltip: { trigger: 'item' },
    legend: { bottom: 0 },
    series: [{
      type: 'pie',
      radius: ['40%', '70%'],
      data: stats.value.statusBreakdown.map(item => ({ name: item.label, value: item.count }))
    }]
  })
}

const renderTrendChart = () => {
  trendChart = echarts.init(trendChartRef.value)
  trendChart.setOption({
    tooltip: { trigger: 'axis' },
    grid: { left: 40, right: 20, top: 20, bottom: 40 },
    xAxis: { type: 'category', data: stats.value.recentTrend.map(item => item.date.slice(5)) },
    yAxis: { type: 'value', minInterval: 1 },
    series: [{ type: 'line', smooth: true, areaStyle: {}, data: stats.value.recentTrend.map(item => item.count) }]
  })
}

const renderOrgChart = () => {
  orgChart = echarts.init(orgChartRef.value)
  const data = [...stats.value.orgBreakdown].reverse()
  orgChart.setOption({
    tooltip: { trigger: 'axis' },
    grid: { left: 120, right: 20, top: 10, bottom: 10 },
    xAxis: { type: 'value', minInterval: 1 },
    yAxis: { type: 'category', data: data.map(item => item.label) },
    series: [{ type: 'bar', data: data.map(item => item.count) }]
  })
}

const renderDisciplineChart = () => {
  disciplineChart = echarts.init(disciplineChartRef.value)
  const data = [...stats.value.disciplineBreakdown].reverse()
  disciplineChart.setOption({
    tooltip: { trigger: 'axis' },
    grid: { left: 100, right: 20, top: 10, bottom: 10 },
    xAxis: { type: 'value', minInterval: 1 },
    yAxis: { type: 'category', data: data.map(item => item.label) },
    series: [{ type: 'bar', data: data.map(item => item.count) }]
  })
}

const loadData = async () => {
  const res = await getDashboardStats()
  if (res.code === 0) {
    stats.value = res.data
    await nextTick()
    renderStatusChart()
    renderTrendChart()
    renderOrgChart()
    renderDisciplineChart()
  }
}

onMounted(loadData)
onUnmounted(() => {
  statusChart && statusChart.dispose()
  trendChart && trendChart.dispose()
  orgChart && orgChart.dispose()
  disciplineChart && disciplineChart.dispose()
})
</script>

<style scoped>
.expert-dashboard {
  padding: 4px;
}
.stat-cards .stat-label {
  color: #909399;
  font-size: 13px;
}
.stat-cards .stat-value {
  font-size: 28px;
  font-weight: 600;
  margin-top: 8px;
}
.chart-row {
  margin-top: 16px;
}
.chart-title {
  font-size: 14px;
  font-weight: 500;
  margin-bottom: 8px;
}
.chart-box {
  width: 100%;
  height: 280px;
}
</style>
