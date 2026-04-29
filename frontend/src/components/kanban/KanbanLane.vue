<template>
  <div class="kanban-lane">
    <div class="lane-title">
      <span class="lane-name" :class="{ 'lane-name-work': lane.id.startsWith('work-project-') }">
        {{ lane.title }}
      </span>
      <span class="lane-count">{{ lane.todos.length }}</span>
    </div>

    <div class="lane-columns">
      <div
        v-for="column in columns"
        :key="column.id"
        class="lane-column"
        :class="{ 'drop-target': isDropTarget(lane.id, column.id) }"
        @dragover="emit('dragover', lane.id, column.id, $event)"
        @drop="emit('drop', lane, column.id, $event)"
      >
        <KanbanCard
          v-for="todo in getTodosByStatus(lane.todos, column.id)"
          :key="todo.id"
          :todo="todo"
          :lane-id="lane.id"
          @dragstart="emit('dragstart', todo, lane, $event)"
          @dragend="emit('dragend')"
          @edit="emit('edit', todo)"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import KanbanCard from './KanbanCard.vue'
import type { Lane, Todo } from 'components/models'

defineProps<{
  lane: Lane
  columns: ReadonlyArray<{
    id: string
    title: string
    icon: string
    iconClass: string
  }>
  getTodosByStatus: (todos: Todo[], status: string) => Todo[]
  isDropTarget: (laneId: string, status: string) => boolean
}>()

const emit = defineEmits<{
  dragstart: [todo: Todo, lane: Lane, event: DragEvent]
  dragover: [laneId: string, status: string, event: DragEvent]
  dragend: []
  drop: [lane: Lane, status: string, event: DragEvent]
  edit: [todo: Todo]
}>()
</script>

<style scoped>
.kanban-lane {
  display: flex;
  flex-direction: column;
  border-bottom: 1px solid #e1e4e8;
}

.kanban-lane:last-child {
  border-bottom: none;
}

.lane-title {
  padding: 12px 16px;
  background: #e5e8e8;
  border-bottom: 1px solid #e1e4e8;
  display: flex;
  align-items: center;
  gap: 12px;
}

.lane-columns {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  min-height: 120px;
}

.lane-name {
  font-size: 14px;
  font-weight: 600;
  color: #24292f;
}

.lane-name-work {
  color: #14b8a6;
}

.lane-count {
  font-size: 12px;
  color: #57606a;
  background: #ffffff;
  padding: 2px 8px;
  border-radius: 12px;
  font-weight: 500;
  border: 1px solid #e1e4e8;
}

.lane-column {
  padding: 12px;
  border-right: 1px solid #e1e4e8;
  min-width: 0;
  overflow: hidden;
  transition:
    background-color 0.2s ease,
    border-color 0.2s ease;
}

.lane-column:last-child {
  border-right: none;
}

.lane-column.drop-target {
  background-color: #e8f5f7;
  border: 2px dashed #4fa1ac;
  border-right: 2px dashed #4fa1ac;
}

.lane-column::-webkit-scrollbar {
  width: 6px;
}

.lane-column::-webkit-scrollbar-track {
  background: transparent;
}

.lane-column::-webkit-scrollbar-thumb {
  background: #d0d7de;
  border-radius: 3px;
}

.lane-column::-webkit-scrollbar-thumb:hover {
  background: #8c959f;
}
</style>
