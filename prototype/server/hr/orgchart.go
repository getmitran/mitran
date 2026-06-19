package hr

import (
	"encoding/json"
	"net/http"
	"github.com/getmitran/mitran/server/apierr"
)

type OrgNode struct {
	Employee      Employee  `json:"employee"`
	DirectReports []OrgNode `json:"direct_reports"`
}

func BuildOrgChart(dir *Directory, rootID string) *OrgNode {
	emps := dir.Employees()
	return buildTree(emps, rootID)
}

func buildTree(emps map[string]Employee, rootID string) *OrgNode {
	emp, ok := emps[rootID]
	if !ok {
		return nil
	}
	node := &OrgNode{Employee: emp}
	for id, e := range emps {
		if e.ManagerID == rootID && id != rootID {
			if child := buildTree(emps, id); child != nil {
				node.DirectReports = append(node.DirectReports, *child)
			}
		}
	}
	return node
}

func GetTeamSize(dir *Directory, managerID string) int {
	emps := dir.Employees()
	return countReports(emps, managerID)
}

func countReports(emps map[string]Employee, managerID string) int {
	count := 0
	for id, e := range emps {
		if e.ManagerID == managerID && id != managerID {
			count += 1 + countReports(emps, id)
		}
	}
	return count
}

func GetManagementChain(dir *Directory, empID string) []Employee {
	emps := dir.Employees()
	var chain []Employee
	cur, ok := emps[empID]
	for ok && cur.ManagerID != "" && cur.ManagerID != cur.ID {
		mgr, exists := emps[cur.ManagerID]
		if !exists {
			break
		}
		chain = append(chain, mgr)
		cur, ok = mgr, true
	}
	return chain
}

func RegisterOrgChartRoutes(mux *http.ServeMux, dir *Directory) {
	mux.HandleFunc("GET /api/v1/hr/org", func(w http.ResponseWriter, r *http.Request) {
		root := r.URL.Query().Get("root")
		if root == "" {
			apierr.WriteError(w, apierr.BadRequest("root param required"))
			return
		}
		json.NewEncoder(w).Encode(BuildOrgChart(dir, root))
	})
	mux.HandleFunc("GET /api/v1/hr/org/chain/{id}", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(GetManagementChain(dir, r.PathValue("id")))
	})
	mux.HandleFunc("GET /api/v1/hr/org/team-size/{id}", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]int{"team_size": GetTeamSize(dir, r.PathValue("id"))})
	})
}
