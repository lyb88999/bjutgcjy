<template>
  <div>
    <div class="gva-table-box" style="margin-bottom: 16px;">
      <el-button icon="back" @click="router.back()">返回</el-button>
      <el-divider content-position="left">背景信息</el-divider>
      <el-descriptions :column="3" border>
        <el-descriptions-item label="姓名">{{ profile.name }}</el-descriptions-item>
        <el-descriptions-item label="性别">{{ profile.gender }}</el-descriptions-item>
        <el-descriptions-item label="民族">{{ profile.ethnicity }}</el-descriptions-item>
        <el-descriptions-item label="政治面貌">{{ profile.politicalStatus }}</el-descriptions-item>
        <el-descriptions-item label="所在单位">{{ profile.unitName }}</el-descriptions-item>
        <el-descriptions-item label="院系/部门">{{ profile.department }}</el-descriptions-item>
        <el-descriptions-item label="行政职务">{{ profile.adminTitle }}</el-descriptions-item>
        <el-descriptions-item label="专业技术职称">{{ profile.techTitle }}</el-descriptions-item>
        <el-descriptions-item label="办公电话">{{ profile.phone }}</el-descriptions-item>
        <el-descriptions-item label="手机号码">{{ profile.mobile }}</el-descriptions-item>
        <el-descriptions-item label="电子邮箱">{{ profile.email }}</el-descriptions-item>
        <el-descriptions-item label="通讯地址">{{ profile.address }}</el-descriptions-item>
        <el-descriptions-item label="最高学历">{{ profile.highestEducation }}</el-descriptions-item>
        <el-descriptions-item label="最高学位">{{ profile.highestDegree }}</el-descriptions-item>
        <el-descriptions-item label="毕业院校">{{ profile.graduateSchool }}</el-descriptions-item>
        <el-descriptions-item label="所学专业">{{ profile.major }}</el-descriptions-item>
        <el-descriptions-item label="人才计划">{{ profile.talentProgram }}</el-descriptions-item>
        <el-descriptions-item label="荣誉称号">{{ profile.honorTitle }}</el-descriptions-item>
      </el-descriptions>

      <el-divider content-position="left">学科领域 / 研究主题</el-divider>
      <el-descriptions :column="2" border>
        <el-descriptions-item label="一级学科">{{ profile.disciplineL1 }}</el-descriptions-item>
        <el-descriptions-item label="二级学科">{{ profile.disciplineL2 }}</el-descriptions-item>
        <el-descriptions-item label="交叉学科领域">{{ profile.crossDiscipline }}</el-descriptions-item>
        <el-descriptions-item label="学科平台/研究基地">{{ profile.disciplinePlatform }}</el-descriptions-item>
        <el-descriptions-item label="重点研究方向" :span="2">{{ profile.researchDirections }}</el-descriptions-item>
        <el-descriptions-item label="研究关键词" :span="2">{{ profile.researchKeywords }}</el-descriptions-item>
        <el-descriptions-item label="所属政策领域">{{ profile.policyFields }}</el-descriptions-item>
        <el-descriptions-item label="主要研究对象">{{ profile.researchObjects }}</el-descriptions-item>
        <el-descriptions-item label="当前关注议题">{{ profile.focusTopics }}</el-descriptions-item>
        <el-descriptions-item label="区域/国别研究专长">{{ profile.regionExpertise }}</el-descriptions-item>
        <el-descriptions-item label="研究方法专长" :span="2">{{ profile.methodExpertise }}</el-descriptions-item>
      </el-descriptions>

      <el-divider content-position="left">审核与推荐排序</el-divider>
      <el-descriptions :column="3" border>
        <el-descriptions-item label="审核状态">{{ statusLabel(profile.status) }}</el-descriptions-item>
        <el-descriptions-item label="成果分">{{ profile.achievementScore }}</el-descriptions-item>
        <el-descriptions-item label="决策影响分">{{ profile.influenceScore }}</el-descriptions-item>
        <el-descriptions-item label="社会贡献分">{{ profile.socialScore }}</el-descriptions-item>
        <el-descriptions-item label="综合排序得分" :span="2">{{ profile.compositeScore }}</el-descriptions-item>
      </el-descriptions>
    </div>

    <div class="gva-table-box">
      <el-tabs v-model="activeTab">
        <el-tab-pane label="研究成果" name="achievement">
          <div v-if="!isReadOnlyReviewer" class="gva-btn-list">
            <el-button type="primary" icon="plus" @click="openAchievementDialog">新增</el-button>
          </div>
          <el-table :data="achievementData" style="width: 100%">
            <el-table-column align="left" label="成果类型" prop="achievementType" width="120" />
            <el-table-column align="left" label="成果名称" prop="title" min-width="200" show-overflow-tooltip />
            <el-table-column align="left" label="成果级别" prop="level" width="100" />
            <el-table-column align="left" label="发表/立项单位" prop="publishOrg" width="160" />
            <el-table-column v-if="!isReadOnlyReviewer" align="left" label="操作" fixed="right" width="150">
              <template #default="scope">
                <el-button type="primary" link icon="edit" @click="openAchievementDialog(scope.row)">编辑</el-button>
                <el-button type="primary" link icon="delete" @click="deleteAchievement(scope.row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <el-tab-pane label="决策影响记录" name="adoption">
          <div v-if="!isReadOnlyReviewer" class="gva-btn-list">
            <el-button type="primary" icon="plus" @click="openAdoptionDialog">新增</el-button>
          </div>
          <el-table :data="adoptionData" style="width: 100%">
            <el-table-column align="left" label="采纳类型" prop="adoptionType" width="180" />
            <el-table-column align="left" label="采纳/批示单位级别" prop="adoptingUnitLevel" width="140" />
            <el-table-column align="left" label="采纳/批示单位名称" prop="adoptingUnitName" min-width="180" show-overflow-tooltip />
            <el-table-column v-if="!isReadOnlyReviewer" align="left" label="操作" fixed="right" width="150">
              <template #default="scope">
                <el-button type="primary" link icon="edit" @click="openAdoptionDialog(scope.row)">编辑</el-button>
                <el-button type="primary" link icon="delete" @click="deleteAdoption(scope.row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <el-tab-pane label="学术兼职" name="position">
          <div v-if="!isReadOnlyReviewer" class="gva-btn-list">
            <el-button type="primary" icon="plus" @click="openPositionDialog">新增</el-button>
          </div>
          <el-table :data="positionData" style="width: 100%">
            <el-table-column align="left" label="兼职类型" prop="positionType" width="180" />
            <el-table-column align="left" label="任职机构名称" prop="organizationName" min-width="180" show-overflow-tooltip />
            <el-table-column align="left" label="担任职务" prop="positionTitle" width="140" />
            <el-table-column v-if="!isReadOnlyReviewer" align="left" label="操作" fixed="right" width="150">
              <template #default="scope">
                <el-button type="primary" link icon="edit" @click="openPositionDialog(scope.row)">编辑</el-button>
                <el-button type="primary" link icon="delete" @click="deletePosition(scope.row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <el-tab-pane label="标签" name="tag">
          <p style="color: #909399; font-size: 13px; margin-bottom: 12px;">
            标签会直接影响"专家检索推荐"里的相关性排序，建议按学科/关键词/政策领域/研究方法等打全。
          </p>
          <el-select
            v-model="selectedTagIds"
            multiple filterable
            :disabled="isReadOnlyReviewer"
            placeholder="选择标签"
            style="width: 100%; max-width: 640px;"
          >
            <el-option-group v-for="group in tagGroups" :key="group.type" :label="group.label">
              <el-option v-for="item in group.options" :key="item.ID" :label="item.tagValue" :value="item.ID" />
            </el-option-group>
          </el-select>
          <div v-if="!isReadOnlyReviewer" style="margin-top: 16px;">
            <el-button type="primary" @click="submitTags">保存标签</el-button>
          </div>
        </el-tab-pane>
      </el-tabs>
    </div>

    <!-- 研究成果 弹窗 -->
    <el-dialog v-model="achievementDialogVisible" :title="achievementForm.ID ? '编辑成果' : '添加成果'" width="50%">
      <el-form ref="achievementFormRef" :model="achievementForm" label-width="120px">
        <el-form-item label="成果类型" prop="achievementType" required>
          <el-select v-model="achievementForm.achievementType" placeholder="请选择" style="width:100%">
            <el-option v-for="item in achievementTypeOptions" :key="item" :label="item" :value="item" />
          </el-select>
        </el-form-item>
        <el-form-item label="成果名称" prop="title" required><el-input v-model="achievementForm.title" /></el-form-item>
        <el-form-item label="成果级别" prop="level"><el-input v-model="achievementForm.level" placeholder="如：国家级/省部级/厅局级/一般级" /></el-form-item>
        <el-form-item label="级别认定来源" prop="levelSource"><el-input v-model="achievementForm.levelSource" /></el-form-item>
        <el-form-item label="发表/出版/立项单位" prop="publishOrg"><el-input v-model="achievementForm.publishOrg" /></el-form-item>
        <el-form-item label="发表/出版/立项时间" prop="publishDate">
          <el-date-picker v-model="achievementForm.publishDate" type="date" style="width:100%" />
        </el-form-item>
        <el-form-item label="关键词" prop="keywords"><el-input v-model="achievementForm.keywords" placeholder="用于相关性检索" /></el-form-item>
        <el-form-item label="备注" prop="remark"><el-input v-model="achievementForm.remark" type="textarea" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="achievementDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitAchievement">确定</el-button>
      </template>
    </el-dialog>

    <!-- 决策影响记录 弹窗 -->
    <el-dialog v-model="adoptionDialogVisible" :title="adoptionForm.ID ? '编辑决策影响记录' : '添加决策影响记录'" width="50%">
      <el-form ref="adoptionFormRef" :model="adoptionForm" label-width="140px">
        <el-form-item label="采纳类型" prop="adoptionType" required>
          <el-select v-model="adoptionForm.adoptionType" placeholder="请选择" style="width:100%">
            <el-option v-for="item in adoptionTypeOptions" :key="item" :label="item" :value="item" />
          </el-select>
        </el-form-item>
        <el-form-item label="采纳/批示单位级别" prop="adoptingUnitLevel"><el-input v-model="adoptionForm.adoptingUnitLevel" placeholder="如：国家级/省部级/厅局级/区县级" /></el-form-item>
        <el-form-item label="采纳/批示单位名称" prop="adoptingUnitName"><el-input v-model="adoptionForm.adoptingUnitName" /></el-form-item>
        <el-form-item label="采纳/批示时间" prop="adoptionDate">
          <el-date-picker v-model="adoptionForm.adoptionDate" type="date" style="width:100%" />
        </el-form-item>
        <el-form-item label="佐证材料描述" prop="evidenceDesc"><el-input v-model="adoptionForm.evidenceDesc" type="textarea" /></el-form-item>
        <el-form-item label="佐证材料附件" prop="attachmentUrl"><el-input v-model="adoptionForm.attachmentUrl" placeholder="附件链接" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="adoptionDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitAdoption">确定</el-button>
      </template>
    </el-dialog>

    <!-- 学术兼职 弹窗 -->
    <el-dialog v-model="positionDialogVisible" :title="positionForm.ID ? '编辑学术兼职' : '添加学术兼职'" width="50%">
      <el-form ref="positionFormRef" :model="positionForm" label-width="120px">
        <el-form-item label="兼职类型" prop="positionType"><el-input v-model="positionForm.positionType" placeholder="如：学术团体任职/专业委员会任职" /></el-form-item>
        <el-form-item label="任职机构名称" prop="organizationName" required><el-input v-model="positionForm.organizationName" /></el-form-item>
        <el-form-item label="担任职务" prop="positionTitle"><el-input v-model="positionForm.positionTitle" /></el-form-item>
        <el-form-item label="任职起始时间" prop="startDate">
          <el-date-picker v-model="positionForm.startDate" type="date" style="width:100%" />
        </el-form-item>
        <el-form-item label="任职结束时间" prop="endDate">
          <el-date-picker v-model="positionForm.endDate" type="date" placeholder="在任请留空" style="width:100%" />
        </el-form-item>
        <el-form-item label="备注" prop="remark"><el-input v-model="positionForm.remark" type="textarea" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="positionDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitPosition">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { findExpertProfile } from '@/api/expertProfile'
