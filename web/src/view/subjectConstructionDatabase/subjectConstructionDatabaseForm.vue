<template>
  <div>
    <div class="gva-form-box">
      <el-form :model="formData" ref="elFormRef" label-position="right" :rules="rule" label-width="80px">
        <el-form-item label="流动站名称:" prop="postdocStationName">
          <el-input v-model="formData.postdocStationName" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="学科门类:" prop="postdocSubjectCategory">
          <el-input v-model="formData.postdocSubjectCategory" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="博士后流动站数:" prop="postdocStationCount">
          <el-input v-model.number="formData.postdocStationCount" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="学科代码:" prop="masterSubjectCode">
          <el-input v-model="formData.masterSubjectCode" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="学科门类:" prop="masterSubjectCategory">
          <el-input v-model="formData.masterSubjectCategory" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="软科排名:" prop="masterSoftScienceRanking">
          <el-input v-model.number="formData.masterSoftScienceRanking" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="招生人数:" prop="masterAdmissionNumber">
          <el-input v-model.number="formData.masterAdmissionNumber" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="类型:" prop="masterType">
        <el-select v-model="formData.masterType" placeholder="请选择" style="width:100%" :clearable="true">
          <el-option v-for="item in ['学术学位','专业学位']" :key="item" :label="item" :value="item" />
        </el-select>
       </el-form-item>
        <el-form-item label="是否一流学科:" prop="masterIsFirstClass">
          <el-switch v-model="formData.masterIsFirstClass" active-color="#13ce66" inactive-color="#ff4949" active-text="是" inactive-text="否" clearable ></el-switch>
       </el-form-item>
        <el-form-item label="所在学校:" prop="masterUniversityName">
          <el-input v-model="formData.masterUniversityName" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="学科代码:" prop="phdSubjectCode">
          <el-input v-model="formData.phdSubjectCode" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="学科门类:" prop="phdSubjectCategory">
          <el-input v-model="formData.phdSubjectCategory" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="软科排名:" prop="phdSoftScienceRanking">
          <el-input v-model.number="formData.phdSoftScienceRanking" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="招生人数:" prop="phdAdmissionNumber">
          <el-input v-model.number="formData.phdAdmissionNumber" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="类型:" prop="phdType">
        <el-select v-model="formData.phdType" placeholder="请选择" style="width:100%" :clearable="true">
          <el-option v-for="item in ['学术学位','专业学位']" :key="item" :label="item" :value="item" />
        </el-select>
       </el-form-item>
        <el-form-item label="是否一流学科:" prop="phdIsFirstClass">
          <el-switch v-model="formData.phdIsFirstClass" active-color="#13ce66" inactive-color="#ff4949" active-text="是" inactive-text="否" clearable ></el-switch>
       </el-form-item>
        <el-form-item label="所在学校:" prop="phdUniversityName">
          <el-input v-model="formData.phdUniversityName" :clearable="true" placeholder="请输入" />
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
  createSubjectConstructionDatabase,
  updateSubjectConstructionDatabase,
  findSubjectConstructionDatabase
} from '@/api/subjectConstructionDatabase'

defineOptions({
    name: 'SubjectConstructionDatabaseForm'
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
            postdocStationName: '',
            postdocSubjectCategory: '',
            postdocStationCount: 0,
            masterSubjectCode: '',
            masterSubjectCategory: '',
            masterSoftScienceRanking: 0,
            masterAdmissionNumber: 0,
            masterIsFirstClass: false,
            masterUniversityName: '',
            phdSubjectCode: '',
            phdSubjectCategory: '',
            phdSoftScienceRanking: 0,
            phdAdmissionNumber: 0,
            phdIsFirstClass: false,
            phdUniversityName: '',
        })
// 验证规则
const rule = reactive({
})

const elFormRef = ref()

// 初始化方法
const init = async () => {
 // 建议通过url传参获取目标数据ID 调用 find方法进行查询数据操作 从而决定本页面是create还是update 以下为id作为url参数示例
    if (route.query.id) {
      const res = await findSubjectConstructionDatabase({ ID: route.query.id })
      if (res.code === 0) {
        formData.value = res.data.reSCD
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
               res = await createSubjectConstructionDatabase(formData.value)
               break
             case 'update':
               res = await updateSubjectConstructionDatabase(formData.value)
               break
             default:
               res = await createSubjectConstructionDatabase(formData.value)
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
