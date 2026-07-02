# Mitran — Low-Level Design (LLD)

## 1. Core Engine (Go)

### 1.1 Module Structure
```
cmd/
├── mitran/              # CLI entrypoint
│   └── main.go
internal/
├── engine/
│   ├── dag.go           # DAG scheduler (topological sort, cycle detection)
│   ├── queue.go         # Priority queue (heap-based, human-reorderable)
│   ├── executor.go      # Task executor (sync by default, parallel on disjoint)
│   └── lock.go          # Resource lock manager (path-level mutex)
├── agent/
│   ├── runtime.go       # Agent process lifecycle (spawn, health, kill)
│   ├── acl.go           # Path ACL enforcement
│   ├── manifest.go      # Parse agent manifest.yaml
│   └── pool.go          # Agent process pool (Python subprocess management)
├── checkpoint/
│   ├── store.go         # Checkpoint persistence (SQLite/Postgres)
│   ├── diff.go          # Diff generation for human review
│   └── approval.go      # Approval gate logic
├── app/
│   ├── registry.go      # App install/enable/disable/uninstall
│   ├── local.go         # Local app resolution
│   └── marketplace.go   # Remote marketplace client (Phase 2)
├── integration/
│   ├── github.go        # GitHub API client (repos, PRs, Actions)
│   ├── slack.go         # Slack bot + slash commands
│   ├── grafana.go       # Grafana HTTP API client
│   └── prometheus.go    # Prometheus config generator
├── event/
│   ├── bus.go           # NATS-backed event bus
│   ├── webhook.go       # Incoming webhook handler
│   └── notify.go        # Human notification dispatch
├── api/
│   ├── grpc.go          # gRPC server (agent ↔ core)
│   ├── http.go          # REST API (dashboard ↔ core)
│   └── ws.go            # WebSocket (real-time updates to dashboard)
├── storage/
│   ├── sqlite.go        # Local storage implementation
│   ├── postgres.go      # Enterprise storage implementation
│   └── migrations/      # Schema migrations
└── config/
    └── config.go        # Configuration loading (YAML/env/flags)
```

### 1.2 Key Data Structures

```go
// Task represents a unit of work in the priority queue
type Task struct {
    ID          string        `json:"id"`
    Title       string        `json:"title"`
    Description string        `json:"description"`
    AgentType   AgentType     `json:"agent_type"`
    Priority    int           `json:"priority"` // 1=highest
    Status      TaskStatus    `json:"status"`   // queued|running|checkpoint|approved|failed
    DependsOn   []string      `json:"depends_on"` // task IDs
    Resources   []string      `json:"resources"`  // paths this task reads/writes
    CreatedAt   time.Time     `json:"created_at"`
    CreatedBy   string        `json:"created_by"` // "human" or "agent:<id>"
    Checkpoint  *Checkpoint   `json:"checkpoint,omitempty"`
}

type TaskStatus string
const (
    TaskQueued     TaskStatus = "queued"
    TaskRunning    TaskStatus = "running"
    TaskCheckpoint TaskStatus = "checkpoint" // waiting for human
    TaskApproved   TaskStatus = "approved"
    TaskRejected   TaskStatus = "rejected"
    TaskFailed     TaskStatus = "failed"
)

// Checkpoint represents an agent's output awaiting human review
type Checkpoint struct {
    ID        string       `json:"id"`
    TaskID    string       `json:"task_id"`
    AgentID   string       `json:"agent_id"`
    Output    string       `json:"output"`     // markdown summary
    Diffs     []FileDiff   `json:"diffs"`      // file changes
    Questions []string     `json:"questions"`  // agent's open questions
    Decision  *Decision    `json:"decision,omitempty"`
    CreatedAt time.Time    `json:"created_at"`
}

type Decision struct {
    Action   string    `json:"action"` // "approve"|"reject"|"modify"
    Feedback string    `json:"feedback,omitempty"`
    By       string    `json:"by"` // human user ID
    At       time.Time `json:"at"`
}

// AgentManifest defines what an agent can do
type AgentManifest struct {
    Name         string   `yaml:"name"`
    Version      string   `yaml:"version"`
    Description  string   `yaml:"description"`
    Permissions  Perms    `yaml:"permissions"`
    Integrations []string `yaml:"integrations"` // which APIs it needs
    Model        string   `yaml:"model"`        // preferred LLM
    Tools        []string `yaml:"tools"`        // tool names it uses
}

type Perms struct {
    Read    []string `yaml:"read"`    // glob patterns
    Write   []string `yaml:"write"`   // glob patterns
    Execute []string `yaml:"execute"` // allowed commands
}
```

