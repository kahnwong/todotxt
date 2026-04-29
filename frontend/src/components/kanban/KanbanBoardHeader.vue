<template>
  <div class="board-header">
    <div v-for="column in columns" :key="column.id" class="column-header">
      <Icon :icon="column.icon" class="column-icon" :class="column.iconClass" />
      <h3 class="column-title">{{ column.title }}</h3>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue'

defineProps<{
  columns: ReadonlyArray<{
    id: string
    title: string
    icon: string
    iconClass: string
  }>
}>()
</script>

<style scoped>
.board-header {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  border-bottom: 2px solid #e1e4e8;
  background: #f6f8fa;
  position: sticky;
  top: 0;
  z-index: 10;
}

.column-header {
  height: 50px;
  padding: 0 16px;
  text-align: center;
  border-right: 1px solid #e1e4e8;
  display: flex;
  align-items: center;
  justify-content: center;
}

.column-header:last-child {
  border-right: none;
}

.column-icon {
  font-size: 18px;
  margin-right: 8px;
  color: #57606a;
}

.column-icon-backlog {
  color: #6c757d;
}

.column-icon-progress {
  color: #f59e0b;
}

.column-icon-stuck {
  color: #ef4444;
}

.column-icon-done {
  color: #10b981;
}

.column-title {
  margin: 0;
  font-size: 14px;
  font-weight: 600;
  color: #24292f;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}
</style>
