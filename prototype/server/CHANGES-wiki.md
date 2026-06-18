# Wiki App - Route Registration

Add to `main.go` (or wherever routes are registered):

```go
wikiHandler := &handlers.WikiHandler{Store: store}
mux.Handle("/api/v1/wiki/", wikiHandler)
```

## Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | /api/v1/wiki/pages | List all wiki pages |
| POST | /api/v1/wiki/pages | Create a new page |
| GET | /api/v1/wiki/pages/:id | Get page by ID |
| PUT | /api/v1/wiki/pages/:id | Update page |
| DELETE | /api/v1/wiki/pages/:id | Delete page |
| GET | /api/v1/wiki/search?q=query | Full-text search title+content |
| GET | /api/v1/wiki/tree | Nested page tree structure |

## Dashboard

Import `WikiApp` from `components/WikiApp.tsx` and add route:

```tsx
import WikiApp from './components/WikiApp';

// In router:
<Route path="/wiki" element={<WikiApp />} />
```
