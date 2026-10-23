<template>
  <li class="card">
    <div class="avatar">{{ (item.name || '?').slice(0, 1) }}</div>
    <div class="card-body">
      <div class="card-top">
        <span
          v-if="rank"
          class="rank-num"
        >第 {{ rank }} 名</span>
        <a
          class="card-name"
          @click="viewDetail"
        >{{ item.name }}</a>
        <span
          v-if="item.techTitle"
          class="badge-title"
        >{{ item.techTitle }}</span>
      </div>
      <div class="card-meta">
        <span>{{ item.unitName || '—' }}</span>
        <span
          v-if="item.department"
          class="sep"
        >·</span>
        <span v-if="item.department">{{ item.department }}</span>
      </div>
      <div
        v-if="item.disciplineL1"
        class="tag-row"
      >
        <el-tag
          size="small"
          effect="plain"
        >{{ item.disciplineL1 }}</el-tag>
      </div>
      <div
        v-if="item.researchDirections"
        class="snippet"
        v-html="highlighted"
      />
      <div
        v-if="item.matchReason"
        class="match-reason"
      >命中：{{ item.matchReason }}</div>
      <div class="card-bottom">
        <div class="score-chips">
          <span
            class="score-chip"
            title="职称分"
          ><b>职称</b>{{ fmt(item.titleScore) }}</span>
          <span
            v-if="scoreColumnAvailability.hasAchievement"
            class="score-chip"
            title="成果分"
          ><b>成果</b>{{ fmt(item.achievementScore) }}</span>
          <span
            v-if="scoreColumnAvailability.hasInfluence"
            class="score-chip"
            title="决策影响分"
          ><b>影响</b>{{ fmt(item.influenceScore) }}</span>
          <span
            v-if="scoreColumnAvailability.hasSocial"
            class="score-chip"
            title="社会贡献分"
          ><b>贡献</b>{{ fmt(item.socialScore) }}</span>
          <span
            v-if="hasKeyword"
            class="score-chip relevance"
            title="这条结果跟检索词的匹配程度——综合得分就是这个百分比乘出来的，匹配度越低，职称/成果分再高也会被打相应折扣"
          ><b>匹配度</b>{{ relevancePercent }}%</span>
          <span class="score-chip total"><b>{{ scoreLabel }}</b>{{ item.realtimeScore.toFixed(2) }}</span>
        </div>
        <a
          class="link-detail"
          @click="viewDetail"
        >查看详情 →</a>
      </div>
    </div>
  </li>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'

const props = defineProps({
  item: {
    type: Object,
    required: true
  },
  scoreColumnAvailability: {
    type: Object,
    default: () => ({ hasAchievement: true, hasInfluence: true, hasSocial: true })
  },
  // 上一次真正提交检索时的关键词原文，用来在研究方向摘要里高亮命中片段；浏览模式（没填关键词）
  // 传空字符串，不高亮
  keyword: {
    type: String,
    default: ''
  },
  // 综合实力/综合得分这一项的标签文字，浏览模式和关键词模式文案不同，由父组件算一次统一传下来，
  // 不在每张卡片里各自判断一遍 hasKeyword
  scoreLabel: {
    type: String,
    default: '综合实力'
  },
  // 上一次真正提交的是不是关键词类检索（简单检索或高级检索都算）——决定要不要显示"匹配度"这个
  // chip。不能直接拿 keyword 字符串是否非空来判断：高级检索如果唯一一行是"不含"，lastKeyword
  // 会是空字符串（不含的词不该被当成高亮词），但确实是一次关键词检索，浏览模式（真的没检索）
  // 才应该不显示——所以专门要父组件传一个显式的布尔值，跟 scoreLabel 是同一个 hasKeyword 来源，
  // 不在每张卡片里各自重新判断一遍
  hasKeyword: {
    type: Boolean,
    default: false
  },
  // 全局排名（跨页连续），不传就不展示
  rank: {
    type: Number,
    default: null
  }
})

const router = useRouter()
const viewDetail = () => {
  router.push({ name: 'expertProfileDetail', params: { id: props.item.ID }})
}

const fmt = (v) => Math.round((v || 0) * 10) / 10
// 匹配度：relevance 是 0~1 的浮点数（语义相似度/字面命中的保底值），页面上用百分比更直观，
// 跟综合得分卡片上"匹配度 X% × 基础分 = 综合得分"这个因果关系对得上
const relevancePercent = computed(() => Math.round((props.item.relevance || 0) * 100))

const escapeHtml = (s) =>
  s.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;')

// 把研究方向摘要里命中当前检索词的片段包一层 <mark>——纯展示层的辅助高亮，不是检索算法本身
// （真正命中依据看 matchReason）。先整体转义再在转义后的文本上做替换，保证插进 v-html 的内容
// 不会因为专家档案里的自由文本夹带 < > & 之类字符而产生非预期标签
const highlighted = computed(() => {
  const text = props.item.researchDirections || ''
  const escapedText = escapeHtml(text)
  const terms = (props.keyword || '')
    .split(/[,，、/;；\s]+/)
    .map((t) => t.trim())
    .filter(Boolean)
  if (!terms.length) return escapedText
  const escapedTerms = terms
    .map((t) => escapeHtml(t))
    .sort((a, b) => b.length - a.length) // 长词优先，避免短词先命中把长词从中间切开
    .map((t) => t.replace(/[.*+?^${}()|[\]\\]/g, '\\$&'))
  const re = new RegExp('(' + escapedTerms.join('|') + ')', 'gi')
  return escapedText.replace(re, '<mark>$1</mark>')
})
</script>

