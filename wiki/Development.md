# Development

Guide to building from source, testing, and contributing to Safecast New Map.

[[Home|← Back to Home]]

---

## Development Environment

### Prerequisites

**Required:**
- Go 1.21 or later
- Git
- PostgreSQL (optional, for development)

**Optional:**
- Docker
- PostgreSQL with PostGIS
- DuckDB
- Node.js (for frontend development)

### Install Go

**Linux:**
```bash
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

**macOS:**
```bash
brew install go
```

**Windows:**
Download from [go.dev](https://go.dev/dl/)

---

## Build from Source

### Clone Repository

```bash
git clone https://github.com/Safecast/safecast-new-map.git
cd safecast-new-map
```

### Basic Build

```bash
# Build unified server
go build -o safecast-new-map ./cmd/unified-server

# Run
./safecast-new-map
```

### Build with DuckDB Support

```bash
# Requires CGO
CGO_ENABLED=1 go build -tags duckdb -o safecast-new-map ./cmd/unified-server
```

### Cross-Compilation

Build for multiple platforms:

```bash
# Run cross-compile script
go run scripts/crosscompile/crosscompile.go
```

This builds binaries for:
- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64)

### Build MCP Server (Legacy)

```bash
# Standalone MCP server
go build -o mcp-server ./cmd/mcp-server

# Standalone web chat
go build -o web-chat ./cmd/web-chat
```

**Note:** The unified server includes all functionality in a single binary.

---

## Run Tests

### Run All Tests

```bash
go test ./...
```

### Run Specific Package Tests

```bash
# Test auth package
go test ./pkg/auth/...

# Test database package
go test ./pkg/database/...

# Test spectrum package
go test ./pkg/spectrum/...
```

### Run Tests with Coverage

```bash
go test -cover ./...
```

### Run Tests with Verbose Output

```bash
go test -v ./...
```

### Run Benchmarks

```bash
go test -bench=. ./...
```

---

## Development Workflow

### Run with Database

**SQLite (no setup):**
```bash
go run ./cmd/unified-server -db-type sqlite -db-path ./dev.db
```

**PostgreSQL:**
```bash
# Create development database
createdb safecast_dev

# Run migrations
psql -d safecast_dev -f migrations/create_users_table.sql

# Run server
DATABASE_URL="postgres://localhost/safecast_dev" \
go run ./cmd/unified-server
```

### Run with Hot Reload

**Air (Go hot reload):**
```bash
# Install air
go install github.com/cosmtrek/air@latest

# Run with hot reload
air
```

### Debugging

**Delve debugger:**
```bash
# Install delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug
dlv debug ./cmd/unified-server
```

**Print debugging:**
```go
import "log"

log.Printf("Debug: value=%v", value)
```

---

## Project Structure

```
safecast-new-map/
├── cmd/
│   ├── unified-server/    # Main server (Map + MCP + Chat)
│   │   ├── main.go
│   │   ├── doc.go         # Swagger documentation
│   │   └── hints/         # MCP model hints
│   ├── mcp-server/        # Standalone MCP server (legacy)
│   ├── web-chat/          # Standalone web chat (legacy)
│   └── tools/             # Maintenance utilities
├── pkg/
│   ├── auth/              # Authentication handlers
│   ├── database/          # Database interface and drivers
│   │   └── drivers/       # Driver registration
│   ├── httpapi/           # HTTP routes and handlers
│   ├── safecast-fetcher/  # Safecast API client
│   ├── safecast-realtime/ # Real-time sensor polling
│   ├── spectrum/          # Spectrum parsing and analysis
│   ├── selfupgrade/       # Auto-update system
│   ├── jsonarchive/       # JSON archive generation
│   ├── email/             # SMTP email sending
│   └── logger/            # Logging utilities
├── migrations/            # Database migrations
├── docs/                  # Documentation and diagrams
├── web/                   # Frontend assets (HTML, CSS, JS)
├── scripts/
│   └── crosscompile/      # Cross-compilation script
└── test/                  # Test files and examples
```

---

## Code Style

### Go Style Guide

Follow standard Go conventions:
- Use `gofmt` or `goimports` for formatting
- Follow [Effective Go](https://go.dev/doc/effective_go)
- Use meaningful variable names
- Add comments for exported functions
- Write tests for new functionality

See `.windsurf/rules/go_style_guide.md` for project-specific style guide.

### Linting

```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
golangci-lint run
```

### Formatting

```bash
# Format code
gofmt -w .

