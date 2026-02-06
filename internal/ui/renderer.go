// Package ui holds ui renders and headers ONLY visuals
package ui

import (
	"io"
)

type PlainRenderer struct {
	Out io.Writer
}

type Renderer interface {
	PackageMenuRenderer
	ProgressRenderer
}
