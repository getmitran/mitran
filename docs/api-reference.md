# Mitran API Reference

Base URL: `http://localhost:7780`

## Projects
| Method | Path | Description |
|--------|------|-------------|
| POST | /api/v1/init | Initialize project |
| GET | /api/v1/project | Get project info |

## Tasks
| Method | Path | Description |
|--------|------|-------------|
| GET | /api/v1/tasks | List all tasks |
| POST | /api/v1/tasks | Create task |
| PATCH | /api/v1/tasks/:id | Update task |
| DELETE | /api/v1/tasks/:id | Delete task |

## Sessions
| Method | Path | Description |
|--------|------|-------------|
| POST | /api/v1/sessions | Create session |
| GET | /api/v1/sessions | List sessions |
| GET | /api/v1/sessions/:id | Get session |
| POST | /api/v1/sessions/:id/messages | Add message |
| DELETE | /api/v1/sessions/:id | End session |

## Auth
| Method | Path | Description |
|--------|------|-------------|
| POST | /api/v1/auth/login | Login (returns JWT) |
| POST | /api/v1/auth/register | Register user |
| GET | /api/v1/auth/me | Current user |

## Agents
| Method | Path | Description |
|--------|------|-------------|
| GET | /api/v1/agents/status | Agent fleet status |

## Checkpoints
| Method | Path | Description |
|--------|------|-------------|
| GET | /api/v1/checkpoints | List checkpoints |
| GET | /api/v1/checkpoints/pending | Pending only |
| POST | /api/v1/checkpoints/:id/approve | Approve |
| POST | /api/v1/checkpoints/:id/reject | Reject |

## Cron
| Method | Path | Description |
|--------|------|-------------|
| GET | /api/v1/cron | List jobs |
| POST | /api/v1/cron | Create job |
| DELETE | /api/v1/cron/:id | Remove job |
| POST | /api/v1/cron/:id/pause | Pause |
| POST | /api/v1/cron/:id/resume | Resume |
| POST | /api/v1/cron/:id/trigger | Trigger now |

## Events
| Method | Path | Description |
|--------|------|-------------|
| GET | /api/v1/events | Recent events |
| GET | /api/v1/events/stream | SSE stream |

## WebSocket
| Path | Description |
|------|-------------|
| /api/v1/ws | Real-time events |

## Chat (Worker)
Base URL: `http://localhost:8888`
| Method | Path | Description |
|--------|------|-------------|
| POST | /execute | Execute agent task |
| POST | /stream | Stream LLM response (SSE) |
| POST | /chat | Chat with agent |
