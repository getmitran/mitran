package billing

import (
	"sync"
	"time"
)

type WorkspaceUsage struct {
	WorkspaceID   string
	Month         string
	TaskCount     int
	TokensUsed    int
	EstimatedCost float64
}

var (
	mu    sync.RWMutex
	store = make(map[string]*WorkspaceUsage) // key: workspaceID+month
)

func key(workspaceID, month string) string {
	return workspaceID + ":" + month
}

func RecordUsage(workspaceID string, tokens int, cost float64) {
	month := time.Now().Format("2006-01")
	mu.Lock()
	defer mu.Unlock()
	k := key(workspaceID, month)
	u, ok := store[k]
	if !ok {
		u = &WorkspaceUsage{WorkspaceID: workspaceID, Month: month}
		store[k] = u
	}
	u.TaskCount++
	u.TokensUsed += tokens
	u.EstimatedCost += cost
}

func GetUsage(workspaceID, month string) *WorkspaceUsage {
	mu.RLock()
	defer mu.RUnlock()
	return store[key(workspaceID, month)]
}

func GetAllUsage(month string) []WorkspaceUsage {
	mu.RLock()
	defer mu.RUnlock()
	var result []WorkspaceUsage
	for _, u := range store {
		if u.Month == month {
			result = append(result, *u)
		}
	}
	return result
}
