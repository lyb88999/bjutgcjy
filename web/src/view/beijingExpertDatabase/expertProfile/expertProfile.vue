<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline" @keyup.enter="onSubmit">
        <el-form-item label="姓名" prop="name">
          <el-input v-model="searchInfo.name" placeholder="搜索条件" />
        </el-form-item>
        <el-form-item label="所在单位" prop="unitName">
          <el-input v-model="searchInfo.unitName" placeholder="搜索条件" />
        </el-form-item>
        <el-form-item label="职称" prop="techTitle">
          <el-input v-model="searchInfo.techTitle" placeholder="搜索条件" />
        </el-form-item>
        <el-form-item label="一级学科" prop="disciplineL1">
          <el-input v-model="searchInfo.disciplineL1" placeholder="搜索条件" />
        </el-form-item>
        <el-form-item label="研究关键词" prop="researchKeywords">
          <el-input v-model="searchInfo.researchKeywords" placeholder="搜索条件" />
        </el-form-item>
        <el-form-item label="审核状态" prop="status">
          <el-select v-model="searchInfo.status" clearable placeholder="请选择">
            <el-option v-for="item in statusOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item v-if="!isReadOnlyReviewer" label="单位归属" prop="orgUnmatched">
          <el-checkbox v-model="searchInfo.orgUnmatched">只看未关联单位</el-checkbox>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <div v-if="!isReadOnlyReviewer" class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="openDialog">新增</el-button>
        <el-button icon="upload" style="margin-left: 10px;" @click="openImportDialog">批量导入</el-button>
        <el-button icon="download" style="margin-left: 10px;" :loading="exporting" @click="handleExport">导出当前结果</el-button>
        <el-button icon="delete" style="margin-left: 10px;" :disabled="!multipleSelection.length" @click="onDelete">删除</el-button>
      </div>
      <el-table ref="multipleTable" style="width: 100%" tooltip-effect="dark" :data="tableData" row-key="ID"
        @selection-change="handleSelectionChange">
        <el-table-column v-if="!isReadOnlyReviewer" type="selection" fixed width="55" />
        <el-table-column align="left" label="姓名" prop="name" width="100" />
        <el-table-column align="left" label="所在单位" prop="unitName" width="180">
          <template #default="scope">
            {{ scope.row.unitName }}
            <el-tooltip v-if="!scope.row.orgId" content="单位没能关联到「单位管理」里的记录，编辑时可以在下拉里重新选一下" placement="top">
              <el-tag type="warning" size="small" style="margin-left: 4px;">未关联</el-tag>
            </el-tooltip>
          </template>
        </el-table-column>
        <el-table-column align="left" label="专业技术职称" prop="techTitle" width="120" />
        <el-table-column align="left" label="一级学科" prop="disciplineL1" width="120" />
        <el-table-column align="left" label="研究关键词" prop="researchKeywords" width="180" show-overflow-tooltip />
        <el-table-column align="left" label="审核状态" prop="status" width="150">
          <template #default="scope">
            {{ statusLabel(scope.row.status) }}
            <el-tag v-if="scope.row.reviewBypassed" type="warning" size="small" style="margin-left: 4px;">免审核发布</el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="综合排序得分" prop="compositeScore" width="110" />
        <el-table-column align="left" label="创建日期" width="180">
          <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>
        <el-table-column align="left" label="操作" fixed="right" min-width="140">
          <template #default="scope">
            <el-button type="primary" link class="table-button" @click="getDetails(scope.row)">详情</el-button>
            <template v-if="!isReadOnlyReviewer">
              <el-button type="primary" link icon="edit" @click="updateExpertProfileFunc(scope.row)">编辑</el-button>
              <el-button
                v-if="['draft', 'org_rejected', 'city_rejected'].includes(scope.row.status)"
                type="primary" link @click="submitForReview(scope.row)"
              >提交审核</el-button>
              <el-button type="primary" link icon="delete" @click="deleteRow(scope.row)">删除</el-button>
            </template>
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

    <el-dialog v-model="dialogFormVisible" :before-close="closeDialog" :title="type === 'create' ? '添加专家' : '编辑专家'" width="60%">
      <el-form ref="elFormRef" :model="formData" :rules="rule" label-position="right" label-width="110px">
        <el-divider content-position="left">背景信息</el-divider>
        <el-row :gutter="20">
          <el-col :span="8"><el-form-item label="姓名" prop="name"><el-input v-model="formData.name" placeholder="请输入姓名" /></el-form-item></el-col>
          <el-col :span="8"><el-form-item label="性别" prop="gender"><el-input v-model="formData.gender" placeholder="请输入性别" /></el-form-item></el-col>
          <el-col :span="8"><el-form-item label="民族" prop="ethnicity"><el-input v-model="formData.ethnicity" placeholder="请输入民族" /></el-form-item></el-col>
          <el-col :span="8"><el-form-item label="政治面貌" prop="politicalStatus"><el-input v-model="formData.politicalStatus" /></el-form-item></el-col>
          <el-col :span="8">
            <el-form-item label="所在单位" prop="unitName">
              <el-select
                v-model="formData.unitName" filterable allow-create default-first-option
                placeholder="输入或选择所在单位" style="width: 100%;" @change="onUnitNameChange"
              >
                <el-option v-for="org in orgOptions" :key="org.id" :label="org.name" :value="org.name" />
              </el-select>
              <div class="unit-hint">单位列表里找不到？可以先直接输入，之后再去"单位管理"里补建并关联</div>
            </el-form-item>
          </el-col>
          <el-col :span="8"><el-form-item label="院系/部门" prop="department"><el-input v-model="formData.department" /></el-form-item></el-col>
          <el-col :span="8"><el-form-item label="行政职务" prop="adminTitle"><el-input v-model="formData.adminTitle" /></el-form-item></el-col>
          <el-col :span="8"><el-form-item label="专业技术职称" prop="techTitle"><el-input v-model="formData.techTitle" /></el-form-item></el-col>
          <el-col :span="8"><el-form-item label="办公电话" prop="phone"><el-input v-model="formData.phone" /></el-form-item></el-col>
          <el-col :span="8"><el-form-item label="手机号码" prop="mobile"><el-input v-model="formData.mobile" /></el-form-item></el-col>
          <el-col :span="8"><el-form-item label="电子邮箱" prop="email"><el-input v-model="formData.email" /></el-form-item></el-col>
          <el-col :span="8"><el-form-item label="通讯地址" prop="address"><el-input v-model="formData.address" /></el-form-item></el-col>
          <el-col :span="8"><el-form-item label="最高学历" prop="highestEducation"><el-input v-model="formData.highestEducation" /></el-form-item></el-col>
          <el-col :span="8"><el-form-item label="最高学位" prop="highestDegree"><el-input v-model="formData.highestDegree" /></el-form-item></el-col>
          <el-col :span="8"><el-form-item label="毕业院校" prop="graduateSchool"><el-input v-model="formData.graduateSchool" /></el-form-item></el-col>
          <el-col :span="8"><el-form-item label="所学专业" prop="major"><el-input v-model="formData.major" /></el-form-item></el-col>
          <el-col :span="8"><el-form-item label="人才计划" prop="talentProgram"><el-input v-model="formData.talentProgram" /></el-form-item></el-col>
          <el-col :span="8"><el-form-item label="荣誉称号" prop="honorTitle"><el-input v-model="formData.honorTitle" /></el-form-item></el-col>
        </el-row>

        <el-divider content-position="left">学科领域</el-divider>
        <el-row :gutter="20">
          <el-col :span="8"><el-form-item label="一级学科" prop="disciplineL1"><el-input v-model="formData.disciplineL1" /></el-form-item></el-col>
          <el-col :span="8"><el-form-item label="二级学科" prop="disciplineL2"><el-input v-model="formData.disciplineL2" /></el-form-item></el-col>
          <el-col :span="8"><el-form-item label="交叉学科领域" prop="crossDiscipline"><el-input v-model="formData.crossDiscipline" /></el-form-item></el-col>
          <el-col :span="24"><el-form-item label="学科平台/研究基地" prop="disciplinePlatform"><el-input v-model="formData.disciplinePlatform" /></el-form-item></el-col>
        </el-row>

        <el-divider content-position="left">研究主题</el-divider>
        <el-row :gutter="20">
          <el-col :span="24"><el-form-item label="重点研究方向" prop="researchDirections"><el-input v-model="formData.researchDirections" type="textarea" /></el-form-item></el-col>
          <el-col :span="24"><el-form-item label="研究关键词" prop="researchKeywords"><el-input v-model="formData.researchKeywords" placeholder="多个关键词用逗号分隔" /></el-form-item></el-col>
          <el-col :span="12"><el-form-item label="所属政策领域" prop="policyFields"><el-input v-model="formData.policyFields" /></el-form-item></el-col>
          <el-col :span="12"><el-form-item label="主要研究对象" prop="researchObjects"><el-input v-model="formData.researchObjects" /></el-form-item></el-col>
          <el-col :span="12"><el-form-item label="当前关注议题" prop="focusTopics"><el-input v-model="formData.focusTopics" /></el-form-item></el-col>
          <el-col :span="12"><el-form-item label="区域/国别研究专长" prop="regionExpertise"><el-input v-model="formData.regionExpertise" /></el-form-item></el-col>
          <el-col :span="24"><el-form-item label="研究方法专长" prop="methodExpertise"><el-input v-model="formData.methodExpertise" /></el-form-item></el-col>
        </el-row>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closeDialog">取消</el-button>
          <el-button type="primary" @click="enterDialog">确定</el-button>
        </div>
      </template>
    </el-dialog>

    <el-dialog v-model="importDialogVisible" title="批量导入专家" width="600px" :before-close="closeImportDialog">
      <div class="import-step">
        <p>1. 下载模板，按格式填写"专家背景信息"和"研究成果"两个 sheet（研究成果表按"专家姓名"关联背景信息表，姓名要能对上）</p>
        <el-button icon="download" @click="handleDownloadTemplate">下载模板</el-button>
      </div>
      <div class="import-step" style="margin-top: 20px;">
        <p>2. 上传填好的 .xlsx 文件</p>
        <el-upload
          drag
          :auto-upload="false"
          :limit="1"
          accept=".xlsx"
          :on-change="handleFileChange"
          :on-remove="() => { importFile = null }"
        >
          <el-icon class="el-icon--upload"><upload-filled /></el-icon>
          <div class="el-upload__text">拖到这里，或<em>点击选择文件</em></div>
        </el-upload>
      </div>
      <div v-if="importResult" style="margin-top: 20px;">
        <el-alert
          v-if="importResult.success"
          type="success" :closable="false"
          :title="`导入成功：新建专家 ${importResult.createdProfiles} 位，复用已有档案 ${importResult.reusedProfiles} 位，新增研究成果 ${importResult.createdAchievements} 条，标题重复跳过 ${importResult.skippedAchievements} 条`"
        />
        <el-alert
          v-if="importResult.success && importResult.unmatchedUnits && importResult.unmatchedUnits.length"
          type="warning" :closable="false" style="margin-top: 10px;"
          title="以下单位名在「单位管理」里找不到对应记录，已按空单位导入（不影响本次导入结果），建议去单位管理核实是不是打法不一致或者需要新建"
        >
          <div>{{ importResult.unmatchedUnits.join('、') }}</div>
        </el-alert>
        <template v-else>
          <el-alert type="error" :closable="false" title="校验未通过，以下问题需要改完重新上传，本次没有导入任何数据" />
          <el-table :data="importResult.errors" size="small" style="margin-top: 10px; max-height: 260px; overflow-y: auto;">
            <el-table-column label="sheet" prop="sheet" width="120" />
            <el-table-column label="行号" prop="row" width="70" />
            <el-table-column label="问题" prop="message" />
          </el-table>
        </template>
      </div>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closeImportDialog">关闭</el-button>
          <el-button type="primary" :disabled="!importFile" :loading="importing" @click="handleImportSubmit">开始导入</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import {
  createExpertProfile,
  deleteExpertProfile,
  deleteExpertProfileByIds,
  updateExpertProfile,
  findExpertProfile,
  getExpertProfileList,
  downloadImportTemplate,
  importExpertBatch,
  exportExpertProfiles
} from '@/api/expertProfile'
import { submitExpertProfile } from '@/api/expertApproval'
import { getSysOrganizationTree } from '@/api/sysOrganization'

