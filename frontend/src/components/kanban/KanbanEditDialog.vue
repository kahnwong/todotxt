<template>
  <q-dialog :model-value="modelValue" @update:model-value="emit('update:modelValue', $event)">
    <q-card style="min-width: 500px">
      <q-card-section>
        <div class="text-h6">Edit Todo</div>
      </q-card-section>

      <q-card-section>
        <div class="editor-container">
          <div class="syntax-highlight" v-html="highlightedText"></div>
          <textarea
            :value="editedText"
            class="editor-input"
            rows="5"
            @input="handleInput"
          ></textarea>
        </div>
      </q-card-section>

      <q-card-actions align="right">
        <q-btn flat label="Cancel" color="primary" v-close-popup />
        <q-btn flat label="Save" color="primary" @click="emit('save')" />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script setup lang="ts">
defineProps<{
  modelValue: boolean
  editedText: string
  highlightedText: string
}>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  updateText: [value: string]
  save: []
}>()

const handleInput = (event: Event) => {
  const target = event.target as HTMLTextAreaElement
  emit('updateText', target.value)
}
</script>

<style scoped>
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

:deep(.highlight-context) {
  color: #4fa1ac;
  font-weight: 600;
}

:deep(.highlight-project) {
  color: #c8aa6f;
  font-weight: 600;
}

:deep(.highlight-status) {
  color: #8250df;
  font-weight: 600;
}

:deep(.highlight-due) {
  color: #cf222e;
  font-weight: 600;
}

:deep(.highlight-date) {
  color: #0969da;
  font-weight: 600;
}
</style>
