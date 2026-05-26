package response

import (
	"encoding/json"
	"net/http"
)

type envelope struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

type errEnvelope struct {
	Success bool      `json:"success"`
	Error   *errBody  `json:"error"`
}

type errBody struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

type Meta struct {
	Page    int   `json:"page"`
	PerPage int   `json:"per_page"`
	Total   int64 `json:"total"`
}

func JSON(w http.ResponseWriter, status int, data interface{}) {
	write(w, status, &envelope{Success: true, Data: data})
}

func JSONMsg(w http.ResponseWriter, status int, message string) {
	write(w, status, &envelope{Success: true, Message: message})
}

func JSONPaged(w http.ResponseWriter, status int, data interface{}, meta *Meta) {
	write(w, status, &envelope{Success: true, Data: data, Meta: meta})
}

func Err(w http.ResponseWriter, status int, code, message string) {
	write(w, status, &errEnvelope{Success: false, Error: &errBody{Code: code, Message: message}})
}

func ErrDetails(w http.ResponseWriter, status int, code, message string, details interface{}) {
	write(w, status, &errEnvelope{Success: false, Error: &errBody{Code: code, Message: message, Details: details}})
}

func ValidationErr(w http.ResponseWriter, details interface{}) {
	ErrDetails(w, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "Validation failed", details)
}

func write(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
