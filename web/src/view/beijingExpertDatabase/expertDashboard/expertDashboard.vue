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
          <span class="panel-sub">全库 {{ statusTotal }} 条</span>
        </div>
        <div
          ref="statusChartRef"
          class="chart-box"
        />
      </div>
      <div class="panel">
        <div class="panel-header">
          <span class="panel-title">最近30天新增趋势</span>
          <span class="panel-sub">共 {{ trendTotal }} 条</span>
        </div>
        <div
          ref="trendChartRef"
          class="chart-box"
        />
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
        <RankingList
          :items="achievementLevelsOrdered"
          :show-rank="false"
          :bar-color="levelColor"
        />
      </div>
      <div class="panel">
        <div class="panel-header">
          <span class="panel-title">决策影响单位级别分布</span>
          <span class="panel-sub">已发布专家名下全部决策影响记录</span>
        </div>
        <RankingList
          :items="adoptionLevelsOrdered"
          :show-rank="false"
          :bar-color="levelColor"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { getDashboardStats } from '@/api/expertDashboard'
import * as echarts from 'echarts'
import { computed, nextTick, onMounted, onUnmounted, ref } from 'vue'
import { useWindowResize } from '@/hooks/use-windows-resize'
import {
  Medal,
  OfficeBuilding,
  CircleCheck,
  Warning,
  CollectionTag,
  TrendCharts,
  Trophy
} from '@element-plus/icons-vue'
import RankingList from './RankingList.vue'

defineOptions({
  name: 'ExpertDashboard'
})

const stats = ref({
  totalPublished: 0,
  totalDraft: 0,
  pendingOrgReview: 0,
  pendingCityReview: 0,
  reviewBypassed: 0,
  statusBreakdown: [],
  orgBreakdown: [],
  disciplineBreakdown: [],
  recentTrend: [],
  coveredOrgCount: 0,
  taggedRatio: 0,
  avgCompositeScore: 0,
  seniorTitleRatio: 0,
  achievementLevelBreakdown: [],
  adoptionLevelBreakdown: []
})

const round1 = (v) => Math.round((v || 0) * 10) / 10

const statusTotal = computed(() =>
  stats.value.statusBreakdown.reduce((sum, s) => sum + s.count, 0)
)
const publishedRatio = computed(() => {
  if (!statusTotal.value) return 0
  return ((stats.value.totalPublished / statusTotal.value) * 100).toFixed(1)
})
const trendTotal = computed(() =>
  stats.value.recentTrend.reduce((sum, p) => sum + p.count, 0)
)

// ---- 配色 ----
// 审核状态是"状态"而不是普通类别，按语义配色：草稿=灰紫（中性）、待审=琥珀、退回=红系、
// 已发布=绿。顺序与后端 statusOrder 固定一致。明暗两套是分别选定并跑过
// 色觉障碍/亮度带/对比度校验的（dataviz 校验脚本），不是简单反色
const STATUS_COLORS_LIGHT = ['#6272B3', '#D9940A', '#F5222D', '#1677FF', '#C41D7F', '#389E0D']
const STATUS_COLORS_DARK = ['#7B8AD9', '#C77F06', '#D6264D', '#1677FF', '#C41D7F', '#389E0D']

// 成果/采纳级别是序数量级：同一蓝色从深到浅对应国家级→一般级，未标注用灰色（数值旁边都印着
// 具体数字和占比，颜色只是冗余编码）
const LEVEL_ORDER = ['国家级', '省部级', '厅局级', '一般级', '未标注']
const LEVEL_COLORS = ['#0958D9', '#1677FF', '#69B1FF', '#ADC6FF', '#C0C4CC']

const isDark = () => document.documentElement.classList.contains('dark')
const statusPalette = () => (isDark() ? STATUS_COLORS_DARK : STATUS_COLORS_LIGHT)
const cssVar = (name, fallback) => {
  const v = getComputedStyle(document.documentElement).getPropertyValue(name).trim()
  return v || fallback
}

const orderByLevel = (items) => {
  const rank = (label) => {
    const i = LEVEL_ORDER.indexOf(label)
    return i === -1 ? LEVEL_ORDER.length : i
  }
  return [...(items || [])].sort((a, b) => rank(a.label) - rank(b.label))
}
const achievementLevelsOrdered = computed(() => orderByLevel(stats.value.achievementLevelBreakdown))
const adoptionLevelsOrdered = computed(() => orderByLevel(stats.value.adoptionLevelBreakdown))
const levelColor = (item) => {
  const i = LEVEL_ORDER.indexOf(item.label)
  return LEVEL_COLORS[i === -1 ? LEVEL_COLORS.length - 1 : i]
}

