import { useState, useEffect } from 'react'
import { api } from '../api'

interface Checkpoint {
  id: string
  task_id: string
  agent: string
  status: 'pending' | 'approved' | 'rejected'
  description: string
  files_changed: string[]
  created_at: string
}

export function CheckpointViewer() {
  const [checkpoints, setCheckpoints] = useState<Checkpoint[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    api.get('/api/v1/checkpoints').then(r => {
      setCheckpoints(r.data || [])
      setLoading(false)
    }).catch(() => setLoading(false))
  }, [])

  const handleAction = (id: string, action: 'approve' | 'reject') => {
    api.post(`/api/v1/checkpoints/${id}/${action}`).then(() => {
      setCheckpoints(prev => prev.map(cp =>
        cp.id === id ? { ...cp, status: action === 'approve' ? 'approved' : 'rejected' } : cp
      ))
    }).catch(() => {})
  }

  if (loading) return <div className="p-6 text-gray-400">Loading checkpoints...</div>

  return (
    <div className="p-6">
      <h2 className="text-xl font-semibold text-white mb-6">Checkpoints</h2>
      {checkpoints.length === 0 ? (
        <div className="text-gray-500">No checkpoints yet. Checkpoints appear when agents reach approval gates.</div>
      ) : (
        <div className="space-y-4">
          {checkpoints.map(cp => (
            <div key={cp.id} className="bg-gray-800 rounded-lg p-4 border border-gray-700">
              <div className="flex items-center justify-between">
                <div>
                  <span className="text-white font-medium">{cp.description}</span>
                  <div className="text-xs text-gray-500 mt-1">Agent: {cp.agent} | Task: {cp.task_id}</div>
                </div>
                <span className={`text-xs px-2 py-1 rounded ${
                  cp.status === 'approved' ? 'bg-green-900 text-green-300' :
                  cp.status === 'rejected' ? 'bg-red-900 text-red-300' :
                  'bg-yellow-900 text-yellow-300'
                }`}>{cp.status}</span>
              </div>
              {cp.files_changed?.length > 0 && (
                <div className="mt-2">
                  <span className="text-xs text-gray-400">Files:</span>
                  <div className="flex flex-wrap gap-1 mt-1">
                    {cp.files_changed.map(f => (
                      <span key={f} className="text-[10px] px-1.5 py-0.5 bg-gray-700 text-gray-300 rounded font-mono">{f}</span>
                    ))}
                  </div>
                </div>
              )}
              {cp.status === 'pending' && (
                <div className="flex gap-2 mt-3">
                  <button onClick={() => handleAction(cp.id, 'approve')}
                    className="px-3 py-1 text-xs bg-green-600 text-white rounded hover:bg-green-700">Approve</button>
                  <button onClick={() => handleAction(cp.id, 'reject')}
                    className="px-3 py-1 text-xs bg-red-600 text-white rounded hover:bg-red-700">Reject</button>
                </div>
              )}
              <div className="text-xs text-gray-600 mt-2">
                {new Date(cp.created_at).toLocaleString('en-IN', { timeZone: 'Asia/Kolkata' })}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
