// Package webassets embeds the compiled frontend (web/dist, copied here by
// `make web-build` or the api.Dockerfile) so the api binary serves both the
// UI and the REST endpoints from one process.
package webassets

import (
	"embed"
	"io/fs"
)

//go:embed dist
var embedded embed.FS

func FS() fs.FS {
	sub, err := fs.Sub(embedded, "dist")
	if err != nil {
		panic(err)
	}
	return sub
}
