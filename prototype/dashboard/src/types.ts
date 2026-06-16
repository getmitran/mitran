export type Priority = 1 | 2 | 3 | 4 | 5

export type TaskStatus = 'queued' | 'in-progress' | 'review' | 'done' | 'blocked'

export interface Task {
  id: string
  title: string
  description?: string
  agent_type?: string
  assignedAgent: string
  status: TaskStatus
  priority: Priority
  created_at?: string
  createdAt: string
}

export type AgentState = 'idle' | 'working' | 'error'

export interface Agent {
  id: string
  name: string
  type?: string
  state: AgentState
  lastActivity: string
  last_heartbeat?: string
  currentTask?: string
}

export interface Checkpoint {
  id: string
  taskId: string
  task_id?: string
  agentName: string
  agent_name?: string
  output: string
  fileDiffs: FileDiff[]
  file_diffs?: FileDiff[]
  status: 'pending' | 'approved' | 'rejected'
  createdAt: string
  created_at?: string
}

export interface FileDiff {
  path: string
  additions: string[]
  deletions: string[]
}
