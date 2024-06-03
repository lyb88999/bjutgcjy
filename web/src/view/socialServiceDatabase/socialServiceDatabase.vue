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

        <el-form-item label="学校名称" prop="universityName">
          <el-input v-model="searchInfo.universityName" placeholder="搜索条件" />

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
          <exportExcel template-id="SSD" :condition="searchInfo"/>
          <el-button type="primary" @click="goTossdVis">自对比</el-button>
      </div>
      <el-table ref="multipleTable" style="width: 100%" tooltip-effect="dark" :data="tableData" row-key="ID"
        @selection-change="handleSelectionChange" @sort-change="sortChange">
        <el-table-column type="index" fixed width="50" />
        <el-table-column type="selection" fixed width="55" />
        <el-table-column align="left" label="日期" width="180">
          <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>

        <el-table-column sortable align="left" label="学校名称" fixed prop="universityName" width="120" />
        <el-table-column label="产学研基地" align="center" width="400">
          <el-table-column sortable align="left" label="行业协同创新中心数或产学研基地数" prop="industryInnovationCenterCount"
            width="200" />
          <el-table-column sortable align="left" label="区域协同创新中心数或产学研基地数" prop="regionalInnovationCenterCount"
            width="200" />
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
              @click="updateSocialServiceDatabaseFunc(scope.row)">变更</el-button>
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
      <el-scrollbar height="180px">
        <el-form :model="formData" label-position="right" ref="elFormRef" :rules="rule" label-width="50%">
          <div style="text-align: center; font-weight: bold; font-size: 16px; margin-top: 20px; margin-bottom: 20px;">
            产学研基地
          </div>
          <el-form-item label="学校名称:" prop="universityName">
            <el-input v-model="formData.universityName" :clearable="true" placeholder="请输入学校名称" style="width: 200px;" />
          </el-form-item>
          <el-form-item label="行业协同创新中心数或产学研基地数:" prop="industryInnovationCenterCount">
            <el-input v-model.number="formData.industryInnovationCenterCount" :clearable="true"
              placeholder="请输入行业协同创新中心数或产学研基地数" style="width: 200px;" />
          </el-form-item>
          <el-form-item label="区域协同创新中心数或产学研基地数:" prop="regionalInnovationCenterCount">
            <el-input v-model.number="formData.regionalInnovationCenterCount" :clearable="true"
              placeholder="请输入区域协同创新中心数或产学研基地数" style="width: 200px;" />
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
      <el-scrollbar height="200px">
        <el-descriptions :column="1" border>

          <el-descriptions-item label="学校名称">
            {{ formData.universityName }}
          </el-descriptions-item>
          <el-descriptions-item label="">
            <el-tag type="success" class="custom-tag">产学研基地</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="行业协同创新中心数或产学研基地数">
            {{ formData.industryInnovationCenterCount }}
          </el-descriptions-item>
          <el-descriptions-item label="区域协同创新中心数或产学研基地数">
            {{ formData.regionalInnovationCenterCount }}
          </el-descriptions-item>
        </el-descriptions>
      </el-scrollbar>
    </el-dialog>
  </div>
</template>

<script setup>
import {
  createSocialServiceDatabase,
  deleteSocialServiceDatabase,
  deleteSocialServiceDatabaseByIds,
  updateSocialServiceDatabase,
  findSocialServiceDatabase,
  getSocialServiceDatabaseList
} from '@/api/socialServiceDatabase'

// 全量引入格式化工具 请按需保留
import { getDictFunc, formatDate, formatBoolean, filterDict, ReturnArrImg, onDownloadFile } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'
import exportExcel from '@/components/exportExcel/exportExcel.vue'
import { useRouter } from 'vue-router'



defineOptions({
  name: 'SocialServiceDatabase'
})

const router = useRouter();

// 定义跳转到 VisBid 的函数
const goTossdVis = () => {
  router.push('/layout/visual/ssdVis');
};

// 自动化生成的字典（可能为空）以及字段
const formData = ref({
  universityName: '',
  industryInnovationCenterCount: 0,
  regionalInnovationCenterCount: 0,
})


// 验证规则
const rule = reactive({
  universityName: [{
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
  industryInnovationCenterCount: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
  ],
  regionalInnovationCenterCount: [{
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
    universityName: 'university_name',
    industryInnovationCenterCount: 'industry_innovation_center_count',
    regionalInnovationCenterCount: 'regional_innovation_center_count',
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
  const table = await getSocialServiceDatabaseList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
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
    deleteSocialServiceDatabaseFunc(row)
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
    const res = await deleteSocialServiceDatabaseByIds({ IDs })
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
const updateSocialServiceDatabaseFunc = async (row) => {
  const res = await findSocialServiceDatabase({ ID: row.ID })
  type.value = 'update'
  if (res.code === 0) {
    formData.value = res.data.reSSD
    dialogFormVisible.value = true
  }
}


// 删除行
const deleteSocialServiceDatabaseFunc = async (row) => {
  const res = await deleteSocialServiceDatabase({ ID: row.ID })
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
  const res = await findSocialServiceDatabase({ ID: row.ID })
  if (res.code === 0) {
    formData.value = res.data.reSSD
    openDetailShow()
  }
}


// 关闭详情弹窗
const closeDetailShow = () => {
  detailShow.value = false
  formData.value = {
    universityName: '',
    industryInnovationCenterCount: 0,
    regionalInnovationCenterCount: 0,
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
    universityName: '',
    industryInnovationCenterCount: 0,
    regionalInnovationCenterCount: 0,
  }
}
// 弹窗确定
const enterDialog = async () => {
  elFormRef.value?.validate(async (valid) => {
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
      closeDialog()
      getTableData()
    }
  })
}

</script>

<style></style>
