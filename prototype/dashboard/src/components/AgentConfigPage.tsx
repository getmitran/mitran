import { useState, useEffect } from 'react';

interface AgentConfig {
  name: string;
  model: string;
  temperature: number;
  system_prompt: string;
  tools: string[];
}

const MODELS = ['claude-sonnet-4-20250514', 'claude-haiku-4-20250414', 'claude-opus-4-20250514', 'gpt-4o', 'gpt-4o-mini'];

export default function AgentConfigPage() {
  const [agents, setAgents] = useState<AgentConfig[]>([]);
  const [selected, setSelected] = useState<AgentConfig | null>(null);
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    fetch('http://localhost:7780/api/v1/agents/config')
      .then(r => r.json())
      .then(data => setAgents(Array.isArray(data) ? data : data.agents || []))
      .catch(() => setAgents([]));
  }, []);

  const save = async () => {
    if (!selected) return;
    setSaving(true);
    await fetch('http://localhost:7780/api/v1/agents/config', {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(selected),
    });
    setAgents(prev => prev.map(a => a.name === selected.name ? selected : a));
    setSaving(false);
  };

  return (
    <div className="flex h-full">
      <div className="w-72 border-r border-gray-700 overflow-y-auto">
        <h2 className="p-4 text-lg font-semibold">Agents</h2>
        {agents.map(a => (
          <div
            key={a.name}
            onClick={() => setSelected({ ...a })}
            className={`px-4 py-3 cursor-pointer border-b border-gray-800 hover:bg-gray-800 ${selected?.name === a.name ? 'bg-gray-800' : ''}`}
          >
            <div className="font-medium">{a.name}</div>
            <div className="text-xs text-gray-400">{a.model} · t={a.temperature}</div>
            <div className="text-xs text-gray-500">{a.tools.length} tools</div>
          </div>
        ))}
      </div>

      <div className="flex-1 p-6 overflow-y-auto">
        {selected ? (
          <div className="max-w-2xl space-y-5">
            <h2 className="text-xl font-semibold">{selected.name}</h2>

            <div>
              <label className="block text-sm text-gray-400 mb-1">Model</label>
              <select
                value={selected.model}
                onChange={e => setSelected({ ...selected, model: e.target.value })}
                className="w-full bg-gray-800 border border-gray-600 rounded px-3 py-2 text-sm"
              >
                {MODELS.map(m => <option key={m} value={m}>{m}</option>)}
              </select>
            </div>

            <div>
              <label className="block text-sm text-gray-400 mb-1">
                Temperature: {selected.temperature.toFixed(2)}
              </label>
              <input
                type="range"
                min="0"
                max="1"
                step="0.05"
                value={selected.temperature}
                onChange={e => setSelected({ ...selected, temperature: parseFloat(e.target.value) })}
                className="w-full"
              />
            </div>

            <div>
              <label className="block text-sm text-gray-400 mb-1">System Prompt</label>
              <textarea
                value={selected.system_prompt}
                onChange={e => setSelected({ ...selected, system_prompt: e.target.value })}
                rows={10}
                className="w-full bg-gray-800 border border-gray-600 rounded px-3 py-2 text-sm font-mono"
              />
            </div>

            <div>
              <label className="block text-sm text-gray-400 mb-1">Tools</label>
              <div className="flex flex-wrap gap-1">
                {selected.tools.map(t => (
                  <span key={t} className="px-2 py-0.5 bg-gray-700 rounded text-xs">{t}</span>
                ))}
              </div>
            </div>

            <button
              onClick={save}
              disabled={saving}
              className="px-4 py-2 bg-blue-600 hover:bg-blue-700 rounded text-sm font-medium disabled:opacity-50"
            >
              {saving ? 'Saving...' : 'Save'}
            </button>
          </div>
        ) : (
          <div className="text-gray-500 mt-20 text-center">Select an agent to configure</div>
        )}
      </div>
    </div>
  );
}