<style scoped>
/* 这几个色值跟登录页（src/view/login/index.vue）用的是同一套：$navy/$red/$gold——登录页
   已经把这套配色定成了这个系统的品牌识别，检索页作为用户打开频率最高的页面之一，理应接上，
   不能自己一套默认 Element Plus 蓝、跟登录页看着不像一个系统。只在这三个新写的分面检索组件里
   局部用，不动全局 el-color-primary，不影响别的页面 */
.card {
  --brand-navy: #16407a;
  --brand-navy-dark: #0c2a54;
  --brand-red: #b5272d;
  --brand-gold: #d4a017;
  list-style: none;
  background: var(--el-bg-color);
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 3px;
  padding: 16px 18px;
  display: flex;
  gap: 14px;
  transition: border-color 120ms ease, box-shadow 120ms ease;
}
.card:hover {
  border-color: var(--brand-red);
  box-shadow: 0 2px 12px rgba(181, 39, 45, 0.1);
}

.avatar {
  width: 42px;
  height: 42px;
  border-radius: 50%;
  flex-shrink: 0;
  display: grid;
  place-items: center;
  font-family: "Songti SC", "Noto Serif SC", "SimSun", serif;
  font-size: 17px;
  font-weight: 700;
  color: #fff;
  background: linear-gradient(155deg, var(--brand-navy), var(--brand-red));
}

.card-body {
  flex: 1;
  min-width: 0;
}
.card-top {
  display: flex;
  align-items: baseline;
  gap: 8px;
  flex-wrap: wrap;
  margin-bottom: 2px;
}
.rank-num {
  font-size: 11px;
  color: var(--el-text-color-placeholder);
  font-variant-numeric: tabular-nums;
}
.card-name {
  font-family: "Songti SC", "Noto Serif SC", "SimSun", serif;
  font-size: 17px;
  font-weight: 700;
  color: var(--el-text-color-primary);
  cursor: pointer;
  text-decoration: none;
}
.card-name:hover {
  color: var(--brand-red);
}
.badge-title {
  font-size: 11.5px;
  background: rgba(212, 160, 23, 0.14);
  color: var(--brand-gold);
  padding: 1px 8px;
  border-radius: 999px;
  font-weight: 600;
}

.card-meta {
  font-size: 12.5px;
  color: var(--el-text-color-secondary);
  margin-bottom: 8px;
}
.card-meta .sep {
  margin: 0 6px;
  color: var(--el-text-color-placeholder);
}

.tag-row {
  margin-bottom: 8px;
}

.snippet {
  font-size: 13px;
  color: var(--el-text-color-regular);
  line-height: 1.7;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  margin-bottom: 6px;
}
.snippet :deep(mark) {
  background: rgba(212, 160, 23, 0.28);
  color: var(--el-text-color-primary);
  border-radius: 2px;
  padding: 0 1px;
}
.match-reason {
  font-size: 12px;
  color: var(--brand-navy);
  margin-bottom: 8px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.card-bottom {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  flex-wrap: wrap;
}
.score-chips {
  display: flex;
  gap: 6px;
  align-items: center;
  flex-wrap: wrap;
}
.score-chip {
  font-size: 11.5px;
  color: var(--el-text-color-regular);
  background: var(--el-fill-color-light);
  border-radius: 3px;
  padding: 2px 8px;
  font-variant-numeric: tabular-nums;
}
.score-chip b {
  font-weight: 400;
  color: var(--el-text-color-secondary);
  margin-right: 3px;
}
.score-chip.total {
  background: var(--brand-navy);
  color: rgba(255, 255, 255, 0.85);
}
.score-chip.total b {
  /* 总分数字用金色跟其它分项区分出"这是最重要的那个数"，跟设计原型
     （expert-search-cnki-prototype.html 的 .score-chip.total b）一致 */
  color: var(--brand-gold);
  font-size: 13px;
}
.score-chip.relevance {
  /* 用红色边框跟职称/成果/影响/贡献那几个纯灰底分项区分开——匹配度不是"档案本身的分"，
     是"这次检索特有的、解释排名由来的数"，边框描红呼应检索页整体"红=检索相关的强调色"这个
     颜色分工（跟排序 tab 选中态/已选条件标签同一套逻辑） */
  background: rgba(181, 39, 45, 0.06);
  border: 1px solid rgba(181, 39, 45, 0.25);
  padding: 1px 7px;
}
.score-chip.relevance b {
  color: var(--brand-red);
}

.link-detail {
  font-size: 13px;
  font-weight: 600;
  color: var(--brand-red);
  text-decoration: none;
  cursor: pointer;
  white-space: nowrap;
}
.link-detail:hover {
  text-decoration: underline;
}
</style>