### 1.3 DAG Scheduler Algorithm

```go
func (e *Engine) Schedule() (*Task, error) {
    e.mu.Lock()
    defer e.mu.Unlock()

    // 1. Get all queued tasks sorted by human-set priority
    candidates := e.queue.GetByStatus(TaskQueued)

    for _, task := range candidates {
        // 2. Check all dependencies are approved
        if !e.allDepsApproved(task) {
            continue
        }
        // 3. Check no resource conflicts with running tasks
        if e.hasResourceConflict(task) {
            continue
        }
        // 4. Check agent is available
        if !e.agentPool.IsAvailable(task.AgentType) {
            continue
        }
        // 5. Acquire resource locks
        e.lockResources(task)
        return task, nil
    }
    return nil, ErrNoExecutableTask
}
```

### 1.4 Resource Locking

```go
// Prevents parallel access to same files
type ResourceLockManager struct {
    locks map[string]*ResourceLock // path → lock
    mu    sync.RWMutex
}

type ResourceLock struct {
    Path     string
    HeldBy   string // task ID
    Mode     LockMode // read|write
    AcquiredAt time.Time
}

// Conflict rules:
// - Write locks are exclusive (no other read or write)
// - Read locks are shared (multiple readers OK)
// - Glob patterns expand to specific paths at lock time
```

## 2. Agent SDK (Python)

### 2.1 SDK Interface

```python
# mitran_sdk/agent.py

from abc import ABC, abstractmethod
from dataclasses import dataclass

@dataclass
class TaskContext:
    task_id: str
    title: str
    description: str
    workspace_path: str
    history: list[dict]       # prior checkpoints for this task
    integrations: dict        # authenticated API clients

class MitranAgent(ABC):
    """Base class for all Mitran agents."""

    @abstractmethod
    def execute(self, ctx: TaskContext) -> AgentOutput:
        """Execute the task. Return output for checkpoint."""
        pass

    def ask_human(self, question: str) -> str:
        """Block execution and ask human a question."""
        return self._checkpoint(questions=[question])

    def report_progress(self, message: str, percent: int):
        """Send progress update to dashboard (non-blocking)."""
        self._emit_event("progress", {"message": message, "percent": percent})

@dataclass
class AgentOutput:
    summary: str                    # markdown summary for human
    files_changed: list[FileChange] # diffs for review
    questions: list[str] = None     # open questions for human
    next_tasks: list[dict] = None   # suggest follow-up tasks

@dataclass  
class FileChange:
    path: str
    action: str  # "create"|"modify"|"delete"
    content: str # new content (for create/modify)
    diff: str    # unified diff (for modify)
```

### 2.2 Example Agent Implementation

```python
# agents/dev_agent/agent.py

from mitran_sdk import MitranAgent, TaskContext, AgentOutput, FileChange
from mitran_sdk.tools import github, llm

class DevAgent(MitranAgent):
    def execute(self, ctx: TaskContext) -> AgentOutput:
        # 1. Understand the task
        plan = llm.generate(
            f"Plan implementation for: {ctx.description}\n"
            f"Codebase context: {self._read_context(ctx)}"
        )

        # 2. Generate code
        changes = []
        for file_plan in plan.files:
            code = llm.generate(f"Write {file_plan.path}: {file_plan.spec}")
            changes.append(FileChange(
                path=file_plan.path,
                action="create",
                content=code,
                diff=self._make_diff(file_plan.path, code)
            ))

        # 3. Run tests
        test_result = self._run_tests(ctx.workspace_path)

        # 4. Return for human review
        return AgentOutput(
            summary=f"## Implementation Complete\n{plan.summary}\n\nTests: {test_result}",
            files_changed=changes,
            next_tasks=[{"title": "Write docs for new feature", "agent": "docs"}]
        )
```

### 2.3 Tool Registry

```python
# Tools available to agents via SDK

class ToolRegistry:
    builtin_tools = {
        # File operations (subject to ACL)
        "read_file": ReadFileTool,
        "write_file": WriteFileTool,
        "search_files": SearchFilesTool,

        # Git operations
        "git_commit": GitCommitTool,
        "git_branch": GitBranchTool,
        "git_diff": GitDiffTool,

        # Integration tools
        "github_create_pr": GitHubCreatePRTool,
        "github_create_repo": GitHubCreateRepoTool,
        "slack_notify": SlackNotifyTool,
        "grafana_create_dashboard": GrafanaDashboardTool,

        # LLM tools
        "llm_generate": LLMGenerateTool,
        "llm_analyze": LLMAnalyzeTool,

        # Shell (restricted to manifest.execute list)
        "shell_exec": ShellExecTool,
    }
```

