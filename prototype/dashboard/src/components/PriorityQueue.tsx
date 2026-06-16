import { useState } from 'react'
import { Task } from '../types'
import { api } from '../api'
import TaskCard from './TaskCard'

interface Props { tasks: Task[]; refetch: () => void }

export default function PriorityQueue({ tasks, refetch }: Props) {
  const [showNew, setShowNew] = useState(false)
  const [title, setTitle] = useState('')
  const sorted = [...tasks].sort((a, b) => a.priority - b.priority)

  const create = async () => {
    if (!title.trim()) return
    await api.createTask({ title, priority: 3 })
    setTitle('')
    setShowNew(false)
    refetch()
  }

  const reorder = async (id: string, newPriority: number) => {
    await api.reorderTask(id, newPriority)
    refetch()
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-5">
        <div>
          <h2 className="text-lg font-semibold text-gray-100">Priority Queue</h2>
          <p className="text-xs text-gray-500 mt-0.5">Drag to reorder • {sorted.length} tasks</p>
        </div>
        <button onClick={() => setShowNew(true)}
          className="text-xs bg-accent/10 text-accent px-3 py-1.5 rounded-md hover:bg-accent/20 transition-colors">
          + New Task
        </button>
      </div>
      {showNew && (
        <div className="flex gap-2 mb-3">
          <input value={title} onChange={e => setTitle(e.target.value)} placeholder="Task title..."
            className="flex-1 bg-base border border-gray-800 rounded-md px-3 py-1.5 text-sm text-gray-200 placeholder-gray-600 focus:outline-none focus:border-accent/50"
            onKeyDown={e => e.key === 'Enter' && create()} />
          <button onClick={create} className="text-xs bg-accent text-black px-3 py-1.5 rounded-md">Add</button>
          <button onClick={() => setShowNew(false)} className="text-xs text-gray-500 px-2">Cancel</button>
        </div>
      )}
      <div className="flex flex-col gap-2">
        {sorted.map((task, i) => (
          <TaskCard key={task.id} task={task}
            onMoveUp={i > 0 ? () => reorder(task.id, sorted[i - 1].priority) : undefined}
            onMoveDown={i < sorted.length - 1 ? () => reorder(task.id, sorted[i + 1].priority) : undefined} />
        ))}
      </div>
    </div>
  )
}
