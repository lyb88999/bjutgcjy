<template>
  <div>
    <div class="gva-form-box">
      <el-form :model="formData" ref="elFormRef" label-position="right" :rules="rule" label-width="80px">
        <el-form-item label="省区市名称:" prop="regionName">
          <el-input v-model="formData.regionName" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="区划代码:" prop="divisionCode">
          <el-input v-model="formData.divisionCode" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="地理区域:" prop="geographicArea">
          <el-input v-model="formData.geographicArea" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="地区生产总值:" prop="regionGDP">
          <el-input-number v-model="formData.regionGDP" :precision="2" :clearable="true"></el-input-number>
       </el-form-item>
        <el-form-item label="学校名称:" prop="schoolName">
          <el-input v-model="formData.schoolName" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="学校代码:" prop="schoolCode">
          <el-input v-model="formData.schoolCode" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="学校类型:" prop="schoolType">
          <el-input v-model="formData.schoolType" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="学校类别:" prop="schoolCategory">
          <el-input v-model="formData.schoolCategory" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="学校层次:" prop="educationLevel">
          <el-input v-model="formData.educationLevel" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="是否双一流:" prop="isDoubleFirstClass">
          <el-switch v-model="formData.isDoubleFirstClass" active-color="#13ce66" inactive-color="#ff4949" active-text="是" inactive-text="否" clearable ></el-switch>
       </el-form-item>
        <el-form-item label="所在地区:" prop="locatedRegion">
          <el-input v-model="formData.locatedRegion" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="在校本科生数:" prop="undergraduateNumber">
          <el-input v-model.number="formData.undergraduateNumber" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="在校硕士生数:" prop="graduateStudentNumber">
          <el-input v-model.number="formData.graduateStudentNumber" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="在校博士生数:" prop="doctoralStudentNumber">
          <el-input v-model.number="formData.doctoralStudentNumber" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="在校留学生数:" prop="internationalStudentNumber">
          <el-input v-model.number="formData.internationalStudentNumber" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="在校教职员工数:" prop="facultyStaffNumber">
          <el-input v-model.number="formData.facultyStaffNumber" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="在校专任教师数:" prop="fullTimeTeacherNumber">
          <el-input v-model.number="formData.fullTimeTeacherNumber" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="本科专业数:" prop="undergraduateProgramNumber">
          <el-input v-model.number="formData.undergraduateProgramNumber" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="硕士授权点数:" prop="mastersAutorizationPointerNumber">
          <el-input v-model.number="formData.mastersAutorizationPointerNumber" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="博士授权点数:" prop="doctoralAuthorizationPointerNumber">
          <el-input v-model.number="formData.doctoralAuthorizationPointerNumber" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="硕士专业学位类别数:" prop="mastersDegreeCategoryNumber">
          <el-input v-model.number="formData.mastersDegreeCategoryNumber" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="博士专业学位类别数:" prop="doctoralDegreeCategoryNumber">
          <el-input v-model.number="formData.doctoralDegreeCategoryNumber" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="博士后流动站数:" prop="postDoctoralMobileStationNumber">
          <el-input v-model.number="formData.postDoctoralMobileStationNumber" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="软科排名:" prop="softScienceRanking">
          <el-input v-model.number="formData.softScienceRanking" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="save">保存</el-button>
          <el-button type="primary" @click="back">返回</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import {
  createBasicInfomationDatabase,
  updateBasicInfomationDatabase,
  findBasicInfomationDatabase
} from '@/api/basicInformationDatabase'

defineOptions({
    name: 'BasicInfomationDatabaseForm'
})

// 自动获取字典
import { getDictFunc } from '@/utils/format'
import { useRoute, useRouter } from "vue-router"
import { ElMessage } from 'element-plus'
import { ref, reactive } from 'vue'

const route = useRoute()
const router = useRouter()

const type = ref('')
const formData = ref({
            regionName: '',
            divisionCode: '',
            geographicArea: '',
            regionGDP: 0,
            schoolName: '',
            schoolCode: '',
            schoolType: '',
            schoolCategory: '',
            educationLevel: '',
            isDoubleFirstClass: false,
            locatedRegion: '',
            undergraduateNumber: 0,
            graduateStudentNumber: 0,
            doctoralStudentNumber: 0,
            internationalStudentNumber: 0,
            facultyStaffNumber: 0,
            fullTimeTeacherNumber: 0,
            undergraduateProgramNumber: 0,
            mastersAutorizationPointerNumber: 0,
            doctoralAuthorizationPointerNumber: 0,
            mastersDegreeCategoryNumber: 0,
            doctoralDegreeCategoryNumber: 0,
            postDoctoralMobileStationNumber: 0,
            softScienceRanking: 0,
        })
// 验证规则
const rule = reactive({
               regionName : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               divisionCode : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               geographicArea : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               regionGDP : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               schoolName : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               schoolCode : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               schoolType : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               schoolCategory : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               educationLevel : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               isDoubleFirstClass : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               locatedRegion : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               undergraduateNumber : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               graduateStudentNumber : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               doctoralStudentNumber : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               internationalStudentNumber : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               facultyStaffNumber : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               fullTimeTeacherNumber : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               undergraduateProgramNumber : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               mastersAutorizationPointerNumber : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               doctoralAuthorizationPointerNumber : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               mastersDegreeCategoryNumber : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               doctoralDegreeCategoryNumber : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               postDoctoralMobileStationNumber : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               softScienceRanking : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
})

const elFormRef = ref()

// 初始化方法
const init = async () => {
 // 建议通过url传参获取目标数据ID 调用 find方法进行查询数据操作 从而决定本页面是create还是update 以下为id作为url参数示例
    if (route.query.id) {
      const res = await findBasicInfomationDatabase({ ID: route.query.id })
      if (res.code === 0) {
        formData.value = res.data.reBID
        type.value = 'update'
      }
    } else {
      type.value = 'create'
    }
}

init()
// 保存按钮
const save = async() => {
      elFormRef.value?.validate( async (valid) => {
         if (!valid) return
            let res
           switch (type.value) {
             case 'create':
               res = await createBasicInfomationDatabase(formData.value)
               break
             case 'update':
               res = await updateBasicInfomationDatabase(formData.value)
               break
             default:
               res = await createBasicInfomationDatabase(formData.value)
               break
           }
           if (res.code === 0) {
             ElMessage({
               type: 'success',
               message: '创建/更改成功'
             })
           }
       })
}

// 返回按钮
const back = () => {
    router.go(-1)
}

</script>

<style>
</style>
