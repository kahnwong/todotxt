<template>
  <q-layout view="lHh Lpr lFf">
    <q-header elevated>
      <q-toolbar>
        <q-toolbar-title>
          <router-link to="/" style="text-decoration: none; color: inherit"> todo.txt </router-link>
        </q-toolbar-title>
        <q-space />
        <q-space />
        <q-toolbar-title shrink class="text-subtitle1">
          <router-link to="kanban/" style="text-decoration: none; color: inherit">
            Kanban View
          </router-link>
        </q-toolbar-title>
      </q-toolbar>
    </q-header>
    <q-page-container>
      <router-view />
    </q-page-container>
  </q-layout>
</template>

<script setup lang="ts"></script>
