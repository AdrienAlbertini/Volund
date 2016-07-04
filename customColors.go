package main

import (
	"github.com/fatih/color"
)

var boldRed *color.Color
var boldGreen *color.Color
var boldYellow *color.Color
var boldCyan *color.Color

func initCustomColors() {
	boldRed = color.New(color.FgRed, color.Bold)
	boldGreen = color.New(color.FgGreen, color.Bold)
	boldYellow = color.New(color.FgYellow, color.Bold)
	boldCyan = color.New(color.FgCyan, color.Bold)
}
