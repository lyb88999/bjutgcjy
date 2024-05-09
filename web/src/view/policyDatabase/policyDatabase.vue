<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline" :rules="searchRule"
        @keyup.enter="onSubmit">
        <el-form-item label="创建日期" prop="createdAt">
          <template #label>
            <span>
              创建日期
              <el-tooltip content="搜索范围是开始日期（包含）至结束日期（不包含）">
                <el-icon>
                  <QuestionFilled />
                </el-icon>
              </el-tooltip>
            </span>
          </template>
          <el-date-picker v-model="searchInfo.startCreatedAt" type="datetime" placeholder="开始日期"
            :disabled-date="time => searchInfo.endCreatedAt ? time.getTime() > searchInfo.endCreatedAt.getTime() : false"></el-date-picker>
          —
          <el-date-picker v-model="searchInfo.endCreatedAt" type="datetime" placeholder="结束日期"
            :disabled-date="time => searchInfo.startCreatedAt ? time.getTime() < searchInfo.startCreatedAt.getTime() : false"></el-date-picker>
        </el-form-item>

        <!-- <el-form-item label="政策发布时间" prop="policyReleaseDate">

          <template #label>
            <span>
              政策发布时间
              <el-tooltip content="搜索范围是开始日期（包含）至结束日期（不包含）">
                <el-icon>
                  <QuestionFilled />
                </el-icon>
              </el-tooltip>
            </span>
          </template>
          <el-date-picker v-model="searchInfo.startPolicyReleaseDate" type="datetime" placeholder="开始日期"
            :disabled-date="time => searchInfo.endPolicyReleaseDate ? time.getTime() > searchInfo.endPolicyReleaseDate.getTime() : false"></el-date-picker>
          —
          <el-date-picker v-model="searchInfo.endPolicyReleaseDate" type="datetime" placeholder="结束日期"
            :disabled-date="time => searchInfo.startPolicyReleaseDate ? time.getTime() < searchInfo.startPolicyReleaseDate.getTime() : false"></el-date-picker>

        </el-form-item> -->
        <el-form-item label="政策国家或地区" prop="policyCountryOrRegion">
          <el-input v-model="searchInfo.policyCountryOrRegion" placeholder="搜索条件" />

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
        <el-button icon="delete" style="margin-left: 10px;" :disabled="!multipleSelection.length"
          @click="onDelete">删除</el-button>
          <exportExcel template-id="PD" :condition="searchInfo"/>
      </div>
      <el-table ref="multipleTable" style="width: 100%" tooltip-effect="dark" :data="tableData" row-key="ID"
        @selection-change="handleSelectionChange" @sort-change="sortChange">
        <el-table-column type="selection" width="55" />

        <el-table-column align="left" label="日期" width="180">
          <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>

        <el-table-column label="政策法规" align="center" width="1040">
          <el-table-column align="left" label="政策名称" prop="policyName" width="120" />
          <el-table-column align="left" label="政策发布时间" prop="policyReleaseDate" width="180">
            <!-- <template #default="scope">{{ formatDate(scope.row.policyReleaseDate) }}</template> -->
          </el-table-column>
          <el-table-column align="left" label="政策国家或地区" prop="policyCountryOrRegion" width="120" />
          <el-table-column align="left" label="数据来源" prop="dataSource" width="120" />
          <el-table-column label="政策内容" width="200">
            详情点击右侧查看详情
          </el-table-column>
          <el-table-column align="left" label="数据采集时间" prop="additionalRemarks" width="180">
            <template #default="scope">{{ formatDatemini(scope.row.acquisitionTime) }}</template>
          </el-table-column>
          <el-table-column align="left" label="其他备注" prop="additionalRemarks" width="120" />
        </el-table-column>
        <el-table-column align="left" label="操作" fixed="right" min-width="240">
          <template #default="scope">
            <el-button type="primary" link class="table-button" @click="getDetails(scope.row)">
              <el-icon style="margin-right: 5px">
                <InfoFilled />
              </el-icon>
              查看详情
            </el-button>
            <el-button type="primary" link icon="edit" class="table-button"
              @click="updatePolicyDatabaseFunc(scope.row)">变更</el-button>
            <el-button type="primary" link icon="delete" @click="deleteRow(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination layout="total, sizes, prev, pager, next, jumper" :current-page="page" :page-size="pageSize"
          :page-sizes="[10, 30, 50, 100]" :total="total" @current-change="handleCurrentChange"
          @size-change="handleSizeChange" />
      </div>
    </div>
    <el-dialog v-model="dialogFormVisible" :before-close="closeDialog" :title="type === 'create' ? '添加' : '修改'"
      destroy-on-close width="600px">
      <el-scrollbar height="500px">
        <el-form :model="formData" label-position="right" ref="elFormRef" :rules="rule" label-width="130px">
          <div style="text-align: center; font-weight: bold; font-size: 16px; margin-top: 20px; margin-bottom: 20px;">
            政策法规
          </div>
          <el-form-item label="政策名称:" prop="policyName">
            <el-input v-model="formData.policyName" :clearable="true" placeholder="请输入政策名称" />
          </el-form-item>
          <el-form-item label="政策发布时间:" prop="policyReleaseDate">
            <!-- <el-date-picker v-model="formData.policyReleaseDate" type="date" style="width:100%" placeholder="选择日期"
              :clearable="true" /> -->
            <el-input v-model="formData.policyReleaseDate" :clearable="true" placeholder="请输入日期" />
          </el-form-item>
          <el-form-item label="政策国家或地区:" prop="policyCountryOrRegion">
            <el-input v-model="formData.policyCountryOrRegion" :clearable="true" placeholder="请输入政策国家或地区" />
          </el-form-item>
          <el-form-item label="数据来源:" prop="dataSource">
            <el-input v-model="formData.dataSource" :clearable="true" placeholder="请输入数据来源" />
          </el-form-item>
          <el-form-item label="政策内容:" prop="policyContent">
            <RichEdit v-model="formData.policyContent" />
          </el-form-item>
          <el-form-item label="数据采集时间:" prop="acquisitionTime">
            <el-date-picker v-model="formData.acquisitionTime" type="date" style="width:100%" placeholder="选择日期"
              :clearable="true" />
            <!-- <el-input v-model="formData.policyReleaseDate" :clearable="true" placeholder="请输入日期" /> -->
          </el-form-item>
          <el-form-item label="其他备注:" prop="additionalRemarks">
            <el-input v-model="formData.additionalRemarks" :clearable="true" placeholder="请输入其他备注" />
          </el-form-item>
        </el-form>
      </el-scrollbar>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closeDialog">取 消</el-button>
          <el-button type="primary" @click="enterDialog">确 定</el-button>
        </div>
      </template>
    </el-dialog>

    <el-dialog v-model="detailShow" style="width: 800px" lock-scroll :before-close="closeDetailShow" title="查看详情"
      destroy-on-close>
      <el-scrollbar height="450px">
        <el-descriptions :column="1" border>
          <el-descriptions-item label="">
            <el-tag type="success" class="custom-tag">政策法规</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="政策名称">
            {{ formData.policyName }}
          </el-descriptions-item>
          <el-descriptions-item label="政策发布时间">
            <!-- {{ formatDate(formData.policyReleaseDate) }} -->
            {{ formData.policyReleaseDate }}
          </el-descriptions-item>
          <el-descriptions-item label="政策国家或地区">
            {{ formData.policyCountryOrRegion }}
          </el-descriptions-item>
          <el-descriptions-item label="数据来源">
            {{ formData.dataSource }}
          </el-descriptions-item>
          <el-descriptions-item label="政策内容">
            <div v-html="formData.policyContent"></div>
          </el-descriptions-item>
          <el-descriptions-item label="数据采集时间">
            {{ formatDatemini(formData.acquisitionTime) }}
          </el-descriptions-item>
          <el-descriptions-item label="其他备注">
            {{ formData.additionalRemarks }}
          </el-descriptions-item>
        </el-descriptions>
      </el-scrollbar>
    </el-dialog>
  </div>
