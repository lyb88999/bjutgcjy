<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo" @keyup.enter="onSubmit">
        <el-form-item label="单位名称">
          <el-input v-model="searchInfo.keyword" placeholder="搜索条件" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="openDialog()">新增单位</el-button>
      </div>
      <el-table :data="tableData" row-key="ID" default-expand-all style="width: 100%">
        <el-table-column align="left" label="单位名称" prop="name" min-width="220" />
        <el-table-column align="left" label="单位级别" prop="level" width="120" />
        <el-table-column align="left" label="行政区划代码" prop="regionCode" width="140" />
        <el-table-column align="left" label="启用状态" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.status ? 'success' : 'info'">{{ scope.row.status ? '启用' : '停用' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="创建日期" width="180">
          <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>
        <el-table-column align="left" label="操作" fixed="right" min-width="200">
          <template #default="scope">
            <el-button type="primary" link @click="openDialog(scope.row)">新增下级</el-button>
            <el-button type="primary" link icon="edit" @click="openDialog(null, scope.row)">编辑</el-button>
            <el-button type="primary" link icon="delete" @click="deleteRow(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <el-dialog v-model="dialogVisible" :title="formData.ID ? '编辑单位' : '新增单位'" width="500px" :before-close="closeDialog">
      <el-form ref="formRef" :model="formData" :rules="rules" label-width="100px">
        <el-form-item label="单位名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入单位名称" />
        </el-form-item>
        <el-form-item label="单位级别" prop="level">
          <el-select v-model="formData.level" placeholder="请选择" style="width: 100%">
            <el-option label="市级" value="市级" />
            <el-option label="区级" value="区级" />
            <el-option label="单位级" value="单位级" />
          </el-select>
        </el-form-item>
        <el-form-item label="上级单位" prop="parentId">
          <el-tree-select
            v-model="formData.parentId"
            :data="parentOptions"
            :props="{ label: 'name', value: 'ID', children: 'children' }"
            node-key="ID"
            check-strictly
            clearable
            placeholder="不选表示顶级单位"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="行政区划代码" prop="regionCode">
          <el-input v-model="formData.regionCode" placeholder="可留空" />
        </el-form-item>
        <el-form-item label="启用状态" prop="status">
          <el-switch v-model="formData.status" />
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
import {
  createSysOrganization,
  deleteSysOrganization,
  updateSysOrganization,
  getSysOrganizationTree
} from '@/api/sysOrganization'
import { formatDate } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'

defineOptions({
  name: 'SysOrganization'
})

const searchInfo = reactive({ keyword: '' })
const tableData = ref([])
const parentOptions = ref([])

const getTableData = async () => {
  const res = await getSysOrganizationTree()
  if (res.code === 0) {
    tableData.value = res.data.tree || []
    parentOptions.value = res.data.tree || []
  }
}
getTableData()

const onSubmit = () => {
  if (!searchInfo.keyword) {
    getTableData()
    return
  }
  // 树形结构下按关键字过滤只在顶层做简单实现：把匹配到的节点连同其祖先展开展示，
  // 这里先用最直接的方式——过滤出名称匹配的节点本身（不含层级），够用且不复杂
  const flatten = (nodes, acc = []) => {
    nodes.forEach(n => {
      if (n.name.includes(searchInfo.keyword)) acc.push({ ...n, children: undefined })
      if (n.children) flatten(n.children, acc)
    })
    return acc
  }
  getSysOrganizationTree().then(res => {
    if (res.code === 0) {
      tableData.value = flatten(res.data.tree || [])
    }
  })
}
const onReset = () => {
  searchInfo.keyword = ''
  getTableData()
}

const initFormData = () => ({ ID: undefined, name: '', level: '单位级', parentId: undefined, regionCode: '', status: true })
const formData = ref(initFormData())
const formRef = ref()
const rules = {
  name: [{ required: true, message: '请输入单位名称', trigger: 'blur' }],
  level: [{ required: true, message: '请选择单位级别', trigger: 'change' }]
}
const dialogVisible = ref(false)

const openDialog = (parentRow, editRow) => {
  if (editRow) {
    formData.value = { ...editRow }
  } else {
    formData.value = initFormData()
    if (parentRow) {
      formData.value.parentId = parentRow.ID
      formData.value.level = '单位级'
    }
  }
  dialogVisible.value = true
}
const closeDialog = () => {
  dialogVisible.value = false
  formData.value = initFormData()
}
const enterDialog = () => {
  formRef.value?.validate(async (valid) => {
    if (!valid) return
    const payload = { ...formData.value }
    const res = payload.ID ? await updateSysOrganization(payload) : await createSysOrganization(payload)
    if (res.code === 0) {
      ElMessage({ type: 'success', message: '保存成功' })
      closeDialog()
      getTableData()
    }
  })
}

const deleteRow = (row) => {
  ElMessageBox.confirm('确定要删除该单位吗？如果下面还有下级单位或者有用户关联在这个单位下，会删除失败。', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    const res = await deleteSysOrganization({ id: row.ID })
    if (res.code === 0) {
      ElMessage({ type: 'success', message: '删除成功' })
      getTableData()
    }
  })
}
</script>

<style></style>
