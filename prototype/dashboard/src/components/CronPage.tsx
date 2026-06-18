import { useState, useEffect } from 'react';
import { Play, Pause, Zap, Plus, Clock, Trash2 } from 'lucide-react';

interface CronJob {
  id: string;
  name: string;
  schedule: string;
  agent: string;
  message: string;
  status: 'active' | 'paused';
  last_run?: string;
  next_run?: string;
}

const API = 'http://localhost:7780/api/v1/cron';

export default function CronPage() {
  const [jobs, setJobs] = useState<CronJob[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [form, setForm] = useState({ name: '', schedule: '', agent: '', message: '' });

  const fetchJobs = async () => {
    try {
      const res = await fetch(API);
      if (res.ok) setJobs(await res.json());
    } catch { /* ignore */ }
    setLoading(false);
  };

  useEffect(() => { fetchJobs(); }, []);

  const action = async (id: string, act: 'pause' | 'resume' | 'trigger' | 'delete') => {
    const method = act === 'delete' ? 'DELETE' : 'POST';
    const url = act === 'delete' ? `${API}/${id}` : `${API}/${id}/${act}`;
    await fetch(url, { method });
    fetchJobs();
  };

  const createJob = async (e: React.FormEvent) => {
    e.preventDefault();
    await fetch(API, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(form),
    });
    setForm({ name: '', schedule: '', agent: '', message: '' });
    setShowForm(false);
    fetchJobs();
  };

  if (loading) return <div className="p-6 text-gray-400">Loading cron jobs...</div>;

  return (
    <div className="p-6 max-w-5xl mx-auto">
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold text-white flex items-center gap-2">
          <Clock size={24} /> Cron Manager
        </h1>
        <button
          onClick={() => setShowForm(!showForm)}
          className="flex items-center gap-1 px-3 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-md text-sm"
        >
          <Plus size={16} /> New Job
        </button>
      </div>

      {showForm && (
        <form onSubmit={createJob} className="mb-6 p-4 bg-gray-800 rounded-lg border border-gray-700 space-y-3">
          <div className="grid grid-cols-2 gap-3">
            <input
              placeholder="Job name"
              value={form.name}
              onChange={e => setForm({ ...form, name: e.target.value })}
              required
              className="px-3 py-2 bg-gray-900 border border-gray-700 rounded text-sm text-white"
            />
            <input
              placeholder="Schedule (cron expr or interval)"
              value={form.schedule}
              onChange={e => setForm({ ...form, schedule: e.target.value })}
              required
              className="px-3 py-2 bg-gray-900 border border-gray-700 rounded text-sm text-white"
            />
            <input
              placeholder="Agent (optional)"
              value={form.agent}
              onChange={e => setForm({ ...form, agent: e.target.value })}
              className="px-3 py-2 bg-gray-900 border border-gray-700 rounded text-sm text-white"
            />
            <input
              placeholder="Message / prompt"
              value={form.message}
              onChange={e => setForm({ ...form, message: e.target.value })}
              required
              className="px-3 py-2 bg-gray-900 border border-gray-700 rounded text-sm text-white"
            />
          </div>
          <div className="flex gap-2">
            <button type="submit" className="px-4 py-2 bg-green-600 hover:bg-green-700 text-white rounded text-sm">Create</button>
            <button type="button" onClick={() => setShowForm(false)} className="px-4 py-2 bg-gray-700 hover:bg-gray-600 text-white rounded text-sm">Cancel</button>
          </div>
        </form>
      )}

      {jobs.length === 0 ? (
        <div className="text-center text-gray-500 py-12">No cron jobs configured</div>
      ) : (
        <div className="space-y-2">
          {jobs.map(job => (
            <div key={job.id} className="flex items-center justify-between p-4 bg-gray-800 rounded-lg border border-gray-700">
              <div className="flex-1 min-w-0">
                <div className="flex items-center gap-2">
                  <span className="font-medium text-white truncate">{job.name}</span>
                  <span className={`text-xs px-2 py-0.5 rounded ${job.status === 'active' ? 'bg-green-900 text-green-300' : 'bg-yellow-900 text-yellow-300'}`}>
                    {job.status}
                  </span>
                </div>
                <div className="text-xs text-gray-400 mt-1">
                  <span className="mr-4">⏱ {job.schedule}</span>
                  {job.agent && <span className="mr-4">🤖 {job.agent}</span>}
                  {job.next_run && <span>Next: {new Date(job.next_run).toLocaleString()}</span>}
                </div>
                <div className="text-xs text-gray-500 mt-0.5 truncate">{job.message}</div>
              </div>
              <div className="flex items-center gap-1 ml-4">
                {job.status === 'active' ? (
                  <button onClick={() => action(job.id, 'pause')} title="Pause" className="p-2 hover:bg-gray-700 rounded text-yellow-400"><Pause size={16} /></button>
                ) : (
                  <button onClick={() => action(job.id, 'resume')} title="Resume" className="p-2 hover:bg-gray-700 rounded text-green-400"><Play size={16} /></button>
                )}
                <button onClick={() => action(job.id, 'trigger')} title="Trigger now" className="p-2 hover:bg-gray-700 rounded text-blue-400"><Zap size={16} /></button>
                <button onClick={() => action(job.id, 'delete')} title="Delete" className="p-2 hover:bg-gray-700 rounded text-red-400"><Trash2 size={16} /></button>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
