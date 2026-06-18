import { useState, useEffect } from 'react';

interface Project {
  id: string;
  name: string;
  description: string;
  created_at: string;
}

export default function ProjectsPage() {
  const [projects, setProjects] = useState<Project[]>([]);
  const [newName, setNewName] = useState('');

  const fetchProjects = () => {
    fetch('/api/v1/projects')
      .then(r => r.json())
      .then(data => setProjects(data.projects || []))
      .catch(() => {});
  };

  useEffect(() => { fetchProjects(); }, []);

  const createProject = () => {
    if (!newName.trim()) return;
    fetch('/api/v1/projects', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ name: newName.trim() }),
    })
      .then(r => r.json())
      .then(() => { setNewName(''); fetchProjects(); })
      .catch(() => {});
  };

  return (
    <div className="p-6">
      <div className="flex items-center gap-3 mb-6">
        <input
          value={newName}
          onChange={e => setNewName(e.target.value)}
          onKeyDown={e => e.key === 'Enter' && createProject()}
          placeholder="Project name"
          className="px-3 py-2 border rounded text-sm"
        />
        <button onClick={createProject} className="px-4 py-2 bg-blue-600 text-white rounded text-sm hover:bg-blue-700">
          New Project
        </button>
      </div>
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
        {projects.map(p => (
          <div key={p.id} className="border rounded-lg p-4">
            <h3 className="font-semibold text-lg">{p.name}</h3>
            <p className="text-sm text-gray-600 mt-1">{p.description || 'No description'}</p>
            <p className="text-xs text-gray-400 mt-2">{new Date(p.created_at).toLocaleDateString()}</p>
          </div>
        ))}
      </div>
    </div>
  );
}
