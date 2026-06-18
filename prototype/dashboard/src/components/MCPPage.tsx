import { useState, useEffect } from 'react';

interface MCPTool {
  name: string;
  description: string;
}

interface MCPServer {
  id: string;
  name: string;
  url: string;
  status: 'connected' | 'disconnected';
  tools: MCPTool[];
}

export default function MCPPage() {
  const [servers, setServers] = useState<MCPServer[]>([]);
  const [expandedServer, setExpandedServer] = useState<string | null>(null);
  const [newName, setNewName] = useState('');
  const [newUrl, setNewUrl] = useState('');
  const [loading, setLoading] = useState(true);

  const fetchServers = async () => {
    try {
      const res = await fetch('http://localhost:7780/api/v1/mcp/servers');
      if (res.ok) setServers(await res.json());
    } catch {
      setServers([]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => { fetchServers(); }, []);

  const addServer = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newName || !newUrl) return;
    try {
      await fetch('http://localhost:7780/api/v1/mcp/servers', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name: newName, url: newUrl }),
      });
      setNewName('');
      setNewUrl('');
      fetchServers();
    } catch {}
  };

  const removeServer = async (id: string) => {
    try {
      await fetch(`http://localhost:7780/api/v1/mcp/servers/${id}`, { method: 'DELETE' });
      fetchServers();
    } catch {}
  };

  if (loading) return <div className="p-6 text-gray-400">Loading...</div>;

  return (
    <div className="p-6 max-w-4xl mx-auto">
      <h1 className="text-2xl font-bold text-white mb-6">MCP Registry</h1>

      {/* Add Server Form */}
      <form onSubmit={addServer} className="mb-8 p-4 bg-gray-800 rounded-lg border border-gray-700">
        <h2 className="text-sm font-semibold text-gray-300 mb-3">Add Server</h2>
        <div className="flex gap-3">
          <input
            value={newName}
            onChange={e => setNewName(e.target.value)}
            placeholder="Server name"
            className="flex-1 px-3 py-2 bg-gray-900 border border-gray-600 rounded text-sm text-white placeholder-gray-500"
          />
          <input
            value={newUrl}
            onChange={e => setNewUrl(e.target.value)}
            placeholder="URL (e.g. stdio:///path/to/server)"
            className="flex-[2] px-3 py-2 bg-gray-900 border border-gray-600 rounded text-sm text-white placeholder-gray-500"
          />
          <button type="submit" className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white text-sm rounded">
            Add
          </button>
        </div>
      </form>

      {/* Server List */}
      {servers.length === 0 ? (
        <p className="text-gray-500 text-sm">No MCP servers connected.</p>
      ) : (
        <div className="space-y-3">
          {servers.map(server => (
            <div key={server.id} className="bg-gray-800 border border-gray-700 rounded-lg overflow-hidden">
              <div
                className="flex items-center justify-between p-4 cursor-pointer hover:bg-gray-750"
                onClick={() => setExpandedServer(expandedServer === server.id ? null : server.id)}
              >
                <div className="flex items-center gap-3">
                  <span className={`w-2.5 h-2.5 rounded-full ${server.status === 'connected' ? 'bg-green-500' : 'bg-red-500'}`} />
                  <span className="text-white font-medium">{server.name}</span>
                  <span className="text-gray-500 text-xs">{server.url}</span>
                </div>
                <div className="flex items-center gap-3">
                  <span className="text-xs text-gray-400">{server.tools?.length || 0} tools</span>
                  <button
                    onClick={e => { e.stopPropagation(); removeServer(server.id); }}
                    className="text-red-400 hover:text-red-300 text-xs"
                  >
                    Remove
                  </button>
                </div>
              </div>

              {expandedServer === server.id && server.tools?.length > 0 && (
                <div className="border-t border-gray-700 p-4">
                  <h3 className="text-xs font-semibold text-gray-400 mb-2 uppercase">Available Tools</h3>
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-2">
                    {server.tools.map(tool => (
                      <div key={tool.name} className="p-2 bg-gray-900 rounded text-sm">
                        <div className="text-blue-400 font-mono text-xs">{tool.name}</div>
                        {tool.description && (
                          <div className="text-gray-500 text-xs mt-0.5 truncate">{tool.description}</div>
                        )}
                      </div>
                    ))}
                  </div>
                </div>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
