import { mockTasks } from '../mock-data'
import TaskCard from './TaskCard'

export default function PriorityQueue() {
  const sorted = [...mockTasks].sort((a, b) => a.priority - b.priority)

  return (
    <div>
      <div className="flex items-center justify-between mb-5">
        <div>
          <h2 className="text-lg font-semibold text-gray-100">Priority Queue</h2>
          <p className="text-xs text-gray-500 mt-0.5">Drag to reorder • {sorted.length} tasks</p>
        </div>
        <button className="text-xs bg-accent/10 text-accent px-3 py-1.5 rounded-md hover:bg-accent/20 transition-colors">
          + New Task
        </button>
      </div>
      <div className="flex flex-col gap-2">
        {sorted.map(task => (
          <TaskCard key={task.id} task={task} />
        ))}
      </div>
    </div>
  )
}
