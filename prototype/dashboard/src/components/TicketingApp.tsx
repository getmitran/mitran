import { useState, useEffect, useCallback } from 'react'

const BASE = 'http://localhost:7780/api/v1'

interface Ticket {
  id: string; title: string; description: string; status: string
  priority: string; assignee: string; labels: string[]
  sprint_id: string; sla_due_at?: string; created_by: string
  created_at: string; updated_at: string
}
interface Sprint { id: string; name: string; start_date: string; end_date: string }
interface Comment { id: string; ticket_id: string; author: string; body: string; created_at: string }

const COLUMNS = ['open', 'in_progress', 'review', 'done'] as const
const PRIORITIES = ['critical', 'high', 'medium', 'low'] as const
const PRIORITY_COLORS: Record<string, string> = {
  critical: '#ef4444', high: '#f97316', medium: '#eab308', low: '#22c55e'
}

async function req<T>(path: string, opts?: RequestInit): Promise<T> {
  const res = await fetch(`${BASE}${path}`, { headers: { 'Content-Type': 'application/json' }, ...opts })
  if (!res.ok) throw new Error(`${res.status}`)
  return res.json()
}

export default function TicketingApp() {
  const [tickets, setTickets] = useState<Ticket[]>([])
  const [sprints, setSprints] = useState<Sprint[]>([])
  const [selectedSprint, setSelectedSprint] = useState('')
  const [filterAssignee, setFilterAssignee] = useState('')
  const [filterPriority, setFilterPriority] = useState('')
  const [filterLabel, setFilterLabel] = useState('')
  const [selected, setSelected] = useState<Ticket | null>(null)
  const [comments, setComments] = useState<Comment[]>([])
  const [showNew, setShowNew] = useState(false)
  const [dragId, setDragId] = useState<string | null>(null)

  const load = useCallback(async () => {
    const params = new URLSearchParams()
    if (selectedSprint) params.set('sprint_id', selectedSprint)
    if (filterAssignee) params.set('assignee', filterAssignee)
    if (filterLabel) params.set('labels', filterLabel)
    const t = await req<Ticket[]>(`/tickets?${params}`)
    setTickets(t)
    setSprints(await req<Sprint[]>('/sprints'))
  }, [selectedSprint, filterAssignee, filterLabel])

  useEffect(() => { load(); const i = setInterval(load, 5000); return () => clearInterval(i) }, [load])

  const updateStatus = async (id: string, status: string) => {
    await req(`/tickets/${id}`, { method: 'PUT', body: JSON.stringify({ status }) })
    load()
  }

  const openDetail = async (t: Ticket) => {
    setSelected(t)
    setComments(await req<Comment[]>(`/tickets/${t.id}/comments`))
  }

  const filtered = tickets.filter(t => !filterPriority || t.priority === filterPriority)

  const slaCountdown = (due?: string) => {
    if (!due) return null
    const diff = new Date(due).getTime() - Date.now()
    if (diff < 0) return <span style={{ color: '#ef4444', fontSize: 11 }}>OVERDUE</span>
    const h = Math.floor(diff / 3600000)
    return <span style={{ color: h < 4 ? '#f97316' : '#6b7280', fontSize: 11 }}>{h}h left</span>
  }

  return (
    <div style={{ padding: 24, background: '#0a0e1a', minHeight: '100vh', color: '#e2e8f0' }}>
      {/* Header */}
      <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 16 }}>
        <h1 style={{ fontSize: 20, fontWeight: 700 }}>Ticketing</h1>
        <button onClick={() => setShowNew(true)} style={btnStyle}>+ New Ticket</button>
      </div>

      {/* Filters */}
      <div style={{ display: 'flex', gap: 12, marginBottom: 16, flexWrap: 'wrap' }}>
        <select value={selectedSprint} onChange={e => setSelectedSprint(e.target.value)} style={selStyle}>
          <option value="">All Sprints</option>
          {sprints.map(s => <option key={s.id} value={s.id}>{s.name}</option>)}
        </select>
        <input placeholder="Filter assignee" value={filterAssignee} onChange={e => setFilterAssignee(e.target.value)} style={inputStyle} />
        <select value={filterPriority} onChange={e => setFilterPriority(e.target.value)} style={selStyle}>
          <option value="">All Priorities</option>
          {PRIORITIES.map(p => <option key={p} value={p}>{p}</option>)}
        </select>
        <input placeholder="Filter label" value={filterLabel} onChange={e => setFilterLabel(e.target.value)} style={inputStyle} />
      </div>

      {/* Board */}
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(4, 1fr)', gap: 16 }}>
        {COLUMNS.map(col => (
          <div key={col}
            onDragOver={e => e.preventDefault()}
            onDrop={() => { if (dragId) { updateStatus(dragId, col); setDragId(null) } }}
            style={{ background: '#141b2d', borderRadius: 8, padding: 12, minHeight: 300 }}>
            <div style={{ fontWeight: 600, textTransform: 'capitalize', marginBottom: 12, fontSize: 13, color: '#94a3b8' }}>
              {col.replace('_', ' ')} ({filtered.filter(t => t.status === col).length})
            </div>
            {filtered.filter(t => t.status === col).map(t => (
              <div key={t.id} draggable onDragStart={() => setDragId(t.id)} onClick={() => openDetail(t)}
                style={{ background: '#1e293b', borderRadius: 6, padding: 10, marginBottom: 8, cursor: 'pointer', border: '1px solid #334155' }}>
                <div style={{ fontSize: 13, fontWeight: 500, marginBottom: 4 }}>{t.title}</div>
                <div style={{ display: 'flex', gap: 6, alignItems: 'center', flexWrap: 'wrap' }}>
                  <span style={{ background: PRIORITY_COLORS[t.priority], color: '#fff', fontSize: 10, padding: '1px 6px', borderRadius: 4 }}>{t.priority}</span>
                  {t.assignee && <span style={{ fontSize: 11, color: '#94a3b8' }}>@{t.assignee}</span>}
                  {t.labels?.map(l => <span key={l} style={{ fontSize: 10, background: '#0e7490', color: '#fff', padding: '1px 5px', borderRadius: 3 }}>{l}</span>)}
                  {slaCountdown(t.sla_due_at)}
                </div>
              </div>
            ))}
          </div>
        ))}
      </div>

      {/* Detail Panel */}
      {selected && <DetailPanel ticket={selected} comments={comments} onClose={() => setSelected(null)} onSave={async (updates) => {
        await req(`/tickets/${selected.id}`, { method: 'PUT', body: JSON.stringify(updates) })
        load(); setSelected(null)
      }} onComment={async (body) => {
        await req(`/tickets/${selected.id}/comments`, { method: 'POST', body: JSON.stringify({ author: 'user', body }) })
        setComments(await req<Comment[]>(`/tickets/${selected.id}/comments`))
      }} />}

      {/* New Ticket Modal */}
      {showNew && <NewTicketModal sprints={sprints} onClose={() => setShowNew(false)} onCreate={async (t) => {
        await req('/tickets', { method: 'POST', body: JSON.stringify(t) })
        setShowNew(false); load()
      }} />}
    </div>
  )
}

