<template>
  <div class="q-pa-md" style="max-width: 400px">
    <q-list bordered separator>
      <q-item no-ripple v-for="todo in todos" :key="todo.id">
        <q-item-section>
          <div>
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
</script>

<style>
.todo-context {
  color: #4fa1ac;
}
.todo-project {
  color: #c8aa6f;
}
</style>
