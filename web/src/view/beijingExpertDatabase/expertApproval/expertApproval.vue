<template>
  <div class="gva-table-box">
    <el-tabs v-model="activeTab" @tab-change="onTabChange">
      <el-tab-pane label="我发起的" name="mine">
        <el-table :data="mineData" style="width: 100%">
          <el-table-column align="left" label="姓名" prop="name" width="120" />
          <el-table-column align="left" label="所在单位" prop="unitName" width="160" />
          <el-table-column align="left" label="审核状态" prop="status" width="140">
            <template #default="scope">{{ statusLabel(scope.row.status) }}</template>
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
        <el-divider content-position="left">待本单位审核</el-divider>
        <el-table :data="orgData" style="width: 100%">
          <el-table-column align="left" label="姓名" prop="name" width="120" />
          <el-table-column align="left" label="所在单位" prop="unitName" width="160" />
          <el-table-column align="left" label="专业技术职称" prop="techTitle" width="120" />
          <el-table-column align="left" label="操作" fixed="right" width="260">
            <template #default="scope">
              <el-button type="primary" link @click="viewDetail(scope.row)">查看详情</el-button>
              <el-button type="primary" link @click="orgApprove(scope.row)">通过</el-button>
              <el-button type="danger" link @click="openReject(scope.row, 'org')">退回</el-button>
              <el-button type="primary" link @click="openLog(scope.row)">审核记录</el-button>
            </template>
          </el-table-column>
        </el-table>
        <el-divider content-position="left">待市级审核</el-divider>
        <el-table :data="cityData" style="width: 100%">
          <el-table-column align="left" label="姓名" prop="name" width="120" />
          <el-table-column align="left" label="所在单位" prop="unitName" width="160" />
          <el-table-column align="left" label="专业技术职称" prop="techTitle" width="120" />
          <el-table-column align="left" label="操作" fixed="right" width="260">
            <template #default="scope">
              <el-button type="primary" link @click="viewDetail(scope.row)">查看详情</el-button>
              <el-button type="primary" link @click="cityApprove(scope.row)">通过</el-button>
              <el-button type="danger" link @click="openReject(scope.row, 'city')">退回</el-button>
              <el-button type="primary" link @click="openLog(scope.row)">审核记录</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>
    </el-tabs>

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
  orgRejectExpertProfile,
  cityApproveExpertProfile,
  cityRejectExpertProfile,
  getMyDrafts,
  getPendingOrgReview,
  getPendingCityReview,
  getExpertApprovalLogList
} from '@/api/expertApproval'

import { formatDate } from '@/utils/format'
import { ElMessage } from 'element-plus'
import { ref } from 'vue'
import { useRouter } from 'vue-router'

defineOptions({
  name: 'ExpertApproval'
})

const router = useRouter()
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
const fetchOrgPending = async () => {
  const res = await getPendingOrgReview({ page: 1, pageSize: 50 })
  if (res.code === 0) {
    orgData.value = res.data.list
  }
}

const cityData = ref([])
const fetchCityPending = async () => {
  const res = await getPendingCityReview({ page: 1, pageSize: 50 })
  if (res.code === 0) {
    cityData.value = res.data.list
  }
}

const onTabChange = (name) => {
  if (name === 'mine') fetchMine()
  if (name === 'pending') {
    fetchOrgPending()
    fetchCityPending()
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

<style></style>
