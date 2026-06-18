import { useState, useEffect } from 'react';

const API = 'http://localhost:7780';

type Tab = 'facts' | 'episodes' | 'corrections';

interface MemoryItem {
  id: string;
  key?: string;
  value?: string;
  content?: string;
  rule?: string;
  category?: string;
  timestamp?: string;
}

export default function MemoryPage() {
  const [tab, setTab] = useState<Tab>('facts');
  const [items, setItems] = useState<MemoryItem[]>([]);
  const [query, setQuery] = useState('');
  const [loading, setLoading] = useState(false);

  const fetchItems = async () => {
    setLoading(true);
    try {
      const url = query
        ? `${API}/api/v1/memory/search?q=${encodeURIComponent(query)}`
        : `${API}/api/v1/memory/${tab}`;
      const res = await fetch(url);
      if (res.ok) {
        const data = await res.json();
        setItems(Array.isArray(data) ? data : data.items || []);
      } else {
        setItems([]);
      }
    } catch {
      setItems([]);
    }
    setLoading(false);
  };

  useEffect(() => { fetchItems(); }, [tab]);

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    fetchItems();
  };

  const handleDelete = async (id: string) => {
    if (!confirm('Delete this item?')) return;
    await fetch(`${API}/api/v1/memory/${tab}/${id}`, { method: 'DELETE' });
    fetchItems();
  };

  const handleEdit = async (item: MemoryItem) => {
    const field = tab === 'facts' ? 'value' : tab === 'corrections' ? 'rule' : 'content';
    const current = item[field as keyof MemoryItem] || '';
    const updated = prompt('Edit:', current as string);
    if (updated === null || updated === current) return;
    await fetch(`${API}/api/v1/memory/${tab}/${item.id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ [field]: updated }),
    });
    fetchItems();
  };

  const getDisplay = (item: MemoryItem) => {
    if (tab === 'facts') return { title: item.key || item.id, body: item.value || '' };
    if (tab === 'corrections') return { title: item.category || 'correction', body: item.rule || '' };
    return { title: item.timestamp || item.id, body: item.content || '' };
  };

  const tabs: { key: Tab; label: string }[] = [
    { key: 'facts', label: 'Facts' },
    { key: 'episodes', label: 'Episodes' },
    { key: 'corrections', label: 'Corrections' },
  ];

  return (
    <div className="p-6 max-w-4xl mx-auto">
      <h1 className="text-2xl font-bold mb-4">Memory</h1>

      <form onSubmit={handleSearch} className="flex gap-2 mb-4">
        <input
          type="text"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          placeholder="Search memory..."
          className="flex-1 px-3 py-2 border border-gray-300 rounded-md bg-white dark:bg-gray-800 dark:border-gray-600"
        />
        <button type="submit" className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700">
          Search
        </button>
        {query && (
          <button type="button" onClick={() => { setQuery(''); fetchItems(); }} className="px-3 py-2 text-gray-500 hover:text-gray-700">
            Clear
          </button>
        )}
      </form>

      <div className="flex gap-1 mb-4 border-b border-gray-200 dark:border-gray-700">
        {tabs.map((t) => (
          <button
            key={t.key}
            onClick={() => { setTab(t.key); setQuery(''); }}
            className={`px-4 py-2 text-sm font-medium rounded-t-md ${
              tab === t.key
                ? 'bg-blue-50 text-blue-700 border-b-2 border-blue-600 dark:bg-gray-700 dark:text-blue-400'
                : 'text-gray-500 hover:text-gray-700 dark:text-gray-400'
            }`}
          >
            {t.label}
          </button>
        ))}
      </div>

      {loading ? (
        <p className="text-gray-500">Loading...</p>
      ) : items.length === 0 ? (
        <p className="text-gray-500">No items found.</p>
      ) : (
        <ul className="space-y-2">
          {items.map((item) => {
            const { title, body } = getDisplay(item);
            return (
              <li key={item.id} className="p-3 border border-gray-200 dark:border-gray-700 rounded-md bg-white dark:bg-gray-800">
                <div className="flex justify-between items-start">
                  <div className="flex-1 min-w-0">
                    <p className="text-sm font-medium text-gray-900 dark:text-gray-100 truncate">{title}</p>
                    <p className="text-sm text-gray-600 dark:text-gray-400 mt-1 whitespace-pre-wrap">{body}</p>
                  </div>
                  <div className="flex gap-1 ml-2 shrink-0">
                    <button onClick={() => handleEdit(item)} className="text-xs px-2 py-1 text-blue-600 hover:bg-blue-50 rounded dark:hover:bg-gray-700">
                      Edit
                    </button>
                    <button onClick={() => handleDelete(item.id)} className="text-xs px-2 py-1 text-red-600 hover:bg-red-50 rounded dark:hover:bg-gray-700">
                      Delete
                    </button>
                  </div>
                </div>
              </li>
            );
          })}
        </ul>
      )}
    </div>
  );
}
