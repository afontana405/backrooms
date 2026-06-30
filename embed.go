package main

import "embed"

// Embed all application assets into the binary.
// The engine reads from this FS instead of the OS filesystem.
// Add new directories here as your app grows.

//go:embed pages public chefscript-dictionary.md
var embeddedAssets embed.FS
