<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" @keyup.enter="onSubmit">
        <el-form-item label="专家ID" prop="expertId">
          <el-input v-model="searchInfo.expertId" placeholder="专家ID" />
        </el-form-item>
        <el-form-item label="采纳类型" prop="adoptionType">
          <el-input v-model="searchInfo.adoptionType" placeholder="搜索条件" />
        </el-form-item>
        <el-form-item label="采纳单位级别" prop="adoptingUnitLevel">
          <el-input v-model="searchInfo.adoptingUnitLevel" placeholder="搜索条件" />
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
        <el-table-column align="left" label="采纳类型" prop="adoptionType" width="160" />
        <el-table-column align="left" label="采纳/批示单位级别" prop="adoptingUnitLevel" width="140" />
        <el-table-column align="left" label="采纳/批示单位名称" prop="adoptingUnitName" min-width="180" show-overflow-tooltip />
        <el-table-column align="left" label="佐证材料描述" prop="evidenceDesc" min-width="180" show-overflow-tooltip />
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

    <el-dialog v-model="dialogFormVisible" :before-close="closeDialog" :title="type === 'create' ? '添加决策影响记录' : '编辑决策影响记录'" width="50%">
      <el-form ref="elFormRef" :model="formData" :rules="rule" label-width="140px">
        <el-form-item label="专家ID" prop="expertId"><el-input-number v-model="formData.expertId" :min="1" style="width:100%" /></el-form-item>
        <el-form-item label="关联成果ID" prop="achievementId"><el-input-number v-model="formData.achievementId" :min="1" style="width:100%" /></el-form-item>
        <el-form-item label="采纳类型" prop="adoptionType">
          <el-select v-model="formData.adoptionType" placeholder="请选择" style="width:100%">
            <el-option v-for="item in adoptionTypeOptions" :key="item" :label="item" :value="item" />
          </el-select>
        </el-form-item>
        <el-form-item label="采纳/批示单位级别" prop="adoptingUnitLevel"><el-input v-model="formData.adoptingUnitLevel" placeholder="如：国家级/省部级/厅局级/区县级" /></el-form-item>
        <el-form-item label="采纳/批示单位名称" prop="adoptingUnitName"><el-input v-model="formData.adoptingUnitName" /></el-form-item>
        <el-form-item label="采纳/批示时间" prop="adoptionDate">
          <el-date-picker v-model="formData.adoptionDate" type="date" placeholder="请选择时间" style="width:100%" />
        </el-form-item>
        <el-form-item label="佐证材料描述" prop="evidenceDesc"><el-input v-model="formData.evidenceDesc" type="textarea" /></el-form-item>
        <el-form-item label="佐证材料附件" prop="attachmentUrl"><el-input v-model="formData.attachmentUrl" placeholder="附件链接" /></el-form-item>
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
  createExpertAdoptionRecord,
  deleteExpertAdoptionRecord,
  deleteExpertAdoptionRecordByIds,
  updateExpertAdoptionRecord,
  findExpertAdoptionRecord,
  getExpertAdoptionRecordList
} from '@/api/expertAdoptionRecord'

import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'

defineOptions({
  name: 'ExpertAdoptionRecord'
})

const adoptionTypeOptions = ['内参/专报采用', '部门采纳', '领导批示', '进入政策文件', '参与政策起草/咨询论证']

const initFormData = () => ({
  expertId: undefined,
  achievementId: undefined,
  adoptionType: '',
  adoptingUnitLevel: '',
  adoptingUnitName: '',
  adoptionDate: null,
  evidenceDesc: '',
  attachmentUrl: ''
})

const formData = ref(initFormData())

const rule = reactive({
  expertId: [{ required: true, message: '请输入专家ID', trigger: ['input', 'blur'] }],
  adoptionType: [{ required: true, message: '请选择采纳类型', trigger: ['change', 'blur'] }]
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
  const table = await getExpertAdoptionRecordList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
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
    const res = await deleteExpertAdoptionRecord({ ID: row.ID })
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
    const res = await deleteExpertAdoptionRecordByIds({ IDs })
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
  const res = await findExpertAdoptionRecord({ ID: row.ID })
  type.value = 'update'
  if (res.code === 0) {
    formData.value = res.data.reExpertAdoptionRecord
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
      res = await updateExpertAdoptionRecord(formData.value)
    } else {
      res = await createExpertAdoptionRecord(formData.value)
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
