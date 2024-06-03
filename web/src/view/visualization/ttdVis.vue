<template>
    <div class="chart-container">
      <!-- Ensure charts render only after data is loaded -->
      <v-chart v-if="schools.length > 0 && selectedChart === 'students'" class="chart" :option="studentsOption" />
      <v-chart v-else-if="schools.length > 0 && selectedChart === 'educationQuality'" class="chart" :option="educationQualityOption" />
      <v-chart v-else-if="schools.length > 0 && selectedChart === 'textbookAndAwards'" class="chart" :option="textbookAndAwardsOption" />
      <v-chart v-else-if="schools.length > 0 && selectedChart === 'entrepreneurshipAndTransfer'" class="chart" :option="entrepreneurshipAndTransferOption" />
  
      <label for="chart-select">选择图表：</label>
      <select v-model="selectedChart" id="chart-select">
        <option value="students">学生数量与就业率</option>
        <option value="educationQuality">教育质量及成就</option>
        <option value="textbookAndAwards">教材建设与学生获奖</option>
        <option value="entrepreneurshipAndTransfer">学生创业项目与转专业</option>
      </select>
    </div>
  </template>
  
  <script setup>
  import { use } from 'echarts/core';
  import { CanvasRenderer } from 'echarts/renderers';
  import { BarChart, LineChart } from 'echarts/charts';
  import { TitleComponent, TooltipComponent, LegendComponent, GridComponent } from 'echarts/components';
  import VChart, { THEME_KEY } from 'vue-echarts';
  import { ref, provide, onMounted, computed } from 'vue';
  import { getTalentTrainingDatabaseList } from '@/api/talentTrainingDatabase';
  
  use([
    CanvasRenderer,
    BarChart,
    LineChart,
    TitleComponent,
    TooltipComponent,
    LegendComponent,
    GridComponent
  ]);
  
  provide(THEME_KEY, 'black');
  
  const schools = ref([]);
  
  // Function to fetch data from the backend
  const fetchSchoolData = async () => {
    try {
      const res = await getTalentTrainingDatabaseList({ page: 1, pageSize: 100 });
      if (res.code === 0) {
        schools.value = res.data.list.map(school => ({
          universityName: school.universityName,
          undergraduateCount: school.undergraduateCount,
          masterStudentCount: school.masterStudentCount,
          doctoralStudentCount: school.doctoralStudentCount,
          nationalFirstClassCourseCount: school.nationalFirstClassCourseCount,
          provincialFirstClassCourseCount: school.provincialFirstClassCourseCount,
          nationalTeachingDemonstrationCenterCount: school.ationalTeachingDemonstrationCenterCount,
          provincialTeachingDemonstrationCenterCount: school.provincialTeachingDemonstrationCenterCount,
          nationalPlanningTextbookCount: school.nationalPlanningTextbookCount,
          nationalExcellentTextbookCount: school.nationalExcellentTextbookCount,
          provincialExcellentTextbookCount: school.provincialExcellentTextbookCount,
          internationalAwardsCount: school.internationalAwardsCount,
          nationalAwardsCount: school.nationalAwardsCount,
          provincialAwardsCount: school.provincialAwardsCount,
          undergraduateEmploymentRate: school.undergraduateEmploymentRate,
          masterEmploymentRate: school.masterEmploymentRate,
          doctoralEmploymentRate: school.doctoralEmploymentRate,
          projectCount: school.projectCount,
          participatingStudentCount: school.participatingStudentCount,
          transferredIntoEngineeringCount: school.transferredIntoEngineeringCount,
          transferredOutOfEngineeringCount: school.transferredOutOfEngineeringCount
        }));
      } else {
        console.error('Error fetching data: ', res.message);
      }
    } catch (error) {
      console.error('Error during data fetch: ', error);
    }
  };
  
  onMounted(fetchSchoolData);
  
  const selectedChart = ref('students');
  
  // 学生数量与就业率
  const studentsOption = computed(() => ({
    title: { text: '学生数量与就业率', left: 'center', top: '5%' },
    tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } },
    legend: { left: 'right', data: ['本科生', '硕士生', '博士生', '本科生就业率', '硕士生就业率', '博士生就业率'] },
    xAxis: {
      type: 'category',
      data: schools.value.map(s => s.universityName)
    },
    yAxis: [
      { type: 'value', name: '学生数量' },
      { type: 'value', name: '就业率', axisLabel: { formatter: '{value} %' }, max: 100 }
    ],
    series: [
      { name: '本科生', type: 'bar', data: schools.value.map(s => s.undergraduateCount) },
      { name: '硕士生', type: 'bar', data: schools.value.map(s => s.masterStudentCount) },
      { name: '博士生', type: 'bar', data: schools.value.map(s => s.doctoralStudentCount) },
      { name: '本科生就业率', type: 'line', yAxisIndex: 1, data: schools.value.map(s => s.undergraduateEmploymentRate * 100) },
      { name: '硕士生就业率', type: 'line', yAxisIndex: 1, data: schools.value.map(s => s.masterEmploymentRate * 100) },
      { name: '博士生就业率', type: 'line', yAxisIndex: 1, data: schools.value.map(s => s.doctoralEmploymentRate * 100) }
    ]
  }));
  
  // 教育质量及成就
  const educationQualityOption = computed(() => ({
    title: { text: '教育质量及成就', left: 'center', top: '0%' },
    tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } },
    legend: { left: 'right', top: '6%', data: ['国家级一流课程', '省级一流课程', '国家级教育示范中心', '省级教育示范中心'] },
    xAxis: {
      type: 'category',
      data: schools.value.map(s => s.universityName)
    },
    yAxis: { type: 'value' },
    series: [
      { name: '国家级一流课程', type: 'bar', data: schools.value.map(s => s.nationalFirstClassCourseCount) },
      { name: '省级一流课程', type: 'bar', data: schools.value.map(s => s.provincialFirstClassCourseCount) },
      { name: '国家级教育示范中心', type: 'bar', data: schools.value.map(s => s.nationalTeachingDemonstrationCenterCount) },
      { name: '省级教育示范中心', type: 'bar', data: schools.value.map(s => s.provincialTeachingDemonstrationCenterCount) }
    ]
  }));
  
  // 教材建设与学生获奖
  const textbookAndAwardsOption = computed(() => ({
    title: { text: '教材建设与学生获奖', left: 'center', top: '0%' },
    tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } },
    legend: { left: 'right', top: '5%', data: ['国家级规划教材数', '全国优秀教材数', '省级优秀教材数', '学生获国际奖项数', '学生获国家级奖项数', '学生获省部级奖项数'] },
    xAxis: {
      type: 'category',
      data: schools.value.map(s => s.universityName)
    },
    yAxis: { type: 'value' },
    series: [
      { name: '国家级规划教材数', type: 'bar', data: schools.value.map(s => s.nationalPlanningTextbookCount) },
      { name: '全国优秀教材数', type: 'bar', data: schools.value.map(s => s.nationalExcellentTextbookCount) },
      { name: '省级优秀教材数', type: 'bar', data: schools.value.map(s => s.provincialExcellentTextbookCount) },
      { name: '学生获国际奖项数', type: 'bar', data: schools.value.map(s => s.internationalAwardsCount) },
      { name: '学生获国家级奖项数', type: 'bar', data: schools.value.map(s => s.nationalAwardsCount) },
      { name: '学生获省部级奖项数', type: 'bar', data: schools.value.map(s => s.provincialAwardsCount) }
    ]
  }));
  
  // 学生创业项目与转专业
  const entrepreneurshipAndTransferOption = computed(() => ({
    title: { text: '学生创业项目与转专业', left: 'center', top: '0%' },
    tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } },
    legend: { left: 'right', top: '5%', data: ['项目数', '参与学生数', '转入工科学生数', '转出工科学生数'] },
    xAxis: {
      type: 'category',
      data: schools.value.map(s => s.universityName)
    },
    yAxis: { type: 'value' },
    series: [
      { name: '项目数', type: 'bar', data: schools.value.map(s => s.projectCount) },
      { name: '参与学生数', type: 'bar', data: schools.value.map(s => s.participatingStudentCount) },
      { name: '转入工科学生数', type: 'bar', data: schools.value.map(s => s.transferredIntoEngineeringCount) },
      { name: '转出工科学生数', type: 'bar', data: schools.value.map(s => s.transferredOutOfEngineeringCount) }
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
  