import {
  createExpertAchievement,
  updateExpertAchievement,
  deleteExpertAchievement,
  getExpertAchievementList
} from '@/api/expertAchievement'
import {
  createExpertAdoptionRecord,
  updateExpertAdoptionRecord,
  deleteExpertAdoptionRecord,
  getExpertAdoptionRecordList
} from '@/api/expertAdoptionRecord'
import {
  createExpertAcademicPosition,
  updateExpertAcademicPosition,
  deleteExpertAcademicPosition,
  getExpertAcademicPositionList
} from '@/api/expertAcademicPosition'
import {
  getExpertTagList,
  getExpertTagsByExpertId,
  setExpertTagRelations
} from '@/api/expertTag'

import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/pinia/modules/user'

defineOptions({
  name: 'ExpertProfileDetail'
})

const route = useRoute()
const router = useRouter()
const expertId = Number(route.params.id)
const userStore = useUserStore()
// 单位审核员/市级审核员在后端只有只读权限，写操作一定会被 Casbin 拒绝——
// 前端直接不展示这些按钮，避免审核员填完一整张表单才发现白填了
const isReadOnlyReviewer = computed(() => [9002, 9003].includes(userStore.userInfo.authorityId))

const statusMap = {
  draft: '草稿',
  pending_org_review: '待单位审核',
  org_rejected: '单位已退回',
  pending_city_review: '待市级审核',
  city_rejected: '市级已退回',
  published: '已发布'
}
const statusLabel = (value) => statusMap[value] || value

