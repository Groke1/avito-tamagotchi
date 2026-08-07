import { DefaultTaskIcon, taskIconMap } from '../model/consts'
import type { Task } from '../model/types'

export const TaskMiniCard = ({ title, description }: Task) => {
  const Icon = taskIconMap[title] ?? DefaultTaskIcon

  return (
    <li className="p-4 rounded-xl border border-outline-variant bg-surface-bg">
      <div className="flex items-center gap-4">
        <Icon size={40} />
        {description}
      </div>
    </li>
  )
}
