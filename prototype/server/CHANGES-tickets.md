# Ticketing Module - Route Wiring

Add these to `main.go` where routes are registered:

```go
// Ticketing routes
ticketH := &handlers.TicketHandler{Store: store}
sprintH := &handlers.SprintHandler{Store: store}
mux.Handle("/api/v1/tickets", ticketH)
mux.Handle("/api/v1/tickets/", ticketH)
mux.Handle("/api/v1/sprints", sprintH)
mux.Handle("/api/v1/sprints/", sprintH)
```

## Endpoints Added

### Tickets
| Method | Path | Description |
|--------|------|-------------|
| GET | /api/v1/tickets | List tickets (filter: status, assignee, sprint_id, labels) |
| POST | /api/v1/tickets | Create ticket |
| GET | /api/v1/tickets/:id | Get ticket by ID |
| PUT | /api/v1/tickets/:id | Update ticket fields |
| DELETE | /api/v1/tickets/:id | Delete ticket |
| POST | /api/v1/tickets/:id/comments | Add comment |
| GET | /api/v1/tickets/:id/comments | List comments |

### Sprints
| Method | Path | Description |
|--------|------|-------------|
| GET | /api/v1/sprints | List all sprints |
| POST | /api/v1/sprints | Create sprint |
| PUT | /api/v1/sprints/:id | Update sprint |

## Dashboard

Import `TicketingApp` from `components/TicketingApp` and add a route/sidebar entry:

```tsx
import TicketingApp from './components/TicketingApp'
// In router: <Route path="/tickets" element={<TicketingApp />} />
// In Sidebar: { name: 'Tickets', path: '/tickets', icon: '🎫' }
```
