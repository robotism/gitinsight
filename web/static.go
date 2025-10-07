package web

import (
	"embed"
)

//go:embed all:dist
var WebDistFs embed.FS

func init() {
	// xres.New(WebDistFs).DumpAll()
}
