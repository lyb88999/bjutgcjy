<template>
  <article class="agent-expert-card">
    <div class="expert-card-heading">
      <div class="expert-card-identity">
        <span class="expert-rank">{{ index + 1 }}</span>
        <router-link
          :to="{ name: 'expertProfileDetail', params: { id: item.id } }"
          class="expert-name"
        >
          {{ item.name || '未填写姓名' }}
        </router-link>
        <el-tag
          v-if="item.techTitle"
          size="small"
          effect="plain"
        >{{ item.techTitle }}</el-tag>
      </div>
      <span
        v-if="item.score !== null"
        class="expert-score"
      >排序得分 {{ item.score.toFixed(2) }}</span>
    </div>
    <p class="expert-meta">{{ [item.unitName, item.department, item.disciplineL1].filter(Boolean).join(' · ') || '单位与学科信息待补充' }}</p>
    <p
      v-if="item.reason"
      class="expert-reason"
    >{{ item.reason }}</p>
    <p
      v-if="item.matchReason"
      class="expert-match"
    ><span>命中依据</span>{{ item.matchReason }}</p>
    <p
      v-if="item.researchDirections"
      class="expert-directions"
    >{{ item.researchDirections }}</p>
    <details
      v-if="item.evidence.length"
      class="expert-evidence"
    >
      <summary>核对依据（{{ item.evidence.length }} 条）</summary>
      <ul>
        <li
          v-for="(entry, evidenceIndex) in item.evidence"
          :key="evidenceIndex"
        >
          <span class="evidence-kind">{{ evidenceKind(entry.kind) }}</span>
          <span>{{ entry.title || '记录未提供标题' }}</span>
          <span
            v-if="entry.sourceId"
            class="evidence-id"
          >#{{ entry.sourceId }}</span>
          <a
            v-if="entry.url"
            :href="entry.url"
            target="_blank"
            rel="noopener noreferrer"
          >来源 ↗</a>
        </li>
      </ul>
    </details>
    <div class="expert-card-links">
      <router-link :to="{ name: 'expertProfileDetail', params: { id: item.id } }">查看专家档案 →</router-link>
      <a
        v-if="item.profileLink"
        :href="item.profileLink"
        target="_blank"
        rel="noopener noreferrer"
      >档案来源 ↗</a>
    </div>
  </article>
</template>

<script setup>
defineProps({
  item: { type: Object, required: true },
  index: { type: Number, default: 0 }
})

const evidenceKind = (kind) => ({
  profile: '专家档案',
  research: '研究方向',
  research_direction: '研究方向',
  researchDirection: '研究方向',
  achievement: '研究成果',
  adoption: '决策影响',
  position: '学术兼职',
  tag: '专家标签'
}[kind] || kind || '库内记录')
</script>

<style scoped>
.agent-expert-card {
  padding: 18px;
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 5px;
  background: var(--el-bg-color);
}
.expert-card-heading,
.expert-card-identity,
.expert-card-links {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}
.expert-card-heading { justify-content: space-between; }
.expert-rank {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  display: inline-grid;
  place-items: center;
  background: var(--el-fill-color-light);
  color: var(--el-text-color-secondary);
  font-size: 12px;
}
.expert-name { font-size: 17px; font-weight: 600; }
.expert-name,
.expert-card-links a,
.expert-evidence a { color: var(--el-color-primary); text-decoration: none; }
a:hover { text-decoration: underline; }
.expert-score { font-size: 12px; color: var(--el-text-color-secondary); }
.expert-meta { margin: 9px 0 0; color: var(--el-text-color-secondary); font-size: 13px; line-height: 1.7; }
.expert-reason { margin: 12px 0 0; color: var(--el-text-color-primary); font-size: 14px; line-height: 1.8; white-space: pre-wrap; overflow-wrap: anywhere; }
.expert-match { margin: 10px 0 0; font-size: 12px; line-height: 1.8; color: var(--el-text-color-regular); overflow-wrap: anywhere; }
.expert-match > span { margin-right: 8px; color: var(--el-text-color-secondary); }
.expert-directions { margin: 8px 0 0; color: var(--el-text-color-secondary); font-size: 13px; line-height: 1.8; white-space: pre-wrap; overflow-wrap: anywhere; }
.expert-evidence { margin-top: 14px; font-size: 12px; color: var(--el-text-color-regular); }
.expert-evidence summary { cursor: pointer; color: var(--el-text-color-primary); }
.expert-evidence ul { margin: 10px 0 0; padding-left: 18px; }
.expert-evidence li { margin: 7px 0; line-height: 1.8; overflow-wrap: anywhere; }
.evidence-kind { display: inline-block; margin-right: 8px; color: var(--el-text-color-secondary); }
.evidence-id { margin-left: 7px; color: var(--el-text-color-placeholder); }
.expert-evidence a { margin-left: 8px; }
.expert-card-links { margin-top: 14px; padding-top: 12px; border-top: 1px solid var(--el-border-color-extra-light); font-size: 12px; }
@media (max-width: 640px) {
  .agent-expert-card { padding: 14px; }
  .expert-card-heading { align-items: flex-start; }
}
</style>
