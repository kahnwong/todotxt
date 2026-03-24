import { ref } from 'vue'
import type { Todo, Lane, DragData } from '../components/models'
import { LANE_IDS, CONTEXT } from '../constants/kanban'
import axios from 'axios'

/**
 * Composable for managing drag-and-drop functionality
 */
export function useDragDrop(
  fetchTodoToday: () => Promise<void>,
  fetchTodoWork: () => Promise<void>,
  fetchTodoTinkering: () => Promise<void>,
) {
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
    const sourceProject = sourceLane.id === LANE_IDS.TODAY ? '' : sourceLane.title
    const targetProject = targetLane.id === LANE_IDS.TODAY ? '' : targetLane.title

    if (todo.status === targetStatus && sourceProject === targetProject) {
      dragData = null
      return
    }

    // Determine context: add @tinkering or @work when moving from Today to a project lane
    let context = todo.context || ''
    if (sourceLane.id === LANE_IDS.TODAY && targetLane.id !== LANE_IDS.TODAY) {
      // Moving from Today to a project lane
      if (targetLane.id.startsWith('work-project-')) {
        // Moving to work project lane - add @work
        context = CONTEXT.WORK
      } else if (targetLane.id.startsWith('project-')) {
        // Moving to tinkering project lane - add @tinkering
        context = CONTEXT.TINKERING
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
      await fetchTodoWork()
      await fetchTodoTinkering()
    } catch (error) {
      console.error('Failed to update todo:', error)
    }

    dragData = null
  }

  const isDropTarget = (laneId: string, status: string): boolean => {
    return dropTarget.value === `${laneId}-${status}`
  }

  return {
    dropTarget,
    handleDragStart,
    handleDragOver,
    handleDragEnd,
    handleDrop,
    isDropTarget,
  }
}
