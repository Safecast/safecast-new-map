package main

import "embed"

// embeddedHintsFS carries the on-disk hint JSON templates into the compiled
// binary. Used as the default source when MCP_HINTS_DIR is unset (e.g. in
// production where the binary lives at /usr/local/bin/ without a neighbouring
// hints/ directory).
//
//go:embed hints/*.json
var embeddedHintsFS embed.FS
