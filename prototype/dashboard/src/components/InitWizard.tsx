import { useState } from 'react'
import { X, Rocket } from 'lucide-react'
import { api } from '../api'

interface Props { open: boolean; onClose: () => void }

export default function InitWizard({ open, onClose }: Props) {
  const [form, setForm] = useState({ company_name: '', description: '', languages: '', team_size: 5 })
  const [status, setStatus] = useState<'idle' | 'loading' | 'done' | 'error'>('idle')
  const [error, setError] = useState('')

  if (!open) return null

  const submit = async () => {
    setStatus('loading')
    try {
      await api.triggerInit({
        ...form,
        languages: form.languages.split(',').map(l => l.trim()).filter(Boolean),
      })
      setStatus('done')
    } catch (e: any) {
      setError(e.message)
      setStatus('error')
    }
  }

  return (
    <div className="fixed inset-0 bg-black/60 flex items-center justify-center z-50">
      <div className="bg-surface border border-gray-800 rounded-xl w-full max-w-md p-6">
        <div className="flex items-center justify-between mb-5">
          <h2 className="text-lg font-semibold text-gray-100 flex items-center gap-2">
            <Rocket size={18} className="text-accent" /> Initialize Company
          </h2>
          <button onClick={onClose} className="text-gray-500 hover:text-gray-300"><X size={18} /></button>
        </div>

        {status === 'done' ? (
          <p className="text-emerald-400 text-sm">✓ Initialization triggered. Tasks are being created.</p>
        ) : (
          <>
            <div className="flex flex-col gap-3 mb-4">
              <input value={form.company_name} onChange={e => setForm(f => ({ ...f, company_name: e.target.value }))}
                placeholder="Company name" className="bg-base border border-gray-800 rounded-md px-3 py-2 text-sm text-gray-200 placeholder-gray-600 focus:outline-none focus:border-accent/50" />
              <textarea value={form.description} onChange={e => setForm(f => ({ ...f, description: e.target.value }))}
                placeholder="What does this company build?" className="bg-base border border-gray-800 rounded-md px-3 py-2 text-sm text-gray-200 placeholder-gray-600 resize-none h-20 focus:outline-none focus:border-accent/50" />
              <input value={form.languages} onChange={e => setForm(f => ({ ...f, languages: e.target.value }))}
                placeholder="Languages (comma-separated, e.g. Go, Python)" className="bg-base border border-gray-800 rounded-md px-3 py-2 text-sm text-gray-200 placeholder-gray-600 focus:outline-none focus:border-accent/50" />
              <input type="number" value={form.team_size} onChange={e => setForm(f => ({ ...f, team_size: +e.target.value }))}
                placeholder="Team size" className="bg-base border border-gray-800 rounded-md px-3 py-2 text-sm text-gray-200 placeholder-gray-600 focus:outline-none focus:border-accent/50" />
            </div>
            {status === 'error' && <p className="text-red-400 text-xs mb-2">{error}</p>}
            <button onClick={submit} disabled={status === 'loading' || !form.company_name}
              className="w-full bg-accent hover:bg-accent/80 disabled:opacity-50 text-black font-medium text-sm py-2 rounded-md transition-colors">
              {status === 'loading' ? 'Initializing...' : 'Start Initialization'}
            </button>
          </>
        )}
      </div>
    </div>
  )
}
