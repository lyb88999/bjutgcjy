<template>
  <div class="expert-detail-page">
    <div class="hero-card">
      <div class="hero-main">
        <el-button
          class="hero-back"
          icon="back"
          text
          @click="router.back()"
        >返回</el-button>
        <div class="hero-identity">
          <div class="hero-avatar">{{ (profile.name || '?').slice(0, 1) }}</div>
          <div>
            <div class="hero-name-row">
              <span class="hero-name">{{ profile.name }}</span>
              <el-tag
                :type="statusTagType(profile.status)"
                size="small"
              >{{ statusLabel(profile.status) }}</el-tag>
              <el-tag
                v-if="profile.reviewBypassed"
                type="warning"
                size="small"
              >免审核发布</el-tag>
            </div>
            <div class="hero-sub">
              {{ [profile.unitName, profile.department, profile.techTitle].filter(Boolean).join(' · ') || '单位/职称信息待补充' }}
            </div>
            <div class="hero-chips">
              <el-tag
                v-if="profile.disciplineL1"
                effect="plain"
                size="small"
              >{{ profile.disciplineL1 }}</el-tag>
              <el-tag
                v-if="profile.disciplineL2"
                effect="plain"
                size="small"
                type="info"
              >{{ profile.disciplineL2 }}</el-tag>
              <el-tag
                v-if="profile.talentProgram"
                effect="plain"
                size="small"
                type="success"
              >{{ profile.talentProgram }}</el-tag>
            </div>
          </div>
        </div>
      </div>
      <div class="hero-scores">
        <div class="hero-score">
          <div class="hero-score-value">{{ profile.titleScore }}</div>
          <div class="hero-score-label">职称分</div>
        </div>
        <div
          v-if="scoreColumnAvailability.hasAchievement"
          class="hero-score"
        >
          <div class="hero-score-value">{{ profile.achievementScore }}</div>
          <div class="hero-score-label">成果分</div>
        </div>
        <div
          v-if="scoreColumnAvailability.hasInfluence"
          class="hero-score"
        >
          <div class="hero-score-value">{{ profile.influenceScore }}</div>
          <div class="hero-score-label">决策影响分</div>
        </div>
        <div
          v-if="scoreColumnAvailability.hasSocial"
          class="hero-score"
        >
          <div class="hero-score-value">{{ profile.socialScore }}</div>
          <div class="hero-score-label">社会贡献分</div>
        </div>
        <div class="hero-score hero-score-primary">
          <div class="hero-score-value">{{ profile.compositeScore }}</div>
          <div class="hero-score-label">综合排序得分</div>
        </div>
      </div>
    </div>

    <div
      class="gva-table-box"
      style="margin-bottom: 16px"
    >
      <el-divider content-position="left">背景信息</el-divider>
      <div class="info-grid">
        <div
          v-for="group in backgroundGroups"
          :key="group.title"
          class="info-card"
        >
          <div class="info-card-title">{{ group.title }}</div>
          <div
            v-for="[label, key] in group.fields"
            :key="key"
            class="info-row"
          >
            <span class="info-label">{{ label }}</span>
            <span class="info-value">{{ dv(profile[key]) }}</span>
          </div>
        </div>
      </div>

      <el-divider content-position="left">学科领域 / 研究主题</el-divider>
      <div class="info-grid">
        <div
          v-for="group in disciplineGroups"
          :key="group.title"
          class="info-card"
        >
          <div class="info-card-title">{{ group.title }}</div>
          <div
            v-for="[label, key] in group.fields"
            :key="key"
            class="info-row"
          >
            <span class="info-label">{{ label }}</span>
            <span class="info-value">{{ dv(profile[key]) }}</span>
          </div>
        </div>

        <div class="info-card info-card-wide">
          <div class="info-card-title">研究方向</div>
          <div class="info-block">
            <span class="info-label">重点研究方向</span>
            <p class="info-paragraph">{{ dv(profile.researchDirections) }}</p>
          </div>
          <div class="info-block">
            <span class="info-label">研究关键词</span>
            <div class="kw-tag-wrap">
              <template v-if="keywordList.length">
                <el-tag
                  v-for="kw in keywordList"
                  :key="kw"
                  size="small"
                  effect="plain"
                  class="kw-tag"
                >{{ kw }}</el-tag>
              </template>
              <span
                v-else
                class="info-value"
              >—</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="gva-table-box">
      <el-tabs v-model="activeTab">
        <el-tab-pane
          :label="`研究成果 (${achievementData.length})`"
          name="achievement"
        >
          <div
            v-if="!isReadOnlyReviewer"
            class="gva-btn-list"
          >
            <el-button
              type="primary"
              icon="plus"
              @click="openAchievementDialog"
            >新增</el-button>
          </div>
          <el-table
            :data="achievementData"
            style="width: 100%"
          >
            <el-table-column
              align="left"
              label="成果类型"
              prop="achievementType"
              width="120"
            />
            <el-table-column
              align="left"
              label="成果名称"
              prop="title"
              min-width="200"
              show-overflow-tooltip
            />
            <el-table-column
              align="left"
              label="成果级别"
              prop="level"
              width="100"
            />
            <el-table-column
              align="left"
              label="发表/立项单位"
              prop="publishOrg"
              width="160"
            />
            <el-table-column
              v-if="!isReadOnlyReviewer"
              align="left"
              label="操作"
              fixed="right"
              width="150"
            >
              <template #default="scope">
                <el-button
                  type="primary"
                  link
                  icon="edit"
                  @click="openAchievementDialog(scope.row)"
                >编辑</el-button>
                <el-button
                  type="primary"
                  link
                  icon="delete"
                  @click="deleteAchievement(scope.row)"
                >删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <el-tab-pane
          :label="`学术兼职 (${positionData.length})`"
          name="position"
        >
          <div
            v-if="!isReadOnlyReviewer"
            class="gva-btn-list"
          >
            <el-button
              type="primary"
              icon="plus"
              @click="openPositionDialog"
            >新增</el-button>
          </div>
          <el-table
            :data="positionData"
            style="width: 100%"
          >
            <el-table-column
              align="left"
              label="兼职类型"
              prop="positionType"
              width="180"
            />
            <el-table-column
              align="left"
              label="任职机构名称"
              prop="organizationName"
              min-width="180"
              show-overflow-tooltip
            />
            <el-table-column
              align="left"
              label="担任职务"
              prop="positionTitle"
              width="140"
            />
            <el-table-column
              v-if="!isReadOnlyReviewer"
              align="left"
              label="操作"
              fixed="right"
              width="150"
            >
              <template #default="scope">
                <el-button
                  type="primary"
                  link
                  icon="edit"
                  @click="openPositionDialog(scope.row)"
                >编辑</el-button>
                <el-button
                  type="primary"
                  link
                  icon="delete"
                  @click="deletePosition(scope.row)"
                >删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <el-tab-pane
          :label="`标签 (${selectedTagIds.length})`"
          name="tag"
        >
          <p style="color: #909399; font-size: 13px; margin-bottom: 12px;">
            标签会直接影响"专家检索推荐"里的相关性排序，建议按学科/关键词/政策领域/研究方法等打全。
          </p>
          <el-select
            v-model="selectedTagIds"
            multiple
            filterable
            :disabled="isReadOnlyReviewer"
            placeholder="选择标签"
            style="width: 100%; max-width: 640px;"
          >
            <el-option-group
              v-for="group in tagGroups"
              :key="group.type"
              :label="group.label"
            >
              <el-option
                v-for="item in group.options"
                :key="item.ID"
                :label="item.tagValue"
                :value="item.ID"
              />
            </el-option-group>
          </el-select>
          <div
            v-if="!isReadOnlyReviewer"
            style="margin-top: 16px;"
          >
            <el-button
              type="primary"
              @click="submitTags"
            >保存标签</el-button>
          </div>
        </el-tab-pane>
      </el-tabs>
    </div>

    <!-- 研究成果 弹窗 -->
    <el-dialog
      v-model="achievementDialogVisible"
      :title="achievementForm.ID ? '编辑成果' : '添加成果'"
      width="50%"
    >
      <el-form
        ref="achievementFormRef"
        :model="achievementForm"
        label-width="120px"
      >
        <el-form-item
          label="成果类型"
          prop="achievementType"
          required
        >
          <el-select
            v-model="achievementForm.achievementType"
            placeholder="请选择"
            style="width:100%"
          >
            <el-option
              v-for="item in achievementTypeOptions"
              :key="item"
              :label="item"
              :value="item"
            />
          </el-select>
        </el-form-item>
        <el-form-item
          label="成果名称"
          prop="title"
          required
        ><el-input v-model="achievementForm.title" /></el-form-item>
        <el-form-item
          label="成果级别"
          prop="level"
        ><el-input
          v-model="achievementForm.level"
          placeholder="如：国家级/省部级/厅局级/一般级"
        /></el-form-item>
        <el-form-item
          label="级别认定来源"
          prop="levelSource"
        ><el-input v-model="achievementForm.levelSource" /></el-form-item>
        <el-form-item
          label="发表/出版/立项单位"
          prop="publishOrg"
        ><el-input v-model="achievementForm.publishOrg" /></el-form-item>
        <el-form-item
          label="发表/出版/立项时间"
          prop="publishDate"
        >
          <el-date-picker
            v-model="achievementForm.publishDate"
            type="date"
            style="width:100%"
          />
        </el-form-item>
        <el-form-item
          label="关键词"
          prop="keywords"
        ><el-input
          v-model="achievementForm.keywords"
          placeholder="用于相关性检索"
        /></el-form-item>
        <el-form-item
          label="备注"
          prop="remark"
        ><el-input
          v-model="achievementForm.remark"
          type="textarea"
        /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="achievementDialogVisible = false">取消</el-button>
        <el-button
          type="primary"
          @click="submitAchievement"
        >确定</el-button>
      </template>
    </el-dialog>

    <!-- 学术兼职 弹窗 -->
    <el-dialog
      v-model="positionDialogVisible"
      :title="positionForm.ID ? '编辑学术兼职' : '添加学术兼职'"
      width="50%"
    >
      <el-form
        ref="positionFormRef"
        :model="positionForm"
        label-width="120px"
      >
        <el-form-item
          label="兼职类型"
          prop="positionType"
        ><el-input
          v-model="positionForm.positionType"
          placeholder="如：学术团体任职/专业委员会任职"
        /></el-form-item>
        <el-form-item
          label="任职机构名称"
          prop="organizationName"
          required
        ><el-input v-model="positionForm.organizationName" /></el-form-item>
        <el-form-item
          label="担任职务"
          prop="positionTitle"
        ><el-input v-model="positionForm.positionTitle" /></el-form-item>
        <el-form-item
          label="任职起始时间"
          prop="startDate"
        >
          <el-date-picker
            v-model="positionForm.startDate"
            type="date"
            style="width:100%"
          />
        </el-form-item>
        <el-form-item
          label="任职结束时间"
          prop="endDate"
        >
          <el-date-picker
            v-model="positionForm.endDate"
            type="date"
            placeholder="在任请留空"
            style="width:100%"
          />
        </el-form-item>
        <el-form-item
          label="备注"
          prop="remark"
        ><el-input
          v-model="positionForm.remark"
          type="textarea"
        /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="positionDialogVisible = false">取消</el-button>
        <el-button
          type="primary"
          @click="submitPosition"
        >确定</el-button>
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
import { getScoreColumnAvailability } from '@/api/expertSearch'

