<template>
  <div>
    <div class="gva-form-box">
      <el-form :model="formData" ref="elFormRef" label-position="right" :rules="rule" label-width="80px">
        <el-form-item label="学校名称:" prop="universityName">
          <el-input v-model="formData.universityName" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="行业协同创新中心数或产学研基地数:" prop="industryInnovationCenterCount">
          <el-input v-model.number="formData.industryInnovationCenterCount" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="区域协同创新中心数或产学研基地数:" prop="regionalInnovationCenterCount">
          <el-input v-model.number="formData.regionalInnovationCenterCount" :clearable="true" placeholder="请输入" />
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
  createSocialServiceDatabase,
  updateSocialServiceDatabase,
  findSocialServiceDatabase
} from '@/api/socialServiceDatabase'

defineOptions({
    name: 'SocialServiceDatabaseForm'
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
            industryInnovationCenterCount: 0,
            regionalInnovationCenterCount: 0,
        })
// 验证规则
const rule = reactive({
               universityName : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               industryInnovationCenterCount : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               regionalInnovationCenterCount : [{
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
      const res = await findSocialServiceDatabase({ ID: route.query.id })
      if (res.code === 0) {
        formData.value = res.data.reSSD
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
               res = await createSocialServiceDatabase(formData.value)
               break
             case 'update':
               res = await updateSocialServiceDatabase(formData.value)
               break
             default:
               res = await createSocialServiceDatabase(formData.value)
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
