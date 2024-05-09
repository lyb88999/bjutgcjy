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

        <el-form-item label="所在学校" prop="universityName">
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
          <exportExcel template-id="TTD" :condition="searchInfo"/>
      </div>
      <el-table ref="multipleTable" style="width: 100%" tooltip-effect="dark" :data="tableData" row-key="ID"
        @selection-change="handleSelectionChange" @sort-change="sortChange">
        <el-table-column type="index" fixed width="50" />
        <el-table-column type="selection" fixed width="55" />

        <el-table-column align="left" label="日期" width="180">
          <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>

        <el-table-column align="left" label="学校名称" fixed prop="universityName" width="120" />
        <el-table-column label="在校学生" align="center" width="360">
          <el-table-column sortable align="left" label="本科生数" prop="undergraduateCount" width="120" />
          <el-table-column sortable align="left" label="硕士生数" prop="masterStudentCount" width="120" />
          <el-table-column sortable align="left" label="博士生数" prop="doctoralStudentCount" width="120" />
        </el-table-column>
        <el-table-column label="课堂建设" align="center" width="240">
          <el-table-column sortable align="left" label="国家级一流课程数" prop="nationalFirstClassCourseCount" width="120" />
          <el-table-column sortable align="left" label="省部级一流课程数" prop="provincialFirstClassCourseCount" width="120" />
        </el-table-column>
        <el-table-column label="教学中心" align="center" width="240">
          <el-table-column sortable align="left" label="国家级教学示范中心" prop="ationalTeachingDemonstrationCenterCount"
            width="120" />
          <el-table-column sortable align="left" label="省部级教学示范中心" prop="provincialTeachingDemonstrationCenterCount"
            width="120" />
        </el-table-column>
        <el-table-column label="教材建设" align="center" width="360">
          <el-table-column sortable align="left" label="国家级规划教材数" prop="nationalPlanningTextbookCount" width="120" />
          <el-table-column sortable align="left" label="全国优秀教材数" prop="nationalExcellentTextbookCount" width="120" />
          <el-table-column sortable align="left" label="省级优秀教材数" prop="provincialExcellentTextbookCount" width="120" />
        </el-table-column>
        <el-table-column label="学生获奖" align="center" width="360">
          <el-table-column sortable align="left" label="学生获国际奖项数" prop="internationalAwardsCount" width="120" />
          <el-table-column sortable align="left" label="学生获国家级奖项数" prop="nationalAwardsCount" width="120" />
          <el-table-column sortable align="left" label="学生获省部级奖项数" prop="provincialAwardsCount" width="120" />
        </el-table-column>
        <el-table-column label="学生毕业去向落实率和升学率" align="center" width="360">
          <el-table-column sortable align="left" label="本科生平均毕业去向落实率" prop="undergraduateEmploymentRate" width="120" />
          <el-table-column sortable align="left" label="硕士生平均毕业去向落实率" prop="masterEmploymentRate" width="120" />
          <el-table-column sortable align="left" label="博士生平均毕业去向落实率" prop="doctoralEmploymentRate" width="120" />
        </el-table-column>
        <el-table-column label="在校学生创业项目" align="center" width="240">
          <el-table-column sortable align="left" label="项目数" prop="projectCount" width="120" />
          <el-table-column sortable align="left" label="参与学生数" prop="participatingStudentCount" width="120" />
        </el-table-column>
        <el-table-column label="本科生转专业" align="center" width="240">
          <el-table-column sortable align="left" label="转入工科学生数" prop="transferredIntoEngineeringCount" width="120" />
          <el-table-column sortable align="left" label="转出工科学生数" prop="transferredOutOfEngineeringCount" width="120" />
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
              @click="updateTalentTrainingDatabaseFunc(scope.row)">变更</el-button>
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
          <el-form-item label="所在学校:" prop="universityName">
            <el-input v-model="formData.universityName" :clearable="true" placeholder="请输入所在学校" style="width: 200px;" />
          </el-form-item>
          <div style="text-align: center; font-weight: bold; font-size: 16px; margin-top: 20px; margin-bottom: 20px;">
            在校学生
          </div>
          <el-form-item label="本科生数:" prop="undergraduateCount">
            <el-input v-model.number="formData.undergraduateCount" :clearable="true" placeholder="请输入本科生数"
              style="width: 200px;" />
          </el-form-item>
          <el-form-item label="硕士生数:" prop="masterStudentCount">
            <el-input v-model.number="formData.masterStudentCount" :clearable="true" placeholder="请输入硕士生数"
              style="width: 200px;" />
          </el-form-item>
          <el-form-item label="博士生数:" prop="doctoralStudentCount">
            <el-input v-model.number="formData.doctoralStudentCount" :clearable="true" placeholder="请输入博士生数"
              style="width: 200px;" />
          </el-form-item>
          <div style="text-align: center; font-weight: bold; font-size: 16px; margin-top: 20px; margin-bottom: 20px;">
            课程建设
          </div>
          <el-form-item label="国家级一流课程数:" prop="nationalFirstClassCourseCount">
            <el-input v-model.number="formData.nationalFirstClassCourseCount" :clearable="true"
              placeholder="请输入国家级一流课程数" style="width: 200px;" />
          </el-form-item>
          <el-form-item label="省部级一流课程数:" prop="provincialFirstClassCourseCount">
            <el-input v-model.number="formData.provincialFirstClassCourseCount" :clearable="true"
              placeholder="请输入省部级一流课程数" style="width: 200px;" />
          </el-form-item>
          <div style="text-align: center; font-weight: bold; font-size: 16px; margin-top: 20px; margin-bottom: 20px;">
            教学中心
          </div>
          <el-form-item label="国家级教学示范中心:" prop="ationalTeachingDemonstrationCenterCount">
            <el-input v-model.number="formData.ationalTeachingDemonstrationCenterCount" :clearable="true"
              placeholder="请输入国家级教学示范中心数" style="width: 200px;" />
          </el-form-item>
          <el-form-item label="省部级教学示范中心:" prop="provincialTeachingDemonstrationCenterCount">
            <el-input v-model.number="formData.provincialTeachingDemonstrationCenterCount" :clearable="true"
              placeholder="请输入省部级教学示范中心数" style="width: 200px;" />
          </el-form-item>
          <div style="text-align: center; font-weight: bold; font-size: 16px; margin-top: 20px; margin-bottom: 20px;">
            教材建设
          </div>
          <el-form-item label="国家级规划教材数:" prop="nationalPlanningTextbookCount">
            <el-input v-model.number="formData.nationalPlanningTextbookCount" :clearable="true"
              placeholder="请输入国家级规划教材数" style="width: 200px;" />
          </el-form-item>
          <el-form-item label="全国优秀教材数:" prop="nationalExcellentTextbookCount">
            <el-input v-model.number="formData.nationalExcellentTextbookCount" :clearable="true"
              placeholder="请输入全国优秀教材数" style="width: 200px;" />
          </el-form-item>
          <el-form-item label="省级优秀教材数:" prop="provincialExcellentTextbookCount">
            <el-input v-model.number="formData.provincialExcellentTextbookCount" :clearable="true"
              placeholder="请输入省级优秀教材数" style="width: 200px;" />
          </el-form-item>
          <div style="text-align: center; font-weight: bold; font-size: 16px; margin-top: 20px; margin-bottom: 20px;">
            学生获奖
          </div>
          <el-form-item label="学生获国际奖项数:" prop="internationalAwardsCount">
            <el-input v-model.number="formData.internationalAwardsCount" :clearable="true" placeholder="请输入学生获国际奖项数"
              style="width: 200px;" />
          </el-form-item>
          <el-form-item label="学生获国家级奖项数:" prop="nationalAwardsCount">
            <el-input v-model.number="formData.nationalAwardsCount" :clearable="true" placeholder="请输入学生获国家级奖项数"
              style="width: 200px;" />
          </el-form-item>
          <el-form-item label="学生获省部级奖项数:" prop="provincialAwardsCount">
            <el-input v-model.number="formData.provincialAwardsCount" :clearable="true" placeholder="请输入学生获省部级奖项数"
              style="width: 200px;" />
          </el-form-item>
          <div style="text-align: center; font-weight: bold; font-size: 16px; margin-top: 20px; margin-bottom: 20px;">
            学生毕业去向落实率和升学率
          </div>
          <el-form-item label="本科生平均毕业去向落实率:" prop="undergraduateEmploymentRate">
            <el-input-number v-model="formData.undergraduateEmploymentRate" style="width:200px" :precision="2"
              :clearable="true" />
          </el-form-item>
          <el-form-item label="硕士生平均毕业去向落实率:" prop="masterEmploymentRate">
            <el-input-number v-model="formData.masterEmploymentRate" style="width:200px" :precision="2"
              :clearable="true" />
          </el-form-item>
          <el-form-item label="博士生平均毕业去向落实率:" prop="doctoralEmploymentRate">
            <el-input-number v-model="formData.doctoralEmploymentRate" style="width:200px" :precision="2"
              :clearable="true" />
          </el-form-item>
          <div style="text-align: center; font-weight: bold; font-size: 16px; margin-top: 20px; margin-bottom: 20px;">
            在校学生创业项目
          </div>
          <el-form-item label="项目数:" prop="projectCount">
            <el-input v-model.number="formData.projectCount" :clearable="true" placeholder="请输入项目数"
              style="width: 200px;" />
          </el-form-item>
          <el-form-item label="参与学生数:" prop="participatingStudentCount">
            <el-input v-model.number="formData.participatingStudentCount" :clearable="true" placeholder="请输入参与学生数"
              style="width: 200px;" />
          </el-form-item>
          <div style="text-align: center; font-weight: bold; font-size: 16px; margin-top: 20px; margin-bottom: 20px;">
            本科生转专业
          </div>
          <el-form-item label="转入工科学生数:" prop="transferredIntoEngineeringCount">
            <el-input v-model.number="formData.transferredIntoEngineeringCount" :clearable="true"
              placeholder="请输入转入工科学生数" style="width: 200px;" />
          </el-form-item>
          <el-form-item label="转出工科学生数:" prop="transferredOutOfEngineeringCount">
            <el-input v-model.number="formData.transferredOutOfEngineeringCount" :clearable="true"
              placeholder="请输入转出工科学生数" style="width: 200px;" />
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
          <el-descriptions-item label="所在学校">
            {{ formData.universityName }}
          </el-descriptions-item>
          <el-descriptions-item label="">
            <el-tag type="success" class="custom-tag">在校学生</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="本科生数">
            {{ formData.undergraduateCount }}
          </el-descriptions-item>
          <el-descriptions-item label="硕士生数">
            {{ formData.masterStudentCount }}
          </el-descriptions-item>
          <el-descriptions-item label="博士生数">
            {{ formData.doctoralStudentCount }}
          </el-descriptions-item>
          <el-descriptions-item label="">
            <el-tag type="success" class="custom-tag">课程建设</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="国家级一流课程数">
            {{ formData.nationalFirstClassCourseCount }}
          </el-descriptions-item>
          <el-descriptions-item label="省部级一流课程数">
            {{ formData.provincialFirstClassCourseCount }}
          </el-descriptions-item>
          <el-descriptions-item label="">
            <el-tag type="success" class="custom-tag">教学中心</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="国家级教学示范中心">
            {{ formData.ationalTeachingDemonstrationCenterCount }}
          </el-descriptions-item>
          <el-descriptions-item label="省部级教学示范中心">
            {{ formData.provincialTeachingDemonstrationCenterCount }}
          </el-descriptions-item>
          <el-descriptions-item label="">
            <el-tag type="success" class="custom-tag">教材建设</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="国家级规划教材数">
            {{ formData.nationalPlanningTextbookCount }}
          </el-descriptions-item>
          <el-descriptions-item label="全国优秀教材数">
            {{ formData.nationalExcellentTextbookCount }}
          </el-descriptions-item>
          <el-descriptions-item label="省级优秀教材数">
            {{ formData.provincialExcellentTextbookCount }}
          </el-descriptions-item>
          <el-descriptions-item label="">
            <el-tag type="success" class="custom-tag">学生获奖</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="学生获国际奖项数">
            {{ formData.internationalAwardsCount }}
          </el-descriptions-item>
          <el-descriptions-item label="学生获国家级奖项数">
            {{ formData.nationalAwardsCount }}
          </el-descriptions-item>
          <el-descriptions-item label="学生获省部级奖项数">
            {{ formData.provincialAwardsCount }}
          </el-descriptions-item>
          <el-descriptions-item label="">
            <el-tag type="success" class="custom-tag">学生毕业去向落实率和升学率</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="本科生平均毕业去向落实率">
            {{ formData.undergraduateEmploymentRate }}
          </el-descriptions-item>
          <el-descriptions-item label="硕士生平均毕业去向落实率">
            {{ formData.masterEmploymentRate }}
          </el-descriptions-item>
          <el-descriptions-item label="博士生平均毕业去向落实率">
            {{ formData.doctoralEmploymentRate }}
          </el-descriptions-item>
          <el-descriptions-item label="">
            <el-tag type="success" class="custom-tag">在校学生创业项目</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="项目数">
            {{ formData.projectCount }}
          </el-descriptions-item>
          <el-descriptions-item label="参与学生数">
            {{ formData.participatingStudentCount }}
          </el-descriptions-item>
          <el-descriptions-item label="">
            <el-tag type="success" class="custom-tag">本科生转专业</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="转入工科学生数">
            {{ formData.transferredIntoEngineeringCount }}
          </el-descriptions-item>
          <el-descriptions-item label="转出工科学生数">
            {{ formData.transferredOutOfEngineeringCount }}
          </el-descriptions-item>
        </el-descriptions>
      </el-scrollbar>
    </el-dialog>

  </div>
