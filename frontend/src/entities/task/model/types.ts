export interface Task {
  id: string
  title: TaskTitle
  description: string
  reward_coins: number
  reward_xp: number
  status: TaskStatus
  completed_at: null
}

export type TaskStatus = 'active' | 'completed'
export type TaskTitle = 'Первая продажа месяца' | 'Лояльный продавец' | 'Быстрый ответ покупателю'

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
