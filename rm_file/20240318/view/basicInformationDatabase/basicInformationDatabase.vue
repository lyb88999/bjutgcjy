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
      </div>
      <el-table ref="multipleTable" style="width: 100%" tooltip-effect="dark" :data="tableData" row-key="ID"
        @selection-change="handleSelectionChange" @sort-change="sortChange">
        <el-table-column type="selection" width="55" />

        <el-table-column align="left" label="日期" width="180">
          <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>

        <el-table-column align="left" label="省区市名称" prop="regionName" width="120" />
        <el-table-column align="left" label="区划代码" prop="divisionCode" width="120" />
        <el-table-column align="left" label="地理区域" prop="geographicArea" width="120" />
        <el-table-column align="left" label="地区生产总值" prop="regionGDP" width="120" />
        <el-table-column align="left" label="学校名称" prop="schoolName" width="120" />
        <el-table-column sortable align="left" label="学校代码" prop="schoolCode" width="120" />
        <el-table-column align="left" label="学校类型" prop="schoolType" width="120" />
        <el-table-column align="left" label="学校类别" prop="schoolCategory" width="120" />
        <el-table-column align="left" label="学校层次" prop="educationLevel" width="120" />
        <el-table-column align="left" label="是否双一流" prop="isDoubleFirstClass" width="120">
          <template #default="scope">{{ formatBoolean(scope.row.isDoubleFirstClass) }}</template>
        </el-table-column>
        <el-table-column align="left" label="所在地区" prop="locatedRegion" width="120" />
        <el-table-column align="left" label="在校本科生数" prop="undergraduateNumber" width="120" />
        <el-table-column align="left" label="在校硕士生数" prop="graduateStudentNumber" width="120" />
        <el-table-column align="left" label="在校博士生数" prop="doctoralStudentNumber" width="120" />
        <el-table-column align="left" label="在校留学生数" prop="internationalStudentNumber" width="120" />
        <el-table-column align="left" label="在校教职员工数" prop="facultyStaffNumber" width="120" />
        <el-table-column align="left" label="在校专任教师数" prop="fullTimeTeacherNumber" width="120" />
        <el-table-column align="left" label="本科专业数" prop="undergraduateProgramNumber" width="120" />
        <el-table-column align="left" label="硕士授权点数" prop="mastersAutorizationPointerNumber" width="120" />
        <el-table-column align="left" label="博士授权点数" prop="doctoralAuthorizationPointerNumber" width="120" />
        <el-table-column align="left" label="硕士专业学位类别数" prop="mastersDegreeCategoryNumber" width="120" />
        <el-table-column align="left" label="博士专业学位类别数" prop="doctoralDegreeCategoryNumber" width="120" />
        <el-table-column align="left" label="博士后流动站数" prop="postDoctoralMobileStationNumber" width="120" />
        <el-table-column align="left" label="软科排名" prop="softScienceRanking" width="120" />
        <el-table-column align="left" label="操作" fixed="right" min-width="240">
          <template #default="scope">
            <el-button type="primary" link class="table-button" @click="getDetails(scope.row)">
              <el-icon style="margin-right: 5px">
                <InfoFilled />
              </el-icon>
              查看详情
            </el-button>
            <el-button type="primary" link icon="edit" class="table-button"
              @click="updateBasicInfomationDatabaseFunc(scope.row)">变更</el-button>
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
          <div style="text-align: center; font-weight: bold; font-size: 16px; margin-top: 20px; margin-bottom: 20px;">
            省区市信息
          </div>
          <el-form-item label="省区市名称:" prop="regionName">
            <el-input v-model="formData.regionName" :clearable="true" placeholder="请输入省区市名称" style="width: 200px;" />
          </el-form-item>
          <el-form-item label="区划代码:" prop="divisionCode">
            <el-input v-model="formData.divisionCode" :clearable="true" placeholder="请输入区划代码" style="width: 200px;" />
          </el-form-item>
          <el-form-item label="地理区域:" prop="geographicArea">
            <el-input v-model="formData.geographicArea" :clearable="true" placeholder="请输入地理区域" style="width: 200px;" />
          </el-form-item>
          <el-form-item label="地区生产总值:" prop="regionGDP">
            <el-input-number v-model="formData.regionGDP" style="width:200px" :precision="2" :clearable="true" />
          </el-form-item>
          <div style="text-align: center; font-weight: bold; font-size: 16px; margin-top: 20px; margin-bottom: 20px;">
            学校信息
          </div>
          <el-form-item label="学校名称:" prop="schoolName">
            <el-input v-model="formData.schoolName" :clearable="true" placeholder="请输入学校名称" style="width: 200px;" />
          </el-form-item>
          <el-form-item label="学校代码:" prop="schoolCode">
            <el-input v-model="formData.schoolCode" :clearable="true" placeholder="请输入学校代码" style="width: 200px;" />
          </el-form-item>
          <el-form-item label="学校类型:" prop="schoolType">
            <el-input v-model="formData.schoolType" :clearable="true" placeholder="请输入学校类型" style="width: 200px;" />
          </el-form-item>
          <el-form-item label="学校类别:" prop="schoolCategory">
            <el-input v-model="formData.schoolCategory" :clearable="true" placeholder="请输入学校类别" style="width: 200px;" />
          </el-form-item>
          <el-form-item label="学校层次:" prop="educationLevel">
            <el-input v-model="formData.educationLevel" :clearable="true" placeholder="请输入学校层次" style="width: 200px;" />
          </el-form-item>
          <el-form-item label="是否双一流:" prop="isDoubleFirstClass">
            <el-switch v-model="formData.isDoubleFirstClass" active-color="#13ce66" inactive-color="#ff4949"
              active-text="是" inactive-text="否" clearable style="width: 200px;"></el-switch>
          </el-form-item>
          <el-form-item label="所在地区:" prop="locatedRegion">
            <el-input v-model="formData.locatedRegion" :clearable="true" placeholder="请输入所在地区" style="width: 200px;" />
          </el-form-item>
          <el-form-item label="在校本科生数:" prop="undergraduateNumber">
            <el-input v-model.number="formData.undergraduateNumber" :clearable="true" placeholder="请输入在校本科生数"
              style="width: 200px;" />
          </el-form-item>
          <el-form-item label="在校硕士生数:" prop="graduateStudentNumber">
            <el-input v-model.number="formData.graduateStudentNumber" :clearable="true" placeholder="请输入在校硕士生数"
              style="width: 200px;" />
          </el-form-item>
          <el-form-item label="在校博士生数:" prop="doctoralStudentNumber">
            <el-input v-model.number="formData.doctoralStudentNumber" :clearable="true" placeholder="请输入在校博士生数"
              style="width: 200px;" />
          </el-form-item>
          <el-form-item label="在校留学生数:" prop="internationalStudentNumber">
            <el-input v-model.number="formData.internationalStudentNumber" :clearable="true" placeholder="请输入在校留学生数"
              style="width: 200px;" />
          </el-form-item>
          <el-form-item label="在校教职员工数:" prop="facultyStaffNumber">
            <el-input v-model.number="formData.facultyStaffNumber" :clearable="true" placeholder="请输入在校教职员工数"
              style="width: 200px;" />
          </el-form-item>
          <el-form-item label="在校专任教师数:" prop="fullTimeTeacherNumber">
            <el-input v-model.number="formData.fullTimeTeacherNumber" :clearable="true" placeholder="请输入在校专任教师数"
              style="width: 200px;" />
          </el-form-item>
          <el-form-item label="本科专业数:" prop="undergraduateProgramNumber">
            <el-input v-model.number="formData.undergraduateProgramNumber" :clearable="true" placeholder="请输入本科专业数"
              style="width: 200px;" />
          </el-form-item>
          <el-form-item label="硕士授权点数:" prop="mastersAutorizationPointerNumber">
            <el-input v-model.number="formData.mastersAutorizationPointerNumber" :clearable="true"
              placeholder="请输入硕士授权点数" style="width: 200px;" />
          </el-form-item>
          <el-form-item label="博士授权点数:" prop="doctoralAuthorizationPointerNumber">
            <el-input v-model.number="formData.doctoralAuthorizationPointerNumber" :clearable="true"
              placeholder="请输入博士授权点数" style="width: 200px;" />
          </el-form-item>
          <el-form-item label="硕士专业学位类别数:" prop="mastersDegreeCategoryNumber">
            <el-input v-model.number="formData.mastersDegreeCategoryNumber" :clearable="true" placeholder="请输入硕士专业学位类别数"
              style="width: 200px;" />
          </el-form-item>
          <el-form-item label="博士专业学位类别数:" prop="doctoralDegreeCategoryNumber">
            <el-input v-model.number="formData.doctoralDegreeCategoryNumber" :clearable="true"
              placeholder="请输入博士专业学位类别数" style="width: 200px;" />
          </el-form-item>
          <el-form-item label="博士后流动站数:" prop="postDoctoralMobileStationNumber">
            <el-input v-model.number="formData.postDoctoralMobileStationNumber" :clearable="true"
              placeholder="请输入博士后流动站数" style="width: 200px;" />
          </el-form-item>
          <el-form-item label="软科排名:" prop="softScienceRanking">
            <el-input v-model.number="formData.softScienceRanking" :clearable="true" placeholder="请输入软科排名"
              style="width: 200px;" />
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
            <el-tag type="success" class="custom-tag">省区市信息</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="省区市名称">
            {{ formData.regionName }}
          </el-descriptions-item>
          <el-descriptions-item label="区划代码">
            {{ formData.divisionCode }}
          </el-descriptions-item>
          <el-descriptions-item label="地理区域">
            {{ formData.geographicArea }}
          </el-descriptions-item>
          <el-descriptions-item label="地区生产总值">
            {{ formData.regionGDP }}
          </el-descriptions-item>
          <el-descriptions-item label="">
            <el-tag type="success" class="custom-tag">学校信息</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="学校名称">
            {{ formData.schoolName }}
          </el-descriptions-item>
          <el-descriptions-item label="学校代码">
            {{ formData.schoolCode }}
          </el-descriptions-item>
          <el-descriptions-item label="学校类型">
            {{ formData.schoolType }}
          </el-descriptions-item>
          <el-descriptions-item label="学校类别">
            {{ formData.schoolCategory }}
          </el-descriptions-item>
          <el-descriptions-item label="学校层次">
            {{ formData.educationLevel }}
          </el-descriptions-item>
          <el-descriptions-item label="是否双一流">
            {{ formatBoolean(formData.isDoubleFirstClass) }}
          </el-descriptions-item>
          <el-descriptions-item label="所在地区">
            {{ formData.locatedRegion }}
          </el-descriptions-item>
          <el-descriptions-item label="在校本科生数">
            {{ formData.undergraduateNumber }}
          </el-descriptions-item>
          <el-descriptions-item label="在校硕士生数">
            {{ formData.graduateStudentNumber }}
          </el-descriptions-item>
          <el-descriptions-item label="在校博士生数">
            {{ formData.doctoralStudentNumber }}
          </el-descriptions-item>
          <el-descriptions-item label="在校留学生数">
            {{ formData.internationalStudentNumber }}
          </el-descriptions-item>
          <el-descriptions-item label="在校教职员工数">
            {{ formData.facultyStaffNumber }}
          </el-descriptions-item>
          <el-descriptions-item label="在校专任教师数">
            {{ formData.fullTimeTeacherNumber }}
          </el-descriptions-item>
          <el-descriptions-item label="本科专业数">
            {{ formData.undergraduateProgramNumber }}
          </el-descriptions-item>
          <el-descriptions-item label="硕士授权点数">
            {{ formData.mastersAutorizationPointerNumber }}
          </el-descriptions-item>
          <el-descriptions-item label="博士授权点数">
            {{ formData.doctoralAuthorizationPointerNumber }}
          </el-descriptions-item>
          <el-descriptions-item label="硕士专业学位类别数">
            {{ formData.mastersDegreeCategoryNumber }}
          </el-descriptions-item>
          <el-descriptions-item label="博士专业学位类别数">
            {{ formData.doctoralDegreeCategoryNumber }}
          </el-descriptions-item>
          <el-descriptions-item label="博士后流动站数">
            {{ formData.postDoctoralMobileStationNumber }}
          </el-descriptions-item>
          <el-descriptions-item label="软科排名">
            {{ formData.softScienceRanking }}
          </el-descriptions-item>
        </el-descriptions>
      </el-scrollbar>
    </el-dialog>
  </div>
