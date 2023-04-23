package main

import (
	"image/color"
)

type LocalColor struct{}

func (lc *LocalColor) blue() color.Color {
	return &color.RGBA{5, 55, 247, 1}
}
