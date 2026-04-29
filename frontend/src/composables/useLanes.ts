import { computed, type Ref } from 'vue'
import type { Todo, Lane } from '../components/models'
import { LANE_IDS, CONTEXT } from '../constants/kanban'

/**
 * Composable for managing kanban lanes
 */
export function useLanes(
  todoToday: Ref<Todo[]>,
  todoTinkering: Ref<Todo[]>,
  todoWork: Ref<Todo[]>,
) {
  const lanes = computed<Lane[]>(() => {
    // Filter today todos to exclude @tinkering and @work tasks with projects
    const filteredTodayTodos = todoToday.value.filter((todo) => {
      const hasTinkering = todo.context === CONTEXT.TINKERING
      const hasWork = todo.context === CONTEXT.WORK
      const hasProject = !!todo.project
      // Exclude if (@tinkering or @work) and project exist
      return !((hasTinkering || hasWork) && hasProject)
    })

    const laneList: Lane[] = [{ id: LANE_IDS.TODAY, title: 'Today', todos: filteredTodayTodos }]

    // Group work todos by project
    const workProjectGroups = new Map<string, Todo[]>()
    todoWork.value.forEach((todo) => {
      const project = todo.project || 'No Project'
      if (!workProjectGroups.has(project)) {
        workProjectGroups.set(project, [])
      }
      workProjectGroups.get(project)!.push(todo)
    })

    // Create a lane for each work project
    workProjectGroups.forEach((todos, project) => {
      laneList.push({
        id: `work-project-${project}`,
        title: project.startsWith('+') ? project.slice(1) : project,
        todos,
      })
    })

    // Group tinkering todos by project
    const tinkeringProjectGroups = new Map<string, Todo[]>()
    todoTinkering.value.forEach((todo) => {
      const project = todo.project || 'No Project'
      if (!tinkeringProjectGroups.has(project)) {
        tinkeringProjectGroups.set(project, [])
      }
      tinkeringProjectGroups.get(project)!.push(todo)
    })

    // Create a lane for each tinkering project
    tinkeringProjectGroups.forEach((todos, project) => {
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
