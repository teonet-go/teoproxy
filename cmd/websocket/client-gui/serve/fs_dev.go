//go:build dev

package main

import (
	"io/fs"
	"os"
)

func getFrontendAssets() fs.FS {
	return os.DirFS("wasm")
}
