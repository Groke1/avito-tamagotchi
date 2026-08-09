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
    <section>
      <Tabs value={activeTab} onValueChange={handleSetTab} className="w-full">
        <TasksTabList
          itemsCount={items.length}
          activeCount={activeCount}
          completedCount={completedCount}
          showCounts={!isLoading && !isError}
        />
        <div className="mt-6">
          {TABS.map(({ id }) => (
            <TasksContentList
              key={id}
              tabId={id}
              tasks={getFilteredTasks(id)}
              isLoading={isLoading}
              isError={isError}
              completingTaskId={completingTaskId}
              onComplete={handleComplete}
              onRetry={refetch}
            />
          ))}
        </div>
      </Tabs>
    </section>
  )
}