function DetailPanel({ ticket, comments, onClose, onSave, onComment }: {
  ticket: Ticket; comments: Comment[]; onClose: () => void
  onSave: (u: Partial<Ticket>) => void; onComment: (body: string) => void
}) {
  const [form, setForm] = useState(ticket)
  const [comment, setComment] = useState('')
  return (
    <div style={{ position: 'fixed', top: 0, right: 0, width: 420, height: '100vh', background: '#141b2d', borderLeft: '1px solid #334155', padding: 20, overflow: 'auto', zIndex: 50 }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 16 }}>
        <h2 style={{ fontSize: 16, fontWeight: 600 }}>Edit Ticket</h2>
        <button onClick={onClose} style={{ background: 'none', border: 'none', color: '#94a3b8', cursor: 'pointer', fontSize: 18 }}>✕</button>
      </div>
      <label style={labelStyle}>Title</label>
      <input value={form.title} onChange={e => setForm({ ...form, title: e.target.value })} style={inputStyle} />
      <label style={labelStyle}>Description</label>
      <textarea value={form.description} onChange={e => setForm({ ...form, description: e.target.value })} style={{ ...inputStyle, height: 80 }} />
      <label style={labelStyle}>Status</label>
      <select value={form.status} onChange={e => setForm({ ...form, status: e.target.value })} style={selStyle}>
        {COLUMNS.map(c => <option key={c} value={c}>{c}</option>)}
        <option value="closed">closed</option>
      </select>
      <label style={labelStyle}>Priority</label>
      <select value={form.priority} onChange={e => setForm({ ...form, priority: e.target.value })} style={selStyle}>
        {PRIORITIES.map(p => <option key={p} value={p}>{p}</option>)}
      </select>
      <label style={labelStyle}>Assignee</label>
      <input value={form.assignee} onChange={e => setForm({ ...form, assignee: e.target.value })} style={inputStyle} />
      <label style={labelStyle}>Labels (comma-separated)</label>
      <input value={form.labels?.join(', ')} onChange={e => setForm({ ...form, labels: e.target.value.split(',').map(s => s.trim()).filter(Boolean) })} style={inputStyle} />
      <button onClick={() => onSave({ title: form.title, description: form.description, status: form.status, priority: form.priority, assignee: form.assignee, labels: form.labels })} style={{ ...btnStyle, marginTop: 12, width: '100%' }}>Save</button>

      <div style={{ marginTop: 20, borderTop: '1px solid #334155', paddingTop: 12 }}>
        <h3 style={{ fontSize: 13, fontWeight: 600, marginBottom: 8 }}>Comments</h3>
        {comments.map(c => (
          <div key={c.id} style={{ marginBottom: 8, fontSize: 12 }}>
            <span style={{ color: '#06b6d4' }}>@{c.author}</span> <span style={{ color: '#64748b' }}>{new Date(c.created_at).toLocaleString()}</span>
            <div style={{ marginTop: 2 }}>{c.body}</div>
          </div>
        ))}
        <div style={{ display: 'flex', gap: 8, marginTop: 8 }}>
          <input value={comment} onChange={e => setComment(e.target.value)} placeholder="Add comment..." style={{ ...inputStyle, flex: 1 }} />
          <button onClick={() => { if (comment) { onComment(comment); setComment('') } }} style={btnStyle}>Post</button>
        </div>
      </div>
    </div>
  )
}

