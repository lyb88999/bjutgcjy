<template>
    <div class="chart-container">
      <!-- 当数据加载完成后渲染图表 -->
      <v-chart v-if="majors.length > 0" class="chart" :option="satisfactionComparisonOption" />
    </div>
  </template>
  
  <script setup>
  import { use } from 'echarts/core';
  import { CanvasRenderer } from 'echarts/renderers';
  import { BarChart } from 'echarts/charts';
  import { TitleComponent, TooltipComponent, LegendComponent, GridComponent } from 'echarts/components';
  import VChart, { THEME_KEY } from 'vue-echarts';
  import { ref, provide, computed } from 'vue';
  
  use([
    CanvasRenderer,
    BarChart,
    TitleComponent,
    TooltipComponent,
    LegendComponent,
    GridComponent
  ]);
  
  provide(THEME_KEY, 'black');
  
  const majors = ref([
    {
      majorName: "机械工程",
      universityName: "清华大学",
      plannedEnrollment: 20,
      isFirstClassMajor: true,
      isNewEngineeringDiscipline: true,
      majorSatisfaction: "满意"
    },
    {
      majorName: "电子科学与技术",
      universityName: "上海交通大学",
      plannedEnrollment: 15,
      isFirstClassMajor: true,
      isNewEngineeringDiscipline: false,
      majorSatisfaction: "非常满意"
    },
    {
      majorName: "计算机科学与技术",
      universityName: "浙江大学",
      plannedEnrollment: 8,
      isFirstClassMajor: false,
      isNewEngineeringDiscipline: false,
      majorSatisfaction: "中等"
    },
    // ...其他模拟数据
  ]);
  
  const satisfactionComparisonOption = computed(() => ({
    title: { text: '专业满意度学校对比', left: 'center', top: '5%' },
    tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } },
    legend: { left: 'right', top: '5%', data: ['满意', '非常满意', '中等'] },
    xAxis: { type: 'category', data: majors.value.map(s => s.universityName) },
    yAxis: { type: 'value' },
    series: [
      {
        name: '满意',
        type: 'bar',
        stack: 'total',
        data: majors.value.map(s => s.majorSatisfaction === '满意' ? 1 : 0)
      },
      {
        name: '非常满意',
        type: 'bar',
        stack: 'total',
        data: majors.value.map(s => s.majorSatisfaction === '非常满意' ? 1 : 0)
      },
      {
        name: '中等',
        type: 'bar',
        stack: 'total',
        data: majors.value.map(s => s.majorSatisfaction === '中等' ? 1 : 0)
      }
    ]
  }));
  
  </script>
  
  <style scoped>
  .chart-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 20px;
    margin-top: 10%;
  }
  
  .chart {
    width: 100%;
    max-width: 800px;
    height: 400px;
  }
  </style>
  