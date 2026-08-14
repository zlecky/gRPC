package web

import "embed"

// FS holds the API console static assets.
//
//go:embed index.html styles.css app.js
var FS embed.FS
