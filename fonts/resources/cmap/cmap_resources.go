// Package cmapresources provides embedded CMap files for PDF font parsing.
package cmapresources

import "embed"

//go:embed *
var Files embed.FS