</template>

<script setup>
import {
  createPolicyDatabase,
  deletePolicyDatabase,
  deletePolicyDatabaseByIds,
  updatePolicyDatabase,
  findPolicyDatabase,
  getPolicyDatabaseList
} from '@/api/policyDatabase'
// 富文本组件
import RichEdit from '@/components/richtext/rich-edit.vue'

// 全量引入格式化工具 请按需保留
import { getDictFunc, formatDate, formatDatemini, formatBoolean, filterDict, ReturnArrImg, onDownloadFile } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'
import exportExcel from '@/components/exportExcel/exportExcel.vue'

defineOptions({
  name: 'PolicyDatabase'
})

// 自动化生成的字典（可能为空）以及字段
const formData = ref({
  policyName: '',
  policyReleaseDate: '',
  policyCountryOrRegion: '',
  dataSource: '',
  policyContent: '',
  additionalRemarks: '',
  acquisitionTime: new Date(),
})


// 验证规则
const rule = reactive({
  policyName: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
  {
    whitespace: true,
    message: '不能只输入空格',
    trigger: ['input', 'blur'],
  }
  ],
  policyReleaseDate: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
  ],
  policyCountryOrRegion: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
  {
    whitespace: true,
    message: '不能只输入空格',
    trigger: ['input', 'blur'],
  }
  ],
  dataSource: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
  {
    whitespace: true,
    message: '不能只输入空格',
    trigger: ['input', 'blur'],
  }
  ],
  policyContent: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
  ],
})

