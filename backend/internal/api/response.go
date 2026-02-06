package api

import (
	"encoding/json"
	"net/http"
)

type responseError struct {
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
}

type responseSuccess struct {
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func writeRespErr(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	resp := responseError{
		Code:    code,
		Message: msg,
	}
	json.NewEncoder(w).Encode(resp)
}

func writeRespWithStatus(w http.ResponseWriter, msg string, data any, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	resp := responseSuccess{
		Code:    code,
		Message: msg,
		Data:    data,
	}
	json.NewEncoder(w).Encode(resp)
}

func writeRespOk(w http.ResponseWriter, msg string, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	resp := responseSuccess{
		Code:    http.StatusOK,
		Message: msg,
		Data:    data,
	}
	json.NewEncoder(w).Encode(resp)
}
