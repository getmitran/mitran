import { useState, useEffect, useCallback, useRef } from 'react';
import ReactMarkdown from 'react-markdown';

const API = '/api/v1/wiki';

interface WikiPage {
  id: string;
  title: string;
  content: string;
  parent_id: string;
  slug: string;
  created_by: string;
  tags: string[];
  created_at: string;
  updated_at: string;
}

interface TreeNode extends WikiPage {
  children: TreeNode[];
}

export default function WikiApp() {
  const [pages, setPages] = useState<WikiPage[]>([]);
  const [tree, setTree] = useState<TreeNode[]>([]);
  const [selected, setSelected] = useState<WikiPage | null>(null);
  const [editing, setEditing] = useState(false);
  const [editTitle, setEditTitle] = useState('');
  const [editContent, setEditContent] = useState('');
  const [editTags, setEditTags] = useState('');
  const [search, setSearch] = useState('');
  const [searchResults, setSearchResults] = useState<WikiPage[] | null>(null);
  const [showNew, setShowNew] = useState(false);
  const [newTitle, setNewTitle] = useState('');
  const [newParent, setNewParent] = useState('');
  const saveTimer = useRef<ReturnType<typeof setTimeout> | null>(null);

  const fetchTree = useCallback(async () => {
    const res = await fetch(`${API}/tree`);
    setTree(await res.json());
  }, []);

  const fetchPages = useCallback(async () => {
    const res = await fetch(`${API}/pages`);
    setPages(await res.json());
  }, []);

  useEffect(() => { fetchTree(); fetchPages(); }, [fetchTree, fetchPages]);

  const selectPage = async (id: string) => {
    const res = await fetch(`${API}/pages/${id}`);
    const p = await res.json();
    setSelected(p);
    setEditTitle(p.title);
    setEditContent(p.content);
    setEditTags((p.tags || []).join(', '));
    setEditing(false);
    setSearchResults(null);
  };

  const doSearch = async () => {
    if (!search.trim()) { setSearchResults(null); return; }
    const res = await fetch(`${API}/search?q=${encodeURIComponent(search)}`);
    setSearchResults(await res.json());
  };

  const savePage = useCallback(async () => {
    if (!selected) return;
    await fetch(`${API}/pages/${selected.id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        title: editTitle,
        content: editContent,
        tags: editTags.split(',').map(t => t.trim()).filter(Boolean),
      }),
    });
    fetchTree();
    fetchPages();
  }, [selected, editTitle, editContent, editTags, fetchTree, fetchPages]);

  const autoSave = useCallback(() => {
    if (saveTimer.current) clearTimeout(saveTimer.current);
    saveTimer.current = setTimeout(savePage, 2000);
  }, [savePage]);

  const handleContentChange = (val: string) => {
    setEditContent(val);
    autoSave();
  };

  const handleTitleChange = (val: string) => {
    setEditTitle(val);
    autoSave();
  };

  const deletePage = async () => {
    if (!selected || !confirm('Delete this page?')) return;
    await fetch(`${API}/pages/${selected.id}`, { method: 'DELETE' });
    setSelected(null);
    fetchTree();
    fetchPages();
  };

  const createPage = async () => {
    if (!newTitle.trim()) return;
    const res = await fetch(`${API}/pages`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ title: newTitle, content: '', parent_id: newParent, tags: [] }),
    });
    const p = await res.json();
    setShowNew(false);
    setNewTitle('');
    setNewParent('');
    fetchTree();
    fetchPages();
    selectPage(p.id);
  };

  const getBreadcrumbs = (pageId: string): WikiPage[] => {
    const crumbs: WikiPage[] = [];
    let cur = pages.find(p => p.id === pageId);
    while (cur) {
      crumbs.unshift(cur);
      cur = cur.parent_id ? pages.find(p => p.id === cur!.parent_id) : undefined;
    }
    return crumbs;
  };

  return (
    <div className="flex h-full bg-gray-900 text-gray-100">
      {/* Sidebar */}
      <div className="w-72 border-r border-gray-700 flex flex-col">
        <div className="p-3 border-b border-gray-700">
          <input
            className="w-full px-3 py-1.5 bg-gray-800 border border-gray-600 rounded text-sm"
            placeholder="Search wiki..."
            value={search}
            onChange={e => setSearch(e.target.value)}
            onKeyDown={e => e.key === 'Enter' && doSearch()}
          />
        </div>
        <div className="p-2">
          <button
            onClick={() => setShowNew(true)}
            className="w-full px-3 py-1.5 bg-blue-600 hover:bg-blue-700 rounded text-sm font-medium"
          >
            + New Page
          </button>
        </div>
        <div className="flex-1 overflow-y-auto p-2">
          {searchResults ? (
            <div>
              <div className="text-xs text-gray-400 mb-2">{searchResults.length} results</div>
              {searchResults.map(p => (
                <div key={p.id} onClick={() => selectPage(p.id)}
                  className="px-2 py-1 rounded cursor-pointer hover:bg-gray-700 text-sm truncate">
                  {p.title}
                </div>
              ))}
            </div>
          ) : (
            <TreeView nodes={tree} onSelect={selectPage} selectedId={selected?.id} />
          )}
        </div>
      </div>

      {/* Main */}
      <div className="flex-1 flex flex-col overflow-hidden">
        {selected ? (
          <>
            {/* Top bar */}
            <div className="flex items-center gap-2 p-3 border-b border-gray-700">
              <div className="flex-1">
                <div className="text-xs text-gray-400 mb-1">
                  {getBreadcrumbs(selected.id).map((b, i) => (
                    <span key={b.id}>
                      {i > 0 && ' / '}
                      <span className="cursor-pointer hover:text-gray-200" onClick={() => selectPage(b.id)}>{b.title}</span>
                    </span>
                  ))}
                </div>
                {editing ? (
                  <input className="bg-gray-800 border border-gray-600 rounded px-2 py-1 text-lg font-semibold w-full"
                    value={editTitle} onChange={e => handleTitleChange(e.target.value)} />
                ) : (
                  <h1 className="text-lg font-semibold">{selected.title}</h1>
                )}
              </div>
              <div className="flex gap-2 text-sm">
                <span className="text-xs text-gray-400">{editTags && `Tags: ${editTags}`}</span>
                <button onClick={() => setEditing(!editing)}
                  className="px-3 py-1 rounded bg-gray-700 hover:bg-gray-600">
                  {editing ? 'View' : 'Edit'}
                </button>
                <button onClick={savePage} className="px-3 py-1 rounded bg-green-700 hover:bg-green-600">Save</button>
                <button onClick={deletePage} className="px-3 py-1 rounded bg-red-700 hover:bg-red-600">Delete</button>
              </div>
            </div>

            {/* Content */}
            <div className="flex-1 overflow-hidden flex">
              {editing ? (
                <>
                  <div className="w-1/2 flex flex-col border-r border-gray-700">
                    <div className="p-2 border-b border-gray-700">
                      <input className="w-full bg-gray-800 border border-gray-600 rounded px-2 py-1 text-xs"
                        placeholder="Tags (comma separated)" value={editTags}
                        onChange={e => { setEditTags(e.target.value); autoSave(); }} />
                    </div>
                    <textarea
                      className="flex-1 p-4 bg-gray-900 resize-none font-mono text-sm focus:outline-none"
                      value={editContent}
                      onChange={e => handleContentChange(e.target.value)}
                    />
                  </div>
                  <div className="w-1/2 overflow-y-auto p-4 prose prose-invert max-w-none">
                    <ReactMarkdown>{editContent}</ReactMarkdown>
                  </div>
                </>
              ) : (
                <div className="flex-1 overflow-y-auto p-6 prose prose-invert max-w-none">
                  <ReactMarkdown>{selected.content}</ReactMarkdown>
                </div>
              )}
            </div>
          </>
        ) : (
          <div className="flex-1 flex items-center justify-center text-gray-500">
            Select a page or create a new one
          </div>
        )}
      </div>

      {/* New page modal */}
      {showNew && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <div className="bg-gray-800 rounded-lg p-6 w-96 border border-gray-600">
            <h2 className="text-lg font-semibold mb-4">New Wiki Page</h2>
            <input className="w-full mb-3 px-3 py-2 bg-gray-900 border border-gray-600 rounded"
              placeholder="Page title" value={newTitle} onChange={e => setNewTitle(e.target.value)} />
            <select className="w-full mb-4 px-3 py-2 bg-gray-900 border border-gray-600 rounded"
              value={newParent} onChange={e => setNewParent(e.target.value)}>
              <option value="">No parent (root level)</option>
              {pages.map(p => <option key={p.id} value={p.id}>{p.title}</option>)}
            </select>
            <div className="flex justify-end gap-2">
              <button onClick={() => setShowNew(false)} className="px-4 py-2 rounded bg-gray-700 hover:bg-gray-600">Cancel</button>
              <button onClick={createPage} className="px-4 py-2 rounded bg-blue-600 hover:bg-blue-700">Create</button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

function TreeView({ nodes, onSelect, selectedId, depth = 0 }: {
  nodes: TreeNode[]; onSelect: (id: string) => void; selectedId?: string; depth?: number;
}) {
  const [collapsed, setCollapsed] = useState<Record<string, boolean>>({});
  return (
    <>
      {nodes.map(node => (
        <div key={node.id}>
          <div
            className={`flex items-center gap-1 px-2 py-1 rounded cursor-pointer text-sm ${
              node.id === selectedId ? 'bg-blue-900/50 text-blue-200' : 'hover:bg-gray-700'
            }`}
            style={{ paddingLeft: `${depth * 16 + 8}px` }}
          >
            {node.children.length > 0 && (
              <span className="text-xs cursor-pointer w-4 text-center"
                onClick={e => { e.stopPropagation(); setCollapsed(p => ({ ...p, [node.id]: !p[node.id] })); }}>
                {collapsed[node.id] ? '▶' : '▼'}
              </span>
            )}
            {node.children.length === 0 && <span className="w-4" />}
            <span className="truncate" onClick={() => onSelect(node.id)}>{node.title}</span>
          </div>
          {!collapsed[node.id] && node.children.length > 0 && (
            <TreeView nodes={node.children} onSelect={onSelect} selectedId={selectedId} depth={depth + 1} />
          )}
        </div>
      ))}
    </>
  );
}
