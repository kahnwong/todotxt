<template>
  <q-page class="kanban-page">
    <div class="kanban-board">
      <KanbanBoardHeader :columns="COLUMN_CONFIG" />

      <KanbanLane
        v-for="lane in lanes"
        :key="lane.id"
        :lane="lane"
        :columns="COLUMN_CONFIG"
        :get-todos-by-status="getTodosByStatus"
        :is-drop-target="isDropTarget"
        @dragstart="handleDragStart"
        @dragover="handleDragOver"
        @dragend="handleDragEnd"
        @drop="handleDrop"
        @edit="openEditDialog"
      />
    </div>

    <KanbanEditDialog
      v-model="showEditDialog"
      :edited-text="editedText"
      :highlighted-text="highlightedText"
      @update-text="handleEditedTextUpdate"
      @save="saveTodo"
    />
  </q-page>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import type { Todo } from 'components/models'
import axios from 'axios'
import KanbanBoardHeader from 'components/kanban/KanbanBoardHeader.vue'
import KanbanEditDialog from 'components/kanban/KanbanEditDialog.vue'
import KanbanLane from 'components/kanban/KanbanLane.vue'
import { COLUMN_CONFIG, UPDATE_INTERVAL_MS } from '../constants/kanban'
import { useLanes } from '../composables/useLanes'
import { useDragDrop } from '../composables/useDragDrop'
import { useTodoEdit } from '../composables/useTodoEdit'

const todoToday = ref<Todo[]>([])
const todoTinkering = ref<Todo[]>([])
const todoWork = ref<Todo[]>([])

// Fetch functions
const fetchTodoToday = async () => {
  const response = await axios.get('/api/todo/today')
  todoToday.value = response.data as Todo[]
}

const fetchTodoTinkering = async () => {
  const response = await axios.get('/api/todo/tinkering')
  todoTinkering.value = response.data as Todo[]
}

const fetchTodoWork = async () => {
  const response = await axios.get('/api/todo/work')
  todoWork.value = response.data as Todo[]
}

// Use composables
const { lanes, getTodosByStatus } = useLanes(todoToday, todoTinkering, todoWork)

const { handleDragStart, handleDragOver, handleDragEnd, handleDrop, isDropTarget } = useDragDrop(
  fetchTodoToday,
  fetchTodoWork,
  fetchTodoTinkering,
)

const { showEditDialog, editedText, highlightedText, openEditDialog, updateHighlight, saveTodo } =
  useTodoEdit(fetchTodoToday, fetchTodoTinkering, fetchTodoWork)

const handleEditedTextUpdate = (value: string) => {
  editedText.value = value
  updateHighlight()
}

// Auto-refresh setup
let updateInterval: number | null = null

onMounted(() => {
  fetchTodoToday()
  fetchTodoTinkering()
  fetchTodoWork()

  updateInterval = setInterval(() => {
    fetchTodoToday()
    fetchTodoTinkering()
    fetchTodoWork()
  }, UPDATE_INTERVAL_MS) as unknown as number
})

onUnmounted(() => {
  if (updateInterval) {
    clearInterval(updateInterval)
  }
})
</script>

<style scoped>
.kanban-page {
  background: #f6f8fa;
  height: 100vh;
  padding: 24px;
  display: flex;
  flex-direction: column;
}

.kanban-board {
  display: flex;
  flex-direction: column;
  gap: 0;
  background: #ffffff;
  border-radius: 8px;
  border: 1px solid #e1e4e8;
  overflow-x: auto;
  overflow-y: auto;
  flex: 1;
  max-height: calc(100vh - 48px);
}
</style>
