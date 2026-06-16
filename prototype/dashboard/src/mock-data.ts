import { Task, Agent, Checkpoint } from './types'

export const mockTasks: Task[] = [
  {
    id: 'task-001',
    title: 'Implement user authentication flow',
    assignedAgent: 'Dev Agent',
    status: 'in-progress',
    priority: 1,
    createdAt: '2026-06-16T09:00:00Z',
  },
  {
    id: 'task-002',
    title: 'Write API documentation for /users endpoint',
    assignedAgent: 'Docs Agent',
    status: 'review',
    priority: 2,
    createdAt: '2026-06-16T08:30:00Z',
  },
  {
    id: 'task-003',
    title: 'Set up CI/CD pipeline for staging',
    assignedAgent: 'CI/CD Agent',
    status: 'queued',
    priority: 3,
    createdAt: '2026-06-16T08:00:00Z',
  },
  {
    id: 'task-004',
    title: 'Create onboarding ticket for new hire',
    assignedAgent: 'HR Agent',
    status: 'done',
    priority: 4,
    createdAt: '2026-06-15T16:00:00Z',
  },
  {
    id: 'task-005',
    title: 'Investigate memory leak in worker process',
    assignedAgent: 'Ops Agent',
    status: 'blocked',
    priority: 2,
    createdAt: '2026-06-16T07:45:00Z',
  },
]

export const mockAgents: Agent[] = [
  { id: 'agent-dev', name: 'Dev Agent', state: 'working', lastActivity: '2m ago', currentTask: 'task-001' },
  { id: 'agent-docs', name: 'Docs Agent', state: 'idle', lastActivity: '5m ago' },
  { id: 'agent-ops', name: 'Ops Agent', state: 'error', lastActivity: '1m ago', currentTask: 'task-005' },
  { id: 'agent-review', name: 'Review Agent', state: 'working', lastActivity: '30s ago', currentTask: 'task-002' },
  { id: 'agent-hr', name: 'HR Agent', state: 'idle', lastActivity: '12m ago' },
  { id: 'agent-cicd', name: 'CI/CD Agent', state: 'idle', lastActivity: '8m ago' },
  { id: 'agent-tickets', name: 'Tickets Agent', state: 'idle', lastActivity: '15m ago' },
  { id: 'agent-wiki', name: 'Wiki Agent', state: 'idle', lastActivity: '20m ago' },
]

export const mockCheckpoints: Checkpoint[] = [
  {
    id: 'cp-001',
    taskId: 'task-002',
    agentName: 'Docs Agent',
    output: `# Users API Documentation\n\n## GET /users\n\nReturns a paginated list of users.\n\n### Parameters\n\n| Name | Type | Required | Description |\n|------|------|----------|-------------|\n| page | int | No | Page number (default: 1) |\n| limit | int | No | Items per page (default: 20) |\n\n### Response\n\n\`\`\`json\n{\n  "users": [...],\n  "total": 142,\n  "page": 1\n}\n\`\`\``,
    fileDiffs: [
      {
        path: 'docs/api/users.md',
        additions: [
          '## GET /users',
          '',
          'Returns a paginated list of users.',
          '',
          '### Parameters',
        ],
        deletions: [
          '## TODO: Document users endpoint',
        ],
      },
    ],
    status: 'pending',
    createdAt: '2026-06-16T09:15:00Z',
  },
]
