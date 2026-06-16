import { useState } from 'react'
import { Check, X } from 'lucide-react'
import { Checkpoint } from '../types'
import { api } from '../api'

interface Props { checkpoints: Checkpoint[]; refetch: () => void }

export default function CheckpointReview({ checkpoints, refetch }: Props) {
  const [feedback, setFeedback] = useState('')
  const [acting, setActing] = useState(false)
  const pending = checkpoints.filter(c => c.status === 'pending')
  const cp = pending[0]

  if (!cp) return <p className="text-gray-500">No checkpoints pending review.</p>

  const approve = async () => {
    setActing(true)
    await api.approveCheckpoint(cp.id, feedback || undefined)
    setFeedback('')
    setActing(false)
    refetch()
  }

  const reject = async () => {
    if (!feedback.trim()) return
    setActing(true)
    await api.rejectCheckpoint(cp.id, feedback)
    setFeedback('')
    setActing(false)
    refetch()
  }

  const agentName = cp.agent_name || cp.agentName
  const taskId = cp.task_id || cp.taskId
  const diffs = cp.file_diffs || cp.fileDiffs || []
  const created = cp.created_at || cp.createdAt

  return (
    <div className="max-w-3xl">
      <h2 className="text-lg font-semibold text-gray-100 mb-1">Checkpoint Review</h2>
      <p className="text-xs text-gray-500 mb-1">{pending.length} pending</p>
      <p className="text-xs text-gray-500 mb-5">
        {agentName} • Task {taskId} • {new Date(created).toLocaleTimeString()}
      </p>

      <div className="bg-card border border-gray-800 rounded-lg p-4 mb-4">
        <h3 className="text-xs font-semibold text-gray-400 uppercase tracking-wider mb-2">Agent Output</h3>
        <pre className="text-sm text-gray-300 whitespace-pre-wrap font-mono leading-relaxed">{cp.output}</pre>
      </div>

      <div className="bg-card border border-gray-800 rounded-lg p-4 mb-4">
        <h3 className="text-xs font-semibold text-gray-400 uppercase tracking-wider mb-2">File Changes</h3>
        {diffs.length === 0 && <p className="text-xs text-gray-600">No file changes</p>}
        {diffs.map((diff, i) => (
          <div key={i} className="mb-3 last:mb-0">
            <p className="text-xs text-gray-500 font-mono mb-1">{diff.path}</p>
            <div className="bg-base rounded border border-gray-800 p-2 font-mono text-xs">
              {diff.deletions.map((line, j) => (
                <div key={`d-${j}`} className="text-red-400 bg-red-950/30 px-2 py-0.5 rounded-sm">- {line}</div>
              ))}
              {diff.additions.map((line, j) => (
                <div key={`a-${j}`} className="text-green-400 bg-green-950/30 px-2 py-0.5 rounded-sm">+ {line}</div>
              ))}
            </div>
          </div>
        ))}
      </div>

      <textarea value={feedback} onChange={e => setFeedback(e.target.value)}
        placeholder="Optional feedback for the agent..."
        className="w-full bg-base border border-gray-800 rounded-lg px-3 py-2 text-sm text-gray-300 placeholder-gray-600 resize-none h-20 focus:outline-none focus:border-accent/50 mb-3" />

      <div className="flex gap-2">
        <button onClick={approve} disabled={acting}
          className="flex items-center gap-1.5 bg-emerald-600 hover:bg-emerald-500 disabled:opacity-50 text-white text-sm px-4 py-2 rounded-md transition-colors">
          <Check size={14} /> Approve
        </button>
        <button onClick={reject} disabled={acting || !feedback.trim()}
          className="flex items-center gap-1.5 bg-red-600 hover:bg-red-500 disabled:opacity-50 text-white text-sm px-4 py-2 rounded-md transition-colors">
          <X size={14} /> Reject
        </button>
      </div>
    </div>
  )
}
