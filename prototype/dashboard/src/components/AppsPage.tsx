import { useEffect, useState } from 'react';

interface App {
  id: string;
  name: string;
  description: string;
  enabled: boolean;
}

export default function AppsPage() {
  const [apps, setApps] = useState<App[]>([]);
  const [loading, setLoading] = useState(true);

  const fetchApps = async () => {
    try {
      const res = await fetch('http://localhost:7780/api/v1/apps');
      if (res.ok) setApps(await res.json());
    } catch { /* ignore */ }
    setLoading(false);
  };

  useEffect(() => { fetchApps(); }, []);

  const toggle = async (id: string, enable: boolean) => {
    await fetch(`http://localhost:7780/api/v1/apps/${id}/${enable ? 'enable' : 'disable'}`, { method: 'POST' });
    fetchApps();
  };

  if (loading) return <div className="p-6 text-gray-400">Loading apps...</div>;

  return (
    <div className="p-6">
      <h1 className="text-2xl font-bold mb-6">Apps</h1>
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {apps.map(app => (
          <div key={app.id} className="border border-gray-700 rounded-lg p-4 bg-gray-800">
            <div className="flex items-center justify-between mb-2">
              <h2 className="text-lg font-semibold">{app.name}</h2>
              <span className={`text-xs px-2 py-0.5 rounded ${app.enabled ? 'bg-green-900 text-green-300' : 'bg-gray-700 text-gray-400'}`}>
                {app.enabled ? 'Enabled' : 'Disabled'}
              </span>
            </div>
            <p className="text-sm text-gray-400 mb-4">{app.description}</p>
            <button
              onClick={() => toggle(app.id, !app.enabled)}
              className={`text-sm px-3 py-1.5 rounded ${app.enabled ? 'bg-red-700 hover:bg-red-600' : 'bg-blue-700 hover:bg-blue-600'} text-white`}
            >
              {app.enabled ? 'Disable' : 'Enable'}
            </button>
          </div>
        ))}
        {apps.length === 0 && <p className="text-gray-500">No apps installed.</p>}
      </div>
    </div>
  );
}
