import { useState } from 'react'

interface MobileNavProps {
  currentView: string
  onNavigate: (view: string) => void
}

const NAV_ITEMS = [
  { id: 'queue', label: 'Queue', icon: '📋' },
  { id: 'kanban', label: 'Kanban', icon: '📊' },
  { id: 'chat', label: 'Chat', icon: '💬' },
  { id: 'agents', label: 'Agents', icon: '🤖' },
  { id: 'tickets', label: 'Tickets', icon: '🎫' },
  { id: 'wiki', label: 'Wiki', icon: '📖' },
  { id: 'sessions', label: 'Sessions', icon: '💾' },
  { id: 'settings', label: 'Settings', icon: '⚙️' },
]

export function MobileNav({ currentView, onNavigate }: MobileNavProps) {
  const [open, setOpen] = useState(false)

  return (
    <div className="md:hidden">
      <button onClick={() => setOpen(!open)}
        className="fixed top-3 left-3 z-50 p-2 bg-gray-800 rounded-md text-white">
        {open ? '✕' : '☰'}
      </button>
      {open && (
        <div className="fixed inset-0 z-40 bg-gray-900/95 flex flex-col items-center justify-center gap-4">
          {NAV_ITEMS.map(item => (
            <button key={item.id}
              onClick={() => { onNavigate(item.id); setOpen(false) }}
              className={`text-lg px-6 py-3 rounded-lg w-48 text-left ${
                currentView === item.id ? 'bg-blue-600 text-white' : 'text-gray-300 hover:bg-gray-800'
              }`}>
              {item.icon} {item.label}
            </button>
          ))}
        </div>
      )}
    </div>
  )
}
