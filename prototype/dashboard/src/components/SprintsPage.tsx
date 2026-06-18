import { useState, useEffect } from 'react';

interface Sprint {
  id: string;
  name: string;
  start_date: string;
  end_date: string;
  status: 'active' | 'completed' | 'planned';
  task_count?: number;
}

export default function SprintsPage() {
  const [sprints, setSprints] = useState<Sprint[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [form, setForm] = useState({ name: '', start_date: '', end_date: '' });

  const fetchSprints = () => {
    setLoading(true);
    fetch('http://localhost:7780/api/v1/sprints')
      .then(r => r.json())
      .then(data => setSprints(Array.isArray(data) ? data : data.sprints || []))
      .catch(() => setSprints([]))
      .finally(() => setLoading(false));
  };

  useEffect(() => { fetchSprints(); }, []);

  const createSprint = (e: React.FormEvent) => {
    e.preventDefault();
    fetch('http://localhost:7780/api/v1/sprints', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(form),
    }).then(() => {
      setForm({ name: '', start_date: '', end_date: '' });
      setShowForm(false);
      fetchSprints();
    });
  };

  const statusColor = (s: string) =>
    s === 'active' ? '#22c55e' : s === 'completed' ? '#6b7280' : '#3b82f6';

  if (loading) return <div style={{ padding: 24 }}>Loading sprints...</div>;

  return (
    <div style={{ padding: 24, maxWidth: 900, margin: '0 auto' }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 20 }}>
        <h1 style={{ fontSize: 24, fontWeight: 700 }}>Sprints</h1>
        <button
          onClick={() => setShowForm(!showForm)}
          style={{ padding: '8px 16px', background: '#3b82f6', color: '#fff', border: 'none', borderRadius: 6, cursor: 'pointer' }}
        >
          {showForm ? 'Cancel' : '+ New Sprint'}
        </button>
      </div>

      {showForm && (
        <form onSubmit={createSprint} style={{ marginBottom: 20, padding: 16, border: '1px solid #e5e7eb', borderRadius: 8 }}>
          <div style={{ display: 'flex', gap: 12, flexWrap: 'wrap' }}>
            <input
              placeholder="Sprint name"
              value={form.name}
              onChange={e => setForm({ ...form, name: e.target.value })}
              required
              style={{ flex: 1, minWidth: 180, padding: '8px 12px', border: '1px solid #d1d5db', borderRadius: 6 }}
            />
            <input
              type="date"
              value={form.start_date}
              onChange={e => setForm({ ...form, start_date: e.target.value })}
              required
              style={{ padding: '8px 12px', border: '1px solid #d1d5db', borderRadius: 6 }}
            />
            <input
              type="date"
              value={form.end_date}
              onChange={e => setForm({ ...form, end_date: e.target.value })}
              required
              style={{ padding: '8px 12px', border: '1px solid #d1d5db', borderRadius: 6 }}
            />
            <button type="submit" style={{ padding: '8px 16px', background: '#22c55e', color: '#fff', border: 'none', borderRadius: 6, cursor: 'pointer' }}>
              Create
            </button>
          </div>
        </form>
      )}

      {sprints.length === 0 ? (
        <p style={{ color: '#6b7280' }}>No sprints found. Create one to get started.</p>
      ) : (
        <div style={{ display: 'flex', flexDirection: 'column', gap: 12 }}>
          {sprints.map(sprint => (
            <div key={sprint.id} style={{ padding: 16, border: '1px solid #e5e7eb', borderRadius: 8, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
              <div>
                <div style={{ fontWeight: 600, fontSize: 16 }}>{sprint.name}</div>
                <div style={{ fontSize: 13, color: '#6b7280', marginTop: 4 }}>
                  {sprint.start_date} → {sprint.end_date}
                </div>
              </div>
              <div style={{ display: 'flex', gap: 12, alignItems: 'center' }}>
                {sprint.task_count !== undefined && (
                  <span style={{ fontSize: 13, color: '#6b7280' }}>{sprint.task_count} tasks</span>
                )}
                <span style={{
                  padding: '4px 10px',
                  borderRadius: 12,
                  fontSize: 12,
                  fontWeight: 600,
                  color: '#fff',
                  background: statusColor(sprint.status),
                }}>
                  {sprint.status}
                </span>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
