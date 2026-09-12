package httpx

import (
	"encoding/json"
	"net/http"
)

type code string

const (
	Error_invalid_id        code = "invalid_id"
	Error_not_found         code = "not_found"
	Error_internal_error    code = "internal_error"
	Error_malformed_json    code = "malformed_json"
	Error_validation_failed code = "validation_failed"
	Error_unauthenticated   code = "authenticated_failed"
	Error_forbdiden         code = "forbidden"
	Error_conflict          code = "conflict"
	Error_rate_limited      code = "rate_limited"
)

type errorEnvelope struct {
	Error errorPayload `json:"error"`
}

type errorPayload struct {
	ErrorMessage string `json:"message"`
	ErrorCode    code   `json:"code"`
	ErrorField   string `json:"field,omitempty"`
}

func Error(w http.ResponseWriter, status int, message string, code code) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(errorEnvelope{
		Error: errorPayload{
			ErrorMessage: message,
			ErrorCode:    code,
		}})
}

func ValidationError(w http.ResponseWriter, status int, message string, code code, field string) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(errorEnvelope{
		Error: errorPayload{
			ErrorMessage: message,
			ErrorCode:    code,
			ErrorField:   field,
		}})
}
