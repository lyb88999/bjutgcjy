<template>
    <div class="chart-container">
      <!-- Ensure charts render only after data is loaded -->
      <v-chart v-if="schools.length > 0 && selectedChart === 'internationalData'" class="chart" :option="internationalOption" />
      <label for="chart-select">选择图表：</label>
      <select v-model="selectedChart" id="chart-select">
        <option value="internationalData">国际化数据分析</option>
      </select>
    </div>
  </template>
  
  <script setup>
  import { use } from 'echarts/core';
  import { CanvasRenderer } from 'echarts/renderers';
  import { BarChart } from 'echarts/charts';
  import { TitleComponent, TooltipComponent, LegendComponent, GridComponent } from 'echarts/components';
  import VChart, { THEME_KEY } from 'vue-echarts';
  import { ref, provide, computed, onMounted } from 'vue';
  import { getInternationalExchangeDatabaseList } from '@/api/internationalExchangeDatabase';
  
  use([
    CanvasRenderer,
    BarChart,
    TitleComponent,
    TooltipComponent,
    LegendComponent,
    GridComponent
  ]);
  
  provide(THEME_KEY, 'black');
  
  const schools = ref([]);
  
  // 使用 onMounted 钩子从后端接口获取数据
  onMounted(async () => {
    try {
      const res = await getInternationalExchangeDatabaseList({ page: 1, pageSize: 100 });
      if (res.code === 0) {
        schools.value = res.data.list.map(school => ({
          universityName: school.universityName,
          internationalConferencesHosted: school.internationalConferencesHosted,
          foreignFacultyCount: school.foreignFacultyCount,
          internationalStudentsCount: school.internationalStudentsCount,
          studentsAbroadExchangeCount: school.studentsAbroadExchangeCount
        }));
      } else {
        console.error('Error fetching data: ', res.message);
      }
    } catch (error) {
      console.error('Error during data fetch: ', error);
    }
  });
  
  const selectedChart = ref('internationalData');
  
  const internationalOption = computed(() => ({
    title: { text: '国际化数据分析', left: 'center', top: '0%' },
    tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } },
    legend: { left: 'right', data: ['举办的国际会议数', '外籍教师数', '来华留学生数', '出国交流学生数'], top: '5%' },
    xAxis: { type: 'category', data: schools.value.map(s => s.universityName) },
    yAxis: { type: 'value' },
    series: [
      {
        name: '举办的国际会议数',
        type: 'bar',
        data: schools.value.map(s => s.internationalConferencesHosted)
      },
      {
        name: '外籍教师数',
        type: 'bar',
        data: schools.value.map(s => s.foreignFacultyCount)
      },
      {
        name: '来华留学生数',
        type: 'bar',
        data: schools.value.map(s => s.internationalStudentsCount)
      },
      {
        name: '出国交流学生数',
        type: 'bar',
        data: schools.value.map(s => s.studentsAbroadExchangeCount)
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
  