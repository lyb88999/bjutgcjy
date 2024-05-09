<template>
  <div>
    <div class="gva-form-box">
      <el-form :model="formData" ref="elFormRef" label-position="right" :rules="rule" label-width="80px">
        <el-form-item label="所在学校:" prop="universityName">
          <el-input v-model="formData.universityName" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="本科生数:" prop="undergraduateCount">
          <el-input v-model.number="formData.undergraduateCount" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="硕士生数:" prop="masterStudentCount">
          <el-input v-model.number="formData.masterStudentCount" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="博士生数:" prop="doctoralStudentCount">
          <el-input v-model.number="formData.doctoralStudentCount" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="国家级一流课程数:" prop="nationalFirstClassCourseCount">
          <el-input v-model.number="formData.nationalFirstClassCourseCount" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="省部级一流课程数:" prop="provincialFirstClassCourseCount">
          <el-input v-model.number="formData.provincialFirstClassCourseCount" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="国家级教学示范中心:" prop="ationalTeachingDemonstrationCenterCount">
          <el-input v-model.number="formData.ationalTeachingDemonstrationCenterCount" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="省部级教学示范中心:" prop="provincialTeachingDemonstrationCenterCount">
          <el-input v-model.number="formData.provincialTeachingDemonstrationCenterCount" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="国家级规划教材数:" prop="nationalPlanningTextbookCount">
          <el-input v-model.number="formData.nationalPlanningTextbookCount" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="全国优秀教材数:" prop="nationalExcellentTextbookCount">
          <el-input v-model.number="formData.nationalExcellentTextbookCount" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="省级优秀教材数:" prop="provincialExcellentTextbookCount">
          <el-input v-model.number="formData.provincialExcellentTextbookCount" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="学生获国际奖项数:" prop="internationalAwardsCount">
          <el-input v-model.number="formData.internationalAwardsCount" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="学生获国家级奖项数:" prop="nationalAwardsCount">
          <el-input v-model.number="formData.nationalAwardsCount" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="学生获省部级奖项数:" prop="provincialAwardsCount">
          <el-input v-model.number="formData.provincialAwardsCount" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="本科生平均毕业去向落实率:" prop="undergraduateEmploymentRate">
          <el-input-number v-model="formData.undergraduateEmploymentRate" :precision="2" :clearable="true"></el-input-number>
       </el-form-item>
        <el-form-item label="硕士生平均毕业去向落实率:" prop="masterEmploymentRate">
          <el-input-number v-model="formData.masterEmploymentRate" :precision="2" :clearable="true"></el-input-number>
       </el-form-item>
        <el-form-item label="博士生平均毕业去向落实率:" prop="doctoralEmploymentRate">
          <el-input-number v-model="formData.doctoralEmploymentRate" :precision="2" :clearable="true"></el-input-number>
       </el-form-item>
        <el-form-item label="项目数:" prop="projectCount">
          <el-input v-model.number="formData.projectCount" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="参与学生数:" prop="participatingStudentCount">
          <el-input v-model.number="formData.participatingStudentCount" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="转入工科学生数:" prop="transferredIntoEngineeringCount">
          <el-input v-model.number="formData.transferredIntoEngineeringCount" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="转出工科学生数:" prop="transferredOutOfEngineeringCount">
          <el-input v-model.number="formData.transferredOutOfEngineeringCount" :clearable="true" placeholder="请输入" />
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
  createTalentTrainingDatabase,
  updateTalentTrainingDatabase,
  findTalentTrainingDatabase
} from '@/api/talentTrainingDatabase'

defineOptions({
    name: 'TalentTrainingDatabaseForm'
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
            universityName: '',
            undergraduateCount: 0,
            masterStudentCount: 0,
            doctoralStudentCount: 0,
            nationalFirstClassCourseCount: 0,
            provincialFirstClassCourseCount: 0,
            ationalTeachingDemonstrationCenterCount: 0,
            provincialTeachingDemonstrationCenterCount: 0,
            nationalPlanningTextbookCount: 0,
            nationalExcellentTextbookCount: 0,
            provincialExcellentTextbookCount: 0,
            internationalAwardsCount: 0,
            nationalAwardsCount: 0,
            provincialAwardsCount: 0,
            undergraduateEmploymentRate: 0,
            masterEmploymentRate: 0,
            doctoralEmploymentRate: 0,
            projectCount: 0,
            participatingStudentCount: 0,
            transferredIntoEngineeringCount: 0,
            transferredOutOfEngineeringCount: 0,
        })
// 验证规则
const rule = reactive({
               universityName : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               undergraduateCount : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               masterStudentCount : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               doctoralStudentCount : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               nationalFirstClassCourseCount : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               provincialFirstClassCourseCount : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               ationalTeachingDemonstrationCenterCount : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               provincialTeachingDemonstrationCenterCount : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               nationalPlanningTextbookCount : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               nationalExcellentTextbookCount : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               provincialExcellentTextbookCount : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               internationalAwardsCount : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               nationalAwardsCount : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               provincialAwardsCount : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               undergraduateEmploymentRate : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               masterEmploymentRate : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               doctoralEmploymentRate : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               projectCount : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               participatingStudentCount : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               transferredIntoEngineeringCount : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               transferredOutOfEngineeringCount : [{
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
      const res = await findTalentTrainingDatabase({ ID: route.query.id })
      if (res.code === 0) {
        formData.value = res.data.reTTD
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
               res = await createTalentTrainingDatabase(formData.value)
               break
             case 'update':
               res = await updateTalentTrainingDatabase(formData.value)
               break
             default:
               res = await createTalentTrainingDatabase(formData.value)
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
