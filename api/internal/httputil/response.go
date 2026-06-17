package httputil

import (
	"encoding/json"
	"net/http"
)

type ErrorBody struct {
	Code    string          `json:"code"`
	Message string          `json:"message"`
	City    json.RawMessage `json:"city,omitempty"`
	Existing json.RawMessage `json:"existing,omitempty"`
}

type ErrorResponse struct {
	Error ErrorBody `json:"error"`
}

type DataResponse struct {
	Data any `json:"data"`
}

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func WriteData(w http.ResponseWriter, status int, data any) {
	WriteJSON(w, status, DataResponse{Data: data})
}

func WriteError(w http.ResponseWriter, status int, code, message string) {
	WriteJSON(w, status, ErrorResponse{
		Error: ErrorBody{Code: code, Message: message},
	})
}

func WriteErrorWithFields(w http.ResponseWriter, status int, body ErrorBody) {
	WriteJSON(w, status, ErrorResponse{Error: body})
}