const statusChartRef = ref(null)
const trendChartRef = ref(null)
let statusChart, trendChart

useWindowResize(() => {
  statusChart && statusChart.resize()
  trendChart && trendChart.resize()
})

const renderStatusChart = () => {
  if (!statusChartRef.value) return
  statusChart = statusChart || echarts.init(statusChartRef.value)
  const textColor = cssVar('--el-text-color-regular', '#606266')
  const mutedColor = cssVar('--el-text-color-secondary', '#909399')
  const surface = cssVar('--el-bg-color', '#ffffff')
  statusChart.setOption({
    color: statusPalette(),
    tooltip: { trigger: 'item', formatter: '{b}：{c} 条（{d}%）' },
    legend: {
      bottom: 0,
      itemWidth: 10,
      itemHeight: 10,
      textStyle: { color: textColor }
    },
    title: {
      text: String(statusTotal.value),
      subtext: '全库档案',
      left: 'center',
      top: '38%',
      itemGap: 4,
      textStyle: { fontSize: 28, fontWeight: 700, color: textColor },
      subtextStyle: { fontSize: 12, color: mutedColor }
    },
    series: [
      {
        type: 'pie',
        radius: ['52%', '72%'],
        avoidLabelOverlap: true,
        itemStyle: { borderRadius: 6, borderColor: surface, borderWidth: 2 },
        label: { formatter: '{b} {d}%', color: textColor },
        labelLine: { lineStyle: { color: mutedColor }},
        data: stats.value.statusBreakdown.map((item) => ({
          name: item.label,
          value: item.count
        }))
      }
    ]
  })
}

const renderTrendChart = () => {
  if (!trendChartRef.value) return
  trendChart = trendChart || echarts.init(trendChartRef.value)
  const textColor = cssVar('--el-text-color-secondary', '#909399')
  const gridColor = cssVar('--el-border-color-lighter', '#ebeef5')
  const brand = isDark() ? '#4096FF' : '#1677FF'
  trendChart.setOption({
    tooltip: {
      trigger: 'axis',
      axisPointer: { type: 'line' },
      formatter: (params) => {
        const p = params[0]
        return `${p.axisValue}<br/>新增 <b>${p.data}</b> 条`
      }
    },
    grid: { left: 40, right: 20, top: 20, bottom: 40 },
    xAxis: {
      type: 'category',
      data: stats.value.recentTrend.map((item) => item.date.slice(5)),
      boundaryGap: false,
      axisLabel: { color: textColor },
      axisLine: { lineStyle: { color: gridColor }}
    },
    yAxis: {
      type: 'value',
      minInterval: 1,
      axisLabel: { color: textColor },
      splitLine: { lineStyle: { type: 'dashed', color: gridColor }}
    },
    series: [
      {
        type: 'line',
        smooth: true,
        symbol: 'circle',
        symbolSize: 7,
        showSymbol: false,
        lineStyle: { width: 2.5, color: brand },
        itemStyle: { color: brand },
        emphasis: { itemStyle: { borderColor: brand, borderWidth: 2 }},
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: isDark() ? 'rgba(64,150,255,0.30)' : 'rgba(22,119,255,0.28)' },
            { offset: 1, color: 'rgba(22,119,255,0.02)' }
          ])
        },
        data: stats.value.recentTrend.map((item) => item.count)
      }
    ]
  })
}

const renderCharts = () => {
  renderStatusChart()
  renderTrendChart()
}

// 明暗主题切换时用当前主题的配色和文字色重画图表，而不是停留在旧主题的颜色上
let themeObserver
const watchTheme = () => {
  themeObserver = new MutationObserver(() => renderCharts())
  themeObserver.observe(document.documentElement, {
    attributes: true,
    attributeFilter: ['class']
  })
}

const loadData = async() => {
  const res = await getDashboardStats()
  if (res.code === 0) {
    stats.value = res.data
    await nextTick()
    renderCharts()
  }
}

onMounted(() => {
  loadData()
  watchTheme()
})
onUnmounted(() => {
  themeObserver && themeObserver.disconnect()
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
  font-variant-numeric: tabular-nums;
}
.kpi-foot {
  font-size: 12px;
  opacity: 0.85;
  margin-top: 8px;
}

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
  font-variant-numeric: tabular-nums;
}
.quality-label {
  font-size: 12px;
  color: var(--el-text-color-secondary);
}

/* ---- 图表面板 ---- */
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
  position: relative;
  padding-left: 10px;
}
.panel-title::before {
  content: '';
  position: absolute;
  left: 0;
  top: 50%;
  transform: translateY(-50%);
  width: 3px;
  height: 14px;
  border-radius: 2px;
  background: linear-gradient(180deg, #1677ff, #69b1ff);
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
