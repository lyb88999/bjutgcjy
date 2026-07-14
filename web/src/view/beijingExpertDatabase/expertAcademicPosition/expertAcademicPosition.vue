<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" @keyup.enter="onSubmit">
        <el-form-item label="专家ID" prop="expertId">
          <el-input v-model="searchInfo.expertId" placeholder="专家ID" />
        </el-form-item>
        <el-form-item label="兼职类型" prop="positionType">
          <el-input v-model="searchInfo.positionType" placeholder="搜索条件" />
        </el-form-item>
        <el-form-item label="任职机构" prop="organizationName">
          <el-input v-model="searchInfo.organizationName" placeholder="搜索条件" />
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
        <el-table-column align="left" label="兼职类型" prop="positionType" width="180" />
        <el-table-column align="left" label="任职机构名称" prop="organizationName" min-width="180" show-overflow-tooltip />
        <el-table-column align="left" label="担任职务" prop="positionTitle" width="140" />
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

    <el-dialog v-model="dialogFormVisible" :before-close="closeDialog" :title="type === 'create' ? '添加学术兼职' : '编辑学术兼职'" width="50%">
      <el-form ref="elFormRef" :model="formData" :rules="rule" label-width="120px">
        <el-form-item label="专家ID" prop="expertId"><el-input-number v-model="formData.expertId" :min="1" style="width:100%" /></el-form-item>
        <el-form-item label="兼职类型" prop="positionType"><el-input v-model="formData.positionType" placeholder="如：学术团体任职/专业委员会任职/政府决策咨询机构兼职" /></el-form-item>
        <el-form-item label="任职机构名称" prop="organizationName"><el-input v-model="formData.organizationName" /></el-form-item>
        <el-form-item label="担任职务" prop="positionTitle"><el-input v-model="formData.positionTitle" /></el-form-item>
        <el-form-item label="任职起始时间" prop="startDate">
          <el-date-picker v-model="formData.startDate" type="date" placeholder="请选择时间" style="width:100%" />
        </el-form-item>
        <el-form-item label="任职结束时间" prop="endDate">
          <el-date-picker v-model="formData.endDate" type="date" placeholder="在任请留空" style="width:100%" />
        </el-form-item>
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
  createExpertAcademicPosition,
  deleteExpertAcademicPosition,
  deleteExpertAcademicPositionByIds,
  updateExpertAcademicPosition,
  findExpertAcademicPosition,
  getExpertAcademicPositionList
} from '@/api/expertAcademicPosition'

import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'

defineOptions({
  name: 'ExpertAcademicPosition'
})

const initFormData = () => ({
  expertId: undefined,
  positionType: '',
  organizationName: '',
  positionTitle: '',
  startDate: null,
  endDate: null,
  remark: ''
})

const formData = ref(initFormData())

const rule = reactive({
  expertId: [{ required: true, message: '请输入专家ID', trigger: ['input', 'blur'] }],
  organizationName: [{ required: true, message: '请输入任职机构名称', trigger: ['input', 'blur'] }]
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
  const table = await getExpertAcademicPositionList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
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
    const res = await deleteExpertAcademicPosition({ ID: row.ID })
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
    const res = await deleteExpertAcademicPositionByIds({ IDs })
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
  const res = await findExpertAcademicPosition({ ID: row.ID })
  type.value = 'update'
  if (res.code === 0) {
    formData.value = res.data.reExpertAcademicPosition
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
      res = await updateExpertAcademicPosition(formData.value)
    } else {
      res = await createExpertAcademicPosition(formData.value)
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
