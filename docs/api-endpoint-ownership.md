# API Endpoint Ownership

This document is the authoritative route ownership map for runtime binaries.
It exists to avoid route drift across `main` files and transport-specific muxes.

## Unified Map Server (main listener)

- Route composer: `pkg/httpapi/register.go`
- Route providers:
  - `pkg/httpapi/server_web.go`
  - `pkg/httpapi/handlers_core.go`
  - `pkg/auth/*` via `RegisterConfig.AuthManager`
- Admin API routes: all `/api/admin/*` must be registered through
  `pkg/httpapi/register.go` using `RegisterConfig`.

## Unified MCP Server (MCP listener)

- Route composer: `cmd/unified-server/mcp_register.go`
- MCP transports:
  - `/mcp-http`
  - `/mcp/sse`
- MCP docs/UI:
  - `/mcp-api/*`
- Policy: the MCP listener also owns the MCP REST mirror (`/api/*`) so
  transport capabilities are consistent with the standalone `cmd/mcp-server`
  binary.

## Standalone MCP Server

- Route composer: `cmd/mcp-server/main.go`
- Routes:
  - MCP transports (`/mcp-http`, `/mcp/sse`)
  - MCP REST mirror (`/api/*`) via `cmd/mcp-server/rest.go`
  - Swagger (`/mcp-api/*`)

## Guardrails

- Avoid adding `/api/*` registrations directly in `cmd/*/main.go` when a
  registry/composer exists.
- When adding an API endpoint:
  1. Add handler implementation in the owning package.
  2. Register via the owner composer file.
  3. Add or update route inventory tests.
