<template>
  <section
    class="expert-agent"
    aria-label="智能遴选专家"
  >
    <div class="agent-heading">
      <div>
        <h2>说清需求，找到合适的专家</h2>
        <p>从专家库中检索、筛选并列出推荐依据，支持继续补充条件。</p>
      </div>
      <el-button
        :disabled="!turns.length && !message"
        @click="newSession"
      >新会话</el-button>
    </div>

    <div
      v-if="statusLoading"
      class="service-notice"
      role="status"
    >正在检查智能遴选服务…</div>
    <div
      v-else-if="configured === false"
      class="service-notice notice-warning"
      role="status"
    >
      <span>智能遴选尚未配置，请联系管理员配置模型服务。您仍可使用「专家检索」查找专家。</span>
      <el-button
        link
        @click="loadStatus"
      >重新检查</el-button>
    </div>
    <div
      v-else-if="statusError"
      class="service-notice notice-warning"
      role="alert"
    >
      <span>{{ statusError }}</span>
      <el-button
        link
        @click="loadStatus"
      >重新检查</el-button>
    </div>

    <div class="agent-layout">
      <div class="conversation-panel">
        <div class="conversation-heading">
          <h3>遴选需求</h3>
          <span
            v-if="configured"
            class="service-ready"
            :title="model ? `当前模型：${model}` : ''"
          >服务已就绪</span>
        </div>
        <div
          v-if="!turns.length"
          class="agent-examples"
        >
          <p>可以从这些需求开始</p>
          <button
            v-for="example in examples"
            :key="example"
            type="button"
            :disabled="loading"
            @click="useExample(example)"
          >
            {{ example }}<span aria-hidden="true">↗</span>
          </button>
        </div>
        <div
          ref="conversationRef"
          class="conversation-history"
          role="log"
          aria-label="遴选对话"
          aria-live="polite"
        >
          <div
            v-for="turn in turns"
            :key="turn.id"
            class="conversation-turn"
          >
            <div class="message-label">我的需求</div>
            <p class="user-message">{{ turn.message }}</p>
            <template v-if="turn.status === 'success'">
              <div class="message-label">智能遴选</div>
              <p class="assistant-message">{{ turn.reply || '已完成本轮遴选，请核对右侧名单和依据。' }}</p>
            </template>
            <p
              v-else-if="turn.status === 'expired'"
              class="turn-note"
            >会话已过期，请新建会话继续。</p>
            <p
              v-else-if="turn.status === 'error'"
              class="turn-note"
            >本轮未完成，可重试或调整需求。</p>
            <p
              v-else-if="turn.status === 'cancelled'"
              class="turn-note"
            >本轮请求已取消。</p>
          </div>
        </div>

        <div
          v-if="loading"
          class="agent-progress"
          role="status"
          aria-live="polite"
        >
          <span
            class="progress-dot"
            aria-hidden="true"
          />
          <div>
            <strong>正在处理本轮需求</strong>
            <p>理解条件、检索专家并整理推荐依据 · {{ elapsedSeconds }} 秒</p>
          </div>
          <el-button
            link
            @click="cancelRequest"
          >取消</el-button>
        </div>
        <div
          v-if="errorText"
          class="request-error"
          role="alert"
        >
          <p>{{ errorText }}</p>
          <el-button
            v-if="sessionExpired"
            link
            @click="recoverExpiredSession"
          >新建会话</el-button>
          <el-button
            v-else-if="lastAttempt && configured"
            link
            :disabled="loading"
            @click="retryLastAttempt"
          >重试这一轮</el-button>
        </div>

        <form
          class="agent-composer"
          @submit.prevent="submitMessage"
        >
          <label
            class="composer-label"
            for="expert-agent-message"
          >{{ sessionId ? '补充或调整条件' : '描述您需要的专家' }}</label>
          <el-input
            id="expert-agent-message"
            ref="messageInputRef"
            v-model="message"
            type="textarea"
            :autosize="{ minRows: 4, maxRows: 8 }"
            :maxlength="2000"
            :disabled="loading"
            placeholder="例如：找5位研究城市更新与社会治理的专家，优先教授，单位尽量多样。"
            @keydown="onComposerKeydown"
          />
          <div class="selection-controls">
            <div>
              <label for="expert-agent-count">推荐人数</label>
              <el-input-number
                id="expert-agent-count"
                v-model="count"
                :min="1"
                :max="30"
                :step="1"
                :precision="0"
                :disabled="loading"
                controls-position="right"
                @change="countEdited = true"
              />
            </div>
            <div>
              <label for="expert-agent-unit-limit">每单位最多</label>
              <el-input-number
                v-if="maxPerUnit !== 0"
                id="expert-agent-unit-limit"
                v-model="maxPerUnit"
                :min="1"
                :max="30"
                :step="1"
                :precision="0"
                :disabled="loading"
                controls-position="right"
                @change="maxPerUnitEdited = true"
              />
              <div
                v-else
                class="unlimited-control"
              >
                <span>不限</span>
                <el-button
                  link
                  :disabled="loading"
                  @click="setUnitLimit"
                >设置上限</el-button>
              </div>
            </div>
          </div>
          <p class="composer-hint">也可在对话中调整人数；手动修改的数字用于下一轮遴选。Ctrl / ⌘ + Enter 发送。</p>
          <el-button
            class="send-button"
            type="primary"
            native-type="submit"
            :disabled="!canSend"
            :loading="loading"
          >
            {{ sessionId ? '按新条件遴选' : '开始遴选' }}
          </el-button>
        </form>
      </div>

      <div
        class="agent-results"
        :aria-busy="loading"
      >
        <div
          v-if="!result"
          class="result-placeholder"
        >
          <span
            class="placeholder-mark"
            aria-hidden="true"
          >选</span>
          <h3>推荐名单与依据将在这里展示</h3>
          <p>输入主题、专家人数和单位要求，结果可查看专家档案、核对库内依据，并导出清单。</p>
          <div class="placeholder-stages"><span>理解需求</span><span>检索与筛选</span><span>核对推荐依据</span></div>
        </div>
        <template v-else>
          <div
            v-if="loading"
            class="previous-result-note"
          >下方保留上一轮结果，新的遴选完成后更新。</div>
          <div class="result-heading">
            <div><h3>本轮遴选结果</h3><p>推荐供筛选参考，请结合专家档案核对。</p></div>
            <el-button
              :disabled="!result.recommendations.length"
              @click="exportRecommendations"
            >导出推荐清单</el-button>
          </div>
          <div class="result-metrics">
            <div
              v-for="metric in metrics"
              :key="metric.label"
            ><strong>{{ metric.value }}</strong><span>{{ metric.label }}</span></div>
          </div>
          <div
            v-if="planLabels.length"
            class="plan-summary"
            aria-label="本轮采用的筛选条件"
          >
            <span class="plan-label">本轮条件</span>
            <el-tag
              v-for="label in planLabels"
              :key="label"
              size="small"
              effect="plain"
            >{{ label }}</el-tag>
          </div>
          <ul
            v-if="result.warnings.length"
            class="result-warnings"
            aria-label="结果提示"
          >
            <li
              v-for="(warning, index) in result.warnings"
              :key="index"
            >{{ warning }}</li>
          </ul>
          <details
            v-if="result.steps.length"
            class="completed-steps"
          >
            <summary>查看本轮处理过程（{{ result.steps.length }} 步）</summary>
            <ol><li
              v-for="(step, index) in result.steps"
              :key="index"
            ><strong>{{ step.label }}</strong><p v-if="step.detail">{{ step.detail }}</p></li></ol>
          </details>

          <div
            v-if="result.recommendations.length && distributions.some((group) => group.items.length)"
            class="distribution-grid"
          >
            <section
              v-for="group in distributions"
              :key="group.title"
              class="distribution-panel"
            >
              <h4>{{ group.title }}<span>本轮推荐</span></h4>
              <div
                v-for="(item, index) in group.items"
                :key="`${item.label}-${index}`"
                class="distribution-row"
              >
                <span
                  class="distribution-label"
                  :title="item.label"
                >{{ item.label || '未标注' }}</span>
                <div class="distribution-track"><span :style="{ width: distributionWidth(item, group.items) }" /></div>
                <span class="distribution-count">{{ item.count }} 人</span>
              </div>
              <p
                v-if="!group.items.length"
                class="distribution-empty"
              >暂无分布数据</p>
            </section>
          </div>

          <div class="result-list-toolbar">
            <el-radio-group
              v-model="resultView"
              size="small"
              aria-label="选择名单范围"
              @change="resultPage = 1"
            >
              <el-radio-button label="recommendations">推荐名单（{{ result.recommendations.length }}）</el-radio-button>
              <el-radio-button label="candidates">候选名单（{{ result.candidates.length }}）</el-radio-button>
            </el-radio-group>
          </div>
          <div
            v-if="visibleExperts.length"
            class="expert-card-list"
          >
            <ExpertAgentCandidate
              v-for="(expert, index) in visibleExperts"
              :key="`${resultView}-${expert.id}-${index}`"
              :item="expert"
              :index="(resultPage - 1) * resultsPerPage + index"
            />
          </div>
          <el-empty
            v-else
            :image-size="70"
            :description="resultView === 'recommendations' ? '本轮没有推荐名单，可继续补充或调整需求。' : '本轮没有候选专家，可放宽筛选条件再试。'"
          />
          <el-pagination
            v-if="currentExperts.length > resultsPerPage"
            v-model:current-page="resultPage"
            class="agent-pagination"
            :page-size="resultsPerPage"
            :total="currentExperts.length"
            layout="prev, pager, next"
            small
          />
        </template>
      </div>
    </div>
  </section>
