package ui

import (
	"github.com/tanema/amore/gfx"
)

var (
	PlayImg, _    = gfx.NewImage("images/play.png")
	PauseImg, _   = gfx.NewImage("images/pause.png")
	StopImg, _    = gfx.NewImage("images/stop.png")
	ForwardImg, _ = gfx.NewImage("images/forward.png")
	RewindImg, _  = gfx.NewImage("images/rewind.png")
	NextImg, _    = gfx.NewImage("images/next.png")
	BackImg, _    = gfx.NewImage("images/back.png")

	Black = gfx.NewColor(0, 0, 0, 255)
	Clear = gfx.NewColor(0, 0, 0, 0)
)
