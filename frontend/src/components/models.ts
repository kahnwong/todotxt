export interface Todo {
  id: number
  context: string
  project: string
  todo: string
  status: string
  created_date: string
}

export interface Lane {
  id: string
  title: string
  todos: Todo[]
}

export interface DragData {
  todo: Todo
  sourceLane: Lane
}