import { formatDate } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/pinia/modules/user'

defineOptions({
  name: 'ExpertProfile'
})

const router = useRouter()
const userStore = useUserStore()
// 单位审核员/市级审核员在后端只有只读权限，写操作一定会被 Casbin 拒绝——
// 前端直接不展示这些按钮，避免审核员填完一整张表单才发现白填了
const isReadOnlyReviewer = computed(() => [9002, 9003].includes(userStore.userInfo.authorityId))

const statusOptions = [
  { label: '草稿', value: 'draft' },
  { label: '待单位审核', value: 'pending_org_review' },
  { label: '单位已退回', value: 'org_rejected' },
  { label: '待市级审核', value: 'pending_city_review' },
  { label: '市级已退回', value: 'city_rejected' },
  { label: '已发布', value: 'published' }
]
const statusLabel = (value) => statusOptions.find(item => item.value === value)?.label || value

const initFormData = () => ({
  name: '',
  gender: '',
  ethnicity: '',
  politicalStatus: '',
  unitName: '',
  orgId: null,
  department: '',
  adminTitle: '',
  techTitle: '',
  phone: '',
  mobile: '',
  email: '',
  address: '',
  highestEducation: '',
  highestDegree: '',
  graduateSchool: '',
  major: '',
  talentProgram: '',
  honorTitle: '',
  disciplineL1: '',
  disciplineL2: '',
  crossDiscipline: '',
  disciplinePlatform: '',
  researchDirections: '',
  researchKeywords: '',
  policyFields: '',
  researchObjects: '',
  focusTopics: '',
  regionExpertise: '',
  methodExpertise: ''
})

