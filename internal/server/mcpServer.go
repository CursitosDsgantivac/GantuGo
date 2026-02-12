package server

import (
	"encoding/json"
	"fmt"
	"log"
	"main/internal/utils"
	"net/http"
	"time"
)

/*
MCP (Model Context Protocol) Server implementation.
This server provides endpoints for MCP-compatible clients to interact with the service.
*/

// MCPRequest represents a standard MCP request structure
type MCPRequest struct {
	Method string                 `json:"method"`
	Params map[string]interface{} `json:"params,omitempty"`
	ID     string                 `json:"id,omitempty"`
}

// MCPResponse represents a standard MCP response structure
type MCPResponse struct {
	Result interface{} `json:"result,omitempty"`
	Error  *MCPError   `json:"error,omitempty"`
	ID     string      `json:"id,omitempty"`
}

// MCPError represents an error in MCP protocol
type MCPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Start initializes and starts the MCP server
func Start() {
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("GET /health", handleHealth)

	// MCP protocol endpoint
	mux.HandleFunc("POST /mcp", handleMCPRequest)

	// Server info endpoint
	mux.HandleFunc("GET /info", handleInfo)

	// Catch-all for undefined routes
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "404 - MCP endpoint not found", http.StatusNotFound)
	})

	port := ":" + utils.LoadPort()

	server := &http.Server{
		Addr:         port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("MCP Server started on %s", port)
	log.Fatal(server.ListenAndServe())
}

// handleHealth handles health check requests
func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "healthy",
		"service": "MCP Server",
		"time":    time.Now().Unix(),
	})
}

// handleInfo provides server information
func handleInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"name":        "GantuGo MCP Server",
		"version":     "1.0.0",
		"protocol":    "MCP",
		"description": "Model Context Protocol server implementation",
	})
}

// handleMCPRequest processes MCP protocol requests
func handleMCPRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req MCPRequest
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	if err := decoder.Decode(&req); err != nil {
		sendError(w, &MCPError{
			Code:    -32700,
			Message: fmt.Sprintf("Parse error: %v", err),
		}, "")
		return
	}

	// Route based on method
	switch req.Method {
	case "ping":
		sendResult(w, map[string]string{"message": "pong"}, req.ID)
	case "echo":
		sendResult(w, req.Params, req.ID)
	case "getContext":
		sendResult(w, map[string]interface{}{
			"context": "MCP server context",
			"params":  req.Params,
		}, req.ID)
	default:
		sendError(w, &MCPError{
			Code:    -32601,
			Message: fmt.Sprintf("Method not found: %s", req.Method),
		}, req.ID)
	}
}

// sendResult sends a successful MCP response
func sendResult(w http.ResponseWriter, result interface{}, id string) {
	response := MCPResponse{
		Result: result,
		ID:     id,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// sendError sends an error MCP response
func sendError(w http.ResponseWriter, err *MCPError, id string) {
	response := MCPResponse{
		Error: err,
		ID:    id,
	}
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(response)
}
