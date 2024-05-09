<template>
  <div>
    <div class="gva-form-box">
      <el-form :model="formData" ref="elFormRef" label-position="right" :rules="rule" label-width="80px">
        <el-form-item label="学校名称:" prop="universityName">
          <el-input v-model="formData.universityName" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="专任教师数:" prop="fullTimeFacultyCount">
          <el-input v-model.number="formData.fullTimeFacultyCount" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="外籍高层次人才数:" prop="foreignHighLevelTalentCount">
          <el-input v-model.number="formData.foreignHighLevelTalentCount" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="院士人数:" prop="academicianCount">
          <el-input v-model.number="formData.academicianCount" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="国家级人才项目获得者人次:" prop="nationalTalentProjectWinnerCount">
          <el-input v-model.number="formData.nationalTalentProjectWinnerCount" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="省部级人才项目获得者人次:" prop="provincialTalentProjectWinnerCount">
          <el-input v-model.number="formData.provincialTalentProjectWinnerCount" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="校外实习、实践、实训基地数:" prop="offCampusInternshipBaseCount">
          <el-input v-model.number="formData.offCampusInternshipBaseCount" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="大型仪器设备共享系统平台数:" prop="largeEquipmentSharingPlatformCount">
          <el-input v-model.number="formData.largeEquipmentSharingPlatformCount" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="国家级实验教学示范中心数:" prop="nationalExperimentalTeachingDemoCenterCount">
          <el-input v-model.number="formData.nationalExperimentalTeachingDemoCenterCount" :clearable="true" placeholder="请输入" />
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
  createConditionalGuaranteeDatabase,
  updateConditionalGuaranteeDatabase,
  findConditionalGuaranteeDatabase
} from '@/api/conditionalGuaranteeDatabase'

defineOptions({
    name: 'ConditionalGuaranteeDatabaseForm'
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
            fullTimeFacultyCount: 0,
            foreignHighLevelTalentCount: 0,
            academicianCount: 0,
            nationalTalentProjectWinnerCount: 0,
            provincialTalentProjectWinnerCount: 0,
            offCampusInternshipBaseCount: 0,
            largeEquipmentSharingPlatformCount: 0,
            nationalExperimentalTeachingDemoCenterCount: 0,
        })
// 验证规则
const rule = reactive({
               universityName : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               fullTimeFacultyCount : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               foreignHighLevelTalentCount : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               academicianCount : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               nationalTalentProjectWinnerCount : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               provincialTalentProjectWinnerCount : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               offCampusInternshipBaseCount : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               largeEquipmentSharingPlatformCount : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               nationalExperimentalTeachingDemoCenterCount : [{
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
      const res = await findConditionalGuaranteeDatabase({ ID: route.query.id })
      if (res.code === 0) {
        formData.value = res.data.reCGD
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
               res = await createConditionalGuaranteeDatabase(formData.value)
               break
             case 'update':
               res = await updateConditionalGuaranteeDatabase(formData.value)
               break
             default:
               res = await createConditionalGuaranteeDatabase(formData.value)
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
