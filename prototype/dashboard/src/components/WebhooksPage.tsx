import { useEffect, useState } from 'react';

interface Webhook {
  id: string;
  url: string;
  event: string;
  secret: string;
}

const EVENTS = ['task.created', 'task.completed', 'agent.error'];
const API = 'http://localhost:7780/api/v1/webhooks';

export default function WebhooksPage() {
  const [webhooks, setWebhooks] = useState<Webhook[]>([]);
  const [url, setUrl] = useState('');
  const [event, setEvent] = useState(EVENTS[0]);
  const [secret, setSecret] = useState('');

  const load = () => fetch(API).then(r => r.json()).then(setWebhooks).catch(() => {});

  useEffect(() => { load(); }, []);

  const add = async (e: React.FormEvent) => {
    e.preventDefault();
    await fetch(API, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ url, event, secret }),
    });
    setUrl(''); setSecret('');
    load();
  };

  const remove = async (id: string) => {
    await fetch(`${API}/${id}`, { method: 'DELETE' });
    load();
  };

  const mask = (s: string) => s ? '•'.repeat(Math.min(s.length, 12)) + s.slice(-4) : '—';

  return (
    <div className="p-6 max-w-4xl mx-auto">
      <h1 className="text-2xl font-bold mb-4">Webhooks</h1>

      <form onSubmit={add} className="flex gap-2 mb-6 flex-wrap">
        <input value={url} onChange={e => setUrl(e.target.value)} placeholder="https://..." required
          className="flex-1 min-w-[200px] px-3 py-2 border rounded bg-transparent" />
        <select value={event} onChange={e => setEvent(e.target.value)}
          className="px-3 py-2 border rounded bg-transparent">
          {EVENTS.map(ev => <option key={ev} value={ev}>{ev}</option>)}
        </select>
        <input value={secret} onChange={e => setSecret(e.target.value)} placeholder="Secret" type="password"
          className="w-40 px-3 py-2 border rounded bg-transparent" />
        <button type="submit" className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700">Add</button>
      </form>

      <table className="w-full text-sm border-collapse">
        <thead>
          <tr className="border-b">
            <th className="text-left p-2">URL</th>
            <th className="text-left p-2">Event</th>
            <th className="text-left p-2">Secret</th>
            <th className="p-2"></th>
          </tr>
        </thead>
        <tbody>
          {webhooks.map(wh => (
            <tr key={wh.id} className="border-b hover:bg-gray-50 dark:hover:bg-gray-800">
              <td className="p-2 font-mono text-xs break-all">{wh.url}</td>
              <td className="p-2">{wh.event}</td>
              <td className="p-2 font-mono">{mask(wh.secret)}</td>
              <td className="p-2">
                <button onClick={() => remove(wh.id)}
                  className="text-red-500 hover:text-red-700 text-xs">Delete</button>
              </td>
            </tr>
          ))}
          {!webhooks.length && (
            <tr><td colSpan={4} className="p-4 text-center text-gray-500">No webhooks configured</td></tr>
          )}
        </tbody>
      </table>
    </div>
  );
}
