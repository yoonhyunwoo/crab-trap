package server

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"
)

type RequestLog struct {
	Timestamp   time.Time              `json:"timestamp"`
	Method      string                 `json:"method"`
	URL         string                 `json:"url"`
	Headers     map[string][]string    `json:"headers"`
	QueryParams map[string][]string    `json:"query_params"`
	Body        string                 `json:"body"`
	RemoteAddr  string                 `json:"remote_addr"`
	UserAgent   string                 `json:"user_agent"`
}

type Handler struct {
	logger *Logger
}

func NewHandler(logger *Logger) *Handler {
	return &Handler{
		logger: logger,
	}
}

func (h *Handler) HandleRequest(w http.ResponseWriter, r *http.Request) {
	log := h.logRequest(r)

	if h.logger != nil {
		h.logger.Log(log)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Request logged successfully",
		"timestamp": log.Timestamp,
	})
}

func (h *Handler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "healthy",
		"time":   time.Now().UTC(),
	})
}

func (h *Handler) HandleLogs(w http.ResponseWriter, r *http.Request) {
	if h.logger == nil {
		http.Error(w, "Logger not configured", http.StatusInternalServerError)
		return
	}

	logs := h.logger.GetAllLogs()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"count":   len(logs),
		"logs":    logs,
	})
}

func (h *Handler) HandleUI(w http.ResponseWriter, r *http.Request) {
	html, err := os.ReadFile("internal/server/ui.html")
	if err != nil {
		http.Error(w, "UI not found", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(html)
}
func (h *Handler) logRequest(r *http.Request) *RequestLog {
	body, _ := io.ReadAll(r.Body)
	r.Body.Close()

	return &RequestLog{
		Timestamp:   time.Now().UTC(),
		Method:      r.Method,
		URL:         r.URL.String(),
		Headers:     r.Header,
		QueryParams: r.URL.Query(),
		Body:        string(body),
		RemoteAddr:  r.RemoteAddr,
		UserAgent:   r.UserAgent(),
	}
}
