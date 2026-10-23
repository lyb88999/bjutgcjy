export const safeSourceUrl = (value) => {
  if (typeof value !== 'string' || !value.trim()) return ''
  try {
    const url = new URL(value)
    if (!['http:', 'https:'].includes(url.protocol) || url.username || url.password) return ''
    return url.href
  } catch {
    return ''
  }
}

const text = (value) => typeof value === 'string' ? value : ''
const list = (value) => Array.isArray(value) ? value : []

const normalizeExpert = (item) => ({
  id: Number(item.id),
  name: text(item.name),
  unitName: text(item.unitName),
  department: text(item.department),
  techTitle: text(item.techTitle),
  disciplineL1: text(item.disciplineL1),
  researchDirections: text(item.researchDirections),
  profileLink: safeSourceUrl(item.profileLink),
  matchReason: text(item.matchReason),
  score: typeof item.score === 'number' && Number.isFinite(item.score) ? item.score : null,
  reason: text(item.reason),
  evidence: list(item.evidence).filter((entry) => entry && typeof entry === 'object').map((entry) => ({
    kind: text(entry.kind),
    title: text(entry.title),
    sourceId: Number.isSafeInteger(Number(entry.sourceId)) && Number(entry.sourceId) > 0 ? Number(entry.sourceId) : null,
    url: safeSourceUrl(entry.url)
  }))
})

export const normalizeAgentResult = (data) => {
  if (!data || typeof data !== 'object' || typeof data.sessionId !== 'string' || !data.sessionId || typeof data.reply !== 'string') {
    throw new Error('服务返回的结果不完整，请重试。')
  }
  const experts = (items) => list(items)
    .filter((item) => item && Number.isSafeInteger(Number(item.id)) && Number(item.id) > 0)
    .map(normalizeExpert)
  const distribution = (items) => list(items)
    .filter((item) => item && typeof item.label === 'string' && Number.isFinite(item.count) && item.count >= 0)
    .map((item) => ({ label: item.label, count: item.count }))
  const candidates = experts(data.candidates)
  const recommendations = experts(data.recommendations)
  const nonNegativeCount = (value, fallback) => Number.isInteger(value) && value >= 0 ? value : fallback
  return {
    sessionId: data.sessionId,
    reply: data.reply,
    plan: data.plan && typeof data.plan === 'object' ? data.plan : {},
    candidates,
    recommendations,
    statistics: {
      candidateCount: nonNegativeCount(data.statistics?.candidateCount, candidates.length),
      selectedCount: nonNegativeCount(data.statistics?.selectedCount, recommendations.length),
      unitCount: nonNegativeCount(data.statistics?.unitCount, new Set(recommendations.map((item) => item.unitName).filter(Boolean)).size),
      disciplineCount: nonNegativeCount(data.statistics?.disciplineCount, new Set(recommendations.map((item) => item.disciplineL1).filter(Boolean)).size),
      units: distribution(data.statistics?.units),
      disciplines: distribution(data.statistics?.disciplines)
    },
    steps: list(data.steps).filter((step) => step && typeof step.label === 'string').map((step) => ({
      label: step.label,
      detail: text(step.detail)
    })),
    warnings: list(data.warnings).filter((warning) => typeof warning === 'string')
  }
}

// 每个单元格都加引号；危险前缀加单引号，防止表格软件把外部文本解释为公式。
export const csvCell = (value) => {
  let cell = String(value ?? '').replace(/\u0000/g, '')
  if (/^[\s\u0000-\u0020\u007f-\u009f\uFEFF]*[=+\-@]/.test(cell) || /^[\t\r\n]/.test(cell)) cell = `'${cell}`
  return `"${cell.replace(/"/g, '""')}"`
}

export const recommendationsCsv = (items) => {
  const rows = [['序号', '专家ID', '姓名', '所在单位', '院系/部门', '职称', '一级学科', '排序得分', '推荐理由', '命中依据', '研究方向', '证据', '来源链接']]
  items.forEach((item, index) => {
    rows.push([
      index + 1, item.id, item.name, item.unitName, item.department, item.techTitle,
      item.disciplineL1, item.score ?? '', item.reason, item.matchReason, item.researchDirections,
      item.evidence.map((entry) => [entry.kind, entry.title, safeSourceUrl(entry.url)].filter(Boolean).join('：')).join('\n'),
      safeSourceUrl(item.profileLink)
    ])
  })
  return '\uFEFF' + rows.map((row) => row.map(csvCell).join(',')).join('\r\n')
}
