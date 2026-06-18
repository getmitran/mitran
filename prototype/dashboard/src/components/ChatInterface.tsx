import { useState, useRef, useEffect } from 'react'
import { Send, Bot, User, Loader2 } from 'lucide-react'
import { api } from '../api'

interface Message {
  id: string
  role: 'user' | 'assistant'
  content: string
  agent?: string
  timestamp: string
}

export default function ChatInterface() {
  const [messages, setMessages] = useState<Message[]>([])
  const [input, setInput] = useState('')
  const [loading, setLoading] = useState(false)
  const [agent, setAgent] = useState('dev')
  const bottomRef = useRef<HTMLDivElement>(null)

  useEffect(() => { bottomRef.current?.scrollIntoView({ behavior: 'smooth' }) }, [messages])

  async function send() {
    if (!input.trim() || loading) return
    const userMsg: Message = { id: Date.now().toString(), role: 'user', content: input, timestamp: new Date().toISOString() }
    setMessages(prev => [...prev, userMsg])
    setInput('')
    setLoading(true)

    try {
      const res = await api.chat(input, agent)
      const assistantMsg: Message = {
        id: (Date.now() + 1).toString(),
        role: 'assistant',
        content: res.response || res.summary || 'No response',
        agent,
        timestamp: new Date().toISOString(),
      }
      setMessages(prev => [...prev, assistantMsg])
    } catch (e: any) {
      setMessages(prev => [...prev, { id: (Date.now() + 1).toString(), role: 'assistant', content: `Error: ${e.message}`, timestamp: new Date().toISOString() }])
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="flex flex-col h-full">
      {/* Header */}
      <div className="flex items-center gap-3 px-4 py-3 border-b border-white/5">
        <h2 className="text-sm font-semibold text-white">Chat</h2>
        <select value={agent} onChange={e => setAgent(e.target.value)}
          className="text-xs px-2 py-1 rounded bg-[#141b2d] border border-white/10 text-gray-300">
          <option value="dev">Dev Agent</option>
          <option value="docs">Docs Agent</option>
          <option value="ops">Ops Agent</option>
          <option value="cicd">CI/CD Agent</option>
          <option value="tickets">Tickets Agent</option>
          <option value="wiki">Wiki Agent</option>
          <option value="hr">HR Agent</option>
          <option value="review">Review Agent</option>
        </select>
      </div>

      {/* Messages */}
      <div className="flex-1 overflow-y-auto px-4 py-4 space-y-4">
        {messages.length === 0 && (
          <div className="flex items-center justify-center h-full text-gray-600 text-sm">
            Start a conversation with any agent...
          </div>
        )}
        {messages.map(msg => (
          <div key={msg.id} className={`flex gap-3 ${msg.role === 'user' ? 'justify-end' : ''}`}>
            {msg.role === 'assistant' && (
              <div className="w-7 h-7 rounded-full bg-cyan-500/20 flex items-center justify-center shrink-0">
                <Bot size={14} className="text-cyan-400" />
              </div>
            )}
            <div className={`max-w-[70%] px-3 py-2 rounded-lg text-sm ${
              msg.role === 'user'
                ? 'bg-cyan-600/20 text-gray-200'
                : 'bg-[#141b2d] border border-white/5 text-gray-300'
            }`}>
              <pre className="whitespace-pre-wrap font-sans">{msg.content}</pre>
              <div className="text-[10px] text-gray-600 mt-1">
                {msg.agent && <span className="text-cyan-500">{msg.agent} • </span>}
                {new Date(msg.timestamp).toLocaleTimeString()}
              </div>
            </div>
            {msg.role === 'user' && (
              <div className="w-7 h-7 rounded-full bg-white/10 flex items-center justify-center shrink-0">
                <User size={14} className="text-gray-400" />
              </div>
            )}
          </div>
        ))}
        {loading && (
          <div className="flex gap-3">
            <div className="w-7 h-7 rounded-full bg-cyan-500/20 flex items-center justify-center">
              <Loader2 size={14} className="text-cyan-400 animate-spin" />
            </div>
            <div className="px-3 py-2 rounded-lg bg-[#141b2d] border border-white/5 text-gray-500 text-sm">
              Thinking...
            </div>
          </div>
        )}
        <div ref={bottomRef} />
      </div>

      {/* Input */}
      <div className="px-4 py-3 border-t border-white/5">
        <div className="flex gap-2">
          <input
            value={input}
            onChange={e => setInput(e.target.value)}
            onKeyDown={e => e.key === 'Enter' && !e.shiftKey && send()}
            placeholder={`Message ${agent} agent...`}
            className="flex-1 px-3 py-2.5 rounded-lg bg-[#141b2d] border border-white/10 text-sm text-gray-200 placeholder-gray-600 focus:outline-none focus:border-cyan-500/50"
          />
          <button onClick={send} disabled={loading || !input.trim()}
            className="px-3 py-2.5 rounded-lg bg-cyan-600 hover:bg-cyan-500 disabled:opacity-30 transition">
            <Send size={16} className="text-white" />
          </button>
        </div>
      </div>
    </div>
  )
}
