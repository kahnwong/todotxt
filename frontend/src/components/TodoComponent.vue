<template>
  <div class="q-pa-md" style="max-width: 400px">
    <q-list bordered separator>
      <q-item no-ripple v-for="todo in todos" :key="todo.id">
        <q-item-section>
          <div>
            <span v-if="todo.priority" :class="['todo-priority', priorityClass(todo.priority)]">
              ({{ todo.priority }})&nbsp;
            </span>
            <span class="todo-context">{{ todo.context }}&nbsp;</span>
            <span class="todo-project">{{ todo.project }}&nbsp;</span>
            <span>{{ todo.todo }}</span>
          </div>
        </q-item-section>
      </q-item>
    </q-list>
  </div>
</template>

<script setup lang="ts">
import type { Todo } from 'components/models'

export interface TodoProps {
  todos: Todo[]
}

withDefaults(defineProps<TodoProps>(), {})

const priorityClass = (priority: string) => {
  switch (priority.toUpperCase()) {
    case 'A':
      return 'priority-a'
    case 'B':
      return 'priority-b'
    case 'C':
      return 'priority-c'
    default:
      return 'priority-other'
  }
}
</script>

<style>
.todo-priority {
  font-weight: 700;
}

.priority-a {
  color: #cf222e;
}

.priority-b {
  color: #ff8c00;
}

.priority-c {
  color: #9a6700;
}

.priority-other {
  color: #825d3a;
}

.todo-context {
  color: #4fa1ac;
}
.todo-project {
  color: #c8aa6f;
}
</style>
