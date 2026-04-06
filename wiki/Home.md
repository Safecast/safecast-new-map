# Safecast New Map - Wiki

Welcome to the **Safecast New Map** wiki! This comprehensive documentation covers all features of the platform.

**Live Demo:** [simplemap.safecast.org](https://simplemap.safecast.org)

**GitHub Repository:** [Safecast/safecast-new-map](https://github.com/Safecast/safecast-new-map)

---

## 📖 Quick Links

| Section | Description |
|---------|-------------|
| [🚀 Getting Started](Getting-Started) | Installation, quick start, and basic configuration |
| [📥 Data Import & Export](Data-Import-Export) | Supported formats, bulk import, and automated sync |
| [🔌 API Documentation](API-Documentation) | REST API endpoints, authentication, and usage |
| [🤖 MCP Server & AI Integration](MCP-Server-AI-Integration) | Model Context Protocol, Claude integration, and AI tools |
| [💾 Database Setup](Database-Setup) | PostgreSQL, DuckDB, SQLite, ClickHouse configuration |
| [🔐 User Authentication](User-Authentication) | Registration, login, API keys, and security |
| [⚙️ Admin Panel](Admin-Panel) | User management, uploads, translations, and analytics |
| [🗺️ Map Features](Map-Features) | Interactive map, spectrum viewer, and visualization |
| [📊 Spectral Analysis](Spectral-Analysis) | Gamma spectroscopy, isotope identification |
| [🚢 Deployment](Deployment) | Production deployment, Docker, HTTPS, CloudFront |
| [💻 Development](Development) | Building from source, testing, and contributing |
| [⚙️ Configuration Reference](Configuration-Reference) | Complete list of flags and options |
| [🛠️ Database Maintenance](Database-Maintenance) | Sequence reset, migrations, and utilities |
| [🌍 Internationalization](Internationalization) | 29 languages, translation management |

---

##  Overview

**Safecast New Map** is a modern, self-hosted radiation monitoring platform that provides environmental radiation data collection, visualization, and analysis. Built to help communities understand radiation levels through open data and transparent mapping.

### Key Capabilities

✅ **Multi-format Data Import** - KML, JSON, CSV, RCTRK, GPX, LOG files  
✅ **Automated Data Sync** - Poll Safecast API and real-time sensors  
✅ **Interactive Map** - Clustered markers, multiple coloring schemes, spectrum viewer  
✅ **REST API** - Complete API with rate limiting and authentication  
✅ **MCP Server** - AI-powered data access via Claude  
✅ **Multi-database Support** - PostgreSQL, DuckDB, SQLite, ClickHouse  
✅ **User Authentication** - Email verification, API keys, session management  
✅ **Admin Panel** - User management, uploads, translations  
✅ **29 Languages** - Full internationalization with admin translation UI  
✅ **Spectral Analysis** - Gamma spectrum parsing and isotope identification  
✅ **Auto-update** - Self-upgrading system with rollback  
✅ **Docker Support** - Easy deployment with containers  

---

## 🏗️ Architecture

The platform consists of three main components unified into a single binary:

1. **Map Server** - Interactive map UI and REST API (port 8765)
2. **MCP Server** - Model Context Protocol for AI integration
3. **Web Chat** - Claude-powered chat interface for data queries

```
┌─────────────────────────────────────────────┐
│         Unified Server (8765)               │
├──────────────┬──────────────┬───────────────┤
│  Map UI +    │  MCP Server  │  Web Chat     │
│  REST API    │  /mcp-http   │  /assistant/  │
│              │  /mcp/sse    │               │
└──────────────┴──────────────┴───────────────┘
                        │
                        ▼
              ┌──────────────────┐
              │   Database       │
              │ (PostgreSQL/     │
              │  DuckDB/SQLite)  │
              └──────────────────┘
```

---

## 📊 Production Scale

The production database at [simplemap.safecast.org](https://simplemap.safecast.org) contains:

- **518.4 million** total markers (with zoom levels)
- **111.0 million** original measurements
- **407.3 million** zoom-optimized markers
- Data spanning **2012-2025**

---

## 🏛️ Community & Support

**Developed by:**
- Matvey Gladkikh (primary designer/developer)
- Rob Oudendijk (developer)
- Safecast volunteers worldwide

**Contributing:** We welcome contributions! Check [GitHub Issues](https://github.com/Safecast/safecast-new-map/issues) for areas needing help.

**License:**
- Code: [Apache 2.0](https://github.com/Safecast/safecast-new-map/blob/main/LICENSE)
- Data: [CC0 1.0 Universal](https://github.com/Safecast/safecast-new-map/blob/main/LICENSE.CC0)

---

## 🔗 External Resources

- [GitHub Repository](https://github.com/Safecast/safecast-new-map)
- [Releases](https://github.com/Safecast/safecast-new-map/releases)
- [GitHub Actions](https://github.com/Safecast/safecast-new-map/actions)
- [Live Demo](https://simplemap.safecast.org)
- [Map API Docs](https://simplemap.safecast.org/map-api/)
- [MCP API Docs](https://simplemap.safecast.org/mcp-api/)
