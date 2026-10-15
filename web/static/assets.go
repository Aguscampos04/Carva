// Package static embeds the web client served by Carva.
package static

import "embed"

// Files contains the web interface assets.
//
//go:embed index.html css js
var Files embed.FS
