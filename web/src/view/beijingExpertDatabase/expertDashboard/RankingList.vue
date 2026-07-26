<template>
  <div
    v-if="!items || !items.length"
    class="ranking-empty"
  >暂无数据</div>
  <div
    v-else
    class="ranking-list"
  >
    <div
      v-for="(item, index) in items"
      :key="item.label + index"
      class="ranking-item"
    >
      <span
        class="rank-badge"
        :class="rankClass(index)"
      >{{ index + 1 }}</span>
      <span
        class="rank-label"
        :title="item.label"
      >{{ item.label }}</span>
      <div class="rank-bar-track">
        <div
          class="rank-bar-fill"
          :style="{ width: (item.count / max * 100) + '%' }"
        />
      </div>
      <span class="rank-count">{{ item.count }}</span>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

defineOptions({
  name: 'RankingList'
})

const props = defineProps({
  items: { type: Array, default: () => [] }
})

const max = computed(() => Math.max(...props.items.map(i => i.count), 1))
const rankClass = (i) => (i === 0 ? 'rank-gold' : i === 1 ? 'rank-silver' : i === 2 ? 'rank-bronze' : 'rank-plain')
</script>

<style scoped>
.ranking-list {
  display: flex;
  flex-direction: column;
  gap: 14px;
  padding: 4px 0;
}
.ranking-item {
  display: flex;
  align-items: center;
  gap: 10px;
}
.rank-badge {
  width: 20px;
  height: 20px;
  min-width: 20px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 700;
  color: #fff;
}
.rank-gold { background: linear-gradient(135deg, #ffd666, #fa8c16); }
.rank-silver { background: linear-gradient(135deg, #d9d9d9, #8c8c8c); }
.rank-bronze { background: linear-gradient(135deg, #ffc069, #ad6800); }
.rank-plain { background: var(--el-fill-color); color: var(--el-text-color-secondary); }
.rank-label {
  width: 150px;
  min-width: 150px;
  font-size: 13px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.rank-bar-track {
  flex: 1;
  height: 8px;
  background: var(--el-fill-color-light);
  border-radius: 4px;
  overflow: hidden;
}
.rank-bar-fill {
  height: 100%;
  border-radius: 4px;
  background: linear-gradient(90deg, #1677ff, #69b1ff);
}
.rank-count {
  min-width: 32px;
  text-align: right;
  font-size: 13px;
  font-weight: 600;
  color: var(--el-text-color-regular);
}
.ranking-empty {
  text-align: center;
  color: var(--el-text-color-placeholder);
  padding: 40px 0;
  font-size: 13px;
}
</style>
