import { ref } from 'vue'
import type { Todo } from '../components/models'
import {
  highlightTodoSyntax,
  reconstructTodoText,
  cleanTodoText,
  prependHttpsToUrls,
} from '../utils/todoText'
import axios from 'axios'

/**
 * Composable for managing todo edit dialog
 */
export function useTodoEdit(
  fetchTodoToday: () => Promise<void>,
  fetchTodoTinkering: () => Promise<void>,
) {
  const showEditDialog = ref(false)
  const editingTodo = ref<Todo | null>(null)
  const editedText = ref('')
  const highlightedText = ref('')

  const openEditDialog = (todo: Todo) => {
    editingTodo.value = todo
    editedText.value = reconstructTodoText(todo)
    updateHighlight()
    showEditDialog.value = true
  }

  const updateHighlight = () => {
    highlightedText.value = highlightTodoSyntax(editedText.value)
  }

  const saveTodo = async () => {
    if (!editingTodo.value) return

    try {
      let textToSave = cleanTodoText(editedText.value)
      textToSave = prependHttpsToUrls(textToSave)

      await axios.put('/api/todo/update-content', {
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

  return {
    showEditDialog,
    editingTodo,
    editedText,
    highlightedText,
    openEditDialog,
    updateHighlight,
    saveTodo,
  }
}
