import { useState, useCallback } from 'react'
import Sidebar from './components/Sidebar'
import PriorityQueue from './components/PriorityQueue'
import { AgentStatusPage } from './components/AgentStatusPage'
import CheckpointReview from './components/CheckpointReview'
import InitWizard from './components/InitWizard'
import KanbanBoard from './components/KanbanBoard'
import SettingsPage from './components/SettingsPage'
import TicketingApp from './components/TicketingApp'
import WikiApp from './components/WikiApp'
import ChatInterface from './components/ChatInterface'
import { SessionsPage } from './components/SessionsPage'
import { api } from './api'
import { usePolling } from './hooks/usePolling'
import { mockTasks, mockAgents, mockCheckpoints } from './mock-data'
import { Task, Agent, Checkpoint } from './types'

type View = 'queue' | 'kanban' | 'chat' | 'sessions' | 'agents' | 'checkpoints' | 'tickets' | 'wiki' | 'settings'

function normalize<T>(data: any[] | null, fallback: T[]): T[] {
  return data ?? fallback
}

export default function App() {
  const [view, setView] = useState<View>('queue')
  const [initOpen, setInitOpen] = useState(false)

  const fetchTasks = useCallback(() => api.fetchTasks(), [])
  const fetchAgents = useCallback(() => api.fetchAgents(), [])
  const fetchCheckpoints = useCallback(() => api.fetchCheckpoints(), [])

  const { data: liveTasks, error: taskErr, refetch: refetchTasks } = usePolling(fetchTasks, 3000)
  const { data: liveAgents, error: agentErr } = usePolling(fetchAgents, 5000)
  const { data: liveCps, error: cpErr, refetch: refetchCps } = usePolling(fetchCheckpoints, 3000)

  const offline = !!(taskErr && agentErr && cpErr)
  const tasks = normalize<Task>(liveTasks, mockTasks)
  const _agents = normalize<Agent>(liveAgents, mockAgents)
  const checkpoints = normalize<Checkpoint>(liveCps, mockCheckpoints)

  return (
    <div className="flex h-screen">
      <Sidebar active={view} onNavigate={setView} onInit={() => setInitOpen(true)} />
      <main className="flex-1 overflow-auto p-6">
        {offline && (
          <div className="bg-amber-900/40 border border-amber-700 text-amber-300 text-xs px-3 py-2 rounded-md mb-4">
            ⚠ Engine offline — showing demo data
          </div>
        )}
        {view === 'queue' && <PriorityQueue tasks={tasks} refetch={refetchTasks} />}
        {view === 'kanban' && <KanbanBoard tasks={tasks} refetch={refetchTasks} />}
        {view === 'chat' && <ChatInterface />}
        {view === 'sessions' && <SessionsPage onResumeSession={(_id, _agent) => { setView('chat') }} />}
        {view === 'agents' && <AgentStatusPage />}
        {view === 'checkpoints' && <CheckpointReview checkpoints={checkpoints} refetch={refetchCps} />}
        {view === 'tickets' && <TicketingApp />}
        {view === 'wiki' && <WikiApp />}
        {view === 'settings' && <SettingsPage />}
      </main>
      <InitWizard open={initOpen} onClose={() => setInitOpen(false)} />
    </div>
  )
}

function _Placeholder({ title }: { title: string }) {
  return (
    <div className="flex items-center justify-center h-full">
      <p className="text-gray-500 text-lg">{title} — coming soon</p>
    </div>
  )
}
