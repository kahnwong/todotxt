<template>
  <q-page class="kanban-page">
    <div class="kanban-board">
      <!-- Column Headers -->
      <div class="board-header">
        <div v-for="column in COLUMN_CONFIG" :key="column.id" class="column-header">
          <Icon :icon="column.icon" class="column-icon" :class="column.iconClass" />
          <h3 class="column-title">{{ column.title }}</h3>
        </div>
      </div>

      <!-- Lanes (Rows) -->
      <div v-for="lane in lanes" :key="lane.id" class="kanban-lane">
        <div class="lane-title">
          <span class="lane-name">{{ lane.title }}</span>
          <span class="lane-count">{{ lane.todos.length }}</span>
        </div>

        <div class="lane-columns">
          <!-- Dynamic Columns -->
          <div
            v-for="column in COLUMN_CONFIG"
            :key="column.id"
            class="lane-column"
            :class="{ 'drop-target': isDropTarget(lane.id, column.id) }"
            @dragover="handleDragOver(lane.id, column.id, $event)"
            @drop="handleDrop(lane, column.id, $event)"
          >
            <div
              v-for="todo in getTodosByStatus(lane.todos, column.id)"
              :key="todo.id"
              class="kanban-card"
              draggable="true"
              @dragstart="handleDragStart(todo, lane, $event)"
              @dragend="handleDragEnd"
              @click="openEditDialog(todo)"
            >
              <div class="card-content">
                <div v-if="lane.id === 'today' && (todo.context || todo.project)" class="card-tags">
                  <span v-if="todo.context" class="card-context">{{ todo.context }}</span>
                  <span v-if="todo.project" class="card-project">{{ todo.project }}</span>
                </div>
                <p class="card-text">{{ todo.todo }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Edit Dialog -->
    <q-dialog v-model="showEditDialog">
      <q-card style="min-width: 500px">
        <q-card-section>
          <div class="text-h6">Edit Todo</div>
        </q-card-section>

        <q-card-section>
          <div class="editor-container">
            <div class="syntax-highlight" v-html="highlightedText"></div>
            <textarea
              v-model="editedText"
              class="editor-input"
              @input="updateHighlight"
              rows="5"
            ></textarea>
          </div>
        </q-card-section>

        <q-card-actions align="right">
          <q-btn flat label="Cancel" color="primary" v-close-popup />
          <q-btn flat label="Save" color="primary" @click="saveTodo" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { Icon } from '@iconify/vue'
import type { Todo } from 'components/models'
import axios from 'axios'
import { COLUMN_CONFIG, UPDATE_INTERVAL_MS } from '../constants/kanban'
import { useLanes } from '../composables/useLanes'
import { useDragDrop } from '../composables/useDragDrop'
import { useTodoEdit } from '../composables/useTodoEdit'

const todoToday = ref<Todo[]>([])
const todoTinkering = ref<Todo[]>([])

// Fetch functions
const fetchTodoToday = async () => {
  const response = await axios.get('/api/todo/today')
  todoToday.value = response.data as Todo[]
}

const fetchTodoTinkering = async () => {
  const response = await axios.get('/api/todo/tinkering')
  todoTinkering.value = response.data as Todo[]
}

// Use composables
const { lanes, getTodosByStatus } = useLanes(todoToday, todoTinkering)

const { handleDragStart, handleDragOver, handleDragEnd, handleDrop, isDropTarget } = useDragDrop(
  fetchTodoToday,
  fetchTodoTinkering,
)

const { showEditDialog, editedText, highlightedText, openEditDialog, updateHighlight, saveTodo } =
  useTodoEdit(fetchTodoToday, fetchTodoTinkering)

// Auto-refresh setup
let updateInterval: number | null = null

onMounted(() => {
  fetchTodoToday()
  fetchTodoTinkering()

  updateInterval = setInterval(() => {
    fetchTodoToday()
    fetchTodoTinkering()
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

/* Scrollbar styling */
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

/* Editor styles */
.editor-container {
  position: relative;
  font-family: 'Courier New', monospace;
  font-size: 14px;
  line-height: 1.5;
}

.syntax-highlight {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  padding: 8px;
  white-space: pre-wrap;
  word-wrap: break-word;
  color: transparent;
  pointer-events: none;
  border: 1px solid transparent;
}

.editor-input {
  position: relative;
  width: 100%;
  padding: 8px;
  font-family: 'Courier New', monospace;
  font-size: 14px;
  line-height: 1.5;
  background: transparent;
  border: 1px solid #d0d7de;
  border-radius: 4px;
  resize: vertical;
  color: #24292f;
}

.editor-input:focus {
  outline: none;
  border-color: #0969da;
}

/* Syntax highlighting colors */
.highlight-context {
  color: #4fa1ac;
  font-weight: 600;
}

.highlight-project {
  color: #c8aa6f;
  font-weight: 600;
}

.highlight-status {
  color: #8250df;
  font-weight: 600;
}

.highlight-due {
  color: #cf222e;
  font-weight: 600;
}

.highlight-date {
  color: #0969da;
  font-weight: 600;
}
</style>
