import { useGetTasksQuery } from '@/entities/task/api/taskApi'
import { Tabs } from '@/shared/ui'
import { TABS } from '../model/consts'
import { useCompleteTask } from '../model/useCompleteTask'
import { useTasksListTabs } from '../model/useTasksListTabs'
import { TasksContentList } from './TasksContentList'
import { TasksTabList } from './TasksTabList'

export const TasksListWidget = () => {
  const { data: tasksData, isLoading, isError, refetch } = useGetTasksQuery()
  const { completingTaskId, handleComplete } = useCompleteTask()
  const { items, activeCount, activeTab, completedCount, getFilteredTasks, handleSetTab } =
    useTasksListTabs(tasksData)

  return (
    <section className="space-y-6">
      <Tabs value={activeTab} onValueChange={handleSetTab} className="w-full">
        <TasksTabList
          itemsCount={items.length}
          activeCount={activeCount}
          completedCount={completedCount}
          showCounts={!isLoading && !isError}
        />
        {TABS.map((tab) => (
          <TasksContentList
            key={tab.id}
            tabId={tab.id}
            tasks={getFilteredTasks(tab.id)}
            isLoading={isLoading}
            isError={isError}
            completingTaskId={completingTaskId}
            onComplete={handleComplete}
            onRetry={refetch}
          />
        ))}
      </Tabs>
    </section>
  )
}
