package dag

import "fmt"

type DAG struct {
	nodes     map[string]bool
	edges     map[string][]string // node -> dependencies
	completed map[string]bool
}

func New() *DAG {
	return &DAG{
		nodes:     make(map[string]bool),
		edges:     make(map[string][]string),
		completed: make(map[string]bool),
	}
}

func (d *DAG) AddNode(id string) {
	d.nodes[id] = true
	if _, ok := d.edges[id]; !ok {
		d.edges[id] = nil
	}
}

func (d *DAG) AddEdge(from, to string) error {
	if !d.nodes[from] || !d.nodes[to] {
		return fmt.Errorf("both nodes must exist: %s -> %s", from, to)
	}
	d.edges[from] = append(d.edges[from], to)
	return nil
}

func (d *DAG) DetectCycle() bool {
	white := make(map[string]bool)
	gray := make(map[string]bool)

	for n := range d.nodes {
		white[n] = true
	}

	var dfs func(string) bool
	dfs = func(n string) bool {
		delete(white, n)
		gray[n] = true
		for _, dep := range d.edges[n] {
			if gray[dep] {
				return true
			}
			if white[dep] {
				if dfs(dep) {
					return true
				}
			}
		}
		delete(gray, n)
		return false
	}

	for n := range d.nodes {
		if white[n] {
			if dfs(n) {
				return true
			}
		}
	}
	return false
}

func (d *DAG) TopologicalSort() ([]string, error) {
	if d.DetectCycle() {
		return nil, fmt.Errorf("graph has a cycle")
	}

	visited := make(map[string]bool)
	var result []string

	var visit func(string)
	visit = func(n string) {
		if visited[n] {
			return
		}
		visited[n] = true
		for _, dep := range d.edges[n] {
			visit(dep)
		}
		result = append(result, n)
	}

	for n := range d.nodes {
		visit(n)
	}

	// reverse
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}
	return result, nil
}

func (d *DAG) MarkCompleted(id string) {
	d.completed[id] = true
}

// GetReady returns nodes whose all dependencies are completed.
func (d *DAG) GetReady() []string {
	var ready []string
	for n := range d.nodes {
		if d.completed[n] {
			continue
		}
		allDone := true
		for _, dep := range d.edges[n] {
			if !d.completed[dep] {
				allDone = false
				break
			}
		}
		if allDone {
			ready = append(ready, n)
		}
	}
	return ready
}
