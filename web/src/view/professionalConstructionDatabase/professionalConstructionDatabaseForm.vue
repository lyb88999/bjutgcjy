<template>
  <div>
    <div class="gva-form-box">
      <el-form :model="formData" ref="elFormRef" label-position="right" :rules="rule" label-width="80px">
        <el-form-item label="专业名称:" prop="majorName">
          <el-input v-model="formData.majorName" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="专业代码:" prop="majorCode">
          <el-input v-model="formData.majorCode" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="授予学位类别:" prop="degreeType">
        <el-select v-model="formData.degreeType" placeholder="请选择" style="width:100%" :clearable="true">
          <el-option v-for="item in ['学士','硕士','博士']" :key="item" :label="item" :value="item" />
        </el-select>
       </el-form-item>
        <el-form-item label="是否新专业:" prop="isNewMajor">
          <el-switch v-model="formData.isNewMajor" active-color="#13ce66" inactive-color="#ff4949" active-text="是" inactive-text="否" clearable ></el-switch>
       </el-form-item>
        <el-form-item label="是否大类招生:" prop="isBroadAdmissionCategory">
          <el-switch v-model="formData.isBroadAdmissionCategory" active-color="#13ce66" inactive-color="#ff4949" active-text="是" inactive-text="否" clearable ></el-switch>
       </el-form-item>
        <el-form-item label="是否一流专业:" prop="isFirstClassMajor">
          <el-switch v-model="formData.isFirstClassMajor" active-color="#13ce66" inactive-color="#ff4949" active-text="是" inactive-text="否" clearable ></el-switch>
       </el-form-item>
        <el-form-item label="是否专业认证:" prop="isMajorAccredited">
          <el-switch v-model="formData.isMajorAccredited" active-color="#13ce66" inactive-color="#ff4949" active-text="是" inactive-text="否" clearable ></el-switch>
       </el-form-item>
        <el-form-item label="是否入选卓越计划:" prop="isSelectedForExcellenceProgram">
          <el-switch v-model="formData.isSelectedForExcellenceProgram" active-color="#13ce66" inactive-color="#ff4949" active-text="是" inactive-text="否" clearable ></el-switch>
       </el-form-item>
        <el-form-item label="是否双学位:" prop="isDoubleDegree">
          <el-switch v-model="formData.isDoubleDegree" active-color="#13ce66" inactive-color="#ff4949" active-text="是" inactive-text="否" clearable ></el-switch>
       </el-form-item>
        <el-form-item label="是否新工科:" prop="isNewEngineeringDiscipline">
          <el-switch v-model="formData.isNewEngineeringDiscipline" active-color="#13ce66" inactive-color="#ff4949" active-text="是" inactive-text="否" clearable ></el-switch>
       </el-form-item>
        <el-form-item label="专业满意度:" prop="majorSatisfaction">
        <el-select v-model="formData.majorSatisfaction" placeholder="请选择" style="width:100%" :clearable="true">
          <el-option v-for="item in ['非常不满意','不满意','中等','满意','非常满意']" :key="item" :label="item" :value="item" />
        </el-select>
       </el-form-item>
        <el-form-item label="所在学校:" prop="universityName">
          <el-input v-model="formData.universityName" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="招生专业:" prop="admissionMajor">
          <el-input v-model="formData.admissionMajor" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="计划招生数:" prop="plannedEnrollment">
          <el-input v-model.number="formData.plannedEnrollment" :clearable="true" placeholder="请输入" />
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
  createProfessionalConstructionDatabase,
  updateProfessionalConstructionDatabase,
  findProfessionalConstructionDatabase
} from '@/api/professionalConstructionDatabase'

defineOptions({
    name: 'ProfessionalConstructionDatabaseForm'
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
            majorName: '',
            majorCode: '',
            isNewMajor: false,
            isBroadAdmissionCategory: false,
            isFirstClassMajor: false,
            isMajorAccredited: false,
            isSelectedForExcellenceProgram: false,
            isDoubleDegree: false,
            isNewEngineeringDiscipline: false,
            universityName: '',
            admissionMajor: '',
            plannedEnrollment: 0,
        })
// 验证规则
const rule = reactive({
               majorName : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               majorCode : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               degreeType : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               isNewMajor : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               isBroadAdmissionCategory : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               isFirstClassMajor : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               isMajorAccredited : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               isSelectedForExcellenceProgram : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               isDoubleDegree : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               isNewEngineeringDiscipline : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               majorSatisfaction : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               universityName : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               admissionMajor : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               plannedEnrollment : [{
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
      const res = await findProfessionalConstructionDatabase({ ID: route.query.id })
      if (res.code === 0) {
        formData.value = res.data.rePCD
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
               res = await createProfessionalConstructionDatabase(formData.value)
               break
             case 'update':
               res = await updateProfessionalConstructionDatabase(formData.value)
               break
             default:
               res = await createProfessionalConstructionDatabase(formData.value)
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
