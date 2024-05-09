<template>
  <div>
    <div class="gva-form-box">
      <el-form :model="formData" ref="elFormRef" label-position="right" :rules="rule" label-width="80px">
        <el-form-item label="政策名称:" prop="policyName">
          <el-input v-model="formData.policyName" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="政策发布时间:" prop="policyReleaseDate">
          <el-date-picker v-model="formData.policyReleaseDate" type="date" placeholder="选择日期" :clearable="true"></el-date-picker>
       </el-form-item>
        <el-form-item label="政策国家或地区:" prop="policyCountryOrRegion">
          <el-input v-model="formData.policyCountryOrRegion" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="数据来源:" prop="dataSource">
          <el-input v-model="formData.dataSource" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="政策内容:" prop="policyContent">
          <RichEdit v-model="formData.policyContent"/>
       </el-form-item>
        <el-form-item label="其他备注:" prop="additionalRemarks">
          <el-input v-model="formData.additionalRemarks" :clearable="true" placeholder="请输入" />
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
  createPolicyDatabase,
  updatePolicyDatabase,
  findPolicyDatabase
} from '@/api/policyDatabase'

defineOptions({
    name: 'PolicyDatabaseForm'
})

// 自动获取字典
import { getDictFunc } from '@/utils/format'
import { useRoute, useRouter } from "vue-router"
import { ElMessage } from 'element-plus'
import { ref, reactive } from 'vue'
// 富文本组件
import RichEdit from '@/components/richtext/rich-edit.vue'

const route = useRoute()
const router = useRouter()

const type = ref('')
const formData = ref({
            policyName: '',
            policyReleaseDate: new Date(),
            policyCountryOrRegion: '',
            dataSource: '',
            policyContent: '',
            additionalRemarks: '',
        })
// 验证规则
const rule = reactive({
               policyName : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               policyReleaseDate : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               policyCountryOrRegion : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               dataSource : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               policyContent : [{
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
      const res = await findPolicyDatabase({ ID: route.query.id })
      if (res.code === 0) {
        formData.value = res.data.rePD
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
               res = await createPolicyDatabase(formData.value)
               break
             case 'update':
               res = await updatePolicyDatabase(formData.value)
               break
             default:
               res = await createPolicyDatabase(formData.value)
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
