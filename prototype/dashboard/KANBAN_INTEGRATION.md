# Kanban Board Integration Guide

Apply these changes to integrate the Kanban board into the existing dashboard.

---

## 1. `src/types.ts` — Add Comment type and extend Task

Add at the end of the file:

```typescript
export interface Comment {
  id: string
  text: string
  author: string
  created_at: string
}
```

Add to the `Task` interface:

```typescript
  labels?: string[]
  comments?: Comment[]
```

---

## 2. `src/api.ts` — Add updateTask method

Add inside the `api` object (after `reorderTask`):

```typescript
  updateTask: (id: string, fields: Partial<{ title: string; description: string; priority: number; status: string; labels: string[] }>) =>
    request(`/tasks/${id}`, { method: 'PATCH', body: JSON.stringify(fields) }),
```

---

## 3. `src/components/Sidebar.tsx` — Add Kanban nav item

Add import:
```typescript
import { ListOrdered, Bot, ShieldCheck, Ticket, BookOpen, Zap, LayoutGrid } from 'lucide-react'
```

Update the `nav` array — insert after 'queue':
```typescript
const nav = [
  { id: 'queue', label: 'Queue', icon: ListOrdered },
  { id: 'kanban', label: 'Kanban', icon: LayoutGrid },
  { id: 'agents', label: 'Agents', icon: Bot },
  { id: 'checkpoints', label: 'Checkpoints', icon: ShieldCheck },
  { id: 'tickets', label: 'Tickets', icon: Ticket },
  { id: 'wiki', label: 'Wiki', icon: BookOpen },
] as const
```

---

## 4. `src/App.tsx` — Add kanban view

Add import:
```typescript
import KanbanBoard from './components/KanbanBoard'
```

Update the View type:
```typescript
type View = 'queue' | 'kanban' | 'agents' | 'checkpoints' | 'tickets' | 'wiki'
```

Add render case (after the queue line):
```typescript
{view === 'kanban' && <KanbanBoard tasks={tasks} refetch={refetchTasks} />}
```
