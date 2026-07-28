<template>
  <div>
    <el-alert
      type="info" :closable="false" show-icon
      title="这里新建的账号固定是「个人申报人」角色，所属单位固定是你自己所在的单位，不能修改——用来给本单位的专家开通申报账号，不用再找超级管理员手工建。"
      style="margin-bottom: 16px;"
    />
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="openDialog">新建账号</el-button>
      </div>
      <el-table :data="tableData" style="width: 100%">
        <el-table-column align="left" label="用户名" prop="username" width="160" />
        <el-table-column align="left" label="昵称" prop="nickName" width="140" />
        <el-table-column align="left" label="手机号" prop="phone" width="140" />
        <el-table-column align="left" label="邮箱" prop="email" min-width="180" />
        <el-table-column align="left" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.enable === 1 ? 'success' : 'info'">{{ scope.row.enable === 1 ? '正常' : '已冻结' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="创建日期" width="180">
          <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>
        <el-table-column align="left" label="操作" fixed="right" width="120">
          <template #default="scope">
            <el-button type="primary" link @click="toggleEnable(scope.row)">
              {{ scope.row.enable === 1 ? '冻结' : '启用' }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination
          :current-page="page"
          :page-size="pageSize"
          :total="total"
          layout="total, prev, pager, next"
          @current-change="v => { page = v; getTableData() }"
        />
      </div>
    </div>

    <el-dialog v-model="dialogVisible" title="新建本单位账号" width="500px" :before-close="closeDialog">
      <el-form ref="formRef" :model="formData" :rules="rules" label-width="90px">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="formData.username" placeholder="登录用的用户名" />
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input v-model="formData.password" type="password" show-password placeholder="至少6位" />
        </el-form-item>
        <el-form-item label="昵称" prop="nickName">
          <el-input v-model="formData.nickName" placeholder="可留空，默认使用用户名" />
        </el-form-item>
        <el-form-item label="手机号" prop="phone">
          <el-input v-model="formData.phone" placeholder="可留空" />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="formData.email" placeholder="可留空" />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closeDialog">取消</el-button>
          <el-button type="primary" @click="enterDialog">确定</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { createOrgApplicant, getOrgUserList, toggleOrgUserEnable } from '@/api/expertOrgUser'
import { formatDate } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'

defineOptions({
  name: 'ExpertOrgUser'
})

const tableData = ref([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

const getTableData = async () => {
  const res = await getOrgUserList({ page: page.value, pageSize: pageSize.value })
  if (res.code === 0) {
    tableData.value = res.data.list || []
    total.value = res.data.total
  }
}
getTableData()

const initFormData = () => ({ username: '', password: '', nickName: '', phone: '', email: '' })
const formData = ref(initFormData())
const formRef = ref()
const rules = reactive({
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }, { min: 6, message: '密码至少6位', trigger: 'blur' }]
})
const dialogVisible = ref(false)

const openDialog = () => {
  formData.value = initFormData()
  dialogVisible.value = true
}
const closeDialog = () => {
  dialogVisible.value = false
}
const enterDialog = () => {
  formRef.value?.validate(async (valid) => {
    if (!valid) return
    const res = await createOrgApplicant(formData.value)
    if (res.code === 0) {
      ElMessage({ type: 'success', message: '创建成功' })
      closeDialog()
      getTableData()
    }
  })
}

const toggleEnable = (row) => {
  const nextEnable = row.enable === 1 ? 2 : 1
  const actionText = nextEnable === 1 ? '启用' : '冻结'
  ElMessageBox.confirm(`确定要${actionText}账号「${row.username}」吗？`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    const res = await toggleOrgUserEnable({ ID: row.ID, enable: nextEnable })
    if (res.code === 0) {
      ElMessage({ type: 'success', message: '操作成功' })
      getTableData()
    }
  })
}
</script>

<style></style>