const formData = ref(initFormData())

// 所在单位下拉的候选项：来自单位管理维护的真实单位树，摊平成一维列表方便过滤/选择。
// 选中已有单位时顺带记下 orgId，写入档案后就能直接进对应单位的审核队列，不用再靠后端
// 按单位名做模糊匹配（这也是导入功能一直存在"单位匹配不上就成孤儿档案"问题的根源）
const orgOptions = ref([])
const flattenOrgTree = (nodes, acc = []) => {
  nodes.forEach(n => {
    // id=1 是顶层的"北京市哲学社会科学规划办公室"，是主管单位本身而不是专家挂靠的实际单位，
    // 跟后端 matchOrgId() 的排除逻辑保持一致，不作为可选单位出现
    if (n.name && n.ID !== 1) acc.push({ id: n.ID, name: n.name })
    if (n.children && n.children.length) flattenOrgTree(n.children, acc)
  })
  return acc
}
const fetchOrgOptions = async () => {
  const res = await getSysOrganizationTree()
  if (res.code === 0) {
    orgOptions.value = flattenOrgTree(res.data.tree || [])
  }
}
fetchOrgOptions()

// 选中的名字如果能在单位列表里找到，就带上 orgId；找不到（新输入的自由文本）就清空 orgId，
// 跟现在后端"org_id 为空时兜底交给市级审核"的逻辑保持一致，不强行卡死流程
const onUnitNameChange = (name) => {
  const matched = orgOptions.value.find(org => org.name === name)
  formData.value.orgId = matched ? matched.id : null
}

