<template>
  <q-page class="kanban-page">
    <div class="kanban-board">
      <!-- Column Headers -->
      <div class="board-header">
        <div class="column-header">
          <Icon icon="mdi:inbox-multiple" class="column-icon column-icon-backlog" />
          <h3 class="column-title">Backlog</h3>
        </div>
        <div class="column-header">
          <Icon icon="mdi:progress-clock" class="column-icon column-icon-progress" />
          <h3 class="column-title">In Progress</h3>
        </div>
        <div class="column-header">
          <Icon icon="mdi:check-circle" class="column-icon column-icon-done" />
          <h3 class="column-title">Done</h3>
        </div>
      </div>

      <!-- Lanes (Rows) -->
      <div v-for="lane in lanes" :key="lane.id" class="kanban-lane">
        <div class="lane-title">
          <span class="lane-name">{{ lane.title }}</span>
          <span class="lane-count">{{ lane.todos.length }}</span>
        </div>

        <div class="lane-columns">
          <!-- Backlog Column -->
          <div
            class="lane-column"
            :class="{ 'drop-target': isDropTarget(lane.id, 'backlog') }"
            @dragover="handleDragOver(lane.id, 'backlog', $event)"
            @drop="handleDrop(lane, 'backlog', $event)"
          >
            <div
              v-for="todo in getTodosByStatus(lane.todos, 'backlog')"
              :key="todo.id"
              class="kanban-card"
              draggable="true"
              @dragstart="handleDragStart(todo, lane, $event)"
              @dragend="handleDragEnd"
              @click="openEditDialog(todo)"
            >
              <div class="card-content">
                <p class="card-text">{{ todo.todo }}</p>
              </div>
            </div>
          </div>

          <!-- In Progress Column -->
          <div
            class="lane-column"
            :class="{ 'drop-target': isDropTarget(lane.id, 'in-progress') }"
            @dragover="handleDragOver(lane.id, 'in-progress', $event)"
            @drop="handleDrop(lane, 'in-progress', $event)"
          >
            <div
              v-for="todo in getTodosByStatus(lane.todos, 'in-progress')"
              :key="todo.id"
              class="kanban-card"
              draggable="true"
              @dragstart="handleDragStart(todo, lane, $event)"
              @dragend="handleDragEnd"
              @click="openEditDialog(todo)"
            >
              <div class="card-content">
                <p class="card-text">{{ todo.todo }}</p>
              </div>
            </div>
          </div>

          <!-- Done Column -->
          <div
            class="lane-column"
            :class="{ 'drop-target': isDropTarget(lane.id, 'done') }"
            @dragover="handleDragOver(lane.id, 'done', $event)"
            @drop="handleDrop(lane, 'done', $event)"
          >
            <div
              v-for="todo in getTodosByStatus(lane.todos, 'done')"
              :key="todo.id"
              class="kanban-card"
              draggable="true"
              @dragstart="handleDragStart(todo, lane, $event)"
              @dragend="handleDragEnd"
              @click="openEditDialog(todo)"
            >
              <div class="card-content">
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
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { Icon } from '@iconify/vue'
import type { Todo } from 'components/models'
import axios from 'axios'

interface Lane {
  id: string
  title: string
  todos: Todo[]
}

const todoToday = ref<Todo[]>([])
const todoTinkering = ref<Todo[]>([])

const lanes = computed<Lane[]>(() => {
  // Filter today todos to exclude @tinkering tasks with projects
  const filteredTodayTodos = todoToday.value.filter((todo) => {
    const hasTinkering = todo.context === '@tinkering'
    const hasProject = !!todo.project
    // Exclude if both @tinkering and project exist
    return !(hasTinkering && hasProject)
  })

  const laneList: Lane[] = [{ id: 'today', title: 'Today', todos: filteredTodayTodos }]

  // Group tinkering todos by project
  const projectGroups = new Map<string, Todo[]>()
  todoTinkering.value.forEach((todo) => {
    const project = todo.project || 'No Project'
    if (!projectGroups.has(project)) {
      projectGroups.set(project, [])
    }
    projectGroups.get(project)!.push(todo)
  })

  // Create a lane for each project
  projectGroups.forEach((todos, project) => {
    laneList.push({
      id: `project-${project}`,
      title: project.startsWith('+') ? project.slice(1) : project,
      todos,
    })
  })

  return laneList
})

const getTodosByStatus = (todos: Todo[], status: string): Todo[] => {
  return todos.filter((todo) => todo.status === status)
}

interface DragData {
  todo: Todo
  sourceLane: Lane
}