const profile = ref({})
const getProfile = async () => {
  const res = await findExpertProfile({ ID: expertId })
  if (res.code === 0) {
    profile.value = res.data.reExpertProfile
  }
}
getProfile()

const activeTab = ref('achievement')

// ============ 研究成果 ============
const achievementData = ref([])
const fetchAchievements = async () => {
  const res = await getExpertAchievementList({ expertId, page: 1, pageSize: 100 })
  if (res.code === 0) achievementData.value = res.data.list
}
fetchAchievements()

const achievementTypeOptions = ['著作', '学术论文', '研究报告', '决策咨询成果', '获奖成果', '课题项目']
const achievementDialogVisible = ref(false)
const achievementFormRef = ref()
const achievementForm = ref({})
const openAchievementDialog = (row) => {
  achievementForm.value = row ? { ...row } : { expertId, achievementType: '', title: '', level: '', levelSource: '', publishOrg: '', publishDate: null, keywords: '', remark: '' }
  achievementDialogVisible.value = true
}
const submitAchievement = async () => {
  const payload = { ...achievementForm.value, expertId }
  const res = payload.ID ? await updateExpertAchievement(payload) : await createExpertAchievement(payload)
  if (res.code === 0) {
    ElMessage({ type: 'success', message: '保存成功' })
    achievementDialogVisible.value = false
    fetchAchievements()
  }
}
const deleteAchievement = (row) => {
  ElMessageBox.confirm('确定要删除吗?', '提示', { type: 'warning' }).then(async () => {
    const res = await deleteExpertAchievement({ ID: row.ID })
    if (res.code === 0) {
      ElMessage({ type: 'success', message: '删除成功' })
      fetchAchievements()
    }
  })
}

