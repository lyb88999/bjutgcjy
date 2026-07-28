<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" @keyup.enter="onSubmit">
        <el-form-item label="专家ID" prop="expertId">
          <el-input v-model="searchInfo.expertId" placeholder="专家ID" />
        </el-form-item>
        <el-form-item label="成果类型" prop="achievementType">
          <el-input v-model="searchInfo.achievementType" placeholder="搜索条件" />
        </el-form-item>
        <el-form-item label="成果级别" prop="level">
          <el-input v-model="searchInfo.level" placeholder="搜索条件" />
        </el-form-item>
        <el-form-item label="成果名称" prop="title">
          <el-input v-model="searchInfo.title" placeholder="搜索条件" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="openDialog">新增</el-button>
        <el-button icon="delete" style="margin-left: 10px;" :disabled="!multipleSelection.length" @click="onDelete">删除</el-button>
      </div>
      <el-table ref="multipleTable" style="width: 100%" tooltip-effect="dark" :data="tableData" row-key="ID"
        @selection-change="handleSelectionChange">
        <el-table-column type="selection" fixed width="55" />
        <el-table-column align="left" label="专家ID" prop="expertId" width="90" />
        <el-table-column align="left" label="成果类型" prop="achievementType" width="120" />
        <el-table-column align="left" label="成果名称" prop="title" min-width="200" show-overflow-tooltip />
        <el-table-column align="left" label="成果级别" prop="level" width="100" />
        <el-table-column align="left" label="级别认定来源" prop="levelSource" width="140" />
        <el-table-column align="left" label="发表/立项单位" prop="publishOrg" width="160" />
        <el-table-column align="left" label="检索相关次数" prop="refCount" width="110" />
        <el-table-column align="left" label="操作" fixed="right" width="150">
          <template #default="scope">
            <el-button type="primary" link icon="edit" @click="updateFunc(scope.row)">编辑</el-button>
            <el-button type="primary" link icon="delete" @click="deleteRow(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination
          :current-page="page"
          :page-size="pageSize"
          :page-sizes="[10, 30, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @current-change="handleCurrentChange"
          @size-change="handleSizeChange"
        />
      </div>
    </div>

    <el-dialog v-model="dialogFormVisible" :before-close="closeDialog" :title="type === 'create' ? '添加成果' : '编辑成果'" width="50%">
      <el-form ref="elFormRef" :model="formData" :rules="rule" label-width="120px">
        <el-form-item label="专家ID" prop="expertId"><el-input-number v-model="formData.expertId" :min="1" style="width:100%" /></el-form-item>
        <el-form-item label="成果类型" prop="achievementType">
          <el-select v-model="formData.achievementType" placeholder="请选择" style="width:100%">
            <el-option v-for="item in achievementTypeOptions" :key="item" :label="item" :value="item" />
          </el-select>
        </el-form-item>
        <el-form-item label="成果名称" prop="title"><el-input v-model="formData.title" /></el-form-item>
        <el-form-item label="成果级别" prop="level"><el-input v-model="formData.level" placeholder="如：国家级/省部级/厅局级/一般级" /></el-form-item>
        <el-form-item label="级别认定来源" prop="levelSource"><el-input v-model="formData.levelSource" /></el-form-item>
        <el-form-item label="发表/出版/立项单位" prop="publishOrg"><el-input v-model="formData.publishOrg" /></el-form-item>
        <el-form-item label="发表/出版/立项时间" prop="publishDate">
          <el-date-picker v-model="formData.publishDate" type="date" placeholder="请选择时间" style="width:100%" />
        </el-form-item>
        <el-form-item label="关键词" prop="keywords"><el-input v-model="formData.keywords" placeholder="用于相关性检索，多个关键词用逗号分隔" /></el-form-item>
        <el-form-item label="检索相关次数" prop="refCount"><el-input-number v-model="formData.refCount" :min="0" style="width:100%" /></el-form-item>
        <el-form-item label="备注" prop="remark"><el-input v-model="formData.remark" type="textarea" /></el-form-item>
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
import {
  createExpertAchievement,
  deleteExpertAchievement,
  deleteExpertAchievementByIds,
  updateExpertAchievement,
  findExpertAchievement,
  getExpertAchievementList
} from '@/api/expertAchievement'

import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'

defineOptions({
  name: 'ExpertAchievement'
})

const achievementTypeOptions = ['著作', '学术论文', '研究报告', '决策咨询成果', '获奖成果', '课题项目']

const initFormData = () => ({
  expertId: undefined,
  achievementType: '',
  title: '',
  level: '',
  levelSource: '',
  publishOrg: '',
  publishDate: null,
  keywords: '',
  refCount: 0,
  remark: ''
})

const formData = ref(initFormData())

const rule = reactive({
  expertId: [{ required: true, message: '请输入专家ID', trigger: ['input', 'blur'] }],
  achievementType: [{ required: true, message: '请选择成果类型', trigger: ['change', 'blur'] }],
  title: [{ required: true, message: '请输入成果名称', trigger: ['input', 'blur'] }]
})

const elFormRef = ref()
const elSearchFormRef = ref()

const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})

const onReset = () => {
  searchInfo.value = {}
  getTableData()
}

const onSubmit = () => {
  page.value = 1
  pageSize.value = 10
  getTableData()
}

const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

const getTableData = async () => {
  const table = await getExpertAchievementList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}

getTableData()

const multipleSelection = ref([])
const handleSelectionChange = (val) => {
  multipleSelection.value = val
}

const deleteRow = (row) => {
  ElMessageBox.confirm('确定要删除吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    const res = await deleteExpertAchievement({ ID: row.ID })
    if (res.code === 0) {
      ElMessage({ type: 'success', message: '删除成功' })
      if (tableData.value.length === 1 && page.value > 1) {
        page.value--
      }
      getTableData()
    }
  })
}

const onDelete = async () => {
  ElMessageBox.confirm('确定要删除吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    if (multipleSelection.value.length === 0) {
      ElMessage({ type: 'warning', message: '请选择要删除的数据' })
      return
    }
    const IDs = multipleSelection.value.map(item => item.ID)
    const res = await deleteExpertAchievementByIds({ IDs })
    if (res.code === 0) {
      ElMessage({ type: 'success', message: '删除成功' })
      if (tableData.value.length === IDs.length && page.value > 1) {
        page.value--
      }
      getTableData()
    }
  })
}

const type = ref('')

const updateFunc = async (row) => {
  const res = await findExpertAchievement({ ID: row.ID })
  type.value = 'update'
  if (res.code === 0) {
    formData.value = res.data.reExpertAchievement
    dialogFormVisible.value = true
  }
}

const dialogFormVisible = ref(false)

const openDialog = () => {
  type.value = 'create'
  dialogFormVisible.value = true
}

const closeDialog = () => {
  dialogFormVisible.value = false
  formData.value = initFormData()
}

const enterDialog = async () => {
  elFormRef.value?.validate(async (valid) => {
    if (!valid) return
    let res
    if (type.value === 'update') {
      res = await updateExpertAchievement(formData.value)
    } else {
      res = await createExpertAchievement(formData.value)
    }
    if (res.code === 0) {
      ElMessage({ type: 'success', message: '创建/更改成功' })
      closeDialog()
      getTableData()
    }
  })
}
</script>

<style></style>
