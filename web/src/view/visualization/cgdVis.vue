<template>
    <div class="chart-container">
        <!-- Ensure charts render only after data is loaded -->
        <v-chart v-if="schools.length > 0 && selectedChart === 'talents'" class="chart" :option="talentsOption" />
        <v-chart v-else-if="schools.length > 0 && selectedChart === 'resources'" class="chart"
            :option="resourcesOption" />

        <label for="chart-select">选择图表：</label>
        <select v-model="selectedChart" id="chart-select">
            <option value="talents">教师与人才数据</option>
            <option value="resources">实习基地与设备共享</option>
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
import { getConditionalGuaranteeDatabaseList } from '@/api/conditionalGuaranteeDatabase';

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

// Fetch data from the backend
const fetchSchoolData = async () => {
    try {
        const res = await getConditionalGuaranteeDatabaseList({ page: 1, pageSize: 100 });
        if (res.code === 0) {
            schools.value = res.data.list.map(school => ({
                universityName: school.universityName,
                fullTimeFacultyCount: school.fullTimeFacultyCount,
                foreignHighLevelTalentCount: school.foreignHighLevelTalentCount,
                academicianCount: school.academicianCount,
                nationalTalentProjectWinnerCount: school.nationalTalentProjectWinnerCount,
                provincialTalentProjectWinnerCount: school.provincialTalentProjectWinnerCount,
                offCampusInternshipBaseCount: school.offCampusInternshipBaseCount,
                largeEquipmentSharingPlatformCount: school.largeEquipmentSharingPlatformCount,
                nationalExperimentalTeachingDemoCenterCount: school.nationalExperimentalTeachingDemoCenterCount
            }));
        } else {
            console.error('Error fetching data: ', res.data.message);
        }
    } catch (error) {
        console.error('Error during data fetch: ', error);
    }
};

onMounted(fetchSchoolData);

const selectedChart = ref('talents');

// 教师与人才数据可视化
const talentsOption = computed(() => ({
    title: { text: '教师与人才数据', left: 'center', top: '5%' },
    tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } },
    legend: { left: 'right', data: ['专任教师数', '外籍高层次人才数', '院士人数', '国家级人才项目获得者人次', '省部级人才项目获得者人次'] },
    xAxis: { type: 'category', data: schools.value.map(s => s.universityName) },
    yAxis: { type: 'value' },
    series: [
        { name: '专任教师数', type: 'bar', data: schools.value.map(s => s.fullTimeFacultyCount) },
        { name: '外籍高层次人才数', type: 'bar', data: schools.value.map(s => s.foreignHighLevelTalentCount) },
        { name: '院士人数', type: 'bar', data: schools.value.map(s => s.academicianCount) },
        { name: '国家级人才项目获得者人次', type: 'bar', data: schools.value.map(s => s.nationalTalentProjectWinnerCount) },
        { name: '省部级人才项目获得者人次', type: 'bar', data: schools.value.map(s => s.provincialTalentProjectWinnerCount) }
    ]
}));

// 实习基地与设备共享可视化
const resourcesOption = computed(() => ({
    title: { text: '实习基地与设备共享', left: 'center', top: '5%' },
    tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } },
    legend: { left: 'right', data: ['校外实习、实践、实训基地数', '大型仪器设备共享系统平台数', '国家级实验教学示范中心数'] },
    xAxis: { type: 'category', data: schools.value.map(s => s.universityName) },
    yAxis: { type: 'value' },
    series: [
        { name: '校外实习、实践、实训基地数', type: 'bar', data: schools.value.map(s => s.offCampusInternshipBaseCount) },
        { name: '大型仪器设备共享系统平台数', type: 'bar', data: schools.value.map(s => s.largeEquipmentSharingPlatformCount) },
        { name: '国家级实验教学示范中心数', type: 'bar', data: schools.value.map(s => s.nationalExperimentalTeachingDemoCenterCount) }
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