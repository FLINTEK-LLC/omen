// Package web embeds OMEN's static frontend (Leaflet map, dashboard,
// watchlist management) so it ships inside the single Go binary.
package web

import (
	"embed"
	"io/fs"
)

//go:embed static
var embedded embed.FS

// Files is the embedded frontend, rooted so that static/index.html is
// served at "/".
var Files fs.FS

func init() {
	sub, err := fs.Sub(embedded, "static")
	if err != nil {
		panic(err)
	}
	Files = sub
}
