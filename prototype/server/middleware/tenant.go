package middleware

import (
	"context"
	"net/http"
)

const tenantContextKey contextKey = "team_id"

// TenantMiddleware extracts team_id from JWT claims or X-Team-ID header
// and injects it into the request context.
func TenantMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		teamID := r.Header.Get("X-Team-ID")

		// If not in header, try to extract from JWT claims in context
		if teamID == "" {
			if claims, ok := r.Context().Value(contextKey("jwt_claims")).(map[string]interface{}); ok {
				if tid, ok := claims["team_id"].(string); ok {
					teamID = tid
				}
			}
		}

		if teamID == "" {
			http.Error(w, `{"error":"missing team_id"}`, http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), tenantContextKey, teamID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// TenantFromContext retrieves the team_id from the request context.
func TenantFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(tenantContextKey).(string); ok {
		return v
	}
	return ""
}
