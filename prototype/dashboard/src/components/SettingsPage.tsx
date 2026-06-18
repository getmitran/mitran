import { useState, useEffect } from 'react'
import { Save, FolderOpen, RefreshCw } from 'lucide-react'
import { api } from '../api'

interface ModulePaths {
  source?: string
  docs?: string
  wiki?: string
  sops?: string
  monitoring?: string
  cicd?: string
  tickets?: string
  hr?: string
}

interface Settings {
  root_dir: string
  projects_dir: string
  project_paths: Record<string, string>
  modules: ModulePaths
  engine_port: string
  worker_port: string
}

export default function SettingsPage() {
  const [settings, setSettings] = useState<Settings | null>(null)
  const [saving, setSaving] = useState(false)
  const [saved, setSaved] = useState(false)
  const [error, setError] = useState('')

  useEffect(() => { loadSettings() }, [])

  async function loadSettings() {
    try {
      const data = await api.fetchSettings()
      setSettings(data)
      setError('')
    } catch (e: any) {
      setError(e.message)
    }
  }

  async function save() {
    if (!settings) return
    setSaving(true)
    setSaved(false)
    try {
      await api.updateSettings(settings)
      setSaved(true)
      setTimeout(() => setSaved(false), 2000)
    } catch (e: any) {
      setError(e.message)
    } finally {
      setSaving(false)
    }
  }

  function update(field: string, value: string) {
    if (!settings) return
    setSettings({ ...settings, [field]: value })
  }

  function updateModule(key: keyof ModulePaths, value: string) {
    if (!settings) return
    setSettings({ ...settings, modules: { ...settings.modules, [key]: value } })
  }

  if (!settings) {
    return <div className="flex items-center justify-center h-full text-gray-500">Loading settings...</div>
  }

  return (
    <div className="max-w-3xl mx-auto space-y-8">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-white">Settings</h1>
        <div className="flex gap-2">
          <button onClick={loadSettings} className="flex items-center gap-1.5 px-3 py-1.5 text-xs rounded-md bg-white/5 border border-white/10 text-gray-300 hover:bg-white/10 transition">
            <RefreshCw size={12} /> Reload
          </button>
          <button onClick={save} disabled={saving} className="flex items-center gap-1.5 px-4 py-1.5 text-xs rounded-md bg-cyan-600 hover:bg-cyan-500 text-white font-medium transition disabled:opacity-50">
            <Save size={12} /> {saving ? 'Saving...' : saved ? 'Saved ✓' : 'Save'}
          </button>
        </div>
      </div>

      {error && <div className="bg-red-900/30 border border-red-700 text-red-300 text-xs px-3 py-2 rounded-md">{error}</div>}

      {/* Global Paths */}
      <section className="p-5 rounded-xl bg-[#141b2d] border border-white/5 space-y-4">
        <h2 className="text-sm font-semibold text-white">Global Paths</h2>
        <Field label="Root Directory" help="Base directory for all Mitran data" value={settings.root_dir} onChange={v => update('root_dir', v)} />
        <Field label="Projects Directory" help="Where project workspaces are created" value={settings.projects_dir} onChange={v => update('projects_dir', v)} />
      </section>

      {/* Module Paths */}
      <section className="p-5 rounded-xl bg-[#141b2d] border border-white/5 space-y-4">
        <h2 className="text-sm font-semibold text-white">Module Output Directories</h2>
        <p className="text-xs text-gray-500">Override where each module saves files. Leave blank for default (&lt;project&gt;/&lt;module&gt;/).</p>
        <div className="grid grid-cols-2 gap-3">
          <Field label="Source Code" value={settings.modules.source || ''} onChange={v => updateModule('source', v)} />
          <Field label="Documentation" value={settings.modules.docs || ''} onChange={v => updateModule('docs', v)} />
          <Field label="Wiki" value={settings.modules.wiki || ''} onChange={v => updateModule('wiki', v)} />
          <Field label="SOPs" value={settings.modules.sops || ''} onChange={v => updateModule('sops', v)} />
          <Field label="Monitoring" value={settings.modules.monitoring || ''} onChange={v => updateModule('monitoring', v)} />
          <Field label="CI/CD" value={settings.modules.cicd || ''} onChange={v => updateModule('cicd', v)} />
          <Field label="Tickets" value={settings.modules.tickets || ''} onChange={v => updateModule('tickets', v)} />
          <Field label="HR" value={settings.modules.hr || ''} onChange={v => updateModule('hr', v)} />
        </div>
      </section>

      {/* Per-Project Overrides */}
      <section className="p-5 rounded-xl bg-[#141b2d] border border-white/5 space-y-4">
        <h2 className="text-sm font-semibold text-white">Project Workspace Overrides</h2>
        <p className="text-xs text-gray-500">Custom workspace paths for specific projects.</p>
        {Object.keys(settings.project_paths || {}).length === 0 ? (
          <p className="text-xs text-gray-600 italic">No custom project paths configured.</p>
        ) : (
          <div className="space-y-2">
            {Object.entries(settings.project_paths).map(([id, path]) => (
              <div key={id} className="flex items-center gap-2">
                <span className="text-xs text-gray-400 font-mono w-28 truncate">{id}</span>
                <input value={path} readOnly className="flex-1 text-xs px-2.5 py-1.5 rounded-md bg-[#0a0e1a] border border-white/10 text-gray-300" />
              </div>
            ))}
          </div>
        )}
      </section>

      {/* Engine Settings */}
      <section className="p-5 rounded-xl bg-[#141b2d] border border-white/5 space-y-4">
        <h2 className="text-sm font-semibold text-white">Engine</h2>
        <div className="grid grid-cols-2 gap-3">
          <Field label="Engine Port" value={settings.engine_port} onChange={v => update('engine_port', v)} />
          <Field label="Worker Port" value={settings.worker_port} onChange={v => update('worker_port', v)} />
        </div>
      </section>
    </div>
  )
}

function Field({ label, help, value, onChange }: { label: string; help?: string; value: string; onChange?: (v: string) => void }) {
  return (
    <div>
      <label className="block text-[11px] text-gray-400 mb-1">{label}</label>
      <div className="flex items-center gap-2">
        <FolderOpen size={14} className="text-gray-600 shrink-0" />
        <input
          value={value}
          onChange={e => onChange?.(e.target.value)}
          placeholder={help || 'Default'}
          className="flex-1 text-sm px-2.5 py-2 rounded-md bg-[#0a0e1a] border border-white/10 text-gray-200 placeholder-gray-600 focus:outline-none focus:border-cyan-500/50"
        />
      </div>
    </div>
  )
}