</template>

<script setup>
import {
  createTalentTrainingDatabase,
  deleteTalentTrainingDatabase,
  deleteTalentTrainingDatabaseByIds,
  updateTalentTrainingDatabase,
  findTalentTrainingDatabase,
  getTalentTrainingDatabaseList
} from '@/api/talentTrainingDatabase'

// 全量引入格式化工具 请按需保留
import { getDictFunc, formatDate, formatBoolean, filterDict, ReturnArrImg, onDownloadFile } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'
import exportExcel from '@/components/exportExcel/exportExcel.vue'

defineOptions({
  name: 'TalentTrainingDatabase'
})

// 自动化生成的字典（可能为空）以及字段
const formData = ref({
  universityName: '',
  undergraduateCount: 0,
  masterStudentCount: 0,
  doctoralStudentCount: 0,
  nationalFirstClassCourseCount: 0,
  provincialFirstClassCourseCount: 0,
  ationalTeachingDemonstrationCenterCount: 0,
  provincialTeachingDemonstrationCenterCount: 0,
  nationalPlanningTextbookCount: 0,
  nationalExcellentTextbookCount: 0,
  provincialExcellentTextbookCount: 0,
  internationalAwardsCount: 0,
  nationalAwardsCount: 0,
  provincialAwardsCount: 0,
  undergraduateEmploymentRate: 0,
  masterEmploymentRate: 0,
  doctoralEmploymentRate: 0,
  projectCount: 0,
  participatingStudentCount: 0,
  transferredIntoEngineeringCount: 0,
  transferredOutOfEngineeringCount: 0,
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
  undergraduateCount: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
  ],
  masterStudentCount: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
  ],
  doctoralStudentCount: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
  ],
  nationalFirstClassCourseCount: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
  ],
  provincialFirstClassCourseCount: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
  ],
  ationalTeachingDemonstrationCenterCount: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
  ],
  provincialTeachingDemonstrationCenterCount: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
  ],
  nationalPlanningTextbookCount: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
  ],
  nationalExcellentTextbookCount: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
  ],
  provincialExcellentTextbookCount: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
  ],
  internationalAwardsCount: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
  ],
  nationalAwardsCount: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
  ],
  provincialAwardsCount: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
  ],
  undergraduateEmploymentRate: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
  ],
  masterEmploymentRate: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
  ],
  doctoralEmploymentRate: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
  ],
  projectCount: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
  ],
  participatingStudentCount: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
  ],
  transferredIntoEngineeringCount: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
  ],
  transferredOutOfEngineeringCount: [{
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
    undergraduateCount: 'undergraduate_count',
    masterStudentCount: 'master_student_count',
    doctoralStudentCount: 'doctoral_student_count',
    nationalFirstClassCourseCount: 'national_first_class_course_count',
    provincialFirstClassCourseCount: 'provincial_first_class_course_count',
    ationalTeachingDemonstrationCenterCount: 'ational_teaching_demonstration_center_count',
    provincialTeachingDemonstrationCenterCount: 'provincial_teaching_demonstration_center_count',
    nationalPlanningTextbookCount: 'national_planning_textbook_count',
    nationalExcellentTextbookCount: 'national_excellent_textbook_count',
    provincialExcellentTextbookCount: 'provincial_excellent_textbook_count',
    internationalAwardsCount: 'international_awards_count',
    nationalAwardsCount: 'national_awards_count',
    provincialAwardsCount: 'provincial_awards_count',
    undergraduateEmploymentRate: 'undergraduate_employment_rate',
    masterEmploymentRate: 'master_employment_rate',
    doctoralEmploymentRate: 'doctoral_employment_rate',
    projectCount: 'project_count',
    participatingStudentCount: 'participating_student_count',
    transferredIntoEngineeringCount: 'transferred_into_engineering_count',
    transferredOutOfEngineeringCount: 'transferred_out_of_engineering_count',
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
  const table = await getTalentTrainingDatabaseList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
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
    deleteTalentTrainingDatabaseFunc(row)
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
    const res = await deleteTalentTrainingDatabaseByIds({ IDs })
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
const updateTalentTrainingDatabaseFunc = async (row) => {
  const res = await findTalentTrainingDatabase({ ID: row.ID })
  type.value = 'update'
  if (res.code === 0) {
    formData.value = res.data.reTTD
    dialogFormVisible.value = true
  }
}


// 删除行
const deleteTalentTrainingDatabaseFunc = async (row) => {
  const res = await deleteTalentTrainingDatabase({ ID: row.ID })
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
  const res = await findTalentTrainingDatabase({ ID: row.ID })
  if (res.code === 0) {
    formData.value = res.data.reTTD
    openDetailShow()
  }
}


// 关闭详情弹窗
const closeDetailShow = () => {
  detailShow.value = false
  formData.value = {
    universityName: '',
    undergraduateCount: 0,
    masterStudentCount: 0,
    doctoralStudentCount: 0,
    nationalFirstClassCourseCount: 0,
    provincialFirstClassCourseCount: 0,
    ationalTeachingDemonstrationCenterCount: 0,
    provincialTeachingDemonstrationCenterCount: 0,
    nationalPlanningTextbookCount: 0,
    nationalExcellentTextbookCount: 0,
    provincialExcellentTextbookCount: 0,
    internationalAwardsCount: 0,
    nationalAwardsCount: 0,
    provincialAwardsCount: 0,
    undergraduateEmploymentRate: 0,
    masterEmploymentRate: 0,
    doctoralEmploymentRate: 0,
    projectCount: 0,
    participatingStudentCount: 0,
    transferredIntoEngineeringCount: 0,
    transferredOutOfEngineeringCount: 0,
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
    undergraduateCount: 0,
    masterStudentCount: 0,
    doctoralStudentCount: 0,
    nationalFirstClassCourseCount: 0,
    provincialFirstClassCourseCount: 0,
    ationalTeachingDemonstrationCenterCount: 0,
    provincialTeachingDemonstrationCenterCount: 0,
    nationalPlanningTextbookCount: 0,
    nationalExcellentTextbookCount: 0,
    provincialExcellentTextbookCount: 0,
    internationalAwardsCount: 0,
    nationalAwardsCount: 0,
    provincialAwardsCount: 0,
    undergraduateEmploymentRate: 0,
    masterEmploymentRate: 0,
    doctoralEmploymentRate: 0,
    projectCount: 0,
    participatingStudentCount: 0,
    transferredIntoEngineeringCount: 0,
    transferredOutOfEngineeringCount: 0,
  }
}
// 弹窗确定
const enterDialog = async () => {
  elFormRef.value?.validate(async (valid) => {
    if (!valid) return
    let res
    switch (type.value) {
      case 'create':
        res = await createTalentTrainingDatabase(formData.value)
        break
      case 'update':
        res = await updateTalentTrainingDatabase(formData.value)
        break
      default:
        res = await createTalentTrainingDatabase(formData.value)
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
