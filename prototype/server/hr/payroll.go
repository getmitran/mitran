package hr

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type SalaryStructure struct {
	EmployeeID       string  `json:"employee_id"`
	Basic            float64 `json:"basic"`
	HRA              float64 `json:"hra"`
	DA               float64 `json:"da"`
	SpecialAllowance float64 `json:"special_allowance"`
	PF               float64 `json:"pf"`
	Tax              float64 `json:"tax"`
	NetPay           float64 `json:"net_pay"`
	Currency         string  `json:"currency"`
}

type Payslip struct {
	ID          string          `json:"id"`
	EmployeeID  string          `json:"employee_id"`
	Month       string          `json:"month"`
	Structure   SalaryStructure `json:"structure"`
	GeneratedAt string          `json:"generated_at"`
}

type PayrollStore struct {
	mu       sync.RWMutex
	salaries map[string]SalaryStructure
	payslips map[string][]Payslip
}

func NewPayrollStore() *PayrollStore {
	return &PayrollStore{salaries: make(map[string]SalaryStructure), payslips: make(map[string][]Payslip)}
}

func (s *PayrollStore) SetSalary(empID string, st SalaryStructure) {
	s.mu.Lock()
	defer s.mu.Unlock()
	st.EmployeeID = empID
	s.salaries[empID] = st
}

func (s *PayrollStore) GetSalary(empID string) *SalaryStructure {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if st, ok := s.salaries[empID]; ok {
		return &st
	}
	return nil
}

func (s *PayrollStore) GeneratePayslip(empID, month string) *Payslip {
	s.mu.Lock()
	defer s.mu.Unlock()
	st, ok := s.salaries[empID]
	if !ok {
		return nil
	}
	st.PF = st.Basic * 0.12
	gross := st.Basic + st.HRA + st.DA + st.SpecialAllowance
	st.Tax = (gross - st.PF - st.HRA) * 0.30
	st.NetPay = gross - st.PF - st.Tax
	p := Payslip{ID: fmt.Sprintf("PS-%s-%s", empID, month), EmployeeID: empID, Month: month, Structure: st, GeneratedAt: time.Now().Format(time.RFC3339)}
	s.payslips[empID] = append(s.payslips[empID], p)
	return &p
}

func (s *PayrollStore) ListPayslips(empID string) []Payslip {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.payslips[empID]
}

func (s *PayrollStore) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/hr/payroll/salary", func(w http.ResponseWriter, r *http.Request) {
		var st SalaryStructure
		json.NewDecoder(r.Body).Decode(&st)
		s.SetSalary(st.EmployeeID, st)
		json.NewEncoder(w).Encode(st)
	})
	mux.HandleFunc("GET /api/v1/hr/payroll/{empId}/payslips", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(s.ListPayslips(r.PathValue("empId")))
	})
	mux.HandleFunc("POST /api/v1/hr/payroll/{empId}/generate", func(w http.ResponseWriter, r *http.Request) {
		p := s.GeneratePayslip(r.PathValue("empId"), r.URL.Query().Get("month"))
		if p == nil {
			http.Error(w, "salary not configured", 404)
			return
		}
		json.NewEncoder(w).Encode(p)
	})
}
