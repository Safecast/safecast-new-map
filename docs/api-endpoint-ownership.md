# API Endpoint Ownership

This document is the authoritative route ownership map for runtime binaries.
It exists to avoid route drift across `main` files and transport-specific muxes.

## Unified Map Server (main listener)

- Route composer: `pkg/httpapi/register.go`
- Route providers:
  - `pkg/httpapi/server_web.go`
  - `pkg/httpapi/handlers_core.go`
  - `cmd/unified-server/rest.go` handlers injected via `RegisterConfig` for:
    - `/api/sensors`
    - `/api/sensors/export`
    - `/api/sensor/*`
    - `/api/feedback`
    - `GET /api/track/{id}/insights`
  - `pkg/auth/*` via `RegisterConfig.AuthManager`
- Admin API routes: all `/api/admin/*` must be registered through
  `pkg/httpapi/register.go` using `RegisterConfig`.
- Legacy non-API routes are registered via `pkg/httpapi/legacy_routes.go`:
  - `/upload`
  - `/upload/progress`
  - `/get_markers`
  - `/stream_markers`
  - `/realtime_history`
  - `/trackid/*`
  - `/tracks/*`
- UI/admin page routes are registered via `pkg/httpapi/page_routes.go`:
  - `/home`
  - `/stories.html`
  - `/data/stories.json`
  - `/`
  - `/profile`
  - `/reset-password`
  - `/admin/users`
  - `/admin/uploads`
  - `/admin/mcp`
  - `/admin/realtime`
  - `/admin/translations`
- Static asset routes are registered via `pkg/httpapi/static_routes.go`:
  - `/static/*`
  - `/js/*`
- System/infrastructure routes are registered via `pkg/httpapi/system_routes.go`:
  - `/selfupgrade/*`

## Unified MCP Surface (same runtime, MCP listener)

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
- Policy: the MCP listener owns the MCP REST mirror (`/api/*`) via
  `pkg/mcpserver/*` + `cmd/unified-server/rest.go`. No standalone MCP runtime exists.

## Guardrails

- Avoid adding `/api/*` registrations directly in `cmd/*/main.go` when a
  registry/composer exists.
- Avoid adding legacy route registrations directly in `cmd/*/main.go`; use
  `pkg/httpapi/legacy_routes.go`.
- Avoid adding page/admin route registrations directly in `cmd/*/main.go`; use
  `pkg/httpapi/page_routes.go`.
- Avoid adding static route registrations directly in `cmd/*/main.go`; use
  `pkg/httpapi/static_routes.go`.
- Avoid adding system route registrations directly in `cmd/*/main.go`; use
  `pkg/httpapi/system_routes.go`.
- Keep `cmd/unified-server/mcp_register.go` companion wiring limited to non-API
  parity routes (for example `/chat`).
- For MCP surface changes (tools/transports/REST/docs):
  1. Update shared registrars under `pkg/mcpserver/*`.
  2. Keep `cmd/unified-server/*` entrypoints as thin wiring only.
- When adding an API endpoint:
  1. Add handler implementation in the owning package.
  2. Register via the owner composer/registrar.
  3. Add or update route inventory tests.