// ============ 决策影响记录 ============
const adoptionData = ref([])
const fetchAdoptionRecords = async () => {
  const res = await getExpertAdoptionRecordList({ expertId, page: 1, pageSize: 100 })
  if (res.code === 0) adoptionData.value = res.data.list
}
fetchAdoptionRecords()

const adoptionTypeOptions = ['内参/专报采用', '部门采纳', '领导批示', '进入政策文件', '参与政策起草/咨询论证']
const adoptionDialogVisible = ref(false)
const adoptionFormRef = ref()
const adoptionForm = ref({})
const openAdoptionDialog = (row) => {
  adoptionForm.value = row ? { ...row } : { expertId, adoptionType: '', adoptingUnitLevel: '', adoptingUnitName: '', adoptionDate: null, evidenceDesc: '', attachmentUrl: '' }
  adoptionDialogVisible.value = true
}
const submitAdoption = async () => {
  const payload = { ...adoptionForm.value, expertId }
  const res = payload.ID ? await updateExpertAdoptionRecord(payload) : await createExpertAdoptionRecord(payload)
  if (res.code === 0) {
    ElMessage({ type: 'success', message: '保存成功' })
    adoptionDialogVisible.value = false
    fetchAdoptionRecords()
  }
}
const deleteAdoption = (row) => {
  ElMessageBox.confirm('确定要删除吗?', '提示', { type: 'warning' }).then(async () => {
    const res = await deleteExpertAdoptionRecord({ ID: row.ID })
    if (res.code === 0) {
      ElMessage({ type: 'success', message: '删除成功' })
      fetchAdoptionRecords()
    }
  })
}

