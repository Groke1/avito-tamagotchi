export interface Task {
  id: string
  title: string
  description: string
  reward_coins: number
  reward_xp: number
  status: TaskStatus
  completed_at: null
  task_type: TaskType
}

export type TaskStatus = 'active' | 'completed'
export type TaskType =
  | 'Отзывы'
  | 'Поиск'
  | 'Сообщения'
  | 'Категории'
  | 'Покупки'
  | 'Избранное'

export interface TasksResponse {
  date: string
  items: Task[]
}

export type TaskResponse = Task

export interface CompleteTaskResponse {
  task: Task
  awarded: {
    coins: number
    xp: number
  }
}