# Or use goimports
goimports -w .
```

---

## Contributing

### Workflow

1. **Fork the repository**
2. **Create feature branch:**
   ```bash
   git checkout -b feature/your-feature
   ```
3. **Make changes**
4. **Write tests**
5. **Run tests:**
   ```bash
   go test ./...
   ```
6. **Commit changes:**
   ```bash
   git commit -m "feat: add your feature"
   ```
7. **Push to fork:**
   ```bash
   git push origin feature/your-feature
   ```
8. **Create Pull Request**

### Commit Message Format

Follow conventional commits:
- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation changes
- `chore:` - Maintenance tasks
- `test:` - Test changes
- `refactor:` - Code refactoring

**Examples:**
```
feat: add spectrum export in CSV format
fix: resolve duplicate key error in markers table
docs: update API documentation
chore: update dependencies
```

### Pull Request Template

See `.github/pull_request_template.md` for PR template.

**Include:**
- Description of changes
- Related issue number
- Testing performed
- Screenshots (if UI changes)

### Code Review

PRs will be reviewed for:
- Code quality
- Test coverage
- Documentation
- Performance impact
- Security considerations

---

## Adding Features

### New API Endpoint

1. **Add route in `pkg/httpapi/`:**
   ```go
   r.HandleFunc("/api/new-endpoint", handler.NewEndpoint).Methods("GET")
   ```

2. **Create handler:**
   ```go
   func NewEndpoint(w http.ResponseWriter, r *http.Request) {
       // Implementation
   }
   ```

3. **Add Swagger documentation:**
   ```go
   // @Summary New endpoint
   // @Description Description of endpoint
   // @Tags Tag
   // @Success 200 {object} ResponseType
   // @Router /api/new-endpoint [get]
   ```

4. **Regenerate docs:**
   ```bash
   cd cmd/unified-server && swag init -g doc.go -o docs/api --parseDependency --parseInternal --parseDependencyLevel 2 --instanceName unifiedapi
   ```

5. **Write tests**

### New Database Driver

1. **Create driver in `pkg/database/drivers/`**
2. **Register driver:**
   ```go
   import _ "your-driver-package"
   ```
3. **Implement database interface**
4. **Add tests**

### New MCP Tool

1. **Add tool in MCP server:**
   ```go
   tools = append(tools, mcp.NewTool("tool_name", ...))
   ```

2. **Implement handler:**
   ```go
   func handleToolName(params map[string]interface{}) (interface{}, error) {
       // Implementation
   }
   ```

3. **Add tool description and parameters**
4. **Write tests**
5. **Update documentation**

---

## Frontend Development

### Web Assets

Located in `web/` directory:
- HTML templates
- CSS stylesheets
- JavaScript modules
- Static assets (images, icons)

### Make Changes

1. **Edit files in `web/`**
2. **Test locally:**
   ```bash
   go run ./cmd/unified-server
   ```
3. **Open browser:**
   ```
   http://localhost:8765
   ```

### Frontend Framework

- **Leaflet** - Interactive maps
- **Canvas 2D** - Spectrum rendering
- **Vanilla JavaScript** - No framework dependency

### Add Translation

1. **Edit `translations.json`**
2. **Add key for all languages:**
   ```json
   {
     "new.key": {
       "en": "English text",
       "ja": "日本語",
       ...
     }
   }
   ```
3. **Run migration:**
   ```bash
   psql -d safecast -f migrations/add_ui_translations.sql
   ```
4. **Test in admin panel**

---

## API Documentation

### Regenerate Swagger Docs

**Map API:**
```bash
cd cmd/unified-server
swag init \
  -g doc.go \
  -o docs/api \
  --parseDependency \
  --parseInternal \
  --parseDependencyLevel 2 \
  --instanceName unifiedapi
