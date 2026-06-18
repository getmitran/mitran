import { useState, useEffect } from 'react'

interface Message {
  id: string
  role: 'user' | 'assistant'
  content: string
  agent?: string
  timestamp: string
}

interface Session {
  id: string
  user_id: string
  agent: string
  messages: Message[]
  created_at: string
  updated_at: string
  active: boolean
}

export function SessionsPage({ onResumeSession }: { onResumeSession?: (sessionId: string, agent: string) => void }) {
  const [sessions, setSessions] = useState<Session[]>([])
  const [selected, setSelected] = useState<Session | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetch('http://localhost:7780/api/v1/sessions').then(r => r.json()).then(r => {
      setSessions(r.data || [])
      setLoading(false)
    }).catch(() => setLoading(false))
  }, [])

  if (loading) return <div className="p-6 text-gray-400">Loading sessions...</div>

  return (
    <div className="flex h-full">
      {/* Session list */}
      <div className="w-80 border-r border-gray-700 overflow-y-auto">
        <div className="p-4 border-b border-gray-700">
          <h2 className="text-lg font-semibold text-white">Chat Sessions</h2>
          <p className="text-sm text-gray-400">{sessions.length} sessions</p>
        </div>
        {sessions.map(s => (
          <div
            key={s.id}
            onClick={() => setSelected(s)}
            className={`p-3 border-b border-gray-800 cursor-pointer hover:bg-gray-800 ${
              selected?.id === s.id ? 'bg-gray-800' : ''
            }`}
          >
            <div className="flex items-center justify-between">
              <span className="text-sm font-medium text-white">{s.agent}</span>
              <span className={`text-xs px-1.5 py-0.5 rounded ${s.active ? 'bg-green-900 text-green-300' : 'bg-gray-700 text-gray-400'}`}>
                {s.active ? 'Active' : 'Ended'}
              </span>
            </div>
            <div className="text-xs text-gray-500 mt-1">
              {new Date(s.updated_at).toLocaleString('en-IN', { timeZone: 'Asia/Kolkata' })}
            </div>
            <div className="text-xs text-gray-400 mt-1 truncate">
              {s.messages?.[s.messages.length - 1]?.content?.slice(0, 60) || 'No messages'}
            </div>
          </div>
        ))}
      </div>

      {/* Session detail */}
      <div className="flex-1 flex flex-col">
        {selected ? (
          <>
            <div className="p-4 border-b border-gray-700 flex items-center justify-between">
              <div>
                <h3 className="text-white font-medium">{selected.agent} -- {selected.id}</h3>
                <p className="text-xs text-gray-500">{selected.messages?.length || 0} messages</p>
              </div>
              {selected.active && onResumeSession && (
                <button
                  onClick={() => onResumeSession(selected.id, selected.agent)}
                  className="px-3 py-1.5 bg-blue-600 text-white text-sm rounded hover:bg-blue-700"
                >
                  Resume
                </button>
              )}
            </div>
            <div className="flex-1 overflow-y-auto p-4 space-y-3">
              {(selected.messages || []).map(m => (
                <div key={m.id} className={`flex ${m.role === 'user' ? 'justify-end' : 'justify-start'}`}>
                  <div className={`max-w-[70%] px-3 py-2 rounded-lg text-sm ${
                    m.role === 'user' ? 'bg-blue-600 text-white' : 'bg-gray-800 text-gray-200'
                  }`}>
                    {m.content}
                  </div>
                </div>
              ))}
            </div>
          </>
        ) : (
          <div className="flex-1 flex items-center justify-center text-gray-500">
            Select a session to view messages
          </div>
        )}
      </div>
    </div>
  )
}