let dragData: DragData | null = null
const dropTarget = ref<string | null>(null)

const handleDragStart = (todo: Todo, lane: Lane, event: DragEvent) => {
  dragData = { todo, sourceLane: lane }
  if (event.dataTransfer) {
    event.dataTransfer.effectAllowed = 'move'
  }
}

const handleDragOver = (laneId: string, status: string, event: DragEvent) => {
  event.preventDefault()
  dropTarget.value = `${laneId}-${status}`
}

const handleDragEnd = () => {
  dropTarget.value = null
}

const handleDrop = async (targetLane: Lane, targetStatus: string, event: DragEvent) => {
  event.preventDefault()
  dropTarget.value = null

  if (!dragData) return

  const { todo, sourceLane } = dragData

  // Check if anything changed
  const sourceProject = sourceLane.id === 'today' ? '' : sourceLane.title
  const targetProject = targetLane.id === 'today' ? '' : targetLane.title

  if (todo.status === targetStatus && sourceProject === targetProject) {
    dragData = null
    return
  }

  // Determine context: add @tinkering when moving from Today to a project lane
  let context = todo.context || ''
  if (sourceLane.id === 'today' && targetLane.id !== 'today') {
    // Moving from Today to a project lane - add @tinkering if not already present
    if (!context.includes('tinkering')) {
      context = '@tinkering'
    }
  }

  try {
    // Update via API
    await axios.put('/api/todo/update', {
      id: todo.id,
      project: targetProject,
      status: targetStatus,
      context: context,
    })

    // Refresh data
    await fetchTodoToday()
    await fetchTodoTinkering()
  } catch (error) {
    console.error('Failed to update todo:', error)
  }

  dragData = null
}

const isDropTarget = (laneId: string, status: string): boolean => {
  return dropTarget.value === `${laneId}-${status}`
}

const fetchTodoToday = async () => {
  const response = await axios.get('/api/todo/today')
  todoToday.value = response.data as Todo[]
}

const fetchTodoTinkering = async () => {
  const response = await axios.get('/api/todo/tinkering')
  todoTinkering.value = response.data as Todo[]
}

// Edit dialog state
const showEditDialog = ref(false)
const editingTodo = ref<Todo | null>(null)
const editedText = ref('')
const highlightedText = ref('')

const openEditDialog = (todo: Todo) => {
  editingTodo.value = todo
  // Reconstruct full todo text with tags in order: @context +project content =status
  const parts = []
  if (todo.context) parts.push(todo.context)
  if (todo.project) parts.push(todo.project)
  parts.push(todo.todo)
  // Don't append =backlog since it's the default state
  if (todo.status && todo.status !== 'backlog') parts.push(`=${todo.status}`)
  editedText.value = parts.join(' ')
  updateHighlight()
  showEditDialog.value = true
}

const updateHighlight = () => {
  let text = editedText.value
  // Escape HTML
  text = text.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;')

  // Highlight @context tags
  text = text.replace(/(@\w+)/g, '<span class="highlight-context">$1</span>')

  // Highlight +project tags
  text = text.replace(/(\+\w+)/g, '<span class="highlight-project">$1</span>')

  // Highlight =status tags
  text = text.replace(/(=[\w-]+)/g, '<span class="highlight-status">$1</span>')

  // Highlight due: tags
  text = text.replace(/(due:\d{4}-\d{2}-\d{2})/g, '<span class="highlight-due">$1</span>')

  highlightedText.value = text
}

const saveTodo = async () => {
  if (!editingTodo.value) return

  try {
    // Remove =backlog since it's the default state
    let textToSave = editedText.value.replace(/=backlog\b/g, '').trim()
    // Clean up multiple spaces
    textToSave = textToSave.replace(/\s+/g, ' ')

    await axios.put('/api/todo/update-text', {
      id: editingTodo.value.id,
      text: textToSave,
    })

    // Refresh data
    await fetchTodoToday()
    await fetchTodoTinkering()

    showEditDialog.value = false
  } catch (error) {
    console.error('Failed to save todo:', error)
  }
}

let updateInterval: number | null = null
const UPDATE_INTERVAL_MS = 10000

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
  min-height: 100vh;
  padding: 24px;
}

.kanban-board {
  display: flex;
  flex-direction: column;
  gap: 0;
  background: #ffffff;
  border-radius: 8px;
  border: 1px solid #e1e4e8;
  overflow-x: auto;
}

.board-header {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
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
  grid-template-columns: repeat(3, 1fr);
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
  overflow-y: auto;
  max-height: 400px;
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
</style>
