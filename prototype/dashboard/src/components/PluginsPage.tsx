import { useEffect, useState } from 'react';

interface Plugin {
  name: string;
  version: string;
  enabled: boolean;
}

export default function PluginsPage() {
  const [plugins, setPlugins] = useState<Plugin[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetch('http://localhost:7780/api/v1/plugins')
      .then(r => r.json())
      .then(data => setPlugins(data.plugins || []))
      .catch(() => setPlugins([]))
      .finally(() => setLoading(false));
  }, []);

  const toggle = (name: string) => {
    setPlugins(ps => ps.map(p => p.name === name ? { ...p, enabled: !p.enabled } : p));
  };

  if (loading) return <div className="p-6 text-gray-400">Loading plugins…</div>;

  return (
    <div className="p-6">
      <div className="flex justify-between items-center mb-4">
        <h1 className="text-xl font-bold">Plugins</h1>
        <button
          onClick={() => {
            const name = prompt('Plugin package name:');
            if (!name) return;
            fetch('http://localhost:7780/api/v1/plugins', {
              method: 'POST',
              headers: { 'Content-Type': 'application/json' },
              body: JSON.stringify({ name }),
            }).then(() => window.location.reload());
          }}
          className="px-3 py-1.5 bg-blue-600 text-white text-sm rounded hover:bg-blue-700"
        >
          Install Plugin
        </button>
      </div>
      <table className="w-full text-sm border-collapse">
        <thead>
          <tr className="border-b border-gray-700 text-left">
            <th className="py-2 px-3">Name</th>
            <th className="py-2 px-3">Version</th>
            <th className="py-2 px-3">Enabled</th>
          </tr>
        </thead>
        <tbody>
          {plugins.length === 0 ? (
            <tr><td colSpan={3} className="py-4 px-3 text-gray-500">No plugins installed</td></tr>
          ) : (
            plugins.map(p => (
              <tr key={p.name} className="border-b border-gray-800">
                <td className="py-2 px-3">{p.name}</td>
                <td className="py-2 px-3">{p.version}</td>
                <td className="py-2 px-3">
                  <button
                    onClick={() => toggle(p.name)}
                    className={`w-10 h-5 rounded-full relative transition-colors ${p.enabled ? 'bg-green-500' : 'bg-gray-600'}`}
                  >
                    <span className={`absolute top-0.5 w-4 h-4 rounded-full bg-white transition-transform ${p.enabled ? 'left-5' : 'left-0.5'}`} />
                  </button>
                </td>
              </tr>
            ))
          )}
        </tbody>
      </table>
    </div>
  );
}
