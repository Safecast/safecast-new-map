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

- Entrypoint: `cmd/unified-server/mcp_register.go`
- Shared MCP composers:
  - `pkg/mcpserver/tools.go`
  - `pkg/mcpserver/transports.go`
  - `pkg/mcpserver/routes.go`
  - `pkg/mcpserver/docs.go`
- MCP transports:
  - `/mcp-http`
  - `/mcp/sse`
- MCP docs/UI:
  - `/mcp-api/*`
- Policy: the MCP listener owns the MCP REST mirror (`/api/*`) using the same
  shared registrars as standalone `cmd/mcp-server`, so capabilities stay aligned.

## Standalone MCP Server

- Entrypoint: `cmd/mcp-server/main.go`
- Shared MCP composers:
  - `pkg/mcpserver/tools.go`
  - `pkg/mcpserver/transports.go`
  - `pkg/mcpserver/routes.go`
  - `pkg/mcpserver/docs.go`
- Routes:
  - MCP transports (`/mcp-http`, `/mcp/sse`)
  - MCP REST mirror (`/api/*`)
  - Swagger (`/mcp-api/*`)

## Guardrails

- Avoid adding `/api/*` registrations directly in `cmd/*/main.go` when a
  registry/composer exists.
- For MCP surface changes (tools/transports/REST/docs):
  1. Update shared registrars under `pkg/mcpserver/*`.
  2. Keep `cmd/*` entrypoints as thin wiring only.
- When adding an API endpoint:
  1. Add handler implementation in the owning package.
  2. Register via the owner composer/registrar.
  3. Add or update route inventory tests.