</template>

<script setup>
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { getExpertAgentStatus, chatWithExpertAgent } from '@/api/expertAgent'
import ExpertAgentCandidate from './ExpertAgentCandidate.vue'
import { normalizeAgentResult, recommendationsCsv } from './expertAgentUtils'

defineOptions({ name: 'ExpertAgent' })
const props = defineProps({ active: { type: Boolean, default: true }})

const examples = [
  '找5位研究城市更新与社会治理的专家，单位尽量多样。',
  '为数字经济课题推荐专家，重点关注产业转型与公共政策。',
  '寻找教育公平领域的教授，希望能看到相关成果依据。'
]
const configured = ref(null)
const model = ref('')
const statusLoading = ref(false)
const statusError = ref('')
const message = ref('')
const sessionId = ref('')
const count = ref(5)
const maxPerUnit = ref(2)
const countEdited = ref(false)
const maxPerUnitEdited = ref(false)
const turns = ref([])
const result = ref(null)
const loading = ref(false)
const errorText = ref('')
const sessionExpired = ref(false)
const elapsedSeconds = ref(0)
const lastAttempt = ref(null)
const resultView = ref('recommendations')
const resultPage = ref(1)
const resultsPerPage = 10
const conversationRef = ref(null)
const messageInputRef = ref(null)
let statusController
let chatController
let progressTimer
let requestSequence = 0
let turnSequence = 0
let disposed = false

