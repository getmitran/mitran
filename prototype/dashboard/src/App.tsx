import { useState } from 'react'
import Sidebar from './components/Sidebar'
import PriorityQueue from './components/PriorityQueue'
import AgentStatus from './components/AgentStatus'
import CheckpointReview from './components/CheckpointReview'

type View = 'queue' | 'agents' | 'checkpoints' | 'tickets' | 'wiki'

export default function App() {
  const [view, setView] = useState<View>('queue')

  return (
    <div className="flex h-screen">
      <Sidebar active={view} onNavigate={setView} />
      <main className="flex-1 overflow-auto p-6">
        {view === 'queue' && <PriorityQueue />}
        {view === 'agents' && <AgentStatus />}
        {view === 'checkpoints' && <CheckpointReview />}
        {view === 'tickets' && <Placeholder title="Tickets" />}
        {view === 'wiki' && <Placeholder title="Wiki" />}
      </main>
    </div>
  )
}

function Placeholder({ title }: { title: string }) {
  return (
    <div className="flex items-center justify-center h-full">
      <p className="text-gray-500 text-lg">{title} — coming soon</p>
    </div>
  )
}
