<template>
  <q-page class="row justify-evenly">
    <div class="q-pa-md example-row-equal-width">
      <div class="row">
        <todo-component title="Today" active :todos="todoToday"></todo-component>
      </div>

      <div class="row">
        <todo-component title="Tinkering" active :todos="todoTinkering"></todo-component>
      </div>
    </div>
  </q-page>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from "vue";
import type { Todo } from "components/models";
import TodoComponent from "components/TodoComponent.vue";
import axios from "axios";

const todoToday = ref<Todo[]>([]);
const todoTinkering = ref<Todo[]>([]);

const fetchTodoToday = async () => {
  const response = await axios.get("/api/todo/today");
  todoToday.value = response.data as Todo[];
};

const fetchTodoTinkering = async () => {
  const response = await axios.get("/api/todo/tinkering");
  todoTinkering.value = response.data as Todo[];
};
//
// onMounted(() => {
//   fetchTodoToday()
//   fetchTodoTinkering()
// })

let updateInterval: number | null = null;
const UPDATE_INTERVAL_MS = 10000;
onMounted(() => {
  fetchTodoToday();
  fetchTodoTinkering();

  updateInterval = setInterval(() => {
    fetchTodoToday();
    // fetchTodoTinkering() // updating both results in a gui glitch, prob race condition
  }, UPDATE_INTERVAL_MS) as unknown as number;
});

onUnmounted(() => {
  if (updateInterval) {
    clearInterval(updateInterval);
  }
});
</script>
