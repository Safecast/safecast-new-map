# Safecast New Map - Wiki Structure

## Overview

This document describes the wiki structure created for the Safecast New Map project on Codeberg.

**Total Pages:** 15 comprehensive wiki pages
**Format:** Codeberg wiki markdown (compatible with Git-based wikis)
**Location:** `/wiki/` directory

---

## Wiki Pages

### 1. Home.md
**Purpose:** Main landing page with overview and quick links
**Content:**
- Project overview and live demo link
- Quick navigation to all wiki pages
- Key capabilities summary
- Architecture diagram
- Production scale statistics
- Community and support links

### 2. Getting-Started.md
**Purpose:** Quick start and installation guide
**Content:**
- 4 deployment options (Binary, Production Data, Docker, Source)
- First configuration guide
- Database setup basics
- Feature enablement
- Production deployment basics
- Verification steps
- Troubleshooting

### 3. Data-Import-Export.md
**Purpose:** Comprehensive data import and export guide
**Content:**
- Supported formats (KML, JSON, CSV, RCTRK, GPX, LOG, spectra)
- File upload via web interface
- Bulk import methods
- Automated data sync (Safecast API fetcher, real-time sensors)
- Data export options
- API-based imports
- Data validation
- Migration tools
- Best practices

### 4. API-Documentation.md
**Purpose:** Complete REST API reference
**Content:**
- Radiation data endpoints
- Real-time sensor endpoints
- Spectroscopy endpoints
- Reference and statistics endpoints
- Authentication endpoints
- Admin endpoints
- Data export endpoints
- Short links
- Rate limiting
- Error responses
- GPT-optimized endpoints
- Authentication methods

### 5. MCP-Server-AI-Integration.md
**Purpose:** MCP server and AI features guide
**Content:**
- Overview and components
- 16 MCP tools documentation
- Claude integration (CLI and web)
- Web chat interface
- Model adapter
- Semantic caching with RAG
- DuckLake analytics
- AI query logging
- MCP transports (HTTP, SSE)
- Best practices

### 6. Database-Setup.md
**Purpose:** Database configuration guide
**Content:**
- PostgreSQL setup with PostGIS
- DuckDB configuration
- SQLite setup
- ClickHouse setup
- Database schema (all tables)
- Migrations
- Performance optimization
- Backup and restore
- Monitoring
- Troubleshooting

### 7. User-Authentication.md
**Purpose:** Authentication and user management
**Content:**
- Enable authentication
- User registration
- Login methods (password, API key)
- API key management
- Password management
- User profile pages
- Session management
- Authentication logging
- Admin user management
- Database schema
- Migration
- Best practices

### 8. Admin-Panel.md
**Purpose:** Admin panel features and usage
**Content:**
- Enable admin panel
- Users management
- Uploads management
- MCP analytics
- Real-time sensors
- Translations management (29 languages)
- Cache management
- Admin API reference
- Security best practices
- Common admin tasks
- Troubleshooting

### 9. Map-Features.md
**Purpose:** Interactive map and visualization features
**Content:**
- Map interface and controls
- Radiation markers and clustering
- Speed-based layers
- Coloring schemes (scientific vs safety)
- Legend and units
- Real-time sensors display
- Spectrum viewer
- Location search
- URL parameters
- Print mode with QR codes
- User profile pages
- Language support
- Performance features
- Accessibility

### 10. Spectral-Analysis.md
**Purpose:** Gamma spectroscopy and isotope identification
**Content:**
- Supported spectrum formats (.spe, .n42, .rctrk)
- Spectrum database schema
- Upload spectrum data
- Spectrum viewer features
- Isotope identification
- Spectrum export formats
- Analysis tools (MCP integration)
- Migration for spectrum support
- Architecture and data flow
- Best practices
- Troubleshooting

### 11. Deployment.md
**Purpose:** Production deployment guide
**Content:**
- Deployment options comparison
- Binary deployment (systemd, supervisor)
- Docker deployment (compose, volumes)
- HTTPS with Let's Encrypt
- Reverse proxy setup (Nginx, Apache, Caddy)
- CloudFront/CDN deployment
- Self-upgrade system
- Production infrastructure example
- Database backup
- Monitoring
- Security hardening
- Performance tuning
- High availability
- Troubleshooting

