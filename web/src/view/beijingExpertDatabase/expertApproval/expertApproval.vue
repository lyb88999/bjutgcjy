<template>
  <div class="gva-table-box">
    <el-tabs v-model="activeTab" @tab-change="onTabChange">
      <el-tab-pane label="我发起的" name="mine">
        <el-table :data="mineData" style="width: 100%">
          <el-table-column align="left" label="姓名" prop="name" width="110" />
          <el-table-column align="left" label="所在单位" prop="unitName" min-width="220" show-overflow-tooltip />
          <el-table-column align="left" label="专业技术职称" prop="techTitle" min-width="130" show-overflow-tooltip />
          <el-table-column align="left" label="一级学科" prop="disciplineL1" min-width="130" show-overflow-tooltip />
          <el-table-column align="left" label="审核状态" prop="status" width="150">
            <template #default="scope">
              {{ statusLabel(scope.row.status) }}
              <el-tag v-if="scope.row.reviewBypassed" type="warning" size="small" style="margin-left: 4px;">免审核发布</el-tag>
            </template>
          </el-table-column>
          <el-table-column align="left" label="创建日期" width="170">
            <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
          </el-table-column>
          <el-table-column align="left" label="操作" fixed="right" width="240">
            <template #default="scope">
              <el-button type="primary" link @click="viewDetail(scope.row)">查看详情</el-button>
              <el-button
                v-if="['draft', 'org_rejected', 'city_rejected'].includes(scope.row.status)"
                type="primary" link @click="doSubmit(scope.row)"
              >提交审核</el-button>
              <el-button type="primary" link @click="openLog(scope.row)">审核记录</el-button>
            </template>
          </el-table-column>
        </el-table>
        <div class="gva-pagination">
          <el-pagination :current-page="minePage" :page-size="minePageSize" :total="mineTotal"
            layout="total, prev, pager, next" @current-change="v => { minePage = v; fetchMine() }" />
        </div>
      </el-tab-pane>

      <el-tab-pane label="待我审核的" name="pending">
        <template v-if="showOrgSection">
          <el-divider content-position="left">待本单位审核</el-divider>
          <div class="batch-bar">
            <el-button
              type="primary" :disabled="!orgSelection.length" @click="openBatchApprove('org')"
            >批量通过（{{ orgSelection.length }}）</el-button>
          </div>
          <el-table :data="orgData" style="width: 100%" @selection-change="v => orgSelection = v">
            <el-table-column type="selection" width="45" />
            <el-table-column align="left" label="姓名" prop="name" width="110" />
            <el-table-column align="left" label="所在单位" prop="unitName" min-width="220" show-overflow-tooltip />
            <el-table-column align="left" label="专业技术职称" prop="techTitle" min-width="130" show-overflow-tooltip />
            <el-table-column align="left" label="一级学科" prop="disciplineL1" min-width="130" show-overflow-tooltip />
            <el-table-column align="left" label="提交日期" width="170">
              <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
            </el-table-column>
            <el-table-column align="left" label="操作" fixed="right" width="260">
              <template #default="scope">
                <el-button type="primary" link @click="viewDetail(scope.row)">查看详情</el-button>
                <el-button type="primary" link @click="orgApprove(scope.row)">通过</el-button>
                <el-button type="danger" link @click="openReject(scope.row, 'org')">退回</el-button>
                <el-button type="primary" link @click="openLog(scope.row)">审核记录</el-button>
              </template>
            </el-table-column>
          </el-table>
          <div class="gva-pagination">
            <el-pagination :current-page="orgPage" :page-size="orgPageSize" :total="orgTotal"
              layout="total, prev, pager, next" @current-change="v => { orgPage = v; fetchOrgPending() }" />
          </div>
        </template>
        <template v-if="showCitySection">
          <el-divider content-position="left">待市级审核</el-divider>
          <div class="batch-bar">
            <el-button
              type="primary" :disabled="!citySelection.length" @click="openBatchApprove('city')"
            >批量通过（{{ citySelection.length }}）</el-button>
          </div>
          <el-table :data="cityData" style="width: 100%" @selection-change="v => citySelection = v">
            <el-table-column type="selection" width="45" />
            <el-table-column align="left" label="姓名" prop="name" width="110" />
            <el-table-column align="left" label="所在单位" prop="unitName" min-width="220" show-overflow-tooltip />
            <el-table-column align="left" label="专业技术职称" prop="techTitle" min-width="130" show-overflow-tooltip />
            <el-table-column align="left" label="一级学科" prop="disciplineL1" min-width="130" show-overflow-tooltip />
            <el-table-column align="left" label="提交日期" width="170">
              <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
            </el-table-column>
            <el-table-column align="left" label="操作" fixed="right" width="260">
              <template #default="scope">
                <el-button type="primary" link @click="viewDetail(scope.row)">查看详情</el-button>
                <el-button type="primary" link @click="cityApprove(scope.row)">通过</el-button>
                <el-button type="danger" link @click="openReject(scope.row, 'city')">退回</el-button>
                <el-button type="primary" link @click="openLog(scope.row)">审核记录</el-button>
              </template>
            </el-table-column>
          </el-table>
          <div class="gva-pagination">
            <el-pagination :current-page="cityPage" :page-size="cityPageSize" :total="cityTotal"
              layout="total, prev, pager, next" @current-change="v => { cityPage = v; fetchCityPending() }" />
          </div>
        </template>
        <el-empty v-if="!showOrgSection && !showCitySection" description="当前角色没有需要审核的记录" />
      </el-tab-pane>
    </el-tabs>

    <el-dialog v-model="batchApproveDialogVisible" title="批量审核通过" width="420px">
      <p>
        确定要批量通过选中的 <strong>{{ batchApproveTarget === 'org' ? orgSelection.length : citySelection.length }}</strong> 条记录吗？
        {{ batchApproveTarget === 'city' ? '通过后会直接发布，进入检索排序范围。' : '通过后会流转到市级审核。' }}
      </p>
      <p style="color: #909399; font-size: 12px;">批量操作跳过逐条核对，请确认这些记录确实不需要单独退回后再继续。</p>
      <template #footer>
        <el-button @click="batchApproveDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="batchApproving" @click="confirmBatchApprove">确定通过</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="rejectDialogVisible" title="退回并填写审核意见" width="40%">
      <el-form :model="rejectForm">
        <el-form-item label="审核意见" required>
          <el-input v-model="rejectForm.opinion" type="textarea" :rows="4" placeholder="请说明退回原因" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="rejectDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmReject">确定退回</el-button>
      </template>
    </el-dialog>

    <el-drawer v-model="logDrawerVisible" size="480" title="审核记录">
      <el-timeline>
        <el-timeline-item v-for="item in logData" :key="item.ID" :timestamp="formatDate(item.CreatedAt)">
          {{ statusLabel(item.fromStatus) }} → {{ statusLabel(item.toStatus) }}
          <div v-if="item.opinion" style="color: #f56c6c;">意见：{{ item.opinion }}</div>
        </el-timeline-item>
      </el-timeline>
      <el-empty v-if="!logData.length" description="暂无审核记录" />
    </el-drawer>
  </div>
