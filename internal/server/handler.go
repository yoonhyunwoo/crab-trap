package server

import (
	"encoding/json"
	"io"
	"net/http"
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
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`
<!DOCTYPE html>
<html>
<head>
    <title>Moltbook Prompt Injector</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .container { max-width: 1200px; }
        .stats { display: grid; grid-template-columns: repeat(4, 1fr); gap: 20px; margin-bottom: 30px; }
        .stat { background: #f4f4f4; padding: 20px; border-radius: 8px; text-align: center; }
        .stat h3 { margin: 0 0 10px 0; color: #666; }
        .stat .value { font-size: 32px; font-weight: bold; color: #333; }
        .logs { background: #fff; border: 1px solid #ddd; border-radius: 8px; padding: 20px; }
        .log-entry { border-bottom: 1px solid #eee; padding: 15px 0; }
        .log-entry:last-child { border-bottom: none; }
        .timestamp { color: #666; font-size: 14px; }
        .method { font-weight: bold; color: #d73a49; }
        .url { font-family: monospace; color: #005cc5; }
        .details { margin-top: 10px; font-size: 13px; color: #444; }
        button { background: #005cc5; color: white; border: none; padding: 10px 20px; border-radius: 5px; cursor: pointer; }
        button:hover { background: #004e8c; }
    </style>
</head>
<body>
    <div class="container">
        <h1>🦞 Moltbook Prompt Injector</h1>
        <div class="stats">
            <div class="stat">
                <h3>Total Requests</h3>
                <div class="value" id="total">0</div>
            </div>
            <div class="stat">
                <h3>Unique IPs</h3>
                <div class="value" id="unique-ips">0</div>
            </div>
            <div class="stat">
                <h3>POST Requests</h3>
                <div class="value" id="post">0</div>
            </div>
            <div class="stat">
                <h3>Last Request</h3>
                <div class="value" id="last">-</div>
            </div>
        </div>
        <button onclick="refreshLogs()">Refresh Logs</button>
        <div class="logs" id="logs"></div>
    </div>
    <script>
        async function loadLogs() {
            const response = await fetch('/logs');
            const data = await response.json();
            
            const logsContainer = document.getElementById('logs');
            logsContainer.innerHTML = '';
            
            let total = 0;
            let postCount = 0;
            const ips = new Set();
            let lastTime = '-';
            
            data.logs.forEach(log => {
                total++;
                ips.add(log.remote_addr.split(':')[0]);
                if (log.method === 'POST') postCount++;
                lastTime = new Date(log.timestamp).toLocaleTimeString();
                
                const entry = document.createElement('div');
                entry.className = 'log-entry';
                entry.innerHTML = \`
                    <div class="timestamp">\${new Date(log.timestamp).toLocaleString()}</div>
                    <div><span class="method">\${log.method}</span> <span class="url">\${log.url}</span></div>
                    <div class="details">Remote: \${log.remote_addr} | User-Agent: \${log.user_agent || 'N/A'}</div>
                    \${log.body ? \`<div class="details">Body: \${log.body.substring(0, 200)}...</div>\` : ''}
                \`;
                logsContainer.appendChild(entry);
            });
            
            document.getElementById('total').textContent = total;
            document.getElementById('unique-ips').textContent = ips.size;
            document.getElementById('post').textContent = postCount;
            document.getElementById('last').textContent = lastTime;
        }
        
        function refreshLogs() {
            loadLogs();
        }
        
        loadLogs();
        setInterval(loadLogs, 10000);
    </script>
</body>
</html>
`))
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
