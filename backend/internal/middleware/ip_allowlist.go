package middleware

import (
	"context"
	"net"
	"net/http"
	"strings"

	"github.com/atenimedia-llc/ncs-online/backend/internal/models"
	"github.com/atenimedia-llc/ncs-online/backend/internal/repository"
	"github.com/atenimedia-llc/ncs-online/backend/internal/response"
)

// EnforceIPAllowlist rejects requests whose source IP isn't on the
// authenticated user's allowlist. If the user hasn't configured any
// entries the middleware is a no-op, so it's safe to wrap broadly.
//
// Must run after Authenticate so models.CtxUserID is set.
func EnforceIPAllowlist(sec *repository.SecurityRepo) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, _ := r.Context().Value(models.CtxUserID).(string)
			if userID == "" {
				next.ServeHTTP(w, r)
				return
			}
			entries, err := sec.ActiveWhitelistFor(r.Context(), userID)
			if err != nil || len(entries) == 0 {
				next.ServeHTTP(w, r)
				return
			}
			ip := net.ParseIP(realIP(r))
			if ip != nil && matchesAllowlist(ip, entries) {
				next.ServeHTTP(w, r)
				return
			}
			response.Err(w, http.StatusForbidden, "IP_BLOCKED",
				"Your account is restricted to specific IP addresses. This network is not allowed.")
		})
	}
}

func matchesAllowlist(ip net.IP, entries []string) bool {
	for _, e := range entries {
		e = strings.TrimSpace(e)
		if strings.Contains(e, "/") {
			_, network, err := net.ParseCIDR(e)
			if err == nil && network.Contains(ip) {
				return true
			}
			continue
		}
		if other := net.ParseIP(e); other != nil && other.Equal(ip) {
			return true
		}
	}
	return false
}

// helper kept package-private so the new file links against the existing
// realIP() defined in middleware.go without re-declaring it.
var _ context.Context = nil
