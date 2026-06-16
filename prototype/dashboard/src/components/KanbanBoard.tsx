import { useState, useRef } from 'react'
import { Plus } from 'lucide-react'
import { Task } from '../types'
import { api } from '../api'
import TaskDetail from './TaskDetail'
import NewTaskModal from './NewTaskModal'

const COLUMNS = [
  { id: 'backlog', label: 'Backlog', status: 'queued' },
  { id: 'in_progress', label: 'In Progress', status: 'running' },
  { id: 'review', label: 'Review', status: 'checkpoint' },
  { id: 'done', label: 'Done', status: 'approved' },
] as const

const AGENT_COLORS: Record<string, string> = {
  'dev-agent': 'bg-blue-500', 'docs-agent': 'bg-green-500', 'cicd-agent': 'bg-purple-500',
  'tickets-agent': 'bg-orange-500', 'wiki-agent': 'bg-teal-500', 'ops-agent': 'bg-red-500',
  'hr-agent': 'bg-pink-500', 'dashboard-agent': 'bg-yellow-500', human: 'bg-gray-400',
}

const PRIORITY_COLORS: Record<number, string> = {
  1: 'bg-red-500', 2: 'bg-orange-500', 3: 'bg-yellow-500', 4: 'bg-blue-500', 5: 'bg-gray-500',
}

function statusToColumn(status: string): string {
  switch (status) {
    case 'queued': return 'backlog'
    case 'running': case 'in-progress': return 'in_progress'
    case 'checkpoint': case 'review': return 'review'
    case 'approved': case 'done': return 'done'
    default: return 'backlog'
  }
}

function columnToApiStatus(col: string): string {
  switch (col) {
    case 'backlog': return 'queued'
    case 'in_progress': return 'running'
    case 'review': return 'checkpoint'
    case 'done': return 'approved'
    default: return 'queued'
  }
}

interface Props { tasks: Task[]; refetch: () => void }

export default function KanbanBoard({ tasks, refetch }: Props) {
  const [selectedTask, setSelectedTask] = useState<Task | null>(null)
  const [newTaskOpen, setNewTaskOpen] = useState(false)
  const [dragOver, setDragOver] = useState<string | null>(null)
  const dragItem = useRef<string | null>(null)

  const grouped = COLUMNS.reduce((acc, col) => {
    acc[col.id] = tasks.filter(t => statusToColumn(t.status) === col.id)
    return acc
  }, {} as Record<string, Task[]>)

  const handleDragStart = (taskId: string) => { dragItem.current = taskId }
  const handleDragEnd = () => { dragItem.current = null; setDragOver(null) }

  const handleDrop = async (colId: string) => {
    setDragOver(null)
    if (!dragItem.current) return
    const newStatus = columnToApiStatus(colId)
    try {
      await api.updateTask(dragItem.current, { status: newStatus })
      refetch()
    } catch { /* silent */ }
  }

  const timeAgo = (d?: string) => {
    if (!d) return ''
    const s = Math.floor((Date.now() - new Date(d).getTime()) / 1000)
    if (s < 60) return 'just now'
    if (s < 3600) return `${Math.floor(s / 60)}m ago`
    if (s < 86400) return `${Math.floor(s / 3600)}h ago`
    return `${Math.floor(s / 86400)}d ago`
  }

  return (
    <div className="h-full flex flex-col">
      <div className="flex items-center justify-between mb-5">
        <h2 className="text-lg font-semibold text-gray-200">Kanban Board</h2>
      </div>

      <div className="flex-1 flex gap-4 overflow-x-auto pb-4">
        {COLUMNS.map(col => (
          <div
            key={col.id}
            className={`flex-1 min-w-[280px] flex flex-col rounded-xl border transition-all duration-200 ${
              dragOver === col.id ? 'ring-2 ring-cyan-500/50 border-cyan-500/30' : 'border-white/5'
            }`}
            style={{ background: '#1e293b' }}
            onDragOver={e => { e.preventDefault(); setDragOver(col.id) }}
            onDragLeave={() => setDragOver(null)}
            onDrop={() => handleDrop(col.id)}
          >
            <div className="flex items-center justify-between px-4 py-3 border-b border-white/5">
              <div className="flex items-center gap-2">
                <span className="text-sm font-medium text-gray-200">{col.label}</span>
                <span className="text-xs bg-white/10 text-gray-400 px-1.5 py-0.5 rounded-full">
                  {grouped[col.id].length}
                </span>
              </div>
              {col.id === 'backlog' && (
                <button onClick={() => setNewTaskOpen(true)}
                  className="p-1 rounded hover:bg-white/10 text-cyan-500 transition-colors">
                  <Plus size={16} />
                </button>
              )}
            </div>

            <div className="flex-1 overflow-y-auto p-3 space-y-2">
              {grouped[col.id].length === 0 ? (
                <p className="text-center text-gray-600 text-xs py-8">No tasks</p>
              ) : (
                grouped[col.id].map(task => (
                  <div
                    key={task.id}
                    draggable
                    onDragStart={() => handleDragStart(task.id)}
                    onDragEnd={handleDragEnd}
                    onClick={() => setSelectedTask(task)}
                    className="p-3 rounded-lg border border-white/5 cursor-pointer hover:border-cyan-500/30 transition-all duration-200 drag:opacity-50 drag:shadow-lg"
                    style={{ background: '#141b2d' }}
                  >
                    <p className="text-sm text-gray-200 font-medium leading-tight mb-2">{task.title}</p>
                    <div className="flex items-center justify-between">
                      <div className="flex items-center gap-2">
                        {task.agent_type && (
                          <span className={`text-[10px] px-1.5 py-0.5 rounded text-white ${AGENT_COLORS[task.agent_type] || 'bg-gray-600'}`}>
                            {task.agent_type.replace('-agent', '')}
                          </span>
                        )}
                        <span className={`w-2 h-2 rounded-full ${PRIORITY_COLORS[task.priority] || 'bg-gray-500'}`} />
                      </div>
                      <span className="text-[10px] text-gray-500">{timeAgo(task.created_at || task.createdAt)}</span>
                    </div>
                  </div>
                ))
              )}
            </div>
          </div>
        ))}
      </div>

      {selectedTask && <TaskDetail task={selectedTask} onClose={() => setSelectedTask(null)} refetch={refetch} />}
      {newTaskOpen && <NewTaskModal onClose={() => setNewTaskOpen(false)} refetch={refetch} />}
    </div>
  )
}
