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
        <exportExcel template-id="CGD" :condition="searchInfo" />
        <el-button type="primary" @click="goToCgdVis">自对比</el-button>
      </div>
      <el-table ref="multipleTable" style="width: 100%" tooltip-effect="dark" :data="tableData" row-key="ID"
        @selection-change="handleSelectionChange" @sort-change="sortChange">
        <el-table-column type="index" fixed width="50" />
        <el-table-column type="selection" fixed width="55" />

        <el-table-column align="left" label="日期" width="180">
          <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>

        <el-table-column align="left" label="学校名称" fixed prop="universityName" width="120" />
        <el-table-column label="专任教师" align="center" width="120">
          <el-table-column sortable align="left" label="专任教师数" prop="fullTimeFacultyCount" width="120">
            <template #default="{ row }">
              {{ formatCount(row.fullTimeFacultyCount) }}
            </template>
          </el-table-column>
        </el-table-column>
        <el-table-column label="高层次人才" align="center" width="480">
          <el-table-column sortable align="left" label="外籍高层次人才数" prop="foreignHighLevelTalentCount" width="120">
            <template #default="{ row }">
              {{ formatCount(row.foreignHighLevelTalentCount) }}
            </template>
          </el-table-column>
          <el-table-column sortable align="left" label="院士人数" prop="academicianCount" width="120">
            <template #default="{ row }">
              {{ formatCount(row.academicianCount) }}
            </template>
          </el-table-column>
          <el-table-column sortable align="left" label="国家级人才项目获得者人次" prop="nationalTalentProjectWinnerCount"
            width="120">
            <template #default="{ row }">
              {{ formatCount(row.nationalTalentProjectWinnerCount) }}
            </template>
          </el-table-column>
          <el-table-column sortable align="left" label="省部级人才项目获得者人次" prop="provincialTalentProjectWinnerCount"
            width="120">
            <template #default="{ row }">
              {{ formatCount(row.provincialTalentProjectWinnerCount) }}
            </template>
          </el-table-column>
        </el-table-column>
        <el-table-column label="基地数" align="center" width="120">
          <el-table-column sortable align="left" label="校外实习、实践、实训基地数" prop="offCampusInternshipBaseCount" width="120">
            <template #default="{ row }">
              {{ formatCount(row.offCampusInternshipBaseCount) }}
            </template>
          </el-table-column>
        </el-table-column>
        <el-table-column label="仪器设备" align="center" width="240">
          <el-table-column sortable align="left" label="大型仪器设备共享系统平台数" prop="largeEquipmentSharingPlatformCount"
            width="120">
            <template #default="{ row }">
              {{ formatCount(row.largeEquipmentSharingPlatformCount) }}
            </template>
          </el-table-column>
          <el-table-column sortable align="left" label="国家级实验教学示范中心数" prop="nationalExperimentalTeachingDemoCenterCount"
            width="120">
            <template #default="{ row }">
              {{ formatCount(row.nationalExperimentalTeachingDemoCenterCount) }}
            </template>
          </el-table-column>
          <el-table-column align="left" label="数据采集时间" prop="additionalRemarks" width="180">
            <template #default="scope">{{ formatDatemini(scope.row.acquisitionTime) }}</template>
          </el-table-column>
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
              @click="updateConditionalGuaranteeDatabaseFunc(scope.row)">变更</el-button>
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
        <el-form :model="formData" label-position="right" ref="elFormRef" :rules="rule" label-width="40%">
          <el-form-item label="学校名称:" prop="universityName">
            <el-input v-model="formData.universityName" :clearable="true" placeholder="请输入学校名称" style="width: 200px;" />
          </el-form-item>
          <div style="text-align: center; font-weight: bold; font-size: 16px; margin-top: 20px; margin-bottom: 20px;">
            专任教师
          </div>
          <el-form-item label="专任教师数:" prop="fullTimeFacultyCount">
            <el-input v-model.number="formData.fullTimeFacultyCount" :clearable="true" placeholder="请输入专任教师数"
              style="width: 200px;" />
          </el-form-item>
          <div style="text-align: center; font-weight: bold; font-size: 16px; margin-top: 20px; margin-bottom: 20px;">
            高层次人才
          </div>
          <el-form-item label="外籍高层次人才数:" prop="foreignHighLevelTalentCount">
            <el-input v-model.number="formData.foreignHighLevelTalentCount" :clearable="true" placeholder="请输入外籍高层次人才数"
              style="width: 200px;" />
          </el-form-item>
          <el-form-item label="院士人数:" prop="academicianCount">
            <el-input v-model.number="formData.academicianCount" :clearable="true" placeholder="请输入院士人数"
              style="width: 200px;" />
          </el-form-item>
          <el-form-item label="国家级人才项目获得者人次:" prop="nationalTalentProjectWinnerCount">
            <el-input v-model.number="formData.nationalTalentProjectWinnerCount" :clearable="true"
              placeholder="请输入国家级人才项目获得者人次" style="width: 200px;" />
          </el-form-item>
          <el-form-item label="省部级人才项目获得者人次:" prop="provincialTalentProjectWinnerCount">
            <el-input v-model.number="formData.provincialTalentProjectWinnerCount" :clearable="true"
              placeholder="请输入省部级人才项目获得者人次" style="width: 200px;" />
          </el-form-item>
          <div style="text-align: center; font-weight: bold; font-size: 16px; margin-top: 20px; margin-bottom: 20px;">
            基地数
          </div>
          <el-form-item label="校外实习、实践、实训基地数:" prop="offCampusInternshipBaseCount">
            <el-input v-model.number="formData.offCampusInternshipBaseCount" :clearable="true"
              placeholder="请输入校外实习、实践、实训基地数" style="width: 200px;" />
          </el-form-item>
          <div style="text-align: center; font-weight: bold; font-size: 16px; margin-top: 20px; margin-bottom: 20px;">
            仪器设备
          </div>
          <el-form-item label="大型仪器设备共享系统平台数:" prop="largeEquipmentSharingPlatformCount">
            <el-input v-model.number="formData.largeEquipmentSharingPlatformCount" :clearable="true"
              placeholder="请输入大型仪器设备共享系统平台数" style="width: 200px;" />
          </el-form-item>
          <el-form-item label="国家级实验教学示范中心数:" prop="nationalExperimentalTeachingDemoCenterCount">
            <el-input v-model.number="formData.nationalExperimentalTeachingDemoCenterCount" :clearable="true"
              placeholder="请输入国家级实验教学示范中心数" style="width: 200px;" />
          </el-form-item>
          <el-form-item label="数据采集时间:" prop="acquisitionTime">
            <el-date-picker v-model="formData.acquisitionTime" type="date" style="width: 200px;" placeholder="选择日期"
              :clearable="true" />
            <!-- <el-input v-model="formData.policyReleaseDate" :clearable="true" placeholder="请输入日期" /> -->
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
      <el-scrollbar height="550px">
        <el-descriptions :column="1" border>
          <el-descriptions-item label="学校名称">
            {{ formData.universityName }}
          </el-descriptions-item>
          <el-descriptions-item label="">
            <el-tag type="success" class="custom-tag">专任教师</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="专任教师数">
            <span v-if="formData.fullTimeFacultyCount < 0">暂无</span>
            <span v-else>{{ formData.fullTimeFacultyCount }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="">
            <el-tag type="success" class="custom-tag">高层次人才</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="外籍高层次人才数">
            <span v-if="formData.foreignHighLevelTalentCount < 0">暂无</span>
            <span v-else>{{ formData.foreignHighLevelTalentCount }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="院士人数">
            <span v-if="formData.academicianCount < 0">暂无</span>
            <span v-else>{{ formData.academicianCount }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="国家级人才项目获得者人次">
            <span v-if="formData.nationalTalentProjectWinnerCount < 0">暂无</span>
            <span v-else>{{ formData.nationalTalentProjectWinnerCount }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="省部级人才项目获得者人次">
            <span v-if="formData.provincialTalentProjectWinnerCount < 0">暂无</span>
            <span v-else>{{ formData.provincialTalentProjectWinnerCount }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="">
            <el-tag type="success" class="custom-tag">基地数</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="校外实习、实践、实训基地数">
            <span v-if="formData.offCampusInternshipBaseCount < 0">暂无</span>
            <span v-else>{{ formData.offCampusInternshipBaseCount }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="">
            <el-tag type="success" class="custom-tag">仪器设备</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="大型仪器设备共享系统平台数">
            <span v-if="formData.largeEquipmentSharingPlatformCount < 0">暂无</span>
            <span v-else>{{ formData.largeEquipmentSharingPlatformCount }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="国家级实验教学示范中心数">
            <span v-if="formData.nationalExperimentalTeachingDemoCenterCount < 0">暂无</span>
            <span v-else>{{ formData.nationalExperimentalTeachingDemoCenterCount }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="数据采集时间">
            {{ formatDatemini(formData.acquisitionTime) }}
          </el-descriptions-item>
        </el-descriptions>
      </el-scrollbar>
    </el-dialog>
  </div>
</template>

<script setup>
import {
  createConditionalGuaranteeDatabase,
  deleteConditionalGuaranteeDatabase,
  deleteConditionalGuaranteeDatabaseByIds,
  updateConditionalGuaranteeDatabase,
  findConditionalGuaranteeDatabase,
  getConditionalGuaranteeDatabaseList
} from '@/api/conditionalGuaranteeDatabase'

// 全量引入格式化工具 请按需保留
import { getDictFunc, formatDate, formatDatemini, formatBoolean, filterDict, ReturnArrImg, onDownloadFile } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'
import exportExcel from '@/components/exportExcel/exportExcel.vue'
import { useRouter } from 'vue-router'

defineOptions({
  name: 'ConditionalGuaranteeDatabase'
})

const formatCount = (value) => {
  return value < 0 ? '暂无' : value;
};


const router = useRouter();

// 定义跳转到 VisBid 的函数
const goToCgdVis = () => {
  router.push('/layout/visual/cgdVis');
};

// 自动化生成的字典（可能为空）以及字段
const formData = ref({
  universityName: '',
  fullTimeFacultyCount: -1,
  foreignHighLevelTalentCount: -1,
  academicianCount: -1,
  nationalTalentProjectWinnerCount: -1,
  provincialTalentProjectWinnerCount: -1,
  offCampusInternshipBaseCount: -1,
  largeEquipmentSharingPlatformCount: -1,
  nationalExperimentalTeachingDemoCenterCount: -1,
  acquisitionTime: new Date(),
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
  fullTimeFacultyCount: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
  ],
  foreignHighLevelTalentCount: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
  ],
  academicianCount: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
  ],
  nationalTalentProjectWinnerCount: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
  ],
  provincialTalentProjectWinnerCount: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
  ],
  offCampusInternshipBaseCount: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
  ],
  largeEquipmentSharingPlatformCount: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
  ],
  nationalExperimentalTeachingDemoCenterCount: [{
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
    fullTimeFacultyCount: 'full_time_faculty_count',
    foreignHighLevelTalentCount: 'foreign_high_level_talent_count',
    academicianCount: 'academician_count',
    nationalTalentProjectWinnerCount: 'national_talent_project_winner_count',
    provincialTalentProjectWinnerCount: 'provincial_talent_project_winner_count',
    offCampusInternshipBaseCount: 'off_campus_internship_base_count',
    largeEquipmentSharingPlatformCount: 'large_equipment_sharing_platform_count',
    nationalExperimentalTeachingDemoCenterCount: 'national_experimental_teaching_demo_center_count',
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
  const table = await getConditionalGuaranteeDatabaseList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
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
    deleteConditionalGuaranteeDatabaseFunc(row)
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
    const res = await deleteConditionalGuaranteeDatabaseByIds({ IDs })
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
const updateConditionalGuaranteeDatabaseFunc = async (row) => {
  const res = await findConditionalGuaranteeDatabase({ ID: row.ID })
  type.value = 'update'
  if (res.code === 0) {
    formData.value = res.data.reCGD
    dialogFormVisible.value = true
  }
}


// 删除行
const deleteConditionalGuaranteeDatabaseFunc = async (row) => {
  const res = await deleteConditionalGuaranteeDatabase({ ID: row.ID })
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
  const res = await findConditionalGuaranteeDatabase({ ID: row.ID })
  if (res.code === 0) {
    formData.value = res.data.reCGD
    openDetailShow()
  }
}


// 关闭详情弹窗
const closeDetailShow = () => {
  detailShow.value = false
  formData.value = {
    universityName: '',
    fullTimeFacultyCount: -1,
    foreignHighLevelTalentCount: -1,
    academicianCount: -1,
    nationalTalentProjectWinnerCount: -1,
    provincialTalentProjectWinnerCount: -1,
    offCampusInternshipBaseCount: -1,
    largeEquipmentSharingPlatformCount: -1,
    nationalExperimentalTeachingDemoCenterCount: -1,
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
    universityName: '',
    fullTimeFacultyCount: -1,
    foreignHighLevelTalentCount: -1,
    academicianCount: -1,
    nationalTalentProjectWinnerCount: -1,
    provincialTalentProjectWinnerCount: -1,
    offCampusInternshipBaseCount: -1,
    largeEquipmentSharingPlatformCount: -1,
    nationalExperimentalTeachingDemoCenterCount: -1,
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
        res = await createConditionalGuaranteeDatabase(formData.value)
        break
      case 'update':
        res = await updateConditionalGuaranteeDatabase(formData.value)
        break
      default:
        res = await createConditionalGuaranteeDatabase(formData.value)
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
