package middleware

import (
	"net/http"
)

// Logger logs the request into the console giving the admin access to the method, route and host address.
func (m *Middleware) Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.log.Info("Received request", "METHOD", r.Method, "ROUTE", r.RequestURI, "HOST", r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}