const searchRule = reactive({
  createdAt: [
    {
      validator: (rule, value, callback) => {
        if (searchInfo.value.startCreatedAt && !searchInfo.value.endCreatedAt) {
          callback(new Error('请填写结束日期'))
        } else if (!searchInfo.value.startCreatedAt && searchInfo.value.endCreatedAt) {
          callback(new Error('请填写开始日期'))
        } else if (searchInfo.value.startCreatedAt && searchInfo.value.endCreatedAt && (searchInfo.value.startCreatedAt.getTime() === searchInfo.value.endCreatedAt.getTime() || searchInfo.value.startCreatedAt.getTime() > searchInfo.value.endCreatedAt.getTime())) {
          callback(new Error('开始日期应当早于结束日期'))
        } else {
          callback()
        }
      }, trigger: 'change'
    }
  ],
  // policyReleaseDate: [{
  //   validator: (rule, value, callback) => {
  //     if (searchInfo.value.startPolicyReleaseDate && !searchInfo.value.endPolicyReleaseDate) {
  //       callback(new Error('请填写结束日期'))
  //     } else if (!searchInfo.value.startPolicyReleaseDate && searchInfo.value.endPolicyReleaseDate) {
  //       callback(new Error('请填写开始日期'))
  //     } else if (searchInfo.value.startPolicyReleaseDate && searchInfo.value.endPolicyReleaseDate && (searchInfo.value.startPolicyReleaseDate.getTime() === searchInfo.value.endPolicyReleaseDate.getTime() || searchInfo.value.startPolicyReleaseDate.getTime() > searchInfo.value.endPolicyReleaseDate.getTime())) {
  //       callback(new Error('开始日期应当早于结束日期'))
  //     } else {
  //       callback()
  //     }
  //   }, trigger: 'change'
  // }],
})

const elFormRef = ref()
const elSearchFormRef = ref()