const canSend = computed(() => configured.value === true && !loading.value && !sessionExpired.value && !!message.value.trim())
const currentExperts = computed(() => result.value?.[resultView.value] || [])
const visibleExperts = computed(() => currentExperts.value.slice((resultPage.value - 1) * resultsPerPage, resultPage.value * resultsPerPage))
const metrics = computed(() => {
  if (!result.value) return []
  const stats = result.value.statistics
  return [
    { label: '候选专家', value: stats.candidateCount },
    { label: '本轮推荐', value: stats.selectedCount },
    { label: '覆盖单位', value: stats.unitCount },
    { label: '覆盖学科', value: stats.disciplineCount }
  ]
})
const distributions = computed(() => [
  { title: '单位分布', items: result.value?.statistics.units || [] },
  { title: '学科分布', items: result.value?.statistics.disciplines || [] }
])
const distributionWidth = (item, items) => `${item.count / Math.max(1, ...items.map((entry) => entry.count)) * 100}%`
const planLabels = computed(() => {
  const plan = result.value?.plan || {}
  const labels = []
  if (Array.isArray(plan.keywords)) labels.push(...plan.keywords.filter((word) => typeof word === 'string' && word).map((word) => `主题：${word}`))
  for (const [key, label] of [['name', '姓名'], ['unitName', '单位'], ['disciplineL1', '学科'], ['titleTier', '职称']]) {
    if (typeof plan[key] === 'string' && plan[key]) labels.push(`${label}：${plan[key]}`)
  }
  if (Number.isInteger(plan.count) && plan.count > 0) labels.push(`目标 ${plan.count} 人`)
  if (Number.isInteger(plan.maxPerUnit) && plan.maxPerUnit > 0) labels.push(`每单位最多 ${plan.maxPerUnit} 人`)
  if (plan.maxPerUnit === 0) labels.push('单位人数不限')
  if (plan.diversify) labels.push('兼顾单位多样性')
  if (Array.isArray(plan.excludedExpertIds) && plan.excludedExpertIds.length) labels.push(`排除 ${plan.excludedExpertIds.length} 位专家`)
  return [...new Set(labels)]
})

