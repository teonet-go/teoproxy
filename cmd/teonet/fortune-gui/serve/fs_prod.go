//go:build !dev

package main

import (
	"embed"
	"io/fs"
)

//go:embed wasm
var embedFrontend embed.FS

func getFrontendAssets() fs.FS {
	f, err := fs.Sub(embedFrontend, "wasm")
	if err != nil {
		panic(err)
	}

	return f
}
