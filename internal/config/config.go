// Package config has config variables
package config

import "github.com/w1lam/Packages/menu"

var Style = menu.Config{
	Width:         60,
	RenderHeaders: true,
	Padding:       3,
	Margin:        10,
}

type Config struct {
	RenderHeaders bool
	Width         int
}