import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { statusLabel, statusTagType } from '../expertStatus'
import { useUserStore } from '@/pinia/modules/user'

defineOptions({
  name: 'ExpertProfileDetail'
})

const route = useRoute()
const router = useRouter()
const expertId = Number(route.params.id)
const userStore = useUserStore()

// 决策影响分/社会贡献分在全库都还没数据支撑时（批量导入的花名册类专家还没配决策影响/学术兼职
// 记录），展示一张清一色 0 的卡片没有意义，查一下全库有没有人非零来决定要不要显示
const scoreColumnAvailability = ref({ hasAchievement: true, hasInfluence: true, hasSocial: true })
getScoreColumnAvailability().then((res) => {
  if (res.code === 0) {
    scoreColumnAvailability.value = res.data
  }
})
// 单位审核员/市级审核员/专家（纯检索用户）在后端只有只读权限，写操作一定会被 Casbin 拒绝——
// 前端直接不展示这些按钮，避免填完一整张表单才发现白填了。个人申报人比较特殊：casbin 层面
// 有写权限（因为要能编辑自己申报的档案），但后端 canModifyProfile 只放行"自己创建的档案"，
// 别人的档案一样会被拒绝——这里补上跟后端一致的归属判断，不然从专家检索点进别人已发布的档案，
// 个人申报人账号还是会看到一整套编辑按钮，点了才发现后端拒绝，体验很差（这正是这次要修的问题：
// 专家检索本该是纯查询入口，点进去却看得到一堆其实点不动的编辑控件）
const isReadOnlyReviewer = computed(() => {
  const authorityId = userStore.userInfo.authorityId
  if ([9002, 9003, 9005].includes(authorityId)) return true
  if (authorityId === 9001) return profile.value.CreatedBy !== userStore.userInfo.ID
  return false
})