function NewTicketModal({ sprints, onClose, onCreate }: { sprints: Sprint[]; onClose: () => void; onCreate: (t: Partial<Ticket>) => void }) {
  const [form, setForm] = useState<Partial<Ticket>>({ status: 'open', priority: 'medium', labels: [] })
  return (
    <div style={{ position: 'fixed', inset: 0, background: 'rgba(0,0,0,.6)', display: 'flex', alignItems: 'center', justifyContent: 'center', zIndex: 60 }}>
      <div style={{ background: '#141b2d', borderRadius: 10, padding: 24, width: 400, border: '1px solid #334155' }}>
        <h2 style={{ fontSize: 16, fontWeight: 600, marginBottom: 16 }}>New Ticket</h2>
        <label style={labelStyle}>Title *</label>
        <input value={form.title || ''} onChange={e => setForm({ ...form, title: e.target.value })} style={inputStyle} />
        <label style={labelStyle}>Description</label>
        <textarea value={form.description || ''} onChange={e => setForm({ ...form, description: e.target.value })} style={{ ...inputStyle, height: 60 }} />
        <label style={labelStyle}>Priority</label>
        <select value={form.priority} onChange={e => setForm({ ...form, priority: e.target.value })} style={selStyle}>
          {PRIORITIES.map(p => <option key={p} value={p}>{p}</option>)}
        </select>
        <label style={labelStyle}>Assignee</label>
        <input value={form.assignee || ''} onChange={e => setForm({ ...form, assignee: e.target.value })} style={inputStyle} />
        <label style={labelStyle}>Labels (comma-separated)</label>
        <input value={form.labels?.join(', ')} onChange={e => setForm({ ...form, labels: e.target.value.split(',').map(s => s.trim()).filter(Boolean) })} style={inputStyle} />
        <label style={labelStyle}>Sprint</label>
        <select value={form.sprint_id || ''} onChange={e => setForm({ ...form, sprint_id: e.target.value })} style={selStyle}>
          <option value="">None</option>
          {sprints.map(s => <option key={s.id} value={s.id}>{s.name}</option>)}
        </select>
        <div style={{ display: 'flex', gap: 8, marginTop: 16 }}>
          <button onClick={() => { if (form.title) onCreate(form) }} style={{ ...btnStyle, flex: 1 }}>Create</button>
          <button onClick={onClose} style={{ ...btnStyle, flex: 1, background: '#334155' }}>Cancel</button>
        </div>
      </div>
    </div>
  )
}

const btnStyle: React.CSSProperties = { background: '#06b6d4', color: '#fff', border: 'none', borderRadius: 6, padding: '6px 14px', cursor: 'pointer', fontSize: 13, fontWeight: 500 }
const inputStyle: React.CSSProperties = { width: '100%', background: '#0f172a', border: '1px solid #334155', borderRadius: 6, padding: '6px 10px', color: '#e2e8f0', fontSize: 13, marginBottom: 8 }
const selStyle: React.CSSProperties = { ...inputStyle, cursor: 'pointer' }
const labelStyle: React.CSSProperties = { fontSize: 11, color: '#94a3b8', display: 'block', marginBottom: 2 }
