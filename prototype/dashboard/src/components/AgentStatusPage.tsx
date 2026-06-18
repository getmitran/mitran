import { useState, useEffect } from 'react'
import { api } from '../api'

interface Agent {
  name: string
  type: string
  status: 'idle' | 'running' | 'error'
  last_task?: string
  last_run?: string
  tasks_completed: number
}

export function AgentStatusPage() {
  const [agents, setAgents] = useState<Agent[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const fetchAgents = () => {
      api.get('/api/v1/agents/status').then(r => {
        setAgents(r.data || [])
        setLoading(false)
      }).catch(() => {
        setAgents([
          { name: 'Dev', type: 'dev', status: 'idle', tasks_completed: 0 },
          { name: 'Docs', type: 'docs', status: 'idle', tasks_completed: 0 },
          { name: 'Ops', type: 'ops', status: 'idle', tasks_completed: 0 },
          { name: 'Review', type: 'review', status: 'idle', tasks_completed: 0 },
          { name: 'HR', type: 'hr', status: 'idle', tasks_completed: 0 },
          { name: 'CI/CD', type: 'cicd', status: 'idle', tasks_completed: 0 },
          { name: 'Tickets', type: 'tickets', status: 'idle', tasks_completed: 0 },
          { name: 'Wiki', type: 'wiki', status: 'idle', tasks_completed: 0 },
        ])
        setLoading(false)
      })
    }
    fetchAgents()
    const interval = setInterval(fetchAgents, 5000)
    return () => clearInterval(interval)
  }, [])

  const statusColor = (s: string) => {
    switch(s) {
      case 'running': return 'bg-blue-500'
      case 'error': return 'bg-red-500'
      default: return 'bg-green-500'
    }
  }

  const statusIcon = (s: string) => {
    switch(s) {
      case 'running': return '⚡'
      case 'error': return '❌'
      default: return '✅'
    }
  }

  if (loading) return <div className="p-6 text-gray-400">Loading agents...</div>

  return (
    <div className="p-6">
      <h2 className="text-xl font-semibold text-white mb-6">Agent Fleet</h2>
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        {agents.map(a => (
          <div key={a.name} className="bg-gray-800 rounded-lg p-4 border border-gray-700">
            <div className="flex items-center justify-between mb-3">
              <h3 className="text-white font-medium">{a.name}</h3>
              <div className={`w-2.5 h-2.5 rounded-full ${statusColor(a.status)}`} />
            </div>
            <div className="space-y-1">
              <div className="text-xs text-gray-400">Type: {a.type}</div>
              <div className="text-xs text-gray-400">Status: {statusIcon(a.status)} {a.status}</div>
              <div className="text-xs text-gray-400">Tasks: {a.tasks_completed}</div>
              {a.last_task && <div className="text-xs text-gray-500 truncate">Last: {a.last_task}</div>}
              {a.last_run && <div className="text-xs text-gray-500">{new Date(a.last_run).toLocaleString('en-IN', { timeZone: 'Asia/Kolkata' })}</div>}
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}