// 空字段统一显示占位符，避免整块描述表里大片空白格看起来像渲染出错
const dv = (v) => (v === undefined || v === null || v === '' ? '—' : v)

// 原来是两张 17/11 字段的巨型 el-descriptions 大表，字段挤在一起不好扫读；按语义分成小卡片
// 网格（跟检索页结果卡片一样的卡片化观感），每组字段数量不多，一眼能扫完。字段本身一个都没删，
// 只是重新分组+换了渲染方式（label/value 行而不是表格单元格）
// 2026-08-25 统计全库23661位已发布专家的字段填充率后，从这个页面去掉了民族/政治面貌/手机号码/
// 荣誉称号/交叉学科领域/所属政策领域/主要研究对象/当前关注议题/区域国别研究专长这9个字段——
// 全部是0%~1%填充，展示出来只有一整屏"—"占位符，没有信息量。字段本身在数据模型和编辑表单
// （expertProfile.vue）里都已经同步去掉，不是只藏了展示，数据本身没删
const backgroundGroups = [
  { title: '基本信息', fields: [['姓名', 'name'], ['性别', 'gender']] },
  { title: '任职信息', fields: [['所在单位', 'unitName'], ['院系/部门', 'department'], ['行政职务', 'adminTitle'], ['专业技术职称', 'techTitle']] },
  { title: '联系方式', fields: [['办公电话', 'phone'], ['电子邮箱', 'email'], ['通讯地址', 'address']] },
  { title: '教育背景', fields: [['最高学历', 'highestEducation'], ['最高学位', 'highestDegree'], ['毕业院校', 'graduateSchool'], ['所学专业', 'major']] },
  { title: '人才计划', fields: [['人才计划', 'talentProgram']] }
]
// 重点研究方向/研究关键词是长文本或标签列表，不适合塞进 label/value 行，单独用一张
// "研究方向"卡片渲染（见模板），不在这个数据驱动的分组列表里。"政策与专长"整组（所属政策领域/
// 主要研究对象/当前关注议题/区域国别研究专长）四个字段全是0%填充，整组一起去掉了
const disciplineGroups = [
  { title: '学科归属', fields: [['一级学科', 'disciplineL1'], ['二级学科', 'disciplineL2'], ['学科平台/研究基地', 'disciplinePlatform']] }
]

