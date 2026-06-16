const BASE = 'http://localhost:7777/api/v1'

async function request<T>(path: string, opts?: RequestInit): Promise<T> {
  const res = await fetch(`${BASE}${path}`, {
    headers: { 'Content-Type': 'application/json' },
    ...opts,
  })
  if (!res.ok) throw new Error(`API ${res.status}: ${await res.text()}`)
  return res.json()
}

export const api = {
  fetchTasks: () => request<any[]>('/tasks'),
  fetchAgents: () => request<any[]>('/agents'),
  fetchCheckpoints: () => request<any[]>('/checkpoints'),
  approveCheckpoint: (id: string, feedback?: string) =>
    request(`/checkpoints/${id}/approve`, { method: 'POST', body: JSON.stringify({ feedback }) }),
  rejectCheckpoint: (id: string, feedback: string) =>
    request(`/checkpoints/${id}/reject`, { method: 'POST', body: JSON.stringify({ feedback }) }),
  createTask: (task: { title: string; description?: string; agent_type?: string; priority?: number }) =>
    request('/tasks', { method: 'POST', body: JSON.stringify(task) }),
  reorderTask: (id: string, priority: number) =>
    request(`/tasks/${id}`, { method: 'PATCH', body: JSON.stringify({ priority }) }),
  triggerInit: (config: { company_name: string; description: string; languages: string[]; team_size: number }) =>
    request('/init', { method: 'POST', body: JSON.stringify(config) }),
}