const rule = reactive({
  name: [
    { required: true, message: '请输入姓名', trigger: ['input', 'blur'] },
    { whitespace: true, message: '不能只输入空格', trigger: ['input', 'blur'] }
  ]
})

const elFormRef = ref()
const elSearchFormRef = ref()

// =========== 表格控制部分 ===========
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
  const table = await getExpertProfileList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
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
  }).then(() => {
    deleteExpertProfileFunc(row)
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
    const res = await deleteExpertProfileByIds({ IDs })
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

const updateExpertProfileFunc = async (row) => {
  const res = await findExpertProfile({ ID: row.ID })
  type.value = 'update'
  if (res.code === 0) {
    formData.value = res.data.reExpertProfile
    dialogFormVisible.value = true
  }
}

const deleteExpertProfileFunc = async (row) => {
  const res = await deleteExpertProfile({ ID: row.ID })
  if (res.code === 0) {
    ElMessage({ type: 'success', message: '删除成功' })
    if (tableData.value.length === 1 && page.value > 1) {
      page.value--
    }
    getTableData()
  }
}

// 改完信息后不用再跑一趟审核台才能提交——编辑完在这里就能直接提交审核
const submitForReview = async (row) => {
  const res = await submitExpertProfile({ expertId: row.ID })
  if (res.code === 0) {
    ElMessage({ type: 'success', message: '提交成功' })
    getTableData()
  }
}

const dialogFormVisible = ref(false)

const getDetails = (row) => {
  router.push({ name: 'expertProfileDetail', params: { id: row.ID } })
}

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
      res = await updateExpertProfile(formData.value)
    } else {
      res = await createExpertProfile(formData.value)
    }
    if (res.code === 0) {
      ElMessage({ type: 'success', message: '创建/更改成功' })
      closeDialog()
      getTableData()
    }
  })
}

