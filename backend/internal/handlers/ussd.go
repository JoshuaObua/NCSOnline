package handlers

import (
	"crypto/subtle"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/atenimedia-llc/ncs-online/backend/internal/services"
)

type USSDHandler struct {
	service        *services.USSDService
	callbackSecret string
}

func (h *USSDHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.callbackSecret != "" {
		supplied := r.Header.Get("X-USSD-Token")
		if subtle.ConstantTimeCompare([]byte(supplied), []byte(h.callbackSecret)) != 1 {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
	}
	var input services.USSDInput
	if strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		if json.NewDecoder(r.Body).Decode(&input) != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
	} else {
		if r.ParseForm() != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		input = services.USSDInput{SessionID: r.FormValue("sessionId"), ServiceCode: r.FormValue("serviceCode"), PhoneNumber: r.FormValue("phoneNumber"), Text: r.FormValue("text")}
	}
	output := h.service.Handle(r.Context(), input)
	prefix := "END "
	if output.Continue {
		prefix = "CON "
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(prefix + output.Message))
}
