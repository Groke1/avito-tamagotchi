import { useCompleteTaskMutation } from '@/entities/task/api/taskApi'
import { useState } from 'react'
import { toast } from 'sonner'

export const useCompleteTask = () => {
  const [completeTask] = useCompleteTaskMutation()
  const [completingTaskId, setCompletingTaskId] = useState<string | null>(null)

  const handleComplete = async (taskId: string) => {
    try {
      setCompletingTaskId(taskId)
      await completeTask(taskId).unwrap()

      toast.success('Задание успешно выполнено! Награда получена 🪙✨')
    } catch (err: unknown) {
      const errorObj = err as { status?: number }
      if (errorObj?.status === 409) {
        toast.error('Награда за это задание уже получена!')
      } else {
        toast.error('Не удалось выполнить задание. Попробуйте позже.')
      }
    } finally {
      setCompletingTaskId(null)
    }
  }

  return { completingTaskId, handleComplete }
}
