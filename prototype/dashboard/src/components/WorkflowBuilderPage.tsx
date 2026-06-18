import { useState, useEffect } from 'react';

interface WorkflowStep {
  name: string;
  agent_id: string;
  action: string;
}

interface Workflow {
  id: string;
  name: string;
  steps: WorkflowStep[];
  created_at: string;
}

export default function WorkflowBuilderPage() {
  const [workflows, setWorkflows] = useState<Workflow[]>([]);
  const [name, setName] = useState('');
  const [steps, setSteps] = useState<WorkflowStep[]>([{ name: '', agent_id: '', action: '' }]);
  const [loading, setLoading] = useState(true);

  const fetchWorkflows = () => {
    fetch('http://localhost:7780/api/v1/workflows')
      .then(r => r.json())
      .then(data => setWorkflows(Array.isArray(data) ? data : data.workflows || []))
      .catch(() => setWorkflows([]))
      .finally(() => setLoading(false));
  };

  useEffect(() => { fetchWorkflows(); }, []);

  const addStep = () => setSteps([...steps, { name: '', agent_id: '', action: '' }]);

  const updateStep = (i: number, field: keyof WorkflowStep, value: string) => {
    const updated = [...steps];
    updated[i] = { ...updated[i], [field]: value };
    setSteps(updated);
  };

  const removeStep = (i: number) => setSteps(steps.filter((_, idx) => idx !== i));

  const createWorkflow = (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim() || steps.some(s => !s.name.trim())) return;
    fetch('http://localhost:7780/api/v1/workflows', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ name, steps }),
    })
      .then(r => r.json())
      .then(() => { setName(''); setSteps([{ name: '', agent_id: '', action: '' }]); fetchWorkflows(); });
  };

  const deleteWorkflow = (id: string) => {
    fetch(`http://localhost:7780/api/v1/workflows/${id}`, { method: 'DELETE' })
      .then(() => fetchWorkflows());
  };

  if (loading) return <div className="p-6 text-gray-400">Loading workflows...</div>;

  return (
    <div className="p-6 max-w-4xl mx-auto">
      <h1 className="text-2xl font-bold mb-6">Workflow Builder</h1>

      <form onSubmit={createWorkflow} className="mb-8 p-4 border border-gray-700 rounded-lg bg-gray-900">
        <h2 className="text-lg font-semibold mb-3">Create Workflow</h2>
        <input
          className="w-full mb-3 px-3 py-2 rounded bg-gray-800 border border-gray-600 text-white"
          placeholder="Workflow name"
          value={name}
          onChange={e => setName(e.target.value)}
          required
        />
        <div className="space-y-2 mb-3">
          {steps.map((step, i) => (
            <div key={i} className="flex gap-2 items-center">
              <input
                className="flex-1 px-2 py-1 rounded bg-gray-800 border border-gray-600 text-white text-sm"
                placeholder="Step name"
                value={step.name}
                onChange={e => updateStep(i, 'name', e.target.value)}
                required
              />
              <input
                className="flex-1 px-2 py-1 rounded bg-gray-800 border border-gray-600 text-white text-sm"
                placeholder="Agent ID"
                value={step.agent_id}
                onChange={e => updateStep(i, 'agent_id', e.target.value)}
              />
              <input
                className="flex-1 px-2 py-1 rounded bg-gray-800 border border-gray-600 text-white text-sm"
                placeholder="Action"
                value={step.action}
                onChange={e => updateStep(i, 'action', e.target.value)}
              />
              {steps.length > 1 && (
                <button type="button" onClick={() => removeStep(i)} className="text-red-400 text-sm hover:text-red-300">✕</button>
              )}
            </div>
          ))}
        </div>
        <div className="flex gap-2">
          <button type="button" onClick={addStep} className="px-3 py-1 text-sm rounded bg-gray-700 hover:bg-gray-600 text-white">+ Add Step</button>
          <button type="submit" className="px-4 py-1 text-sm rounded bg-blue-600 hover:bg-blue-500 text-white font-medium">Create</button>
        </div>
      </form>

      <div className="space-y-2">
        {workflows.length === 0 ? (
          <p className="text-gray-500">No workflows yet.</p>
        ) : (
          workflows.map(w => (
            <div key={w.id} className="flex items-center justify-between p-3 border border-gray-700 rounded bg-gray-900">
              <div>
                <span className="font-medium text-white">{w.name}</span>
                <span className="ml-3 text-sm text-gray-400">{w.steps?.length || 0} steps</span>
                <span className="ml-3 text-xs text-gray-500">{new Date(w.created_at).toLocaleDateString()}</span>
              </div>
              <button onClick={() => deleteWorkflow(w.id)} className="text-red-400 text-sm hover:text-red-300">Delete</button>
            </div>
          ))
        )}
      </div>
    </div>
  );
}
