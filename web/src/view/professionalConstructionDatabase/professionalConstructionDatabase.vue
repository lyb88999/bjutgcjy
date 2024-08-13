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

        <el-form-item label="专业名称" prop="majorName">
          <el-input v-model="searchInfo.majorName" placeholder="搜索条件" />

        </el-form-item>
        <el-form-item label="专业代码" prop="majorCode">
          <el-input v-model="searchInfo.majorCode" placeholder="搜索条件" />

        </el-form-item>
        <el-form-item label="授予学位类别" prop="degreeType">
          <el-input v-model="searchInfo.degreeType" placeholder="搜索条件" />

        </el-form-item>
        <el-form-item label="是否新专业" prop="isNewMajor">
          <el-select v-model="searchInfo.isNewMajor" clearable placeholder="请选择">
            <el-option key="true" label="是" value="true">
            </el-option>
            <el-option key="false" label="否" value="false">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="是否大类招生" prop="isBroadAdmissionCategory">
          <el-select v-model="searchInfo.isBroadAdmissionCategory" clearable placeholder="请选择">
            <el-option key="true" label="是" value="true">
            </el-option>
            <el-option key="false" label="否" value="false">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="是否一流专业" prop="isFirstClassMajor">
          <el-select v-model="searchInfo.isFirstClassMajor" clearable placeholder="请选择">
            <el-option key="true" label="是" value="true">
            </el-option>
            <el-option key="false" label="否" value="false">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="是否专业认证" prop="isMajorAccredited">
          <el-select v-model="searchInfo.isMajorAccredited" clearable placeholder="请选择">
            <el-option key="true" label="是" value="true">
            </el-option>
            <el-option key="false" label="否" value="false">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="是否入选卓越计划" prop="isSelectedForExcellenceProgram">
          <el-select v-model="searchInfo.isSelectedForExcellenceProgram" clearable placeholder="请选择">
            <el-option key="true" label="是" value="true">
            </el-option>
            <el-option key="false" label="否" value="false">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="是否双学位" prop="isDoubleDegree">
          <el-select v-model="searchInfo.isDoubleDegree" clearable placeholder="请选择">
            <el-option key="true" label="是" value="true">
            </el-option>
            <el-option key="false" label="否" value="false">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="是否新工科" prop="isNewEngineeringDiscipline">
          <el-select v-model="searchInfo.isNewEngineeringDiscipline" clearable placeholder="请选择">
            <el-option key="true" label="是" value="true">
            </el-option>
            <el-option key="false" label="否" value="false">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="专业满意度" prop="majorSatisfaction">
          <el-input v-model="searchInfo.majorSatisfaction" placeholder="搜索条件" />

        </el-form-item>
        <el-form-item label="所在学校" prop="universityName">
          <el-input v-model="searchInfo.universityName" placeholder="搜索条件" />

        </el-form-item>
        <el-form-item label="招生专业" prop="admissionMajor">
          <el-input v-model="searchInfo.admissionMajor" placeholder="搜索条件" />

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
        <exportExcel template-id="PCD" :condition="searchInfo" />
      </div>
      <el-table ref="multipleTable" style="width: 100%" tooltip-effect="dark" height="500" :data="tableData"
        row-key="ID" @selection-change="handleSelectionChange" @sort-change="sortChange">
        <el-table-column type="index" fixed width="50" />
        <el-table-column type="selection" fixed width="55" />

        <el-table-column align="left" label="日期" width="180">
          <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>
        <el-table-column align="left" label="学校名称" prop="universityName" fixed width="120" />
        <el-table-column label="专业信息" align="center" width="1440">
          <el-table-column align="left" label="专业名称" prop="majorName" width="120" />
          <el-table-column sortable align="left" label="专业代码" prop="majorCode" width="120" />
          <el-table-column align="left" label="授予学位类别" prop="degreeType" width="120" />
          <el-table-column align="left" label="是否新专业" prop="isNewMajor" width="120">
            <template #default="scope">{{ formatBoolean(scope.row.isNewMajor) }}</template>
          </el-table-column>
          <el-table-column align="left" label="是否大类招生" prop="isBroadAdmissionCategory" width="120">
            <template #default="scope">{{ formatBoolean(scope.row.isBroadAdmissionCategory) }}</template>
          </el-table-column>
          <el-table-column align="left" label="是否一流专业" prop="isFirstClassMajor" width="120">
            <template #default="scope">{{ formatBoolean(scope.row.isFirstClassMajor) }}</template>
          </el-table-column>
          <el-table-column align="left" label="是否专业认证" prop="isMajorAccredited" width="120">
            <template #default="scope">{{ formatBoolean(scope.row.isMajorAccredited) }}</template>
          </el-table-column>
          <el-table-column align="left" label="是否入选卓越计划" prop="isSelectedForExcellenceProgram" width="120">
            <template #default="scope">{{ formatBoolean(scope.row.isSelectedForExcellenceProgram) }}</template>
          </el-table-column>
          <el-table-column align="left" label="是否双学位" prop="isDoubleDegree" width="120">
            <template #default="scope">{{ formatBoolean(scope.row.isDoubleDegree) }}</template>
          </el-table-column>
          <el-table-column align="left" label="是否新工科" prop="isNewEngineeringDiscipline" width="120">
            <template #default="scope">{{ formatBoolean(scope.row.isNewEngineeringDiscipline) }}</template>
          </el-table-column>
          <el-table-column align="left" label="专业满意度" prop="majorSatisfaction" width="120" />
          <el-table-column align="left" label="所在学校" prop="universityName" width="120" />
        </el-table-column>
        <el-table-column label="招生信息" align="center" width="240">
          <el-table-column align="left" label="招生专业" prop="admissionMajor" width="120" />
          <el-table-column align="left" label="计划招生数" prop="plannedEnrollment" width="120">
            <template #default="{ row }">
              {{ formatCount(row.plannedEnrollment) }}
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
              @click="updateProfessionalConstructionDatabaseFunc(scope.row)">变更</el-button>
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
          <!-- <el-form-item label="专业信息" class="form-section-title"></el-form-item> -->
          <div style="text-align: center; font-weight: bold; font-size: 16px; margin-bottom: 20px;">
            专业信息
          </div>
          <el-form-item label="专业名称:" prop="majorName">
            <el-input v-model="formData.majorName" :clearable="true" placeholder="请输入专业名称" style="width: 200px;" />
          </el-form-item>
          <el-form-item label="专业代码:" prop="majorCode">
            <el-input v-model="formData.majorCode" :clearable="true" placeholder="请输入专业代码" style="width: 200px;" />
          </el-form-item>
          <el-form-item label="授予学位类别:" prop="degreeType">
            <el-select v-model="formData.degreeType" placeholder="请选择授予学位类别" style="width:200px" :clearable="true">
              <el-option v-for="item in ['学士', '硕士', '博士']" :key="item" :label="item" :value="item" />
            </el-select>
          </el-form-item>
          <el-form-item label="是否新专业:" prop="isNewMajor">
            <el-switch v-model="formData.isNewMajor" active-color="#13ce66" inactive-color="#ff4949" active-text="是"
              inactive-text="否" clearable></el-switch>
          </el-form-item>
          <el-form-item label="是否大类招生:" prop="isBroadAdmissionCategory">
            <el-switch v-model="formData.isBroadAdmissionCategory" active-color="#13ce66" inactive-color="#ff4949"
              active-text="是" inactive-text="否" clearable></el-switch>
          </el-form-item>
          <el-form-item label="是否一流专业:" prop="isFirstClassMajor">
            <el-switch v-model="formData.isFirstClassMajor" active-color="#13ce66" inactive-color="#ff4949"
              active-text="是" inactive-text="否" clearable></el-switch>
          </el-form-item>
          <el-form-item label="是否专业认证:" prop="isMajorAccredited">
            <el-switch v-model="formData.isMajorAccredited" active-color="#13ce66" inactive-color="#ff4949"
              active-text="是" inactive-text="否" clearable></el-switch>
          </el-form-item>
          <el-form-item label="是否入选卓越计划:" prop="isSelectedForExcellenceProgram">
            <el-switch v-model="formData.isSelectedForExcellenceProgram" active-color="#13ce66" inactive-color="#ff4949"
              active-text="是" inactive-text="否" clearable></el-switch>
          </el-form-item>
          <el-form-item label="是否双学位:" prop="isDoubleDegree">
            <el-switch v-model="formData.isDoubleDegree" active-color="#13ce66" inactive-color="#ff4949" active-text="是"
              inactive-text="否" clearable></el-switch>
          </el-form-item>
          <el-form-item label="是否新工科:" prop="isNewEngineeringDiscipline">
            <el-switch v-model="formData.isNewEngineeringDiscipline" active-color="#13ce66" inactive-color="#ff4949"
              active-text="是" inactive-text="否" clearable></el-switch>
          </el-form-item>
          <el-form-item label="专业满意度:" prop="majorSatisfaction">
            <el-select v-model="formData.majorSatisfaction" placeholder="请选择专业满意度" style="width:200px"
              :clearable="true">
              <el-option v-for="item in ['非常不满意', '不满意', '中等', '满意', '非常满意']" :key="item" :label="item" :value="item" />
            </el-select>
          </el-form-item>
          <el-form-item label="所在学校:" prop="universityName">
            <el-input v-model="formData.universityName" :clearable="true" placeholder="请输入所在学校" style="width: 200px;" />
          </el-form-item>
          <div style="text-align: center; font-weight: bold; font-size: 16px; margin-top: 20px; margin-bottom: 20px;">
            招生信息
          </div>
          <el-form-item label="招生专业:" prop="admissionMajor">
            <el-input v-model="formData.admissionMajor" :clearable="true" placeholder="请输入招生专业" style="width: 200px;" />
          </el-form-item>
          <el-form-item label="计划招生数:" prop="plannedEnrollment">
            <el-input v-model.number="formData.plannedEnrollment" :clearable="true" placeholder="请输入计划招生数"
              style="width: 200px;" />
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
          <el-descriptions-item label="">
            <el-tag type="success" class="custom-tag">专业信息</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="专业名称">
            {{ formData.majorName }}
          </el-descriptions-item>
          <el-descriptions-item label="专业代码">
            {{ formData.majorCode }}
          </el-descriptions-item>
          <el-descriptions-item label="授予学位类别">
            {{ formData.degreeType }}
          </el-descriptions-item>
          <el-descriptions-item label="是否新专业">
            {{ formatBoolean(formData.isNewMajor) }}
          </el-descriptions-item>
          <el-descriptions-item label="是否大类招生">
            {{ formatBoolean(formData.isBroadAdmissionCategory) }}
          </el-descriptions-item>
          <el-descriptions-item label="是否一流专业">
            {{ formatBoolean(formData.isFirstClassMajor) }}
          </el-descriptions-item>
          <el-descriptions-item label="是否专业认证">
            {{ formatBoolean(formData.isMajorAccredited) }}
          </el-descriptions-item>
          <el-descriptions-item label="是否入选卓越计划">
            {{ formatBoolean(formData.isSelectedForExcellenceProgram) }}
          </el-descriptions-item>
          <el-descriptions-item label="是否双学位">
            {{ formatBoolean(formData.isDoubleDegree) }}
          </el-descriptions-item>
          <el-descriptions-item label="是否新工科">
            {{ formatBoolean(formData.isNewEngineeringDiscipline) }}
          </el-descriptions-item>
          <el-descriptions-item label="专业满意度">
            {{ formData.majorSatisfaction }}
          </el-descriptions-item>
          <el-descriptions-item label="所在学校">
            {{ formData.universityName }}
          </el-descriptions-item>
          <el-descriptions-item label="">
            <el-tag type="success" class="custom-tag">招生信息</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="招生专业">
            {{ formData.admissionMajor }}
          </el-descriptions-item>
          <el-descriptions-item label="计划招生数">
            <span v-if="formData.plannedEnrollment < 0">暂无</span>
            <span v-else>{{ formData.plannedEnrollment }}</span>
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
  createProfessionalConstructionDatabase,
  deleteProfessionalConstructionDatabase,
  deleteProfessionalConstructionDatabaseByIds,
  updateProfessionalConstructionDatabase,
  findProfessionalConstructionDatabase,
  getProfessionalConstructionDatabaseList
} from '@/api/professionalConstructionDatabase'

