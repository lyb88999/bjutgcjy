<template>
  <div class="expert-dashboard">
    <div class="kpi-grid">
      <div class="kpi-card kpi-blue">
        <div class="kpi-top">
          <span class="kpi-label">已发布专家数</span>
          <span class="kpi-icon"><el-icon><Medal /></el-icon></span>
        </div>
        <div class="kpi-value">{{ stats.totalPublished }}</div>
        <div class="kpi-foot">占全库 {{ publishedRatio }}%</div>
      </div>
      <div class="kpi-card kpi-orange">
        <div class="kpi-top">
          <span class="kpi-label">待单位审核</span>
          <span class="kpi-icon"><el-icon><OfficeBuilding /></el-icon></span>
        </div>
        <div class="kpi-value">{{ stats.pendingOrgReview }}</div>
        <div class="kpi-foot">等本单位审核员处理</div>
      </div>
      <div class="kpi-card kpi-red">
        <div class="kpi-top">
          <span class="kpi-label">待市级审核</span>
          <span class="kpi-icon"><el-icon><CircleCheck /></el-icon></span>
        </div>
        <div class="kpi-value">{{ stats.pendingCityReview }}</div>
        <div class="kpi-foot">等市级审核员处理</div>
      </div>
      <div class="kpi-card kpi-purple">
        <div class="kpi-top">
          <span class="kpi-label">免审核发布</span>
          <span class="kpi-icon"><el-icon><Warning /></el-icon></span>
        </div>
        <div class="kpi-value">{{ stats.reviewBypassed }}</div>
        <div class="kpi-foot">审核开关关闭期间收录</div>
      </div>
    </div>

    <div class="quality-grid">
      <div class="quality-card">
        <span class="quality-icon quality-icon-cyan"><el-icon><OfficeBuilding /></el-icon></span>
        <div>
          <div class="quality-value">{{ stats.coveredOrgCount }}</div>
          <div class="quality-label">覆盖单位数</div>
        </div>
      </div>
      <div class="quality-card">
        <span class="quality-icon quality-icon-green"><el-icon><CollectionTag /></el-icon></span>
        <div>
          <div class="quality-value">{{ round1(stats.taggedRatio) }}%</div>
          <div class="quality-label">已发布专家打标签比例</div>
        </div>
      </div>
      <div class="quality-card">
        <span class="quality-icon quality-icon-blue"><el-icon><TrendCharts /></el-icon></span>
        <div>
          <div class="quality-value">{{ round1(stats.avgCompositeScore) }}</div>
          <div class="quality-label">平均综合排序得分</div>
        </div>
      </div>
      <div class="quality-card">
        <span class="quality-icon quality-icon-gold"><el-icon><Trophy /></el-icon></span>
        <div>
          <div class="quality-value">{{ round1(stats.seniorTitleRatio) }}%</div>
          <div class="quality-label">教授/研究员占比</div>
        </div>
      </div>
    </div>

    <div class="panel-grid">
      <div class="panel">
        <div class="panel-header">
          <span class="panel-title">审核状态分布</span>
        </div>
        <div ref="statusChartRef" class="chart-box" />
      </div>
      <div class="panel">
        <div class="panel-header">
          <span class="panel-title">最近30天新增趋势</span>
          <span class="panel-sub">共 {{ trendTotal }} 条</span>
        </div>
        <div ref="trendChartRef" class="chart-box" />
      </div>
    </div>

    <div class="panel-grid">
      <div class="panel">
        <div class="panel-header">
          <span class="panel-title">按单位分布</span>
          <span class="panel-sub">已发布 · Top 10</span>
        </div>
        <RankingList :items="stats.orgBreakdown" />
      </div>
      <div class="panel">
        <div class="panel-header">
          <span class="panel-title">按一级学科分布</span>
          <span class="panel-sub">已发布 · Top 10</span>
        </div>
        <RankingList :items="stats.disciplineBreakdown" />
      </div>
    </div>

    <div class="panel-grid">
      <div class="panel">
        <div class="panel-header">
          <span class="panel-title">研究成果级别分布</span>
          <span class="panel-sub">已发布专家名下全部成果</span>
        </div>
        <RankingList :items="stats.achievementLevelBreakdown" />
      </div>
      <div class="panel">
        <div class="panel-header">
          <span class="panel-title">决策影响单位级别分布</span>
          <span class="panel-sub">已发布专家名下全部决策影响记录</span>
        </div>
        <RankingList :items="stats.adoptionLevelBreakdown" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { getDashboardStats } from '@/api/expertDashboard'
import * as echarts from 'echarts'
import { computed, nextTick, onMounted, onUnmounted, ref } from 'vue'
import { useWindowResize } from '@/hooks/use-windows-resize'
import { Medal, OfficeBuilding, CircleCheck, Warning, CollectionTag, TrendCharts, Trophy } from '@element-plus/icons-vue'
import RankingList from './RankingList.vue'

defineOptions({
  name: 'ExpertDashboard'
})

const stats = ref({
  totalPublished: 0, totalDraft: 0, pendingOrgReview: 0, pendingCityReview: 0, reviewBypassed: 0,
  statusBreakdown: [], orgBreakdown: [], disciplineBreakdown: [], recentTrend: [],
  coveredOrgCount: 0, taggedRatio: 0, avgCompositeScore: 0, seniorTitleRatio: 0,
  achievementLevelBreakdown: [], adoptionLevelBreakdown: []
})

const round1 = (v) => Math.round((v || 0) * 10) / 10