const unwrap = (response) => {
  if (response instanceof Error) throw response
  if (!response || response.code !== 0) {
    const error = new Error(response?.msg || '连接中断或服务暂不可用，请稍后重试。')
    error.response = response?.response
    throw error
  }
  return response.data
}
const describeError = (error) => {
  const status = error?.response?.status
  const detail = typeof error?.response?.data?.msg === 'string'
    ? error.response.data.msg
    : (typeof error?.message === 'string' ? error.message : '')
  if (status === 429 || /繁忙|限流|rate.?limit|too many requests/i.test(detail)) return '模型服务繁忙，请稍后重试。'
  if (['ECONNABORTED', 'ETIMEDOUT'].includes(error?.code) || /timeout|超时/i.test(detail)) return '本轮处理超时，请稍后重试，或缩小需求范围。'
  if (/network error|failed to fetch/i.test(detail)) return '暂时无法连接服务，请检查网络后重试。'
  if (status === 401 || status === 403) return '当前账号无法使用智能遴选，请重新登录或联系管理员确认权限。'
  if (status === 404) return '智能遴选服务暂未就绪，请联系管理员检查。'
  return detail ? detail.slice(0, 300) : '本轮遴选暂未完成，请稍后重试。'
}

const loadStatus = async() => {
  statusController?.abort()
  const controller = new AbortController()
  statusController = controller
  statusLoading.value = true
  statusError.value = ''
  try {
    const response = await getExpertAgentStatus({ signal: controller.signal })
    if (disposed || controller.signal.aborted || statusController !== controller) return
    const data = unwrap(response)
    if (typeof data?.configured !== 'boolean') throw new Error('无法确认智能遴选服务状态，请重试。')
    configured.value = data.configured
    model.value = typeof data.model === 'string' ? data.model : ''
  } catch (error) {
    if (!disposed && !controller.signal.aborted && statusController === controller) {
      configured.value = null
      statusError.value = describeError(error)
    }
  } finally {
    if (statusController === controller) {
      statusLoading.value = false
      statusController = null
    }
  }
}

const scrollConversation = async() => {
  await nextTick()
  if (!disposed && conversationRef.value) conversationRef.value.scrollTop = conversationRef.value.scrollHeight
}
const stopProgress = () => {
  clearInterval(progressTimer)
  progressTimer = null
}
const cancelRequest = () => {
  if (!chatController) return
  const controller = chatController
  chatController = null
  requestSequence++
  controller.abort()
  stopProgress()
  loading.value = false
  const turn = turns.value.find((item) => item.id === lastAttempt.value?.turnId)
  if (turn?.status === 'pending') turn.status = 'cancelled'
  errorText.value = result.value ? '已取消本轮请求，上一轮结果已保留。' : '已取消本轮请求，可修改需求或重试。'
}

