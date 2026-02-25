package server

/*
Installed with go get github.com/modelcontextprotocol/go-sdk


claude desktop json

{
  "mcpServers": {
    "my-go-mcp-server": {
      "type": "streamable-http",
      "url": "http://localhost:8080/mcp"
    }
  }
}





*/
import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type GreetInput struct {
	Name string `json:"name" jsonschema:"The name of the person to greet"`
}

// GreetOutput is the typed output for the "greet" tool.
type GreetOutput struct {
	Message string `json:"message" jsonschema:"The greeting message"`
}

func SayHi(
	_ context.Context,
	_ *mcp.CallToolRequest,
	in GreetInput,
) (*mcp.CallToolResult, GreetOutput, error) {
	name := in.Name
	if name == "" {
		name = "World"
	}
	return nil, GreetOutput{Message: "Hello, " + name + "! 👋"}, nil
}

func StartMCPBinary() {
	server := mcp.NewServer(&mcp.Implementation{Name: "greeter", Version: "v1.0.0"}, nil)
	mcp.AddTool(server, &mcp.Tool{Name: "greet", Description: "say hi"}, SayHi)
	// Run the server over stdin/stdout, until the client disconnects.
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}

}

func StartMCP() {
	// Create a server with a single tool.
	server := mcp.NewServer(&mcp.Implementation{Name: "greeter", Version: "v1.0.0"}, nil)
	mcp.AddTool(server, &mcp.Tool{Name: "greet", Description: "say hi"}, SayHi)
	// Run the server over stdin/stdout, until the client disconnects.
	// if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
	// 	log.Fatal(err)
	// }

	// 3. Create a StreamableHTTP handler.
	//    This is the modern MCP transport (replaces the old SSE transport).
	//    A single /mcp endpoint handles both POST (requests) and GET (SSE stream).

	// commented to check if the sse handler works with claude desktop
	// handler := mcp.NewStreamableHTTPHandler(
	// 	func(_ *http.Request) *mcp.Server { return server },
	// 	nil, // use default StreamableHTTPOptions
	// )

	handler := mcp.NewSSEHandler(
		func(_ *http.Request) *mcp.Server { return server },
		nil,
	)

	// 4. Mount the handler and start the HTTP server.
	mux := http.NewServeMux()
	// mux.Handle("/mcp", handler) // to use with streamble http handler
	mux.Handle("/sse", handler) // to use with sse

	// Optional: a simple health-check endpoint.
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	addr := ":8080"
	log.Printf("MCP server listening on http://localhost%s/mcp", addr)
	log.Printf("Health check:           http://localhost%s/health", addr)
	log.Print("check the integration with claude web.ai using ngrook check the claude conversation if I make it to work I would remove the exec docs and the github action")

	if err := http.ListenAndServe(addr, corsMiddleware(mux)); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

// corsMiddleware adds permissive CORS headers so that browser-based MCP
// clients (e.g. Claude.ai web) can reach a locally running server.
// Tighten this for production deployments.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			origin = "*"
		}
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", strings.Join([]string{
			"Content-Type",
			"Accept",
			"Mcp-Session-Id",
			"Last-Event-ID",
		}, ", "))
		w.Header().Set("Access-Control-Expose-Headers", "Mcp-Session-Id")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
