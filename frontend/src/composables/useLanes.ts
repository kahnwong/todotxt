import { computed, type Ref } from 'vue'
import type { Todo, Lane } from '../components/models'
import { LANE_IDS, CONTEXT } from '../constants/kanban'

/**
 * Composable for managing kanban lanes
 */
export function useLanes(todoToday: Ref<Todo[]>, todoTinkering: Ref<Todo[]>) {
  const lanes = computed<Lane[]>(() => {
    // Filter today todos to exclude @tinkering tasks with projects
    const filteredTodayTodos = todoToday.value.filter((todo) => {
      const hasTinkering = todo.context === CONTEXT.TINKERING
      const hasProject = !!todo.project
      // Exclude if both @tinkering and project exist
      return !(hasTinkering && hasProject)
    })

    const laneList: Lane[] = [{ id: LANE_IDS.TODAY, title: 'Today', todos: filteredTodayTodos }]

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

  return {
    lanes,
    getTodosByStatus,
  }
}