</template>

<script setup>
import {
  createBasicInfomationDatabase,
  deleteBasicInfomationDatabase,
  deleteBasicInfomationDatabaseByIds,
  updateBasicInfomationDatabase,
  findBasicInfomationDatabase,
  getBasicInfomationDatabaseList
} from '@/api/basicInformationDatabase'

// 全量引入格式化工具 请按需保留
import { getDictFunc, formatDate, formatBoolean, filterDict, ReturnArrImg, onDownloadFile } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'

defineOptions({
  name: 'BasicInfomationDatabase'
})

// 自动化生成的字典（可能为空）以及字段
const formData = ref({
  regionName: '',
  divisionCode: '',
  geographicArea: '',
  regionGDP: 0,
  schoolName: '',
  schoolCode: '',
  schoolType: '',
  schoolCategory: '',
  educationLevel: '',
  isDoubleFirstClass: false,
  locatedRegion: '',
  undergraduateNumber: 0,
  graduateStudentNumber: 0,
  doctoralStudentNumber: 0,
  internationalStudentNumber: 0,
  facultyStaffNumber: 0,
  fullTimeTeacherNumber: 0,
  undergraduateProgramNumber: 0,
  mastersAutorizationPointerNumber: 0,
  doctoralAuthorizationPointerNumber: 0,
  mastersDegreeCategoryNumber: 0,
  doctoralDegreeCategoryNumber: 0,
  postDoctoralMobileStationNumber: 0,
  softScienceRanking: 0,
})


