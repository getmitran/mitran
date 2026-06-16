import { useState } from 'react'
import { X, Save } from 'lucide-react'
import { Task, Comment } from '../types'
import { api } from '../api'

interface Props { task: Task; onClose: () => void; refetch: () => void }

export default function TaskDetail({ task, onClose, refetch }: Props) {
  const [title, setTitle] = useState(task.title)
  const [description, setDescription] = useState(task.description || '')
  const [priority, setPriority] = useState(task.priority)
  const [comments, setComments] = useState<Comment[]>((task as any).comments || [])
  const [newComment, setNewComment] = useState('')
  const [saving, setSaving] = useState(false)

  const handleSave = async () => {
    setSaving(true)
    try {
      await api.updateTask(task.id, { title, description, priority })
      refetch()
      onClose()
    } catch { /* silent */ } finally { setSaving(false) }
  }

  const addComment = () => {
    if (!newComment.trim()) return
    setComments([...comments, { id: crypto.randomUUID(), text: newComment, author: 'you', created_at: new Date().toISOString() }])
    setNewComment('')
  }

  return (
    <div className="fixed inset-0 z-50 flex justify-end" onClick={onClose}>
      <div className="absolute inset-0 bg-black/50" />
      <div className="relative w-[480px] h-full bg-[#0a0e1a] border-l border-white/10 overflow-y-auto"
        onClick={e => e.stopPropagation()}>
        <div className="flex items-center justify-between px-6 py-4 border-b border-white/5 sticky top-0 bg-[#0a0e1a] z-10">
          <span className="text-xs text-gray-500 font-mono">{task.id.slice(0, 8)}</span>
          <div className="flex gap-2">
            <button onClick={handleSave} disabled={saving}
              className="flex items-center gap-1.5 text-xs bg-cyan-500 hover:bg-cyan-400 text-black font-medium px-3 py-1.5 rounded-md transition-colors disabled:opacity-50">
              <Save size={12} /> Save
            </button>
            <button onClick={onClose} className="p-1.5 rounded hover:bg-white/10 text-gray-400"><X size={16} /></button>
          </div>
        </div>

        <div className="p-6 space-y-5">
          <input value={title} onChange={e => setTitle(e.target.value)}
            className="w-full bg-transparent text-lg font-semibold text-gray-200 border-b border-white/10 pb-2 focus:outline-none focus:border-cyan-500" />

          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="text-[11px] text-gray-500 uppercase tracking-wider">Status</label>
              <p className="text-sm text-gray-300 mt-1 capitalize">{task.status}</p>
            </div>
            <div>
              <label className="text-[11px] text-gray-500 uppercase tracking-wider">Agent</label>
              <p className="text-sm text-gray-300 mt-1">{task.agent_type || task.assignedAgent || '—'}</p>
            </div>
            <div>
              <label className="text-[11px] text-gray-500 uppercase tracking-wider">Priority</label>
              <select value={priority} onChange={e => setPriority(+e.target.value as any)}
                className="mt-1 w-full bg-[#141b2d] text-sm text-gray-300 border border-white/10 rounded px-2 py-1 focus:outline-none focus:border-cyan-500">
                {[1,2,3,4,5].map(p => <option key={p} value={p}>P{p}</option>)}
              </select>
            </div>
            <div>
              <label className="text-[11px] text-gray-500 uppercase tracking-wider">Created</label>
              <p className="text-sm text-gray-300 mt-1">{new Date(task.created_at || task.createdAt).toLocaleDateString()}</p>
            </div>
          </div>

          <div>
            <label className="text-[11px] text-gray-500 uppercase tracking-wider">Description</label>
            <textarea value={description} onChange={e => setDescription(e.target.value)} rows={6}
              className="mt-1 w-full bg-[#141b2d] text-sm text-gray-300 border border-white/10 rounded-md p-3 focus:outline-none focus:border-cyan-500 resize-none font-mono" />
          </div>

          {(task as any).checkpoint_output && (
            <div>
              <label className="text-[11px] text-gray-500 uppercase tracking-wider">Checkpoint Output</label>
              <pre className="mt-1 bg-[#141b2d] text-xs text-gray-400 p-3 rounded-md border border-white/5 overflow-x-auto whitespace-pre-wrap">
                {(task as any).checkpoint_output}
              </pre>
            </div>
          )}

          <div className="border-t border-white/5 pt-4">
            <label className="text-[11px] text-gray-500 uppercase tracking-wider">Comments</label>
            <div className="mt-2 space-y-2 max-h-48 overflow-y-auto">
              {comments.length === 0 && <p className="text-xs text-gray-600">No comments yet</p>}
              {comments.map(c => (
                <div key={c.id} className="bg-[#141b2d] rounded p-2 border border-white/5">
                  <p className="text-xs text-gray-300">{c.text}</p>
                  <p className="text-[10px] text-gray-500 mt-1">{c.author} • {new Date(c.created_at).toLocaleString()}</p>
                </div>
              ))}
            </div>
            <div className="flex gap-2 mt-3">
              <input value={newComment} onChange={e => setNewComment(e.target.value)}
                placeholder="Add a comment..." onKeyDown={e => e.key === 'Enter' && addComment()}
                className="flex-1 bg-[#141b2d] text-xs text-gray-300 border border-white/10 rounded px-3 py-2 focus:outline-none focus:border-cyan-500" />
              <button onClick={addComment} className="text-xs bg-cyan-500/20 text-cyan-400 px-3 py-2 rounded hover:bg-cyan-500/30 transition-colors">Post</button>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
