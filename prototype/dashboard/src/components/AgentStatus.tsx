import { Bot } from 'lucide-react'
import { Agent, AgentState } from '../types'

const stateStyles: Record<AgentState, { dot: string; label: string }> = {
  idle: { dot: 'bg-gray-500', label: 'text-gray-400' },
  working: { dot: 'bg-cyan-400 animate-pulse', label: 'text-cyan-300' },
  error: { dot: 'bg-red-400', label: 'text-red-300' },
}

interface Props { agents: Agent[] }

export default function AgentStatus({ agents }: Props) {
  return (
    <div>
      <h2 className="text-lg font-semibold text-gray-100 mb-5">Agents</h2>
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-3">
        {agents.map(agent => {
          const s = stateStyles[agent.state] || stateStyles.idle
          const heartbeat = agent.last_heartbeat
            ? new Date(agent.last_heartbeat).toLocaleTimeString()
            : agent.lastActivity
          return (
            <div key={agent.id} className="bg-card border border-gray-800 rounded-lg p-4 hover:border-accent/30 transition-colors">
              <div className="flex items-center gap-2 mb-3">
                <Bot size={16} className="text-accent" />
                <span className="text-sm font-medium text-gray-200">{agent.name || agent.type}</span>
              </div>
              <div className="flex items-center gap-2 mb-2">
                <span className={`w-2 h-2 rounded-full ${s.dot}`} />
                <span className={`text-xs capitalize ${s.label}`}>{agent.state}</span>
              </div>
              <p className="text-[11px] text-gray-600">Last active: {heartbeat}</p>
              {agent.currentTask && (
                <p className="text-[11px] text-gray-500 mt-1 truncate">→ {agent.currentTask}</p>
              )}
            </div>
          )
        })}
      </div>
    </div>
  )
}