## 3. Database Schema

```sql
-- Core tables (SQLite/PostgreSQL compatible)

CREATE TABLE tasks (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    agent_type TEXT NOT NULL,
    priority INTEGER DEFAULT 100,
    status TEXT DEFAULT 'queued',
    depends_on JSON DEFAULT '[]',
    resources JSON DEFAULT '[]',
    created_by TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    started_at TIMESTAMP,
    completed_at TIMESTAMP
);

CREATE TABLE checkpoints (
    id TEXT PRIMARY KEY,
    task_id TEXT REFERENCES tasks(id),
    agent_id TEXT NOT NULL,
    output TEXT NOT NULL,  -- markdown
    diffs JSON,           -- [{path, action, content, diff}]
    questions JSON,       -- ["question1", "question2"]
    decision_action TEXT, -- approve|reject|modify
    decision_feedback TEXT,
    decision_by TEXT,
    decision_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE agents (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    type TEXT NOT NULL,    -- dev|docs|cicd|tickets|wiki|ops|hr|dashboard
    status TEXT DEFAULT 'idle',
    manifest JSON NOT NULL,
    last_heartbeat TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE resource_locks (
    path TEXT NOT NULL,
    task_id TEXT REFERENCES tasks(id),
    mode TEXT NOT NULL,    -- read|write
    acquired_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (path, task_id)
);

CREATE TABLE events (
    id TEXT PRIMARY KEY,
    type TEXT NOT NULL,
    source TEXT NOT NULL,   -- agent_id or "human" or "system"
    payload JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE audit_log (
    id TEXT PRIMARY KEY,
    actor TEXT NOT NULL,    -- user or agent
    action TEXT NOT NULL,   -- what happened
    resource TEXT,          -- what was affected
    details JSON,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Ticketing (first-party app)
CREATE TABLE tickets (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    status TEXT DEFAULT 'open',
    priority TEXT DEFAULT 'medium',
    assignee TEXT,
    labels JSON DEFAULT '[]',
    sprint_id TEXT,
    created_by TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Wiki (first-party app)
CREATE TABLE wiki_pages (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT NOT NULL,  -- markdown
    parent_id TEXT REFERENCES wiki_pages(id),
    created_by TEXT NOT NULL,
    updated_by TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- HR (first-party app)
CREATE TABLE pto_requests (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    type TEXT NOT NULL,     -- vacation|sick|personal
    status TEXT DEFAULT 'pending',
    approved_by TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## 4. API Design

### 4.1 gRPC (Agent ↔ Core)

```protobuf
syntax = "proto3";
package mitran.v1;

service AgentService {
    // Agent registers with core on startup
    rpc Register(RegisterRequest) returns (RegisterResponse);

    // Agent receives task assignment
    rpc WaitForTask(WaitRequest) returns (stream TaskAssignment);

    // Agent submits checkpoint for human review
    rpc SubmitCheckpoint(CheckpointRequest) returns (CheckpointResponse);

    // Agent asks human a question (blocking)
    rpc AskHuman(QuestionRequest) returns (QuestionResponse);

    // Agent reports progress
    rpc ReportProgress(ProgressRequest) returns (Empty);

    // Agent requests file operation (goes through ACL)
    rpc FileOp(FileOpRequest) returns (FileOpResponse);
}
```

### 4.2 REST API (Dashboard ↔ Core)

```
GET    /api/v1/tasks                    # List tasks (filterable)
POST   /api/v1/tasks                    # Create task (human or NL)
PUT    /api/v1/tasks/:id/priority       # Reorder priority
POST   /api/v1/tasks/:id/assign         # Assign to specific agent

GET    /api/v1/checkpoints              # List pending checkpoints
POST   /api/v1/checkpoints/:id/approve  # Approve checkpoint
POST   /api/v1/checkpoints/:id/reject   # Reject with feedback

GET    /api/v1/agents                   # List agents + status
POST   /api/v1/agents/:id/restart       # Restart agent

GET    /api/v1/queue                    # Priority queue view
PUT    /api/v1/queue/reorder            # Bulk reorder

POST   /api/v1/init                     # The "mitran init" flow
GET    /api/v1/init/status              # Init progress

# Ticketing
GET    /api/v1/tickets
POST   /api/v1/tickets
PUT    /api/v1/tickets/:id