const sendRequest = async(payload, retryTurnId) => {
  if (loading.value || sessionExpired.value || configured.value !== true || !props.active) return
  const requestId = ++requestSequence
  const controller = new AbortController()
  chatController = controller
  loading.value = true
  errorText.value = ''
  elapsedSeconds.value = 0
  stopProgress()
  progressTimer = setInterval(() => { elapsedSeconds.value++ }, 1000)
  let turn = retryTurnId ? turns.value.find((item) => item.id === retryTurnId) : null
  if (!turn) {
    const turnId = ++turnSequence
    turns.value.push({ id: turnId, message: payload.message, status: 'pending', reply: '' })
    turn = turns.value.find((item) => item.id === turnId)
  } else {
    turn.status = 'pending'
  }
  lastAttempt.value = { payload: { ...payload }, turnId: turn.id }
  scrollConversation()
  try {
    const response = await chatWithExpertAgent(payload, { signal: controller.signal })
    if (disposed || controller.signal.aborted || requestId !== requestSequence) return
    const nextResult = normalizeAgentResult(unwrap(response))
    result.value = nextResult
    sessionId.value = nextResult.sessionId
    turn.status = 'success'
    turn.reply = nextResult.reply
    resultView.value = 'recommendations'
    resultPage.value = 1
    if (Number.isInteger(nextResult.plan.count) && nextResult.plan.count >= 1 && nextResult.plan.count <= 30) count.value = nextResult.plan.count
    if (Number.isInteger(nextResult.plan.maxPerUnit) && nextResult.plan.maxPerUnit >= 0 && nextResult.plan.maxPerUnit <= 30) maxPerUnit.value = nextResult.plan.maxPerUnit
    countEdited.value = false
    maxPerUnitEdited.value = false
    if (message.value.trim() === payload.message) message.value = ''
    lastAttempt.value = null
  } catch (error) {
    if (!disposed && !controller.signal.aborted && requestId === requestSequence) {
      const detail = describeError(error)
      sessionExpired.value = /会话不存在或已过期|session (?:not found|expired)/i.test(detail)
      turn.status = sessionExpired.value ? 'expired' : 'error'
      if (sessionExpired.value) sessionId.value = ''
      errorText.value = detail + (result.value ? ' 上一轮结果已保留。' : '')
    }
  } finally {
    if (requestId === requestSequence) {
      chatController = null
      loading.value = false
      stopProgress()
      scrollConversation()
    }
  }
}

