// Package views provides the HTML templates for the broken application.
package views

import (
	"embed"
	_ "embed"
)

//go:embed list.html
var Content embed.FS