```

**MCP API:**
```bash
cd cmd/mcp-server
swag init -g rest.go
```

### Swagger Annotations

Add to handler:
```go
// @Summary Get radiation data
// @Description Query radiation measurements near coordinates
// @Tags Radiation
// @Accept json
// @Produce json
// @Param lat query float64 true "Latitude"
// @Param lon query float64 true "Longitude"
// @Success 200 {array} models.Marker
// @Failure 400 {object} models.Error
// @Router /api/radiation [get]
func GetRadiation(w http.ResponseWriter, r *http.Request) {
    // Implementation
}
```

---

## Testing Strategies

### Unit Tests

Test individual functions:
```go
func TestParseSpectrum(t *testing.T) {
    data := []byte("...")
    spectrum, err := ParseSpectrum(data)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if len(spectrum.Channels) != 1024 {
        t.Errorf("expected 1024 channels, got %d", len(spectrum.Channels))
    }
}
```

### Integration Tests

Test database interactions:
```go
func TestDatabaseInsert(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()
    
    marker := models.Marker{Lat: 35.6762, Lon: 139.6503}
    err := db.InsertMarker(&marker)
    if err != nil {
        t.Fatalf("insert failed: %v", err)
    }
}
```

### API Tests

Test HTTP endpoints:
```go
func TestGetRadiation(t *testing.T) {
    server := startTestServer()
    defer server.Close()
    
    resp, err := http.Get(server.URL + "/api/radiation?lat=35.6762&lon=139.6503")
    if err != nil {
        t.Fatalf("request failed: %v", err)
    }
    defer resp.Body.Close()
    
    assert.Equal(t, http.StatusOK, resp.StatusCode)
}
```

### Mock Data

Use mock data for testing:
- Mock database responses
- Mock HTTP clients
- Mock external APIs

---

## Debugging Tips

### Enable Debug Logging

```bash
# Set log level
export LOG_LEVEL=debug
go run ./cmd/unified-server
```

### Database Query Logging

**PostgreSQL:**
```sql
-- Enable query logging
ALTER SYSTEM SET log_statement = 'all';
SELECT pg_reload_conf();
```

### Profile Application

```bash
# Enable pprof
go run ./cmd/unified-server -pprof

# Access profiles
curl http://localhost:8765/debug/pprof/profile
curl http://localhost:8765/debug/pprof/heap
curl http://localhost:8765/debug/pprof/goroutine
```

### Memory Profiling

```bash
# Generate heap profile
go tool pprof http://localhost:8765/debug/pprof/heap

# Analyze
(pprof) top
(pprof) list FunctionName
```

---

## Common Development Tasks

### Add Database Migration

1. **Create migration file:**
   ```bash
   touch migrations/your_migration.sql
   ```

2. **Write SQL:**
   ```sql
   ALTER TABLE your_table ADD COLUMN new_column TEXT;
   ```

3. **Test migration:**
   ```bash
   psql -d test_db -f migrations/your_migration.sql
   ```

4. **Document in wiki**

### Update Dependencies

```bash
# Update dependencies
go get -u ./...

# Tidy go.mod
go mod tidy

# Test after update
go test ./...
```

### Fix Duplicate Key Error

After database restore:
```bash
# Fix sequence
go run ./cmd/tools/fix-pg-sequence

# Or manually
psql -d safecast -c "SELECT setval('markers_id_seq', (SELECT MAX(id) FROM markers) + 1);"
```

See [FIX_DUPLICATE_KEY.md](/FIX_DUPLICATE_KEY.md).

---

## See Also

- [API Documentation](API-Documentation) - API reference and endpoints
- [Database Setup](Database-Setup) - Database configuration
- [Configuration Reference](Configuration-Reference) - All flags and options
- [GitHub Actions Guide](/GITHUB_ACTIONS_GUIDE.md) - CI/CD automation
- [Pull Request Template](/.github/pull_request_template.md) - PR template
- [Go Style Guide](/.windsurf/rules/go_style_guide.md) - Code style guide
