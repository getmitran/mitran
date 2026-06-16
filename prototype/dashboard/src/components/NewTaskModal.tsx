import { useState } from 'react'
import { X } from 'lucide-react'
import { api } from '../api'

const AGENT_TYPES = ['dev-agent', 'docs-agent', 'cicd-agent', 'tickets-agent', 'wiki-agent', 'ops-agent', 'hr-agent', 'dashboard-agent', 'human']

interface Props { onClose: () => void; refetch: () => void }

export default function NewTaskModal({ onClose, refetch }: Props) {
  const [title, setTitle] = useState('')
  const [description, setDescription] = useState('')
  const [agentType, setAgentType] = useState('dev-agent')
  const [priority, setPriority] = useState(3)
  const [labels, setLabels] = useState('')
  const [submitting, setSubmitting] = useState(false)

  const handleSubmit = async () => {
    if (!title.trim()) return
    setSubmitting(true)
    try {
      await api.createTask({ title, description, agent_type: agentType, priority, labels: labels.split(',').map(l => l.trim()).filter(Boolean) } as any)
      refetch()
      onClose()
    } catch { /* silent */ } finally { setSubmitting(false) }
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center" onClick={onClose}>
      <div className="absolute inset-0 bg-black/60" />
      <div className="relative w-[440px] bg-[#0a0e1a] border border-white/10 rounded-xl shadow-2xl" onClick={e => e.stopPropagation()}>
        <div className="flex items-center justify-between px-5 py-4 border-b border-white/5">
          <h3 className="text-sm font-semibold text-gray-200">New Task</h3>
          <button onClick={onClose} className="p-1 rounded hover:bg-white/10 text-gray-400"><X size={16} /></button>
        </div>

        <div className="p-5 space-y-4">
          <div>
            <label className="text-[11px] text-gray-500 uppercase tracking-wider">Title *</label>
            <input value={title} onChange={e => setTitle(e.target.value)} autoFocus
              className="mt-1 w-full bg-[#141b2d] text-sm text-gray-200 border border-white/10 rounded-md px-3 py-2 focus:outline-none focus:border-cyan-500" />
          </div>
          <div>
            <label className="text-[11px] text-gray-500 uppercase tracking-wider">Description</label>
            <textarea value={description} onChange={e => setDescription(e.target.value)} rows={3}
              className="mt-1 w-full bg-[#141b2d] text-sm text-gray-300 border border-white/10 rounded-md px-3 py-2 focus:outline-none focus:border-cyan-500 resize-none" />
          </div>
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="text-[11px] text-gray-500 uppercase tracking-wider">Agent Type</label>
              <select value={agentType} onChange={e => setAgentType(e.target.value)}
                className="mt-1 w-full bg-[#141b2d] text-sm text-gray-300 border border-white/10 rounded px-2 py-2 focus:outline-none focus:border-cyan-500">
                {AGENT_TYPES.map(a => <option key={a} value={a}>{a}</option>)}
              </select>
            </div>
            <div>
              <label className="text-[11px] text-gray-500 uppercase tracking-wider">Priority (P{priority})</label>
              <input type="range" min={1} max={5} value={priority} onChange={e => setPriority(+e.target.value)}
                className="mt-2 w-full accent-cyan-500" />
            </div>
          </div>
          <div>
            <label className="text-[11px] text-gray-500 uppercase tracking-wider">Labels (comma-separated)</label>
            <input value={labels} onChange={e => setLabels(e.target.value)} placeholder="frontend, urgent, mvp"
              className="mt-1 w-full bg-[#141b2d] text-sm text-gray-300 border border-white/10 rounded-md px-3 py-2 focus:outline-none focus:border-cyan-500" />
          </div>
        </div>

        <div className="flex justify-end gap-2 px-5 py-4 border-t border-white/5">
          <button onClick={onClose} className="text-xs text-gray-400 hover:text-gray-200 px-3 py-2 transition-colors">Cancel</button>
          <button onClick={handleSubmit} disabled={!title.trim() || submitting}
            className="text-xs bg-cyan-500 hover:bg-cyan-400 text-black font-medium px-4 py-2 rounded-md transition-colors disabled:opacity-50">
            Create Task
          </button>
        </div>
      </div>
    </div>
  )
}
