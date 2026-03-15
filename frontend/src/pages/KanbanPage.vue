<template>
  <q-page class="kanban-page">
    <div class="kanban-board">
      <!-- Column Headers -->
      <div class="board-header">
        <div class="lane-label">Lane</div>
        <div class="column-header">
          <h3 class="column-title">Backlog</h3>
        </div>
        <div class="column-header">
          <h3 class="column-title">In Progress</h3>
        </div>
        <div class="column-header">
          <h3 class="column-title">Done</h3>
        </div>
      </div>

      <!-- Lanes (Rows) -->
      <div v-for="lane in lanes" :key="lane.id" class="kanban-lane">
        <div class="lane-title">
          <span class="lane-name">{{ lane.title }}</span>
          <span class="lane-count">{{ lane.todos.length }}</span>
        </div>

        <!-- Backlog Column -->
        <div class="lane-column">
          <div v-for="todo in lane.todos" :key="todo.id" class="kanban-card">
            <div class="card-content">
              <span class="card-context" v-if="todo.context">{{ todo.context }}</span>
              <span class="card-project" v-if="todo.project && lane.id !== 'today'">{{
                todo.project
              }}</span>
              <p class="card-text">{{ todo.todo }}</p>
            </div>
          </div>
        </div>

        <!-- In Progress Column -->
        <div class="lane-column">
          <!-- Empty for now -->
        </div>

        <!-- Done Column -->
        <div class="lane-column">
          <!-- Empty for now -->
        </div>
      </div>
    </div>
  </q-page>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed } from 'vue'
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
  const laneList: Lane[] = [{ id: 'today', title: 'Today', todos: todoToday.value }]

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
      title: project,
      todos,
    })
  })

  return laneList
})

const fetchTodoToday = async () => {
  const response = await axios.get('/api/todo/today')
  todoToday.value = response.data as Todo[]
}

const fetchTodoTinkering = async () => {
  const response = await axios.get('/api/todo/tinkering')
  todoTinkering.value = response.data as Todo[]
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
  grid-template-columns: 200px repeat(3, 1fr);
  border-bottom: 2px solid #e1e4e8;
  background: #f6f8fa;
  position: sticky;
  top: 0;
  z-index: 10;
}

.lane-label {
  padding: 16px;
  font-size: 14px;
  font-weight: 600;
  color: #24292f;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  border-right: 1px solid #e1e4e8;
}

.column-header {
  padding: 16px;
  text-align: center;
  border-right: 1px solid #e1e4e8;
}

.column-header:last-child {
  border-right: none;
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
  display: grid;
  grid-template-columns: 200px repeat(3, 1fr);
  border-bottom: 1px solid #e1e4e8;
  min-height: 120px;
}

.kanban-lane:last-child {
  border-bottom: none;
}

.lane-title {
  padding: 16px;
  border-right: 1px solid #e1e4e8;
  background: #fafbfc;
  display: flex;
  flex-direction: column;
  gap: 8px;
  align-items: flex-start;
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
}

.lane-column:last-child {
  border-right: none;
}

.kanban-card {
  background: #ffffff;
  border: 1px solid #d0d7de;
  border-radius: 6px;
  padding: 12px;
  margin-bottom: 8px;
  cursor: pointer;
  transition: all 0.15s ease;
}

.kanban-card:hover {
  border-color: #8c959f;
  box-shadow: 0 3px 8px rgba(0, 0, 0, 0.08);
  transform: translateY(-1px);
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
</style>