# Wiki
GET    /api/v1/wiki/pages
POST   /api/v1/wiki/pages
PUT    /api/v1/wiki/pages/:id

# NL interface
POST   /api/v1/chat                     # "Hey Mitran, do X"
```

### 4.3 WebSocket (Real-time Updates)

```
WS /api/v1/ws

Events pushed to dashboard:
- task.created
- task.started
- task.checkpoint    # agent produced output, human must review
- task.approved
- task.failed
- agent.status       # idle|working|error
- progress.update    # percent complete
```

## 5. Agent Specifications

### 5.1 Dev Agent
| Field | Value |
|-------|-------|
| Read | `src/**`, `docs/**`, `config/**`, `tests/**` |
| Write | `src/**`, `tests/**` |
| Execute | `go build`, `go test`, `npm run *`, `cargo *` |
| Integrations | GitHub (PRs, branch protection) |
| LLM tasks | Code generation, refactoring, test writing, PR description |

### 5.2 Docs Agent
| Field | Value |
|-------|-------|
| Read | `src/**`, `docs/**`, `wiki/**` |
| Write | `docs/**` |
| Execute | `mdbook build`, `mkdocs build` |
| Integrations | None |
| LLM tasks | HLD/LLD writing, ADRs, API docs, README generation |

### 5.3 CI/CD Agent
| Field | Value |
|-------|-------|
| Read | `.github/**`, `config/**`, `src/**` |
| Write | `.github/workflows/**`, `config/ci/**` |
| Execute | `gh workflow *`, `act` |
| Integrations | GitHub Actions |
| LLM tasks | Pipeline generation, flaky test analysis, build optimization |

### 5.4 Tickets Agent
| Field | Value |
|-------|-------|
| Read | `tickets/**` |
| Write | `tickets/**` |
| Execute | None |
| Integrations | Slack (notifications) |
| LLM tasks | Triage, priority suggestion, SLA tracking, sprint planning |

### 5.5 Wiki Agent
| Field | Value |
|-------|-------|
| Read | `src/**`, `docs/**`, `wiki/**`, `incidents/**` |
| Write | `wiki/**` |
| Execute | None |
| Integrations | Slack |
| LLM tasks | Runbook generation, onboarding docs, knowledge sync |

### 5.6 Ops Agent
| Field | Value |
|-------|-------|
| Read | `monitoring/**`, `config/**`, `logs/**` |
| Write | `monitoring/**`, `config/alerting/**` |
| Execute | `promtool check`, `amtool check-config` |
| Integrations | Grafana, Prometheus, Alertmanager, Slack |
| LLM tasks | Alert diagnosis, dashboard creation, runbook execution |

### 5.7 HR Agent
| Field | Value |
|-------|-------|
| Read | `hr/**` |
| Write | `hr/**` |
| Execute | None |
| Integrations | Slack |
| LLM tasks | PTO management, onboarding checklists, policy Q&A |

### 5.8 Dashboard Agent
| Field | Value |
|-------|-------|
| Read | `monitoring/**`, `src/**` |
| Write | `monitoring/dashboards/**` |
| Execute | None |
| Integrations | Grafana |
| LLM tasks | Dashboard creation, metric identification, visualization design |

## 6. Deployment & Installation

### 6.1 Single-Command Install
```bash
# Linux/macOS
curl -sSL https://install.mitran.dev | sh

# Verifies: Go (for potential plugins), Python 3.11+, Node.js 18+ (dashboard)
# Downloads: mitran binary (~50MB)
# Installs: Python agent dependencies, Prometheus, Grafana
```

### 6.2 Post-Install Setup
```bash
mitran setup
# → GitHub token
# → Slack bot token
# → LLM API key (Claude/OpenAI/local)
# → Workspace path

mitran init
# → Interactive team description
# → Generates full infrastructure
```

### 6.3 Process Architecture (Self-Hosted)
```
mitran serve (PID 1)
├── Core Engine (goroutines)
├── NATS (embedded)
├── HTTP/WS server (dashboard API)
├── gRPC server (agent communication)
├── Agent Pool:
│   ├── dev-agent (Python subprocess)
│   ├── docs-agent (Python subprocess)
│   ├── cicd-agent (Python subprocess)
│   ├── tickets-agent (Python subprocess)
│   ├── wiki-agent (Python subprocess)
│   ├── ops-agent (Python subprocess)
│   ├── hr-agent (Python subprocess)
│   └── dashboard-agent (Python subprocess)
├── Prometheus (child process, port 9090)
└── Grafana (child process, port 3000)
```
