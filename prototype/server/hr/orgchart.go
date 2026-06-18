package hr

import "net/http"

type OrgNode struct {
	Employee      Employee  `json:"employee"`
	DirectReports []OrgNode `json:"direct_reports"`
}

func BuildOrgChart(dir *Directory, rootID string) *OrgNode {
	emp, ok := dir.Employees[rootID]
	if !ok {
		return nil
	}
	node := &OrgNode{Employee: emp}
	for id, e := range dir.Employees {
		if e.ManagerID == rootID && id != rootID {
			if child := BuildOrgChart(dir, id); child != nil {
				node.DirectReports = append(node.DirectReports, *child)
			}
		}
	}
	return node
}

func GetTeamSize(dir *Directory, managerID string) int {
	count := 0
	for id, e := range dir.Employees {
		if e.ManagerID == managerID && id != managerID {
			count += 1 + GetTeamSize(dir, id)
		}
	}
	return count
}

func GetManagementChain(dir *Directory, empID string) []Employee {
	var chain []Employee
	for cur, ok := dir.Employees[empID]; ok && cur.ManagerID != "" && cur.ManagerID != cur.ID; cur, ok = dir.Employees[cur.ManagerID] {
		if mgr, exists := dir.Employees[cur.ManagerID]; exists {
			chain = append(chain, mgr)
		} else {
			break
		}
	}
	return chain
}

func RegisterOrgChartRoutes(mux *http.ServeMux, dir *Directory) {
	mux.HandleFunc("/api/v1/hr/org", func(w http.ResponseWriter, r *http.Request) {
		root := r.URL.Query().Get("root")
		if root == "" {
			httpError(w, http.StatusBadRequest, "root param required")
			return
		}
		jsonResponse(w, BuildOrgChart(dir, root))
	})
	mux.HandleFunc("/api/v1/hr/org/chain/", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Path[len("/api/v1/hr/org/chain/"):]
		jsonResponse(w, GetManagementChain(dir, id))
	})
	mux.HandleFunc("/api/v1/hr/org/team-size/", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Path[len("/api/v1/hr/org/team-size/"):]
		jsonResponse(w, map[string]int{"team_size": GetTeamSize(dir, id)})
	})
}
