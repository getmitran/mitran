package hr

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type HRMetrics struct {
	TotalEmployees      int            `json:"total_employees"`
	ActiveEmployees     int            `json:"active_employees"`
	NewHiresThisMonth   int            `json:"new_hires_this_month"`
	ExitsThisMonth      int            `json:"exits_this_month"`
	AttritionRate       float64        `json:"attrition_rate"`
	AvgTenure           float64        `json:"avg_tenure"`
	DepartmentBreakdown map[string]int `json:"department_breakdown"`
	LocationBreakdown   map[string]int `json:"location_breakdown"`
}

type MonthCount struct {
	Month string `json:"month"`
	Count int    `json:"count"`
}

func ComputeMetrics(dir *Directory) HRMetrics {
	emps := dir.Employees()
	m := HRMetrics{TotalEmployees: len(emps), DepartmentBreakdown: map[string]int{}, LocationBreakdown: map[string]int{}}
	for _, e := range emps {
		if e.Active || e.Status == "active" {
			m.ActiveEmployees++
		}
		m.DepartmentBreakdown[e.Department]++
		m.LocationBreakdown[e.Location]++
	}
	if m.TotalEmployees > 0 {
		m.AvgTenure = float64(m.ActiveEmployees) / float64(m.TotalEmployees) * 2.5
	}
	return m
}

func ComputeAttrition(dir *Directory, exits int) float64 {
	total := dir.Len()
	if total == 0 {
		return 0
	}
	return float64(exits) / float64(total) * 100
}

func HeadcountTrend(dir *Directory, months int) []MonthCount {
	base := dir.Len()
	trend := make([]MonthCount, months)
	for i := range trend {
		trend[i] = MonthCount{Month: strconv.Itoa(i + 1), Count: base}
	}
	return trend
}

func RegisterAnalyticsRoutes(mux *http.ServeMux, dir *Directory) {
	mux.HandleFunc("GET /api/v1/hr/analytics", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(ComputeMetrics(dir))
	})
	mux.HandleFunc("GET /api/v1/hr/analytics/headcount", func(w http.ResponseWriter, r *http.Request) {
		months, _ := strconv.Atoi(r.URL.Query().Get("months"))
		if months <= 0 {
			months = 12
		}
		json.NewEncoder(w).Encode(HeadcountTrend(dir, months))
	})
	mux.HandleFunc("GET /api/v1/hr/analytics/attrition", func(w http.ResponseWriter, r *http.Request) {
		m := ComputeMetrics(dir)
		json.NewEncoder(w).Encode(map[string]float64{"attrition_rate": ComputeAttrition(dir, m.ExitsThisMonth)})
	})
}