// ============ 批量导入 ============
const importDialogVisible = ref(false)
const importFile = ref(null)
const importing = ref(false)
const importResult = ref(null)

const openImportDialog = () => {
  importFile.value = null
  importResult.value = null
  importDialogVisible.value = true
}
const closeImportDialog = () => {
  importDialogVisible.value = false
}
const handleFileChange = (file) => {
  importFile.value = file.raw
  importResult.value = null
}
const exporting = ref(false)
const handleExport = async () => {
  exporting.value = true
  try {
    // 导出跟当前列表页完全相同的筛选条件（不带分页参数，后端会返回全部匹配结果，不只是当前这一页）
    const res = await exportExpertProfiles(searchInfo.value)
    const blob = new Blob([res], { type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' })
    const url = window.URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = '专家列表导出.xlsx'
    a.click()
    window.URL.revokeObjectURL(url)
  } finally {
    exporting.value = false
  }
}

const handleDownloadTemplate = async () => {
  const res = await downloadImportTemplate()
  const blob = new Blob([res], { type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' })
  const url = window.URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = '专家批量导入模板.xlsx'
  a.click()
  window.URL.revokeObjectURL(url)
}
const handleImportSubmit = async () => {
  if (!importFile.value) return
  importing.value = true
  try {
    const formData = new FormData()
    formData.append('file', importFile.value)
    const res = await importExpertBatch(formData)
    if (res.code === 0) {
      ElMessage({ type: 'success', message: '导入成功' })
      importResult.value = res.data
      getTableData()
    } else {
      importResult.value = res.data
    }
  } finally {
    importing.value = false
  }
}
</script>

<style>
.unit-hint {
  font-size: 12px;
  color: #909399;
  line-height: 1.4;
  margin-top: 4px;
}
</style>
