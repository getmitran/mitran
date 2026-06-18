package middleware

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"
)

type PageParams struct {
	Page     int
	PageSize int
	Offset   int
}

func ParsePagination(r *http.Request) PageParams {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	return PageParams{
		Page:     page,
		PageSize: pageSize,
		Offset:   (page - 1) * pageSize,
	}
}

func PaginatedResponse(w http.ResponseWriter, total, page, pageSize int, items interface{}) {
	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	resp := map[string]interface{}{
		"data": items,
		"meta": map[string]int{
			"total":       total,
			"page":        page,
			"page_size":   pageSize,
			"total_pages": totalPages,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