const publishedRatio = computed(() => {
  const total = stats.value.statusBreakdown.reduce((sum, s) => sum + s.count, 0)
  if (!total) return 0
  return ((stats.value.totalPublished / total) * 100).toFixed(1)
})
const trendTotal = computed(() => stats.value.recentTrend.reduce((sum, p) => sum + p.count, 0))

// 跟 Ant Design 默认图表色板对齐，用在饼图/趋势图上，跟整体的"高级感"配色保持一致
const chartPalette = ['#1677FF', '#52C41A', '#FAAD14', '#F5222D', '#722ED1', '#13C2C2']

const statusChartRef = ref(null)
const trendChartRef = ref(null)
let statusChart, trendChart

useWindowResize(() => {
  statusChart && statusChart.resize()
  trendChart && trendChart.resize()
})

const renderStatusChart = () => {
  statusChart = echarts.init(statusChartRef.value)
  statusChart.setOption({
    color: chartPalette,
    tooltip: { trigger: 'item' },
    legend: { bottom: 0, itemWidth: 10, itemHeight: 10 },
    series: [{
      type: 'pie',
      radius: ['48%', '72%'],
      avoidLabelOverlap: true,
      itemStyle: { borderRadius: 6, borderColor: '#fff', borderWidth: 2 },
      label: { formatter: '{b}\n{d}%' },
      data: stats.value.statusBreakdown.map(item => ({ name: item.label, value: item.count }))
    }]
  })
}

const renderTrendChart = () => {
  trendChart = echarts.init(trendChartRef.value)
  trendChart.setOption({
    color: chartPalette,
    tooltip: { trigger: 'axis' },
    grid: { left: 40, right: 20, top: 20, bottom: 40 },
    xAxis: { type: 'category', data: stats.value.recentTrend.map(item => item.date.slice(5)), boundaryGap: false },
    yAxis: { type: 'value', minInterval: 1, splitLine: { lineStyle: { type: 'dashed' } } },
    series: [{
      type: 'line',
      smooth: true,
      symbol: 'circle',
      symbolSize: 6,
      lineStyle: { width: 3 },
      areaStyle: {
        color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
          { offset: 0, color: 'rgba(22,119,255,0.35)' },
          { offset: 1, color: 'rgba(22,119,255,0.02)' }
        ])
      },
      data: stats.value.recentTrend.map(item => item.count)
    }]
  })
}

const loadData = async () => {
  const res = await getDashboardStats()
  if (res.code === 0) {
    stats.value = res.data
    await nextTick()
    renderStatusChart()
    renderTrendChart()
  }
}

onMounted(loadData)
onUnmounted(() => {
  statusChart && statusChart.dispose()
  trendChart && trendChart.dispose()
})
</script>

<style scoped>
.expert-dashboard {
  padding: 4px;
}

/* ---- KPI 卡片 ---- */
.kpi-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
}
.kpi-card {
  border-radius: 12px;
  padding: 18px 20px;
  color: #fff;
  position: relative;
  overflow: hidden;
  box-shadow: 0 4px 14px rgba(0, 0, 0, 0.08);
  transition: transform 0.2s;
}
.kpi-card:hover {
  transform: translateY(-2px);
}
.kpi-blue { background: linear-gradient(135deg, #3b82f6, #1d4ed8); }
.kpi-orange { background: linear-gradient(135deg, #fb923c, #ea580c); }
.kpi-red { background: linear-gradient(135deg, #f87171, #dc2626); }
.kpi-purple { background: linear-gradient(135deg, #a78bfa, #7c3aed); }
.kpi-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.kpi-label {
  font-size: 13px;
  opacity: 0.9;
}
.kpi-icon {
  width: 30px;
  height: 30px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.22);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
}
.kpi-value {
  font-size: 34px;
  font-weight: 700;
  margin-top: 10px;
  line-height: 1.2;
}
.kpi-foot {
  font-size: 12px;
  opacity: 0.85;
  margin-top: 8px;
}

/* ---- 图表面板 ---- */
/* ---- 内容质量指标 ---- */
.quality-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
  margin-top: 16px;
}
.quality-card {
  display: flex;
  align-items: center;
  gap: 12px;
  background: var(--el-bg-color);
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 12px;
  padding: 16px 18px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}
.quality-icon {
  width: 38px;
  height: 38px;
  min-width: 38px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  color: #fff;
}
.quality-icon-cyan { background: linear-gradient(135deg, #22d3ee, #0891b2); }
.quality-icon-green { background: linear-gradient(135deg, #4ade80, #16a34a); }
.quality-icon-blue { background: linear-gradient(135deg, #60a5fa, #2563eb); }
.quality-icon-gold { background: linear-gradient(135deg, #fbbf24, #d97706); }
.quality-value {
  font-size: 20px;
  font-weight: 700;
  line-height: 1.3;
}
.quality-label {
  font-size: 12px;
  color: var(--el-text-color-secondary);
}

.panel-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
  margin-top: 16px;
}
.panel {
  background: var(--el-bg-color);
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 12px;
  padding: 16px 20px 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}
.panel-header {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  margin-bottom: 12px;
}
.panel-title {
  font-size: 15px;
  font-weight: 600;
}
.panel-sub {
  font-size: 12px;
  color: var(--el-text-color-placeholder);
}
.chart-box {
  width: 100%;
  height: 280px;
}

@media (max-width: 1200px) {
  .kpi-grid { grid-template-columns: repeat(2, 1fr); }
  .quality-grid { grid-template-columns: repeat(2, 1fr); }
  .panel-grid { grid-template-columns: 1fr; }
}
</style>
