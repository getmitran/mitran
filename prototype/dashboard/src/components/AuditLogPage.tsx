import { useState, useEffect } from 'react';

interface AuditEvent {
  timestamp: string;
  actor: string;
  action: string;
  resource: string;
}

const SAMPLE_EVENTS: AuditEvent[] = [
  { timestamp: '2026-06-18T19:00:00Z', actor: 'admin', action: 'create', resource: 'project/mitran-core' },
  { timestamp: '2026-06-18T18:45:00Z', actor: 'ci-bot', action: 'deploy', resource: 'service/agent-worker' },
  { timestamp: '2026-06-18T18:30:00Z', actor: 'ravitejb', action: 'update', resource: 'config/settings.json' },
  { timestamp: '2026-06-18T18:15:00Z', actor: 'system', action: 'restart', resource: 'service/engine' },
  { timestamp: '2026-06-18T18:00:00Z', actor: 'ravitejb', action: 'delete', resource: 'task/MITRAN-42' },
  { timestamp: '2026-06-18T17:30:00Z', actor: 'admin', action: 'invite', resource: 'user/newdev' },
];

export default function AuditLogPage() {
  const [events, setEvents] = useState<AuditEvent[]>([]);
  const [filter, setFilter] = useState('');

  useEffect(() => {
    fetch('http://localhost:7780/api/v1/audit')
      .then(r => r.ok ? r.json() : Promise.reject())
      .then(data => setEvents(data))
      .catch(() => setEvents(SAMPLE_EVENTS));
  }, []);

  const filtered = events.filter(e =>
    `${e.timestamp} ${e.actor} ${e.action} ${e.resource}`.toLowerCase().includes(filter.toLowerCase())
  );

  return (
    <div style={{ padding: '1.5rem' }}>
      <h1 style={{ fontSize: '1.5rem', fontWeight: 600, marginBottom: '1rem' }}>Audit Log</h1>
      <input
        type="text"
        placeholder="Filter events..."
        value={filter}
        onChange={e => setFilter(e.target.value)}
        style={{ width: '100%', maxWidth: 400, padding: '0.5rem', marginBottom: '1rem', border: '1px solid #ccc', borderRadius: 4 }}
      />
      <table style={{ width: '100%', borderCollapse: 'collapse', fontSize: '0.875rem' }}>
        <thead>
          <tr style={{ borderBottom: '2px solid #ddd', textAlign: 'left' }}>
            <th style={{ padding: '0.5rem' }}>Timestamp</th>
            <th style={{ padding: '0.5rem' }}>Actor</th>
            <th style={{ padding: '0.5rem' }}>Action</th>
            <th style={{ padding: '0.5rem' }}>Resource</th>
          </tr>
        </thead>
        <tbody>
          {filtered.map((e, i) => (
            <tr key={i} style={{ borderBottom: '1px solid #eee' }}>
              <td style={{ padding: '0.5rem' }}>{new Date(e.timestamp).toLocaleString()}</td>
              <td style={{ padding: '0.5rem' }}>{e.actor}</td>
              <td style={{ padding: '0.5rem' }}>{e.action}</td>
              <td style={{ padding: '0.5rem', fontFamily: 'monospace' }}>{e.resource}</td>
            </tr>
          ))}
        </tbody>
      </table>
      {filtered.length === 0 && <p style={{ color: '#888', marginTop: '1rem' }}>No matching events.</p>}
    </div>
  );
}
