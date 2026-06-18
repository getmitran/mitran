import { useState, useEffect } from 'react';

interface Artifact {
  slug: string;
  name: string;
  kind: string;
  version: number;
  tags: string[];
  content: string;
  updated_at: string;
}

const API = import.meta.env.VITE_API_URL || 'http://localhost:7780';

export default function ArtifactPage() {
  const [artifacts, setArtifacts] = useState<Artifact[]>([]);
  const [loading, setLoading] = useState(true);
  const [filterKind, setFilterKind] = useState('');
  const [filterTag, setFilterTag] = useState('');
  const [selected, setSelected] = useState<Artifact | null>(null);
  const [previewVersion, setPreviewVersion] = useState<number | null>(null);

  useEffect(() => {
    fetch(`${API}/api/v1/artifacts`)
      .then(r => r.json())
      .then(data => setArtifacts(Array.isArray(data) ? data : data.artifacts || []))
      .catch(() => setArtifacts([]))
      .finally(() => setLoading(false));
  }, []);

  const kinds = [...new Set(artifacts.map(a => a.kind))];
  const allTags = [...new Set(artifacts.flatMap(a => a.tags || []))];

  const filtered = artifacts.filter(a => {
    if (filterKind && a.kind !== filterKind) return false;
    if (filterTag && !(a.tags || []).includes(filterTag)) return false;
    return true;
  });

  const handleSelect = (a: Artifact) => {
    setSelected(a);
    setPreviewVersion(a.version);
  };

  if (loading) return <div className="p-6 text-gray-400">Loading artifacts...</div>;

  return (
    <div className="p-6">
      <h1 className="text-2xl font-bold mb-4">Artifact Library</h1>

      <div className="flex gap-3 mb-6">
        <select
          value={filterKind}
          onChange={e => setFilterKind(e.target.value)}
          className="px-3 py-1.5 rounded bg-gray-800 border border-gray-700 text-sm"
        >
          <option value="">All Kinds</option>
          {kinds.map(k => <option key={k} value={k}>{k}</option>)}
        </select>
        <select
          value={filterTag}
          onChange={e => setFilterTag(e.target.value)}
          className="px-3 py-1.5 rounded bg-gray-800 border border-gray-700 text-sm"
        >
          <option value="">All Tags</option>
          {allTags.map(t => <option key={t} value={t}>{t}</option>)}
        </select>
        <span className="text-sm text-gray-400 self-center">{filtered.length} artifacts</span>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {filtered.map(a => (
          <div
            key={a.slug}
            onClick={() => handleSelect(a)}
            className="p-4 rounded-lg bg-gray-800 border border-gray-700 cursor-pointer hover:border-blue-500 transition-colors"
          >
            <div className="font-semibold truncate">{a.name}</div>
            <div className="text-xs text-gray-400 mt-1">
              {a.kind} &middot; v{a.version}
            </div>
            {a.tags?.length > 0 && (
              <div className="flex flex-wrap gap-1 mt-2">
                {a.tags.map(t => (
                  <span key={t} className="text-xs px-2 py-0.5 rounded bg-gray-700 text-gray-300">{t}</span>
                ))}
              </div>
            )}
            <div className="text-xs text-gray-500 mt-2">
              {new Date(a.updated_at).toLocaleDateString()}
            </div>
          </div>
        ))}
        {filtered.length === 0 && (
          <div className="col-span-full text-gray-500 text-center py-12">No artifacts found</div>
        )}
      </div>

      {selected && (
        <div className="fixed inset-0 bg-black/60 flex items-center justify-center z-50" onClick={() => setSelected(null)}>
          <div className="bg-gray-900 border border-gray-700 rounded-xl w-full max-w-2xl max-h-[80vh] overflow-auto p-6" onClick={e => e.stopPropagation()}>
            <div className="flex justify-between items-start mb-4">
              <div>
                <h2 className="text-xl font-bold">{selected.name}</h2>
                <span className="text-xs text-gray-400">{selected.kind} &middot; {selected.slug}</span>
              </div>
              <div className="flex items-center gap-2">
                <select
                  value={previewVersion ?? selected.version}
                  onChange={e => setPreviewVersion(Number(e.target.value))}
                  className="px-2 py-1 rounded bg-gray-800 border border-gray-700 text-sm"
                >
                  {Array.from({ length: selected.version }, (_, i) => i + 1).reverse().map(v => (
                    <option key={v} value={v}>v{v}</option>
                  ))}
                </select>
                <button onClick={() => setSelected(null)} className="text-gray-400 hover:text-white text-xl">&times;</button>
              </div>
            </div>
            {selected.tags?.length > 0 && (
              <div className="flex gap-1 mb-4">
                {selected.tags.map(t => (
                  <span key={t} className="text-xs px-2 py-0.5 rounded bg-gray-700 text-gray-300">{t}</span>
                ))}
              </div>
            )}
            <div className="bg-gray-950 rounded p-4 text-sm font-mono whitespace-pre-wrap overflow-auto max-h-96 border border-gray-800">
              {selected.content || <span className="text-gray-500">No content</span>}
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
