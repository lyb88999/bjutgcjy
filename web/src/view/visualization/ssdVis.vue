<template>
    <div class="chart-container">
      <!-- 当数据加载完成后渲染图表 -->
      <v-chart v-if="schools.length > 0 && selectedChart === 'innovationData'" class="chart" :option="innovationOption" />
      <label for="chart-select">选择图表：</label>
      <select v-model="selectedChart" id="chart-select">
        <option value="innovationData">协同创新中心分析</option>
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
  import { getSocialServiceDatabaseList } from '@/api/socialServiceDatabase';
  
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

  const fetchSchoolData = async () => {
    try {
      const res = await getSocialServiceDatabaseList({ page: 1, pageSize: 100 });
      if (res.code === 0) {
        schools.value = res.data.list.map(school => ({
          universityName: school.universityName,
          industryInnovationCenterCount: school.industryInnovationCenterCount,
          regionalInnovationCenterCount: school.regionalInnovationCenterCount
        }));
      } else {
        console.error('Error fetching data: ', res.message);
      }
    } catch (error) {
      console.error('Error during data fetch: ', error);
    }
  };
  
  // 使用 onMounted 钩子从后端接口获取数据
  onMounted(fetchSchoolData);
  
  const selectedChart = ref('innovationData');
  
  const innovationOption = computed(() => ({
    title: { text: '协同创新中心分析', left: 'center', top: '0%' },
    tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } },
    legend: { left: 'right', top: '5%', data: ['行业协同创新中心数', '区域协同创新中心数'] },
    xAxis: { type: 'category', data: schools.value.map(s => s.universityName) },
    yAxis: { type: 'value' },
    series: [
      {
        name: '行业协同创新中心数',
        type: 'bar',
        data: schools.value.map(s => s.industryInnovationCenterCount)
      },
      {
        name: '区域协同创新中心数',
        type: 'bar',
        data: schools.value.map(s => s.regionalInnovationCenterCount)
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
  