// 验证规则
const rule = reactive({
  regionName: [{
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
  divisionCode: [{
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
  geographicArea: [{
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
  regionGDP: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
  ],
  schoolName: [{
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
  schoolCode: [{
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
  schoolType: [{
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
  schoolCategory: [{
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
  educationLevel: [{
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
  isDoubleFirstClass: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
  ],
  locatedRegion: [{
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
  undergraduateNumber: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
  ],
  graduateStudentNumber: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
  ],
  doctoralStudentNumber: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
  ],
  internationalStudentNumber: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
  ],
  facultyStaffNumber: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
  ],
  fullTimeTeacherNumber: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
  ],
  undergraduateProgramNumber: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
  ],
  mastersAutorizationPointerNumber: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
  ],
  doctoralAuthorizationPointerNumber: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
  ],
  mastersDegreeCategoryNumber: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
  ],
  doctoralDegreeCategoryNumber: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
  ],
  postDoctoralMobileStationNumber: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
  ],
  softScienceRanking: [{
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
    schoolCode: 'school_code',
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
    if (searchInfo.value.isDoubleFirstClass === "") {
      searchInfo.value.isDoubleFirstClass = null
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
  const table = await getBasicInfomationDatabaseList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
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
    deleteBasicInfomationDatabaseFunc(row)
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
    const res = await deleteBasicInfomationDatabaseByIds({ IDs })
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
const updateBasicInfomationDatabaseFunc = async (row) => {
  const res = await findBasicInfomationDatabase({ ID: row.ID })
  type.value = 'update'
  if (res.code === 0) {
    formData.value = res.data.reBID
    dialogFormVisible.value = true
  }
}


// 删除行
const deleteBasicInfomationDatabaseFunc = async (row) => {
  const res = await deleteBasicInfomationDatabase({ ID: row.ID })
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
  const res = await findBasicInfomationDatabase({ ID: row.ID })
  if (res.code === 0) {
    formData.value = res.data.reBID
    openDetailShow()
  }
}


// 关闭详情弹窗
const closeDetailShow = () => {
  detailShow.value = false
  formData.value = {
    regionName: '',
    divisionCode: '',
    geographicArea: '',
    regionGDP: 0,
    schoolName: '',
    schoolCode: '',
    schoolType: '',
    schoolCategory: '',
    educationLevel: '',
    isDoubleFirstClass: false,
    locatedRegion: '',
    undergraduateNumber: 0,
    graduateStudentNumber: 0,
    doctoralStudentNumber: 0,
    internationalStudentNumber: 0,
    facultyStaffNumber: 0,
    fullTimeTeacherNumber: 0,
    undergraduateProgramNumber: 0,
    mastersAutorizationPointerNumber: 0,
    doctoralAuthorizationPointerNumber: 0,
    mastersDegreeCategoryNumber: 0,
    doctoralDegreeCategoryNumber: 0,
    postDoctoralMobileStationNumber: 0,
    softScienceRanking: 0,
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
    regionName: '',
    divisionCode: '',
    geographicArea: '',
    regionGDP: 0,
    schoolName: '',
    schoolCode: '',
    schoolType: '',
    schoolCategory: '',
    educationLevel: '',
    isDoubleFirstClass: false,
    locatedRegion: '',
    undergraduateNumber: 0,
    graduateStudentNumber: 0,
    doctoralStudentNumber: 0,
    internationalStudentNumber: 0,
    facultyStaffNumber: 0,
    fullTimeTeacherNumber: 0,
    undergraduateProgramNumber: 0,
    mastersAutorizationPointerNumber: 0,
    doctoralAuthorizationPointerNumber: 0,
    mastersDegreeCategoryNumber: 0,
    doctoralDegreeCategoryNumber: 0,
    postDoctoralMobileStationNumber: 0,
    softScienceRanking: 0,
  }
}
// 弹窗确定
const enterDialog = async () => {
  elFormRef.value?.validate(async (valid) => {
    if (!valid) return
    let res
    switch (type.value) {
      case 'create':
        res = await createBasicInfomationDatabase(formData.value)
        break
      case 'update':
        res = await updateBasicInfomationDatabase(formData.value)
        break
      default:
        res = await createBasicInfomationDatabase(formData.value)
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

<style>
.custom-descriptions .description-title .el-descriptions-item__label {
  font-size: 25px;
  /* 或者其他你想要的大小 */
  font-weight: bold;
  /* 加粗 */
  text-align: center;
  /* 居中 */
  width: 100%;
  /* 宽度为100%，确保居中生效 */
}
</style>
