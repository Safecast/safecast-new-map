package mcpserver

import (
	"net/http"

	"github.com/mark3labs/mcp-go/server"
)

type TransportConfig struct {
	Mux        *http.ServeMux
	MCPServer  *server.MCPServer
	BaseURL    string
	Middleware func(http.Handler) http.Handler
}

// RegisterTransports wires /mcp-http and /mcp/sse on a mux.
func RegisterTransports(cfg TransportConfig) {
	if cfg.Mux == nil || cfg.MCPServer == nil {
		return
	}
	sseServer := server.NewSSEServer(cfg.MCPServer,
		server.WithBaseURL(cfg.BaseURL),
		server.WithStaticBasePath("/mcp"),
	)
	httpServer := server.NewStreamableHTTPServer(cfg.MCPServer,
		server.WithEndpointPath("/mcp-http"),
	)

	httpHandler := http.Handler(httpServer)
	sseHandler := http.Handler(sseServer)
	if cfg.Middleware != nil {
		httpHandler = cfg.Middleware(httpHandler)
		sseHandler = cfg.Middleware(sseHandler)
	}

	cfg.Mux.Handle("/mcp-http", httpHandler)
	cfg.Mux.Handle("/mcp/", sseHandler)
}
