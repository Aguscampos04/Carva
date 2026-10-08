// Package migrations embeds the versioned SQL migrations for Carva.
package migrations

import "embed"

// Files contains the versioned SQL migrations in this directory.
//
//go:embed *.sql
var Files embed.FS
