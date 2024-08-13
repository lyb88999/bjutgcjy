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

        <el-form-item label="博士后流动站名称" prop="postdocStationName">
          <el-input v-model="searchInfo.postdocStationName" placeholder="搜索条件" />

        </el-form-item>
        <el-form-item label="博士后学科门类" prop="postdocSubjectCategory">
          <el-input v-model="searchInfo.postdocSubjectCategory" placeholder="搜索条件" />

        </el-form-item>
        <el-form-item label="硕士学位点学科代码" prop="masterSubjectCode">
          <el-input v-model="searchInfo.masterSubjectCode" placeholder="搜索条件" />

        </el-form-item>
        <el-form-item label="硕士学位点学科门类" prop="masterSubjectCategory">
          <el-input v-model="searchInfo.masterSubjectCategory" placeholder="搜索条件" />

        </el-form-item>
        <el-form-item label="硕士学位点类型" prop="masterType">
          <el-input v-model="searchInfo.masterType" placeholder="搜索条件" />

        </el-form-item>
        <el-form-item label="硕士学位点是否一流学科" prop="masterIsFirstClass">
          <el-select v-model="searchInfo.masterIsFirstClass" clearable placeholder="请选择">
            <el-option key="true" label="是" value="true">
            </el-option>
            <el-option key="false" label="否" value="false">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="硕士学位点所在学校" prop="masterUniversityName">
          <el-input v-model="searchInfo.masterUniversityName" placeholder="搜索条件" />

        </el-form-item>
        <el-form-item label="博士学位点学科代码" prop="phdSubjectCode">
          <el-input v-model="searchInfo.phdSubjectCode" placeholder="搜索条件" />

        </el-form-item>
        <el-form-item label="博士学位点学科门类" prop="phdSubjectCategory">
          <el-input v-model="searchInfo.phdSubjectCategory" placeholder="搜索条件" />

        </el-form-item>
        <el-form-item label="博士学位点类型" prop="phdType">
          <el-input v-model="searchInfo.phdType" placeholder="搜索条件" />

        </el-form-item>
        <el-form-item label="博士学位点是否一流学科" prop="phdIsFirstClass">
          <el-select v-model="searchInfo.phdIsFirstClass" clearable placeholder="请选择">
            <el-option key="true" label="是" value="true">
            </el-option>
            <el-option key="false" label="否" value="false">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="博士学位点所在学校" prop="phdUniversityName">
          <el-input v-model="searchInfo.phdUniversityName" placeholder="搜索条件" />

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
        <exportExcel template-id="SCD" :condition="searchInfo" />
      </div>
      <el-table ref="multipleTable" style="width: 100%" tooltip-effect="dark" :data="tableData" row-key="ID"
        @selection-change="handleSelectionChange" @sort-change="sortChange">
        <el-table-column type="index" fixed width="50" />
        <el-table-column type="selection" fixed width="55" />


        <el-table-column align="left" label="日期" width="180">
          <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>
        <el-table-column align="left" label="学校名称" prop="masterUniversityName" fixed width="120" />
        <el-table-column label="博士后流动站" align="center" width="360">
          <el-table-column align="left" label="流动站名称" prop="postdocStationName" width="120" />
          <el-table-column align="left" label="学科门类" prop="postdocSubjectCategory" width="120" />
          <el-table-column align="left" label="流动站数" prop="postdocStationCount" width="120">
            <template #default="{ row }">
              {{ formatCount(row.postdocStationCount) }}
            </template>
          </el-table-column>
        </el-table-column>
        <el-table-column label="硕士学位点" align="center" width="840">
          <el-table-column align="left" label="学科代码" prop="masterSubjectCode" width="120" />
          <el-table-column align="left" label="学科门类" prop="masterSubjectCategory" width="120" />
          <el-table-column sortable align="left" label="软科排名" prop="masterSoftScienceRanking" width="120">
            <template #default="{ row }">
              {{ formatCount(row.masterSoftScienceRanking) }}
            </template>
          </el-table-column>
          <el-table-column sortable align="left" label="招生人数" prop="masterAdmissionNumber" width="120">
            <template #default="{ row }">
              {{ formatCount(row.masterAdmissionNumber) }}
            </template>
          </el-table-column>
          <el-table-column align="left" label="类型" prop="masterType" width="120">
            <template #default="{ row }">
              {{ formatType(row.masterType) }}
            </template>
          </el-table-column>
          <el-table-column align="left" label="是否一流学科" prop="masterIsFirstClass" width="120">
            <template #default="scope">{{ formatBoolean(scope.row.masterIsFirstClass) }}</template>
          </el-table-column>
          <el-table-column align="left" label="所在学校" prop="masterUniversityName" width="120" />
        </el-table-column>
        <el-table-column label="博士学位点" align="center" width="840">
          <el-table-column align="left" label="学科代码" prop="phdSubjectCode" width="120" />
          <el-table-column align="left" label="学科门类" prop="phdSubjectCategory" width="120" />
          <el-table-column sortable align="left" label="软科排名" prop="phdSoftScienceRanking" width="120">
            <template #default="{ row }">
              {{ formatCount(row.phdSoftScienceRanking) }}
            </template>
          </el-table-column>
          <el-table-column sortable align="left" label="招生人数" prop="phdAdmissionNumber" width="120">
            <template #default="{ row }">
              {{ formatCount(row.phdAdmissionNumber) }}
            </template>
          </el-table-column>
          <el-table-column align="left" label="类型" prop="phdType" width="120">
            <template #default="{ row }">
              {{ formatType(row.phdType) }}
            </template>
          </el-table-column>
          <el-table-column align="left" label="是否一流学科" prop="phdIsFirstClass" width="120">
            <template #default="scope">{{ formatBoolean(scope.row.phdIsFirstClass) }}</template>
          </el-table-column>
          <el-table-column align="left" label="所在学校" prop="phdUniversityName" width="120" />
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
              @click="updateSubjectConstructionDatabaseFunc(scope.row)">变更</el-button>
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
        <div style="text-align: center; font-weight: bold; font-size: 16px; margin-top: 20px; margin-bottom: 20px;">
          博士后流动站
        </div>
        <el-form :model="formData" label-position="right" ref="elFormRef" :rules="rule" label-width="40%">
          <el-form-item label="流动站名称:" prop="postdocStationName">
            <el-input v-model="formData.postdocStationName" :clearable="true" placeholder="请输入流动站名称"
              style="width: 200px;" />
          </el-form-item>
          <el-form-item label="学科门类:" prop="postdocSubjectCategory">
            <el-input v-model="formData.postdocSubjectCategory" :clearable="true" placeholder="请输入学科门类"
              style="width: 200px;" />
          </el-form-item>
          <el-form-item label="流动站数:" prop="postdocStationCount">
            <el-input v-model.number="formData.postdocStationCount" :clearable="true" placeholder="请输入流动站数"
              style="width: 200px;" />
          </el-form-item>
          <div style="text-align: center; font-weight: bold; font-size: 16px; margin-top: 20px; margin-bottom: 20px;">
            硕士学位点
          </div>
          <el-form-item label="学科代码:" prop="masterSubjectCode">
            <el-input v-model="formData.masterSubjectCode" :clearable="true" placeholder="请输入学科代码"
              style="width: 200px;" />
          </el-form-item>
          <el-form-item label="学科门类:" prop="masterSubjectCategory">
            <el-input v-model="formData.masterSubjectCategory" :clearable="true" placeholder="请输入学科门类"
              style="width: 200px;" />
          </el-form-item>
          <el-form-item label="软科排名:" prop="masterSoftScienceRanking">
            <el-input v-model.number="formData.masterSoftScienceRanking" :clearable="true" placeholder="请输入软科排名"
              style="width: 200px;" />
          </el-form-item>
          <el-form-item label="招生人数:" prop="masterAdmissionNumber">
            <el-input v-model.number="formData.masterAdmissionNumber" :clearable="true" placeholder="请输入招生人数"
              style="width: 200px;" />
          </el-form-item>
          <el-form-item label="类型:" prop="masterType">
            <el-select v-model="formData.masterType" placeholder="请选择类型" style="width:200px" :clearable="true">
              <el-option v-for="item in ['学术学位', '专业学位']" :key="item" :label="item" :value="item" />
            </el-select>
          </el-form-item>
          <el-form-item label="是否一流学科:" prop="masterIsFirstClass">
            <el-switch v-model="formData.masterIsFirstClass" active-color="#13ce66" inactive-color="#ff4949"
              active-text="是" inactive-text="否" clearable style="width: 200px;"></el-switch>
          </el-form-item>
          <el-form-item label="所在学校:" prop="masterUniversityName">
            <el-input v-model="formData.masterUniversityName" :clearable="true" placeholder="请输入所在学校"
              style="width: 200px;" />
          </el-form-item>
          <div style="text-align: center; font-weight: bold; font-size: 16px; margin-top: 20px; margin-bottom: 20px;">
            博士学位点
          </div>
          <el-form-item label="学科代码:" prop="phdSubjectCode">
            <el-input v-model="formData.phdSubjectCode" :clearable="true" placeholder="请输入学科代码" style="width: 200px;" />
          </el-form-item>
          <el-form-item label="学科门类:" prop="phdSubjectCategory">
            <el-input v-model="formData.phdSubjectCategory" :clearable="true" placeholder="请输入学科门类"
              style="width: 200px;" />
          </el-form-item>
          <el-form-item label="软科排名:" prop="phdSoftScienceRanking">
            <el-input v-model.number="formData.phdSoftScienceRanking" :clearable="true" placeholder="请输入软科排名"
              style="width: 200px;" />
          </el-form-item>
          <el-form-item label="招生人数:" prop="phdAdmissionNumber">
            <el-input v-model.number="formData.phdAdmissionNumber" :clearable="true" placeholder="请输入招生人数"
              style="width: 200px;" />
          </el-form-item>
          <el-form-item label="类型:" prop="phdType">
            <el-select v-model="formData.phdType" placeholder="请选择类型" style="width:200px" :clearable="true">
              <el-option v-for="item in ['学术学位', '专业学位']" :key="item" :label="item" :value="item" />
            </el-select>
          </el-form-item>
          <el-form-item label="是否一流学科:" prop="phdIsFirstClass">
            <el-switch v-model="formData.phdIsFirstClass" active-color="#13ce66" inactive-color="#ff4949"
              active-text="是" inactive-text="否" clearable style="width: 200px;"></el-switch>
          </el-form-item>
          <el-form-item label="所在学校:" prop="phdUniversityName">
            <el-input v-model="formData.phdUniversityName" :clearable="true" placeholder="请输入所在学校"
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
            <el-tag type="success" class="custom-tag">博士后流动站</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="流动站名称">
            {{ formData.postdocStationName }}
          </el-descriptions-item>
          <el-descriptions-item label="学科门类">
            {{ formData.postdocSubjectCategory }}
          </el-descriptions-item>
          <el-descriptions-item label="流动站数">
            <span v-if="formData.postdocStationCount < 0">暂无</span>
            <span v-else>{{ formData.postdocStationCount }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="">
            <el-tag type="success" class="custom-tag">硕士学位点</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="学科代码">
            {{ formData.masterSubjectCode }}
          </el-descriptions-item>
          <el-descriptions-item label="学科门类">
            {{ formData.masterSubjectCategory }}
          </el-descriptions-item>
          <el-descriptions-item label="软科排名">
            <span v-if="formData.masterSoftScienceRanking < 0">暂无</span>
            <span v-else>{{ formData.masterSoftScienceRanking }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="招生人数">
            <span v-if="formData.masterAdmissionNumber < 0">暂无</span>
            <span v-else>{{ formData.masterAdmissionNumber }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="类型">
            <span v-if="formData.masterType == ''">暂无</span>
            <span v-else>{{ formData.masterType }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="是否一流学科">
            {{ formatBoolean(formData.masterIsFirstClass) }}
          </el-descriptions-item>
          <el-descriptions-item label="所在学校">
            {{ formData.masterUniversityName }}
          </el-descriptions-item>
          <el-descriptions-item label="">
            <el-tag type="success" class="custom-tag">博士学位点</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="学科代码">
            {{ formData.phdSubjectCode }}
          </el-descriptions-item>
          <el-descriptions-item label="学科门类">
            {{ formData.phdSubjectCategory }}
          </el-descriptions-item>
          <el-descriptions-item label="软科排名">
            <span v-if="formData.phdSoftScienceRanking < 0">暂无</span>
            <span v-else>{{ formData.phdSoftScienceRanking }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="招生人数">
            <span v-if="formData.phdAdmissionNumber < 0">暂无</span>
            <span v-else>{{ formData.phdAdmissionNumber }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="类型">
            <span v-if="formData.phdType == ''">暂无</span>
            <span v-else>{{ formData.phdType }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="是否一流学科">
            {{ formatBoolean(formData.phdIsFirstClass) }}
          </el-descriptions-item>
          <el-descriptions-item label="所在学校">
            {{ formData.phdUniversityName }}
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
  createSubjectConstructionDatabase,
  deleteSubjectConstructionDatabase,
  deleteSubjectConstructionDatabaseByIds,
  updateSubjectConstructionDatabase,
  findSubjectConstructionDatabase,
  getSubjectConstructionDatabaseList
} from '@/api/subjectConstructionDatabase'
import ExportExcel from '@/components/exportExcel/exportExcel.vue';

// 全量引入格式化工具 请按需保留
import { getDictFunc, formatDate,formatDatemini, formatBoolean, filterDict, ReturnArrImg, onDownloadFile } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'
import exportExcel from '@/components/exportExcel/exportExcel.vue';

defineOptions({
  name: 'SubjectConstructionDatabase'
})

const formatCount = (value) => {
  return value < 0 ? '暂无' : value;
};

const formatType = (value) => {
  return value == '' ? '暂无' : value
}

// 自动化生成的字典（可能为空）以及字段
const formData = ref({
  postdocStationName: '',
  postdocSubjectCategory: '',
  postdocStationCount: -1,
  masterSubjectCode: '',
  masterSubjectCategory: '',
  masterSoftScienceRanking: -1,
  masterAdmissionNumber: -1,
  masterIsFirstClass: false,
  masterUniversityName: '',
  phdSubjectCode: '',
  phdSubjectCategory: '',
  phdSoftScienceRanking: -1,
  phdAdmissionNumber: -1,
  phdIsFirstClass: false,
  phdUniversityName: '',
  acquisitionTime: new Date(),
})


// 验证规则
const rule = reactive({
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
    postdocStationCount: 'postdoc_station_count',
    masterSubjectCode: 'master_subject_code',
    masterSoftScienceRanking: 'master_soft_science_ranking',
    masterAdmissionNumber: 'master_admission_number',
    phdSubjectCode: 'phd_subject_code',
    phdSoftScienceRanking: 'phd_soft_science_ranking',
    phdAdmissionNumber: 'phd_admission_number',
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
    if (searchInfo.value.masterIsFirstClass === "") {
      searchInfo.value.masterIsFirstClass = null
    }
    if (searchInfo.value.phdIsFirstClass === "") {
      searchInfo.value.phdIsFirstClass = null
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
  const table = await getSubjectConstructionDatabaseList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
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
    deleteSubjectConstructionDatabaseFunc(row)
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
    const res = await deleteSubjectConstructionDatabaseByIds({ IDs })
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
const updateSubjectConstructionDatabaseFunc = async (row) => {
  const res = await findSubjectConstructionDatabase({ ID: row.ID })
  type.value = 'update'
  if (res.code === 0) {
    formData.value = res.data.reSCD
    dialogFormVisible.value = true
  }
}


// 删除行
const deleteSubjectConstructionDatabaseFunc = async (row) => {
  const res = await deleteSubjectConstructionDatabase({ ID: row.ID })
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
  const res = await findSubjectConstructionDatabase({ ID: row.ID })
  if (res.code === 0) {
    formData.value = res.data.reSCD
    openDetailShow()
  }
}


// 关闭详情弹窗
const closeDetailShow = () => {
  detailShow.value = false
  formData.value = {
    postdocStationName: '',
    postdocSubjectCategory: '',
    postdocStationCount: -1,
    masterSubjectCode: '',
    masterSubjectCategory: '',
    masterSoftScienceRanking: -1,
    masterAdmissionNumber: -1,
    masterIsFirstClass: false,
    masterUniversityName: '',
    phdSubjectCode: '',
    phdSubjectCategory: '',
    phdSoftScienceRanking: -1,
    phdAdmissionNumber: -1,
    phdIsFirstClass: false,
    phdUniversityName: '',
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
    postdocStationName: '',
    postdocSubjectCategory: '',
    postdocStationCount: -1,
    masterSubjectCode: '',
    masterSubjectCategory: '',
    masterSoftScienceRanking: -1,
    masterAdmissionNumber: -1,
    masterIsFirstClass: false,
    masterUniversityName: '',
    phdSubjectCode: '',
    phdSubjectCategory: '',
    phdSoftScienceRanking: -1,
    phdAdmissionNumber: -1,
    phdIsFirstClass: false,
    phdUniversityName: '',
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
        res = await createSubjectConstructionDatabase(formData.value)
        break
      case 'update':
        res = await updateSubjectConstructionDatabase(formData.value)
        break
      default:
        res = await createSubjectConstructionDatabase(formData.value)
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