// =========== 表格控制部分 ===========
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})
// 排序
const sortChange = ({ prop, order }) => {
  const sortMap = {
    policyReleaseDate: 'policy_release_date',
  }

  let sort = sortMap[prop]
  if (!sort) {
    sort = prop.replace(/[A-Z]/g, match => `_${match.toLowerCase()}`)
  }

  searchInfo.value.sort = sort
  searchInfo.value.order = order
  getTableData()
}

// 重置
const onReset = () => {
  searchInfo.value = {}
  getTableData()
}

// 搜索
const onSubmit = () => {
  elSearchFormRef.value?.validate(async (valid) => {
    if (!valid) return
    page.value = 1
    pageSize.value = 10
    getTableData()
  })
}

// 分页
const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

// 修改页面容量
const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

// 查询
const getTableData = async () => {
  const table = await getPolicyDatabaseList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}

getTableData()

// ============== 表格控制部分结束 ===============

// 获取需要的字典 可能为空 按需保留
const setOptions = async () => {
}

// 获取需要的字典 可能为空 按需保留
setOptions()


// 多选数据
const multipleSelection = ref([])
// 多选
const handleSelectionChange = (val) => {
  multipleSelection.value = val
}

// 删除行
const deleteRow = (row) => {
  ElMessageBox.confirm('确定要删除吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    deletePolicyDatabaseFunc(row)
  })
}

// 多选删除
const onDelete = async () => {
  ElMessageBox.confirm('确定要删除吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    const IDs = []
    if (multipleSelection.value.length === 0) {
      ElMessage({
        type: 'warning',
        message: '请选择要删除的数据'
      })
      return
    }
    multipleSelection.value &&
      multipleSelection.value.map(item => {
        IDs.push(item.ID)
      })
    const res = await deletePolicyDatabaseByIds({ IDs })
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: '删除成功'
      })
      if (tableData.value.length === IDs.length && page.value > 1) {
        page.value--
      }
      getTableData()
    }
  })
}

// 行为控制标记（弹窗内部需要增还是改）
const type = ref('')

// 更新行
const updatePolicyDatabaseFunc = async (row) => {
  const res = await findPolicyDatabase({ ID: row.ID })
  type.value = 'update'
  if (res.code === 0) {
    formData.value = res.data.rePD
    dialogFormVisible.value = true
  }
}


// 删除行
const deletePolicyDatabaseFunc = async (row) => {
  const res = await deletePolicyDatabase({ ID: row.ID })
  if (res.code === 0) {
    ElMessage({
      type: 'success',
      message: '删除成功'
    })
    if (tableData.value.length === 1 && page.value > 1) {
      page.value--
    }
    getTableData()
  }
}

// 弹窗控制标记
const dialogFormVisible = ref(false)


// 查看详情控制标记
const detailShow = ref(false)


// 打开详情弹窗
const openDetailShow = () => {
  detailShow.value = true
}


// 打开详情
const getDetails = async (row) => {
  // 打开弹窗
  const res = await findPolicyDatabase({ ID: row.ID })
  if (res.code === 0) {
    formData.value = res.data.rePD
    openDetailShow()
  }
}


// 关闭详情弹窗
const closeDetailShow = () => {
  detailShow.value = false
  formData.value = {
    policyName: '',
    policyReleaseDate: '',
    policyCountryOrRegion: '',
    dataSource: '',
    additionalRemarks: '',
    acquisitionTime: new Date(),
  }
}


// 打开弹窗
const openDialog = () => {
  type.value = 'create'
  dialogFormVisible.value = true
}

// 关闭弹窗
const closeDialog = () => {
  dialogFormVisible.value = false
  formData.value = {
    policyName: '',
    policyReleaseDate: '',
    policyCountryOrRegion: '',
    dataSource: '',
    additionalRemarks: '',
    acquisitionTime: new Date(),
  }
}
// 弹窗确定
const enterDialog = async () => {
  elFormRef.value?.validate(async (valid) => {
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
      closeDialog()
      getTableData()
    }
  })
}

</script>

<style></style>