const submitMessage = () => {
  if (!canSend.value) return
  const payload = { message: message.value.trim() }
  if (sessionId.value) payload.sessionId = sessionId.value
  if (countEdited.value && Number.isInteger(count.value)) payload.count = count.value
  if (maxPerUnitEdited.value && Number.isInteger(maxPerUnit.value)) {
    if (maxPerUnit.value > 0) {
      payload.maxPerUnit = maxPerUnit.value
    } else if (maxPerUnit.value === 0) {
      // 接口中的0表示沿用条件；新会话保留“不限”时，用可见需求明确表达，不能发送0覆盖。
      const unlimitedInstruction = '每单位人数不限。'
      if (!payload.message.includes(unlimitedInstruction)) payload.message += `\n${unlimitedInstruction}`
      if (payload.message.length > 2000) {
        errorText.value = '保留“不限单位”设置需要补充一行条件，请缩短需求后再发送（最多2000字）。'
        return
      }
      message.value = payload.message
    }
  }
  sendRequest(payload)
}
const retryLastAttempt = () => {
  if (!sessionExpired.value && lastAttempt.value) sendRequest({ ...lastAttempt.value.payload }, lastAttempt.value.turnId)
}
const onComposerKeydown = (event) => {
  if ((event.ctrlKey || event.metaKey) && event.key === 'Enter' && !event.isComposing) {
    event.preventDefault()
    submitMessage()
  }
}
const useExample = (example) => {
  message.value = example
  messageInputRef.value?.focus()
}
const setUnitLimit = () => {
  maxPerUnit.value = 2
  maxPerUnitEdited.value = true
}
const newSession = () => {
  cancelRequest()
  requestSequence++
  turns.value = []
  sessionId.value = ''
  result.value = null
  message.value = ''
  lastAttempt.value = null
  errorText.value = ''
  sessionExpired.value = false
  count.value = 5
  maxPerUnit.value = 2
  countEdited.value = false
  maxPerUnitEdited.value = false
  resultPage.value = 1
  messageInputRef.value?.focus()
}
const recoverExpiredSession = () => {
  const previousResult = result.value
  const previousAttempt = lastAttempt.value?.payload
  const unfinishedMessage = message.value || previousAttempt?.message || ''
  const recoverSetting = (edited, current, attempted, minimum, fallback) =>
    [edited ? current : undefined, attempted, current]
      .find((value) => Number.isInteger(value) && value >= minimum && value <= 30) ?? fallback
  const pendingCount = recoverSetting(countEdited.value, count.value, previousAttempt?.count, 1, 5)
  const pendingMaxPerUnit = recoverSetting(maxPerUnitEdited.value, maxPerUnit.value, previousAttempt?.maxPerUnit, 0, 2)
  newSession()
  // 新会话不携带旧上下文，但保留已完成的名单和未发送需求，供用户核对、补充后再发送。
  result.value = previousResult
  message.value = unfinishedMessage
  count.value = pendingCount
  countEdited.value = true
  maxPerUnit.value = pendingMaxPerUnit
  maxPerUnitEdited.value = true
  errorText.value = '已准备新会话，请完整描述本轮需求后发送。' + (previousResult ? ' 上一轮结果已保留供核对。' : '')
}
const exportRecommendations = () => {
  if (!result.value?.recommendations.length) return
  const blob = new Blob([recommendationsCsv(result.value.recommendations)], { type: 'text/csv;charset=utf-8;' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = '专家智能遴选推荐清单.csv'
  document.body.appendChild(link)
  link.click()
  link.remove()
  setTimeout(() => URL.revokeObjectURL(url), 0)
}

watch(() => props.active, (active) => {
  if (!active) {
    cancelRequest()
    statusController?.abort()
  } else if (configured.value === null) {
    loadStatus()
  }
})
onMounted(() => { if (props.active) loadStatus() })
onBeforeUnmount(() => {
  disposed = true
  statusController?.abort()
  cancelRequest()
  stopProgress()
})
</script>

<style scoped>
.expert-agent { color: var(--el-text-color-primary); }
.agent-heading,
.result-heading,
.conversation-heading { display: flex; align-items: center; justify-content: space-between; gap: 16px; }
.agent-heading { margin: 4px 0 24px; }
.agent-heading h2 { font-size: 21px; font-weight: 600; margin: 0; }
.agent-heading p,
.result-heading p { margin: 8px 0 0; color: var(--el-text-color-secondary); font-size: 13px; line-height: 1.7; }
.agent-layout { display: grid; grid-template-columns: minmax(290px, 350px) minmax(0, 1fr); gap: 24px; align-items: start; }
.conversation-panel { border: 1px solid var(--el-border-color-lighter); border-radius: 6px; padding: 20px; background: var(--el-bg-color); }
.conversation-heading h3,
.result-heading h3 { font-size: 16px; margin: 0; font-weight: 600; }
.service-ready { font-size: 11px; color: var(--el-color-success); }
.service-notice { display: flex; align-items: center; justify-content: space-between; gap: 16px; padding: 12px 16px; margin-bottom: 20px; font-size: 13px; line-height: 1.7; background: var(--el-fill-color-light); border-radius: 4px; }
.notice-warning { color: var(--el-color-warning-dark-2); background: var(--el-color-warning-light-9); }
.agent-examples { margin-top: 20px; }
.agent-examples > p { font-size: 12px; color: var(--el-text-color-secondary); margin: 0 0 10px; }
.agent-examples button { width: 100%; display: flex; align-items: center; gap: 10px; justify-content: space-between; text-align: left; margin: 8px 0; padding: 12px; border: 1px solid var(--el-border-color-lighter); background: var(--el-fill-color-extra-light); color: var(--el-text-color-regular); border-radius: 4px; font: inherit; font-size: 13px; line-height: 1.7; cursor: pointer; }
.agent-examples button:hover { border-color: var(--el-color-primary-light-5); color: var(--el-color-primary); }
.agent-examples button > span { color: var(--el-text-color-placeholder); }
.conversation-history { max-height: 470px; overflow-y: auto; overflow-wrap: anywhere; scrollbar-width: thin; }
.conversation-turn { margin-top: 20px; }
.message-label { color: var(--el-text-color-secondary); font-size: 11px; margin-bottom: 7px; }
.user-message,
.assistant-message { margin: 0 0 14px; white-space: pre-wrap; font-size: 13px; line-height: 1.9; }
.user-message { background: var(--el-color-primary-light-9); border-radius: 4px; padding: 11px 13px; }
.assistant-message { padding: 0 2px; }
.turn-note { margin: 0; font-size: 12px; color: var(--el-text-color-secondary); line-height: 1.7; }
.agent-progress { display: flex; align-items: center; gap: 10px; margin: 18px 0; padding: 12px; background: var(--el-fill-color-light); border-radius: 4px; }
.agent-progress > div { flex: 1; }
.agent-progress strong { font-size: 12px; font-weight: 500; }
.agent-progress p { margin: 5px 0 0; font-size: 11px; line-height: 1.7; color: var(--el-text-color-secondary); }
.progress-dot { width: 8px; height: 8px; flex: 0 0 8px; border-radius: 50%; background: var(--el-color-primary); animation: agent-pulse 1.5s ease-in-out infinite; }
.request-error { margin-top: 16px; padding: 12px; background: var(--el-color-warning-light-9); color: var(--el-color-warning-dark-2); border-radius: 4px; font-size: 12px; line-height: 1.8; overflow-wrap: anywhere; }
.request-error p { margin: 0; }
.agent-composer { margin-top: 22px; padding-top: 18px; border-top: 1px solid var(--el-border-color-extra-light); }
.composer-label { display: block; font-size: 13px; margin-bottom: 10px; }
.selection-controls { display: grid; grid-template-columns: 1fr 1fr; gap: 12px; margin-top: 14px; }
.selection-controls label { display: block; margin-bottom: 8px; font-size: 12px; color: var(--el-text-color-regular); }
.selection-controls :deep(.el-input-number) { width: 100%; }
.unlimited-control { display: flex; align-items: center; justify-content: space-between; min-height: 32px; font-size: 13px; }
.composer-hint { font-size: 11px; color: var(--el-text-color-secondary); margin: 12px 0; line-height: 1.8; }
.send-button { width: 100%; }
.agent-results { min-width: 0; }
.result-placeholder { min-height: 410px; display: flex; flex-direction: column; justify-content: center; align-items: center; text-align: center; padding: 32px; background: var(--el-fill-color-extra-light); border: 1px dashed var(--el-border-color); border-radius: 6px; }
.placeholder-mark { display: grid; place-items: center; height: 48px; width: 48px; margin-bottom: 22px; font-size: 21px; font-family: serif; color: var(--el-color-primary); background: var(--el-color-primary-light-9); border-radius: 50%; }
.result-placeholder h3 { font-size: 17px; font-weight: 500; margin: 0; }
.result-placeholder p { max-width: 390px; color: var(--el-text-color-secondary); font-size: 13px; line-height: 1.9; margin: 12px 0 26px; }
.placeholder-stages { display: flex; flex-wrap: wrap; justify-content: center; gap: 18px; color: var(--el-text-color-placeholder); font-size: 12px; }
.previous-result-note { padding: 10px 12px; margin-bottom: 16px; font-size: 12px; color: var(--el-text-color-secondary); background: var(--el-fill-color-light); border-radius: 4px; }
.result-metrics { display: grid; grid-template-columns: repeat(4, 1fr); gap: 12px; margin: 20px 0 16px; }
.result-metrics > div { padding: 14px 12px; border: 1px solid var(--el-border-color-lighter); border-radius: 4px; }
.result-metrics strong { display: block; font-size: 24px; font-weight: 600; font-variant-numeric: tabular-nums; }
.result-metrics span { display: block; margin-top: 6px; color: var(--el-text-color-secondary); font-size: 12px; }
.plan-summary { display: flex; gap: 7px; flex-wrap: wrap; align-items: center; margin-bottom: 16px; }
.plan-label { font-size: 12px; color: var(--el-text-color-secondary); margin-right: 3px; }
.plan-summary :deep(.el-tag) { height: auto; min-height: 24px; max-width: 100%; white-space: normal; line-height: 1.8; overflow-wrap: anywhere; }
.result-warnings { margin: 0 0 16px; padding: 12px 16px 12px 30px; background: var(--el-color-warning-light-9); color: var(--el-color-warning-dark-2); border-radius: 4px; font-size: 12px; line-height: 1.9; overflow-wrap: anywhere; }
.completed-steps { font-size: 12px; padding: 12px 0; border-top: 1px solid var(--el-border-color-extra-light); border-bottom: 1px solid var(--el-border-color-extra-light); margin-bottom: 16px; }
.completed-steps summary { cursor: pointer; color: var(--el-text-color-regular); }
.completed-steps ol { margin: 16px 0 0; padding-left: 22px; }
.completed-steps li { padding-left: 4px; margin-bottom: 12px; }
.completed-steps strong { font-weight: 500; }
.completed-steps p { margin: 5px 0 0; line-height: 1.8; color: var(--el-text-color-secondary); white-space: pre-wrap; overflow-wrap: anywhere; }
.distribution-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; margin-bottom: 20px; }
.distribution-panel { min-width: 0; padding: 15px; background: var(--el-fill-color-extra-light); border-radius: 4px; }
.distribution-panel h4 { margin: 0 0 14px; font-size: 13px; font-weight: 500; }
.distribution-panel h4 > span { float: right; color: var(--el-text-color-placeholder); font-size: 11px; font-weight: 400; }
.distribution-row { display: flex; gap: 8px; align-items: center; margin: 10px 0; font-size: 11px; }
.distribution-label { width: 44%; max-width: 150px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; color: var(--el-text-color-regular); }
.distribution-track { flex: 1; height: 5px; background: var(--el-fill-color); border-radius: 3px; }
.distribution-track > span { display: block; height: 100%; background: var(--el-color-primary-light-3); border-radius: inherit; }
.distribution-count { min-width: 30px; text-align: right; font-variant-numeric: tabular-nums; color: var(--el-text-color-secondary); }
.distribution-empty { font-size: 12px; color: var(--el-text-color-secondary); }
.result-list-toolbar { display: flex; justify-content: space-between; margin-bottom: 14px; }
.expert-card-list { display: grid; gap: 12px; }
.agent-pagination { justify-content: center; margin-top: 20px; }
@keyframes agent-pulse { 50% { opacity: 0.3; } }
@media (prefers-reduced-motion: reduce) { .progress-dot { animation: none; } }
@media (max-width: 1100px) {
  .agent-layout { grid-template-columns: minmax(260px, 310px) minmax(0, 1fr); gap: 16px; }
  .conversation-panel { padding: 16px; }
  .distribution-grid { grid-template-columns: 1fr; }
  .result-heading { flex-wrap: wrap; }
}
@media (max-width: 800px) {
  .agent-layout { grid-template-columns: 1fr; }
  .agent-heading { align-items: flex-start; }
  .agent-heading h2 { font-size: 18px; }
  .conversation-history { max-height: 320px; }
  .result-placeholder { min-height: 230px; }
}
@media (max-width: 480px) {
  .result-metrics { grid-template-columns: repeat(2, 1fr); }
  .agent-heading { gap: 10px; }
  .service-notice { align-items: flex-start; }
}
</style>