// ============ 学术兼职 ============
const positionData = ref([])
const fetchPositions = async () => {
  const res = await getExpertAcademicPositionList({ expertId, page: 1, pageSize: 100 })
  if (res.code === 0) positionData.value = res.data.list
}
fetchPositions()

const positionDialogVisible = ref(false)
const positionFormRef = ref()
const positionForm = ref({})
const openPositionDialog = (row) => {
  positionForm.value = row ? { ...row } : { expertId, positionType: '', organizationName: '', positionTitle: '', startDate: null, endDate: null, remark: '' }
  positionDialogVisible.value = true
}
const submitPosition = async () => {
  const payload = { ...positionForm.value, expertId }
  const res = payload.ID ? await updateExpertAcademicPosition(payload) : await createExpertAcademicPosition(payload)
  if (res.code === 0) {
    ElMessage({ type: 'success', message: '保存成功' })
    positionDialogVisible.value = false
    fetchPositions()
  }
}
const deletePosition = (row) => {
  ElMessageBox.confirm('确定要删除吗?', '提示', { type: 'warning' }).then(async () => {
    const res = await deleteExpertAcademicPosition({ ID: row.ID })
    if (res.code === 0) {
      ElMessage({ type: 'success', message: '删除成功' })
      fetchPositions()
    }
  })
}

// ============ 标签 ============
// 标签直接参与"专家检索推荐"的相关性排序，之前只有标签库本身的维护页面，没有把标签关联到
// 具体某个专家的入口——加在这里，跟成果/决策影响/学术兼职一样是这个专家详情页的一部分
const tagTypeLabels = {
  discipline_l1: '一级学科', keyword: '关键词', policy_field: '政策领域', region: '区域/国别', method: '研究方法'
}
const allTags = ref([])
const selectedTagIds = ref([])
const tagGroups = computed(() => {
  const byType = {}
  allTags.value.forEach(tag => {
    if (!byType[tag.tagType]) byType[tag.tagType] = []
    byType[tag.tagType].push(tag)
  })
  return Object.keys(byType).map(type => ({
    type, label: tagTypeLabels[type] || type, options: byType[type]
  }))
})
const fetchAllTags = async () => {
  const res = await getExpertTagList({ page: 1, pageSize: 500 })
  if (res.code === 0) allTags.value = res.data.list || []
}
const fetchExpertTags = async () => {
  const res = await getExpertTagsByExpertId({ expertId })
  if (res.code === 0) selectedTagIds.value = (res.data.tags || []).map(t => t.ID)
}
fetchAllTags()
fetchExpertTags()

const submitTags = async () => {
  const res = await setExpertTagRelations({ expertId, tagIds: selectedTagIds.value })
  if (res.code === 0) {
    ElMessage({ type: 'success', message: '保存成功' })
  }
}
</script>

<style></style>
