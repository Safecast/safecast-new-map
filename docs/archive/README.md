# Archived Documentation

This directory contains historical planning documents and implementation plans that have either been completed or are no longer actively pursued.

## Status Legend
- ✅ **IMPLEMENTED** - Feature was built and is now part of the production system
- 🗄️ **ARCHIVED** - Historical document kept for reference
- ⏸️ **DEFERRED** - Plan exists but not currently being pursued

## Archived Documents

### IMPLEMENTATION_PLAN.md
**Status:** 🗄️ ARCHIVED - Not implemented  
**Date Archived:** 2026-02-22  
**Reason:** Plan for integrating Claude AI chat directly into the map interface. Not pursued in favor of the unified architecture where:
- Web chat is served at `/assistant/` on the main map server (port 8765)
- MCP server integration available at `/mcp-http` for external AI clients (Claude Code, Claude.ai)
- See [README.md](../../README.md#mcp-server--ai-integration) for current architecture

### user-login-profile-plan.md
**Status:** ✅ IMPLEMENTED  
**Date Archived:** 2026-02-22  
**Implemented In:** See [README.md](../../README.md#user-authentication--api-keys) sections on:
- User Registration & Login
- API Key Authentication
- User profile pages
- Authentication logging

### mcp-server-plan.md
**Status:** ✅ IMPLEMENTED (evolved)  
**Date Archived:** 2026-02-22  
**Original Implementation:** Separate [safecast-map-MCP](https://github.com/Safecast/safecast-map-MCP) repository  
**Current State:** MCP server functionality merged into unified server binary (`cmd/unified-server/`)  
**Live Endpoint:** https://simplemap.safecast.org/mcp-http  
**Documentation:** See [README.md](../../README.md#mcp-server--ai-integration)

---

**Note:** For current project plans and documentation, see the main [README.md](../../README.md).