### 12. Development.md
**Purpose:** Development workflow and contributing
**Content:**
- Development environment setup
- Build from source
- Testing strategies
- Project structure
- Code style guide
- Contributing workflow
- Adding features (API endpoints, database drivers, MCP tools)
- Frontend development
- API documentation generation
- Debugging tips
- Common development tasks
- See also links

### 13. Configuration-Reference.md
**Purpose:** Complete configuration options reference
**Content:**
- Server configuration flags
- Database configuration flags
- Map defaults
- Authentication flags
- Email configuration
- Data sync flags
- Data import/export flags
- Self-upgrade flags
- Environment variables
- Configuration examples
- Flag precedence
- Database connection strings
- Configuration validation

### 14. Database-Maintenance.md
**Purpose:** Database maintenance and optimization
**Content:**
- PostgreSQL maintenance (sequences, statistics, vacuum, indexes)
- SQLite maintenance
- DuckDB maintenance
- ClickHouse maintenance
- Backup strategies
- Migration tools
- Populate usernames
- Fix recording dates
- Performance monitoring
- Scheduled maintenance
- Troubleshooting

### 15. Internationalization.md
**Purpose:** Multi-language support and translation management
**Content:**
- 29 supported languages
- Language selection priority
- Translation architecture
- Admin panel for translations
- Translated components
- Translation file format
- Incremental seeding
- Migration for translations
- Translation best practices
- Contributing translations
- API reference
- Troubleshooting

---

## Wiki Navigation

All pages include:
- **Back to Home** link at the top
- **See Also** section at the bottom with related pages
- Consistent formatting and structure
- Code examples with proper syntax highlighting
- Tables for reference data
- Cross-references between pages

---

## How to Use This Wiki

### On Codeberg

1. Clone the wiki repository:
   ```bash
   git clone https://codeberg.org/Safecast/safecast-new-map.wiki.git
   ```

2. Add markdown files:
   ```bash
   cd safecast-new-map.wiki
   cp /path/to/wiki/*.md .
   git add *.md
   git commit -m "Add comprehensive wiki documentation"
   git push
   ```

3. The wiki will be available at:
   ```
   https://codeberg.org/Safecast/safecast-new-map/wiki
   ```

### Wiki Features on Codeberg

- Markdown rendering
- Page navigation sidebar
- Search functionality
- Edit button for each page
- Version history
- Clone via HTTPS or SSH

---

## Content Quality

### Features

✅ **Comprehensive coverage** - All project features documented
✅ **Up-to-date** - Based on latest main branch
✅ **Practical examples** - Real-world usage examples
✅ **Code snippets** - Ready-to-use commands and configurations
✅ **Cross-references** - Links between related topics
✅ **Troubleshooting** - Common issues and solutions
✅ **Best practices** - Production-ready recommendations
✅ **API documentation** - Complete REST API reference
✅ **Database schema** - All tables and relationships
✅ **Configuration reference** - All flags and environment variables

### Style Guidelines

- Clear, concise language
- Consistent formatting
- Proper markdown syntax
- Code blocks with language specification
- Tables for structured data
- Bullet points for lists
- Numbered steps for procedures
- Warning/note callouts where needed

---

## Maintenance

### Keeping Wiki Updated

When new features are added:
1. Update relevant wiki pages
2. Add new pages if needed
3. Update cross-references
4. Commit and push changes

### Review Process

Before major releases:
1. Review all pages for accuracy
2. Test code examples
3. Update version-specific information
4. Check cross-references
5. Update Home page if needed

---

## Future Enhancements

Potential additions:
- Video tutorials
- Architecture diagrams (Mermaid)
- API playground examples
- Performance benchmarks
- Migration guides between versions
- FAQ section
- Glossary of terms
- Changelog documentation

---

## Contributing

To contribute to the wiki:
1. Clone the wiki repository
2. Make edits or additions
3. Test markdown rendering
4. Submit pull request or push directly (if you have write access)

---

**Created:** 2026-04-05
**Based on:** Latest main branch (commit: b1937b9)
**Total pages:** 15
**Total features documented:** 50+
**API endpoints documented:** 40+
**MCP tools documented:** 16
**Supported languages:** 29
