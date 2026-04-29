<template>
  <div
    class="kanban-card"
    draggable="true"
    @dragstart="emit('dragstart', $event)"
    @dragend="emit('dragend')"
    @click="emit('edit')"
  >
    <div class="card-content">
      <div v-if="showTags" class="card-tags">
        <span v-if="todo.context" class="card-context">{{ todo.context }}</span>
        <span v-if="todo.project" class="card-project">{{ todo.project }}</span>
      </div>
      <p class="card-text">{{ todo.todo }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Todo } from 'components/models'

const props = defineProps<{
  todo: Todo
  laneId: string
}>()

const emit = defineEmits<{
  dragstart: [event: DragEvent]
  dragend: []
  edit: []
}>()

const showTags = computed(() => {
  return props.laneId === 'today' && (props.todo.context || props.todo.project)
})
</script>

<style scoped>
.kanban-card {
  background: #ffffff;
  border: 1px solid #d0d7de;
  border-radius: 6px;
  padding: 12px;
  margin-bottom: 8px;
  cursor: grab;
  transition: all 0.15s ease;
}

.kanban-card:active {
  cursor: grabbing;
}

.kanban-card:hover {
  border-color: #8c959f;
  box-shadow: 0 3px 8px rgba(0, 0, 0, 0.08);
  transform: translateY(-1px);
}

.kanban-card[draggable='true'] {
  cursor: grab;
}

.kanban-card[draggable='true']:active {
  cursor: grabbing;
  opacity: 0.5;
}

.card-content {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.card-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-bottom: 4px;
}

.card-context {
  font-size: 12px;
  font-weight: 500;
  color: #4fa1ac;
  background: #e8f5f7;
  padding: 2px 6px;
  border-radius: 3px;
  display: inline-block;
  width: fit-content;
}

.card-project {
  font-size: 12px;
  font-weight: 500;
  color: #c8aa6f;
  background: #fef9f0;
  padding: 2px 6px;
  border-radius: 3px;
  display: inline-block;
  width: fit-content;
}

.card-text {
  margin: 0;
  font-size: 14px;
  line-height: 1.5;
  color: #24292f;
  word-wrap: break-word;
  overflow-wrap: break-word;
  word-break: break-word;
  max-width: 100%;
}
</style>