</template>

<script setup>
import {
  submitExpertProfile,
  orgApproveExpertProfile,
  batchOrgApproveExpertProfile,
  orgRejectExpertProfile,
  cityApproveExpertProfile,
  batchCityApproveExpertProfile,
  cityRejectExpertProfile,
  getMyDrafts,
  getPendingOrgReview,
  getPendingCityReview,
  getExpertApprovalLogList
} from '@/api/expertApproval'

import { formatDate } from '@/utils/format'
import { ElMessage } from 'element-plus'
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/pinia/modules/user'

defineOptions({
  name: 'ExpertApproval'
})

const router = useRouter()
const userStore = useUserStore()
// 审核台是三个审核角色共用的同一个页面，"待本单位审核"/"待市级审核"两个区块原来对谁都无条件展示——
// 单位审核员因此能看到跟自己无关的"待市级审核"队列（虽然按不了通过/退回，但看到了不该看的东西）。
// 这里按角色只展示各自负责的那一块；管理员两块都能看，个人申报人两块都不展示。
const isSuperAdmin = computed(() => userStore.userInfo.authorityId === 1)
const showOrgSection = computed(() => isSuperAdmin.value || userStore.userInfo.authorityId === 9002)
const showCitySection = computed(() => isSuperAdmin.value || userStore.userInfo.authorityId === 9003)
const viewDetail = (row) => {
  router.push({ name: 'expertProfileDetail', params: { id: row.ID } })
}

const statusMap = {
  draft: '草稿',
  pending_org_review: '待单位审核',
  org_rejected: '单位已退回',
  pending_city_review: '待市级审核',
  city_rejected: '市级已退回',
  published: '已发布'
}
const statusLabel = (value) => statusMap[value] || value

const activeTab = ref('mine')

const mineData = ref([])
const minePage = ref(1)
const minePageSize = ref(10)
const mineTotal = ref(0)
const fetchMine = async () => {
  const res = await getMyDrafts({ page: minePage.value, pageSize: minePageSize.value })
  if (res.code === 0) {
    mineData.value = res.data.list
    mineTotal.value = res.data.total
  }
}