// 全量引入格式化工具 请按需保留
import { getDictFunc, formatDate, formatDatemini, formatBoolean, filterDict, ReturnArrImg, onDownloadFile } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'
import exportExcel from '@/components/exportExcel/exportExcel.vue'

defineOptions({
  name: 'ProfessionalConstructionDatabase'
})

const formatCount = (value) => {
  return value < 0 ? '暂无' : value;
};

// 自动化生成的字典（可能为空）以及字段
const formData = ref({
  majorName: '',
  majorCode: '',
  isNewMajor: false,
  isBroadAdmissionCategory: false,
  isFirstClassMajor: false,
  isMajorAccredited: false,
  isSelectedForExcellenceProgram: false,
  isDoubleDegree: false,
  isNewEngineeringDiscipline: false,
  universityName: '',
  admissionMajor: '',
  plannedEnrollment: -1,
  acquisitionTime: new Date(),
})


// 验证规则
const rule = reactive({
  majorName: [{
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
  majorCode: [{
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
  degreeType: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
  ],
  isNewMajor: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
  ],
  isBroadAdmissionCategory: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
  ],
  isFirstClassMajor: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
  ],
  isMajorAccredited: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
  ],
  isSelectedForExcellenceProgram: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
  ],
  isDoubleDegree: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
  ],
  isNewEngineeringDiscipline: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
  ],
  majorSatisfaction: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
  ],
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
  admissionMajor: [{
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
  plannedEnrollment: [{
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
    majorCode: 'major_code',
    plannedEnrollment: 'planned_enrollment',
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
    if (searchInfo.value.isNewMajor === "") {
      searchInfo.value.isNewMajor = null
    }
    if (searchInfo.value.isBroadAdmissionCategory === "") {
      searchInfo.value.isBroadAdmissionCategory = null
    }
    if (searchInfo.value.isFirstClassMajor === "") {
      searchInfo.value.isFirstClassMajor = null
    }
    if (searchInfo.value.isMajorAccredited === "") {
      searchInfo.value.isMajorAccredited = null
    }
    if (searchInfo.value.isSelectedForExcellenceProgram === "") {
      searchInfo.value.isSelectedForExcellenceProgram = null
    }
    if (searchInfo.value.isDoubleDegree === "") {
      searchInfo.value.isDoubleDegree = null
    }
    if (searchInfo.value.isNewEngineeringDiscipline === "") {
      searchInfo.value.isNewEngineeringDiscipline = null
    }
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
  const table = await getProfessionalConstructionDatabaseList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
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
    deleteProfessionalConstructionDatabaseFunc(row)
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
    const res = await deleteProfessionalConstructionDatabaseByIds({ IDs })
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
const updateProfessionalConstructionDatabaseFunc = async (row) => {
  const res = await findProfessionalConstructionDatabase({ ID: row.ID })
  type.value = 'update'
  if (res.code === 0) {
    formData.value = res.data.rePCD
    dialogFormVisible.value = true
  }
}


// 删除行
const deleteProfessionalConstructionDatabaseFunc = async (row) => {
  const res = await deleteProfessionalConstructionDatabase({ ID: row.ID })
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
  const res = await findProfessionalConstructionDatabase({ ID: row.ID })
  if (res.code === 0) {
    formData.value = res.data.rePCD
    openDetailShow()
  }
}


// 关闭详情弹窗
const closeDetailShow = () => {
  detailShow.value = false
  formData.value = {
    majorName: '',
    majorCode: '',
    isNewMajor: false,
    isBroadAdmissionCategory: false,
    isFirstClassMajor: false,
    isMajorAccredited: false,
    isSelectedForExcellenceProgram: false,
    isDoubleDegree: false,
    isNewEngineeringDiscipline: false,
    universityName: '',
    admissionMajor: '',
    plannedEnrollment: -1,
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
    majorName: '',
    majorCode: '',
    isNewMajor: false,
    isBroadAdmissionCategory: false,
    isFirstClassMajor: false,
    isMajorAccredited: false,
    isSelectedForExcellenceProgram: false,
    isDoubleDegree: false,
    isNewEngineeringDiscipline: false,
    universityName: '',
    admissionMajor: '',
    plannedEnrollment: -1,
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
        res = await createProfessionalConstructionDatabase(formData.value)
        break
      case 'update':
        res = await updateProfessionalConstructionDatabase(formData.value)
        break
      default:
        res = await createProfessionalConstructionDatabase(formData.value)
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
