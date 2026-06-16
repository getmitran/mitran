import { GripVertical, Bot, Clock, ChevronUp, ChevronDown } from 'lucide-react'
import { Task } from '../types'

const statusColors: Record<string, string> = {
  'queued': 'bg-gray-600 text-gray-200',
  'in-progress': 'bg-cyan-900/60 text-cyan-300',
  'review': 'bg-amber-900/60 text-amber-300',
  'done': 'bg-emerald-900/60 text-emerald-300',
  'blocked': 'bg-red-900/60 text-red-300',
}

const priorityColors: Record<number, string> = {
  1: 'bg-red-500',
  2: 'bg-orange-500',
  3: 'bg-yellow-500',
  4: 'bg-blue-500',
  5: 'bg-gray-500',
}

interface Props { task: Task; onMoveUp?: () => void; onMoveDown?: () => void }

export default function TaskCard({ task, onMoveUp, onMoveDown }: Props) {
  const time = new Date(task.created_at || task.createdAt).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })

  return (
    <div className="group flex items-center gap-3 bg-card border border-gray-800 rounded-lg px-4 py-3 hover:border-accent/40 hover:bg-card-hover transition-all">
      <div className="flex flex-col gap-0.5 shrink-0">
        {onMoveUp && <button onClick={onMoveUp} className="text-gray-600 hover:text-accent"><ChevronUp size={12} /></button>}
        <GripVertical size={14} className="text-gray-600 group-hover:text-gray-400" />
        {onMoveDown && <button onClick={onMoveDown} className="text-gray-600 hover:text-accent"><ChevronDown size={12} /></button>}
      </div>
      <div className={`w-2 h-2 rounded-full shrink-0 ${priorityColors[task.priority] || priorityColors[3]}`} />
      <div className="flex-1 min-w-0">
        <p className="text-sm font-medium text-gray-200 truncate">{task.title}</p>
        <div className="flex items-center gap-3 mt-1">
          <span className="flex items-center gap-1 text-[11px] text-gray-500">
            <Bot size={11} /> {task.agent_type || task.assignedAgent}
          </span>
          <span className="flex items-center gap-1 text-[11px] text-gray-500">
            <Clock size={11} /> {time}
          </span>
        </div>
      </div>
      <span className={`text-[11px] px-2 py-0.5 rounded-full font-medium ${statusColors[task.status] || statusColors.queued}`}>
        {task.status}
      </span>
      <span className="text-[11px] text-gray-600 font-mono w-5 text-right">P{task.priority}</span>
    </div>
  )
}
