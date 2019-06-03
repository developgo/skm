package main

import (
	"github.com/fatih/color"
)

var (
	// Version is the default version of SKM
	Version = "0.5"
	logo    = `

███████╗██╗  ██╗███╗   ███╗
██╔════╝██║ ██╔╝████╗ ████║
███████╗█████╔╝ ██╔████╔██║
╚════██║██╔═██╗ ██║╚██╔╝██║
███████║██║  ██╗██║ ╚═╝ ██║
╚══════╝╚═╝  ╚═╝╚═╝     ╚═╝

SKM V%s
https://github.com/developgo/skm

`
)

func displayLogo() {
	color.Cyan(logo, Version)
}
