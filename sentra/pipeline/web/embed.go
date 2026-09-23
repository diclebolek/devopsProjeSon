package web

import _ "embed"

// IndexHTML is the lab dashboard. It is embedded so the ingest binary has no
// separate frontend build.
//
//go:embed index.html
var IndexHTML []byte
