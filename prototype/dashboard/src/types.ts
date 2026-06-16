export type Priority = 1 | 2 | 3 | 4 | 5

export type TaskStatus = 'queued' | 'in-progress' | 'review' | 'done' | 'blocked'

export interface Task {
  id: string
  title: string
  assignedAgent: string
  status: TaskStatus
  priority: Priority
  createdAt: string
  description?: string
}

export type AgentState = 'idle' | 'working' | 'error'

export interface Agent {
  id: string
  name: string
  state: AgentState
  lastActivity: string
  currentTask?: string
}

export interface Checkpoint {
  id: string
  taskId: string
  agentName: string
  output: string
  fileDiffs: FileDiff[]
  status: 'pending' | 'approved' | 'rejected'
  createdAt: string
}

export interface FileDiff {
  path: string
  additions: string[]
  deletions: string[]
}
