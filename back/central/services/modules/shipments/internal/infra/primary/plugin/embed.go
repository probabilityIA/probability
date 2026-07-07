package plugin

import "embed"

//go:embed probability-shipping
var Files embed.FS

const (
	FolderName = "probability-shipping"
	Version    = "1.4.0"
)
