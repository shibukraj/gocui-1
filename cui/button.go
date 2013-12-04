package cui

import (
	"github.com/nsf/termbox-go"
	"time"
)

type Button struct {
	Base   Controlbase
	Text   string
	Events events
	Box    *Boxchars
	Shadow bool
}

func NewButton() Button {
	b := Button{}
	b.Base = Controlbase{}
	b.Events = NewEventBase()
	b.Box = Box["d"]
	return b
}

func NewButtonB(base Controlbase, text string) *Button {
	b := &Button{}
	b.Text = text
	b.Base = base
	b.Box = Box["d"]
	return b
}

func (b *Button) Draw(win *Window) {

	var fg termbox.Attribute
	var bg termbox.Attribute

	var (
		X = b.Base.X
		Y = b.Base.Y
		W = b.Base.Width
		H = b.Base.Height
	)

	if b == win.SelectedControl {
		fg = win.SelectedForeground
		bg = win.SelectedBackground
	} else {
		fg = win.Foreground
		bg = win.Background
	}

	win.SetCell(X, Y, b.Box.tl, fg, bg)
	win.SetCell(X+W-1, Y, b.Box.tr, fg, bg)
	win.SetCell(X, Y+H-1, b.Box.bl, fg, bg)
	win.SetCell(X+W-1, Y+H-1, b.Box.br, fg, bg)

	for x := 1; x < W-1; x++ {
		win.SetCell(X+x, Y, b.Box.h, fg, bg)
		win.SetCell(X+x, Y+H-1, b.Box.h, fg, bg)
		if timeeffect {
			termbox.Flush()
			time.Sleep(10 * time.Millisecond)
		}
	}

	if b.Shadow {
		for x := 1; x <= W; x++ {
			win.SetCell(X+x-1, Y+H, Box2["shadow"], fg, bg)
		}

		for y := 1; y <= H+1; y++ {
			win.SetCell(X+W, Y+y-1, Box2["shadow"], fg, bg)
			if timeeffect {
				termbox.Flush()
				time.Sleep(10 * time.Millisecond)
			}
		}
	}

	for y := 1; y < H-1; y++ {
		win.SetCell(X, Y+y, b.Box.v, fg, bg)
		win.SetCell(X+W-1, Y+y, b.Box.v, fg, bg)
		if timeeffect {
			termbox.Flush()
			time.Sleep(10 * time.Millisecond)
		}
	}

	for x := 1; x < W-1; x++ {
		for y := 1; y < H-1; y++ {
			win.SetCell(X+x, Y+y, ' ', fg, bg)
			if timeeffect {
				termbox.Flush()
				time.Sleep(10 * time.Millisecond)
			}
		}
	}

	for k, v := range b.Text {
		win.SetCell(X+k+1, Y+1, v, fg, bg)
		if timeeffect {
			termbox.Flush()
			time.Sleep(10 * time.Millisecond)
		}
	}

	termbox.Flush()
}

func (b *Button) B() *Controlbase {
	return &b.Base
}

func (b *Button) E() *events {
	return &b.Events
}
