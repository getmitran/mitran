package middleware

import (
	"log"
	"net/http"
	"runtime/debug"
	"github.com/getmitran/mitran/server/apierr"
)

func RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %v\n%s", err, debug.Stack())
				apierr.WriteError(w, apierr.Internal("Internal Server Error"))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