const keywordList = computed(() =>
  (profile.value.researchKeywords || '')
    .split(/[,，、;；\s]+/)
    .map((s) => s.trim())
    .filter(Boolean)
)

const profile = ref({})
const getProfile = async() => {
  const res = await findExpertProfile({ ID: expertId })
  if (res.code === 0) {
    profile.value = res.data.reExpertProfile
  }
}
getProfile()

const activeTab = ref('achievement')

// ============ 研究成果 ============
const achievementData = ref([])
const fetchAchievements = async() => {
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
const submitAchievement = async() => {
  const payload = { ...achievementForm.value, expertId }
  const res = payload.ID ? await updateExpertAchievement(payload) : await createExpertAchievement(payload)
  if (res.code === 0) {
    ElMessage({ type: 'success', message: '保存成功' })
    achievementDialogVisible.value = false
    fetchAchievements()
  }
}
const deleteAchievement = (row) => {
  ElMessageBox.confirm('确定要删除吗?', '提示', { type: 'warning' }).then(async() => {
    const res = await deleteExpertAchievement({ ID: row.ID })
    if (res.code === 0) {
      ElMessage({ type: 'success', message: '删除成功' })
      fetchAchievements()
    }
  })
}

// ============ 学术兼职 ============
const positionData = ref([])
const fetchPositions = async() => {
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
const submitPosition = async() => {
  const payload = { ...positionForm.value, expertId }
  const res = payload.ID ? await updateExpertAcademicPosition(payload) : await createExpertAcademicPosition(payload)
  if (res.code === 0) {
    ElMessage({ type: 'success', message: '保存成功' })
    positionDialogVisible.value = false
    fetchPositions()
  }
}
const deletePosition = (row) => {
  ElMessageBox.confirm('确定要删除吗?', '提示', { type: 'warning' }).then(async() => {
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
const fetchAllTags = async() => {
  const res = await getExpertTagList({ page: 1, pageSize: 500 })
  if (res.code === 0) allTags.value = res.data.list || []
}
// 标签字典这次导入之后从几十条暴涨到2.7万+条，fetchAllTags 只拉第一页500条（按id倒序），
// 专家自己已选的标签 ID 很可能不落在这最新500条里——这才是"已选标签显示数字ID"的真正原因
// （不是异步时序问题：就算等字典加载完再设置已选标签，只要那个 ID 本来就不在这500条里，
// el-select 一样找不到对应 label）。getExpertTagsByExpertId 本身返回的就是完整标签对象
// （含 tagValue），所以直接把专家已选的这些标签对象补进字典，不管它们是否在字典分页窗口内，
// 保证选中态一定能解析出真实文字；下拉框仍然只是字典的一页，不受影响
const fetchExpertTags = async() => {
  const res = await getExpertTagsByExpertId({ expertId })
  if (res.code === 0) {
    const tags = res.data.tags || []
    const known = new Set(allTags.value.map(t => t.ID))
    const missing = tags.filter(t => !known.has(t.ID))
    if (missing.length) allTags.value = [...allTags.value, ...missing]
    selectedTagIds.value = tags.map(t => t.ID)
  }
}
// 先等字典加载完再取已选标签：fetchAllTags 是整体覆盖式赋值(allTags.value = res.data.list)，
// 如果反过来 fetchExpertTags 先跑完把缺失的标签对象合并进 allTags，字典请求晚一步返回时会把
// 这次合并直接覆盖掉，等于白合并
const initTags = async() => {
  await fetchAllTags()
  await fetchExpertTags()
}
initTags()

const submitTags = async() => {
  const res = await setExpertTagRelations({ expertId, tagIds: selectedTagIds.value })
  if (res.code === 0) {
    ElMessage({ type: 'success', message: '保存成功' })
  }
}
</script>

<style scoped>
/* 跟检索页/结果卡片（expertSearch.vue、ExpertResultCard.vue）同一套品牌色，只做局部
   视觉换肤——信息结构（hero + 两张大表 + tabs + 三个弹窗）完全不动，只换配色/字体/间距，
   不碰全局 el-color-primary，弹窗里的按钮仍用默认 Element 蓝，跟全站操作类控件保持一致 */
.expert-detail-page {
  --brand-navy: #16407a;
  --brand-navy-dark: #0c2a54;
  --brand-red: #b5272d;
  --brand-gold: #d4a017;
}

.hero-card {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 24px;
  flex-wrap: wrap;
  background: var(--el-bg-color);
  border: 1px solid var(--el-border-color-lighter);
  border-top: 3px solid var(--brand-navy);
  border-radius: 12px;
  padding: 20px 24px;
  margin-bottom: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}
.hero-back {
  margin-bottom: 8px;
  padding-left: 0;
}
.hero-identity {
  display: flex;
  align-items: center;
  gap: 16px;
}
.hero-avatar {
  width: 56px;
  height: 56px;
  min-width: 56px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-family: "Songti SC", "Noto Serif SC", "SimSun", serif;
  font-size: 24px;
  font-weight: 700;
  color: #fff;
  background: linear-gradient(155deg, var(--brand-navy), var(--brand-red));
}
.hero-name-row {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}
.hero-name {
  font-family: "Songti SC", "Noto Serif SC", "SimSun", serif;
  font-size: 22px;
  font-weight: 700;
  color: var(--el-text-color-primary);
}
.hero-sub {
  margin-top: 4px;
  font-size: 13px;
  color: var(--el-text-color-secondary);
}
.hero-chips {
  margin-top: 8px;
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
}
.hero-scores {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}
.hero-score {
  min-width: 96px;
  padding: 12px 16px;
  border-radius: 10px;
  background: var(--el-fill-color-light);
  text-align: center;
}
.hero-score-primary {
  background: linear-gradient(135deg, var(--brand-navy-dark), var(--brand-navy));
}
.hero-score-primary .hero-score-value,
.hero-score-primary .hero-score-label {
  color: #fff;
}
.hero-score-value {
  font-size: 22px;
  font-weight: 700;
  color: var(--el-text-color-primary);
  font-variant-numeric: tabular-nums;
  line-height: 1.2;
}
.hero-score-label {
  margin-top: 4px;
  font-size: 12px;
  color: var(--el-text-color-secondary);
}
.kw-tag {
  margin-right: 6px;
  margin-bottom: 4px;
  color: var(--brand-navy);
  background-color: rgba(22, 64, 122, 0.06);
  border-color: rgba(22, 64, 122, 0.22);
}
.kw-tag-wrap {
  display: flex;
  flex-wrap: wrap;
}

/* 背景信息/学科领域两块原来的巨型 el-descriptions，改成小卡片网格——宽屏自动摆 2-3 列，
   窄屏自动收成 1 列，卡片本身比表格单元格更透气，标题用金色小竖条做重点标记，不是纯装饰，
   是跟"重点研究方向/研究关键词"合并的宽卡片（.info-card-wide）做视觉区分 */
.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
  gap: 14px;
  margin-bottom: 8px;
}
.info-card {
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 10px;
  padding: 14px 16px;
  background: var(--el-fill-color-blank);
  transition: border-color 120ms ease, box-shadow 120ms ease;
}
.info-card:hover {
  border-color: var(--brand-navy);
  box-shadow: 0 2px 10px rgba(22, 64, 122, 0.08);
}
.info-card-wide {
  grid-column: 1 / -1;
}
.info-card-title {
  display: flex;
  align-items: center;
  gap: 7px;
  font-family: "Songti SC", "Noto Serif SC", "SimSun", serif;
  font-weight: 700;
  font-size: 14px;
  color: var(--brand-navy);
  margin-bottom: 10px;
}
.info-card-title::before {
  content: '';
  width: 3px;
  height: 13px;
  border-radius: 2px;
  background: var(--brand-gold);
}
.info-row {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  gap: 12px;
  padding: 5px 0;
  border-bottom: 1px dashed var(--el-border-color-lighter);
  font-size: 13px;
}
.info-row:last-child {
  border-bottom: none;
}
.info-label {
  flex-shrink: 0;
  color: var(--el-text-color-secondary);
}
.info-value {
  color: var(--el-text-color-primary);
  text-align: right;
  word-break: break-word;
}
.info-block {
  margin-bottom: 12px;
}
.info-block:last-child {
  margin-bottom: 0;
}
.info-block > .info-label {
  display: block;
  margin-bottom: 4px;
}
.info-paragraph {
  margin: 0;
  font-size: 13px;
  line-height: 1.7;
  color: var(--el-text-color-primary);
}

/* 两块 el-descriptions 前面的分节标题（背景信息/学科领域），跟检索页卡片名同一套宋体，
   颜色改成品牌藏青而不是默认灰黑，让这两大块表格看起来是同一个产品的一部分 */
.expert-detail-page :deep(.el-divider__text) {
  font-family: "Songti SC", "Noto Serif SC", "SimSun", serif;
  font-size: 15px;
  font-weight: 700;
  color: var(--brand-navy);
}

/* tabs 的选中态跟检索页排序 tab（expertSearch.vue .sort-tab.active）用一样的藏青，
   不动全局 el-tabs 默认样式，只在这个页面局部生效 */
.expert-detail-page :deep(.el-tabs__item.is-active) {
  color: var(--brand-navy);
}
.expert-detail-page :deep(.el-tabs__item:hover) {
  color: var(--brand-navy);
}
.expert-detail-page :deep(.el-tabs__active-bar) {
  background-color: var(--brand-navy);
}
</style>
