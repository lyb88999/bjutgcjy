// 专家档案审核状态的前端统一映射：文案、el-tag 颜色、筛选选项。
// 之前列表页/详情页/审核台各自维护一份 statusMap，文案改一处会漏另外两处，收拢到这里
export const STATUS_OPTIONS = [
  { value: 'draft', label: '草稿', tagType: 'info' },
  { value: 'pending_org_review', label: '待单位审核', tagType: 'warning' },
  { value: 'org_rejected', label: '单位已退回', tagType: 'danger' },
  { value: 'pending_city_review', label: '待市级审核', tagType: 'warning' },
  { value: 'city_rejected', label: '市级已退回', tagType: 'danger' },
  { value: 'published', label: '已发布', tagType: 'success' }
]

const byValue = Object.fromEntries(STATUS_OPTIONS.map((s) => [s.value, s]))

export const statusLabel = (value) => byValue[value]?.label || value
export const statusTagType = (value) => byValue[value]?.tagType || 'info'

// 档案编辑表单里的枚举选项。职称不在这里——它直接决定打分，选项以
// expert_title_level 权重字典为准（getDict 取），保证表单里能选的和打分认的完全一致
export const GENDER_OPTIONS = ['男', '女']
export const POLITICAL_STATUS_OPTIONS = [
  '中共党员',
  '中共预备党员',
  '共青团员',
  '民主党派',
  '无党派人士',
  '群众'
]
export const EDUCATION_OPTIONS = ['博士研究生', '硕士研究生', '大学本科', '大专及以下']
export const DEGREE_OPTIONS = ['博士', '硕士', '学士', '无']
