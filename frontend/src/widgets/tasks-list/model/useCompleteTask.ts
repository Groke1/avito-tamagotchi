import { useCompleteTaskMutation } from '@/entities/task/api/taskApi'
import { isApiError, isFetchBaseQueryError } from '@/shared/lib/guards'
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
    } catch (error: unknown) {
      if (isFetchBaseQueryError(error) && isApiError(error.data)) {
        toast.error(error.data.message)
      }
    } finally {
      setCompletingTaskId(null)
    }
  }

  return { completingTaskId, handleComplete }
}