const orgData = ref([])
const orgPage = ref(1)
const orgPageSize = ref(20)
const orgTotal = ref(0)
const orgSelection = ref([])
const fetchOrgPending = async () => {
  const res = await getPendingOrgReview({ page: orgPage.value, pageSize: orgPageSize.value })
  if (res.code === 0) {
    orgData.value = res.data.list
    orgTotal.value = res.data.total
  }
}

const cityData = ref([])
const cityPage = ref(1)
const cityPageSize = ref(20)
const cityTotal = ref(0)
const citySelection = ref([])
const fetchCityPending = async () => {
  const res = await getPendingCityReview({ page: cityPage.value, pageSize: cityPageSize.value })
  if (res.code === 0) {
    cityData.value = res.data.list
    cityTotal.value = res.data.total
  }
}

const onTabChange = (name) => {
  if (name === 'mine') fetchMine()
  if (name === 'pending') {
    if (showOrgSection.value) fetchOrgPending()
    if (showCitySection.value) fetchCityPending()
  }
}

fetchMine()

const doSubmit = async (row) => {
  const res = await submitExpertProfile({ expertId: row.ID })
  if (res.code === 0) {
    ElMessage({ type: 'success', message: '提交成功' })
    fetchMine()
  }
}

const orgApprove = async (row) => {
  const res = await orgApproveExpertProfile({ expertId: row.ID })
  if (res.code === 0) {
    ElMessage({ type: 'success', message: '审核通过' })
    fetchOrgPending()
  }
}

const cityApprove = async (row) => {
  const res = await cityApproveExpertProfile({ expertId: row.ID })
  if (res.code === 0) {
    ElMessage({ type: 'success', message: '审核通过，专家已发布' })
    fetchCityPending()
  }
}

// 批量通过特意加一步二次确认，不做成勾完就静默一键全过——批量操作本来就少了逐条核对，
// 至少让审核员在点下去之前再确认一遍，别手滑批量放过了不该放过的记录
const batchApproveDialogVisible = ref(false)
const batchApproveTarget = ref('')
const batchApproving = ref(false)
const openBatchApprove = (target) => {
  batchApproveTarget.value = target
  batchApproveDialogVisible.value = true
}
const confirmBatchApprove = async () => {
  const isOrg = batchApproveTarget.value === 'org'
  const ids = (isOrg ? orgSelection.value : citySelection.value).map(row => row.ID)
  if (!ids.length) return
  batchApproving.value = true
  try {
    const res = isOrg
      ? await batchOrgApproveExpertProfile({ expertIds: ids })
      : await batchCityApproveExpertProfile({ expertIds: ids })
    if (res.code === 0) {
      const { successCount, failCount, failures } = res.data
      if (failCount > 0) {
        ElMessage({
          type: 'warning',
          message: `成功 ${successCount} 条，失败 ${failCount} 条：${failures.map(f => `#${f.expertId} ${f.message}`).join('；')}`,
          duration: 8000,
          showClose: true
        })
      } else {
        ElMessage({ type: 'success', message: `批量通过成功，共 ${successCount} 条` })
      }
      batchApproveDialogVisible.value = false
      if (isOrg) {
        orgSelection.value = []
        fetchOrgPending()
      } else {
        citySelection.value = []
        fetchCityPending()
      }
    }
  } finally {
    batchApproving.value = false
  }
}

const rejectDialogVisible = ref(false)
const rejectForm = ref({ expertId: null, opinion: '', scope: '' })
const openReject = (row, scope) => {
  rejectForm.value = { expertId: row.ID, opinion: '', scope }
  rejectDialogVisible.value = true
}
const confirmReject = async () => {
  if (!rejectForm.value.opinion) {
    ElMessage({ type: 'warning', message: '请填写审核意见' })
    return
  }
  const payload = { expertId: rejectForm.value.expertId, opinion: rejectForm.value.opinion }
  const res = rejectForm.value.scope === 'org'
    ? await orgRejectExpertProfile(payload)
    : await cityRejectExpertProfile(payload)
  if (res.code === 0) {
    ElMessage({ type: 'success', message: '已退回' })
    rejectDialogVisible.value = false
    fetchOrgPending()
    fetchCityPending()
  }
}

const logDrawerVisible = ref(false)
const logData = ref([])
const openLog = async (row) => {
  const res = await getExpertApprovalLogList({ expertId: row.ID, page: 1, pageSize: 50 })
  if (res.code === 0) {
    logData.value = res.data.list
    logDrawerVisible.value = true
  }
}
</script>

<style>
.batch-bar {
  margin-bottom: 10px;
}
</style>
