package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"number-service/internal/domain"
)

type NumberHandler struct {
	service domain.NumberService
}

func NewNumberHandler(service domain.NumberService) *NumberHandler {
	return &NumberHandler{service: service}
}

type AddNumberRequest struct {
	Number int `json:"number"`
}

type AddNumberResponse struct {
	Numbers []int `json:"numbers"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (h *NumberHandler) AddNumber(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req AddNumberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, fmt.Sprintf("invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	numbers, err := h.service.AddNumber(req.Number)
	if err != nil {
		h.sendError(w, fmt.Sprintf("failed to add number: %v", err), http.StatusInternalServerError)
		return
	}

	h.sendJSON(w, AddNumberResponse{Numbers: numbers}, http.StatusOK)
}

func (h *NumberHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.sendError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (h *NumberHandler) sendJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *NumberHandler) sendError(w http.ResponseWriter, message string, status int) {
	h.sendJSON(w, ErrorResponse{Error: message}, status)
}
