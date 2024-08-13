<template>
  <div class="chart-container">
    <!-- 只有当 schools 数据已加载时，才渲染图表 -->
    <v-chart v-if="schools.length > 0 && selectedChart === 'student'" class="chart" :option="studentOption" />
    <v-chart v-else-if="schools.length > 0 && selectedChart === 'staff'" class="chart" :option="staffOption" />
    <v-chart v-else-if="schools.length > 0 && selectedChart === 'program'" class="chart" :option="programOption" />
    <label for="chart-select">选择图表：</label>
    <select v-model="selectedChart" id="chart-select">
      <option value="student">学生数量</option>
      <option value="staff">教职员工数量</option>
      <option value="program">学术项目</option>
    </select>
  </div>
</template>

<script setup>
import { use } from 'echarts/core';
import { CanvasRenderer } from 'echarts/renderers';
import { BarChart } from 'echarts/charts';
import { TitleComponent, TooltipComponent, LegendComponent, GridComponent } from 'echarts/components';
import VChart, { THEME_KEY } from 'vue-echarts';
import { ref, provide, onMounted, computed } from 'vue';
import { getBasicInfomationDatabaseList } from '@/api/basicInformationDatabase';

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

// 查询
const getTableData = async () => {
  try {
    const res = await getBasicInfomationDatabaseList({ page: 1, pageSize: 100 });
    if (res.code === 0) {
      schools.value = res.data.list.map(school => ({
        schoolName: school.schoolName,
        undergraduateNumber: school.undergraduateNumber,
        graduateStudentNumber: school.graduateStudentNumber,
        doctoralStudentNumber: school.doctoralStudentNumber,
        internationalStudentNumber: school.internationalStudentNumber,
        facultyStaffNumber: school.facultyStaffNumber,
        fullTimeTeacherNumber: school.fullTimeTeacherNumber,
        undergraduateProgramNumber: school.undergraduateProgramNumber,
        mastersAutorizationPointerNumber: school.mastersAutorizationPointerNumber,
        doctoralAuthorizationPointerNumber: school.doctoralAuthorizationPointerNumber,
        mastersDegreeCategoryNumber: school.mastersDegreeCategoryNumber,
        doctoralDegreeCategoryNumber: school.doctoralDegreeCategoryNumber,
        postDoctoralMobileStationNumber: school.postDoctoralMobileStationNumber
      }));
    } else {
      console.error("Error fetching data: ", res.message);
    }
  } catch (error) {
    console.error("Error during data fetch: ", error);
  }
}

onMounted(getTableData);

const selectedChart = ref('student');

// Use computed properties to ensure chart options update with data changes
const studentOption = computed(() => ({
  title: { text: '学生数量', left: 'center', top: '5%' },
  tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } },
  legend: { left: 'right', data: ['本科生', '硕士生', '博士生', '留学生'] },
  grid: { left: '3%', right: '4%', bottom: '15%', containLabel: true }, // 增加底部间距
  xAxis: { 
    type: 'category', 
    data: schools.value.map(s => s.schoolName),
    axisLabel: {
      rotate: 30, // 旋转30度
      fontStyle: 'italic', // 字体斜体
    }
  },
  yAxis: { type: 'value' },
  series: [
    { name: '本科生', type: 'bar', data: schools.value.map(s => s.undergraduateNumber) },
    { name: '硕士生', type: 'bar', data: schools.value.map(s => s.graduateStudentNumber) },
    { name: '博士生', type: 'bar', data: schools.value.map(s => s.doctoralStudentNumber) },
    { name: '留学生', type: 'bar', data: schools.value.map(s => s.internationalStudentNumber) }
  ]
}));

const staffOption = computed(() => ({
  title: { text: '教职员工数量', left: 'center', top: '5%' },
  tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } },
  legend: { left: 'right', data: ['教职员工', '专任教师'] },
  grid: { left: '3%', right: '4%', bottom: '15%', containLabel: true }, // 增加底部间距
  xAxis: { 
    type: 'category', 
    data: schools.value.map(s => s.schoolName),
    axisLabel: {
      rotate: 30, // 旋转30度
      fontStyle: 'italic', // 字体斜体
    }
  },
  yAxis: { type: 'value' },
  series: [
    { name: '教职员工', type: 'bar', data: schools.value.map(s => s.facultyStaffNumber) },
    { name: '专任教师', type: 'bar', data: schools.value.map(s => s.fullTimeTeacherNumber) }
  ]
}));

const programOption = computed(() => ({
  title: { text: '学术项目', left: 'center', top: '5%' },
  tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } },
  legend: { left: 'right', data: ['本科专业', '硕士授权点', '博士授权点', '硕士专业学位', '博士专业学位', '博士后流动站'] },
  grid: { left: '3%', right: '4%', bottom: '15%', containLabel: true }, // 增加底部间距
  xAxis: { 
    type: 'category', 
    data: schools.value.map(s => s.schoolName),
    axisLabel: {
      rotate: 30, // 旋转30度
      fontStyle: 'italic', // 字体斜体
    }
  },
  yAxis: { type: 'value' },
  series: [
    { name: '本科专业', type: 'bar', data: schools.value.map(s => s.undergraduateProgramNumber) },
    { name: '硕士授权点', type: 'bar', data: schools.value.map(s => s.mastersAutorizationPointerNumber) },
    { name: '博士授权点', type: 'bar', data: schools.value.map(s => s.doctoralAuthorizationPointerNumber) },
    { name: '硕士专业学位', type: 'bar', data: schools.value.map(s => s.mastersDegreeCategoryNumber) },
    { name: '博士专业学位', type: 'bar', data: schools.value.map(s => s.doctoralDegreeCategoryNumber) },
    { name: '博士后流动站', type: 'bar', data: schools.value.map(s => s.postDoctoralMobileStationNumber) }
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
  overflow-x: auto; /* 添加水平滚动条 */
}

.chart {
  width: 1200px; /* 设置足够大的宽度 */
  height: 500px; /* 增加图表高度 */
}
</style>
