package cui

import (
	"github.com/nsf/termbox-go"
	"time"
)

type Textbox struct {
	Base      Controlbase
	Text      []rune
	Name      string
	events    events
	Events    events
	Box       *Boxchars
	Cursorpos int
	Maskchar  rune
}

func NewTextbox() *Textbox {
	b := &Textbox{}
	b.Base = Controlbase{}
	b.Events = NewEventBase()
	b.events = NewEventBase()
	b.Cursorpos = 0
	b.events.Onkey = func(ev Event) bool {
		if ev.Termbox.Key == termbox.KeyBackspace || ev.Termbox.Key == termbox.KeyBackspace2 {
			StatusLine("Backspace")
			if len(b.Text) > 0 && b.Cursorpos > 0 {
				//b.Text = b.Text[:len(b.Text)-1]

				copy(b.Text[b.Cursorpos-1:], b.Text[b.Cursorpos-1+1:])
				b.Text[len(b.Text)-1] = 0 // or the zero value of T
				b.Text = b.Text[:len(b.Text)-1]

				b.Cursorpos--
			}
			return true
		} else if ev.Termbox.Key == termbox.KeyDelete {
			StatusLine("DEL")
			if len(b.Text) > 0 && b.Cursorpos < len(b.Text) {
				copy(b.Text[b.Cursorpos:], b.Text[b.Cursorpos+1:])
				b.Text[len(b.Text)-1] = 0 // or the zero value of T
				b.Text = b.Text[:len(b.Text)-1]
			}
			return true
		} else if ev.Termbox.Key == termbox.KeyArrowLeft || ev.Termbox.Key == termbox.KeyArrowRight {
			if ev.Termbox.Key == termbox.KeyArrowLeft && b.Cursorpos > 0 {
				StatusLine("Cursor minused")
				b.Cursorpos--
			} else if ev.Termbox.Key == termbox.KeyArrowRight && b.Cursorpos < len(b.Text) {
				StatusLine("Cursor plussed")
				b.Cursorpos++
			}
			return true
		} else if ev.Termbox.Key == termbox.KeySpace {
			StatusLine("Space: \"" + " " + "\"")
			b.Text = append(b.Text, 0)
			copy(b.Text[b.Cursorpos+1:], b.Text[b.Cursorpos:])
			b.Text[b.Cursorpos] = ' '

			//b.Text = append(append(b.Text[:b.Cursorpos], ev.Termbox.Ch), b.Text[b.Cursorpos:]...)
			b.Cursorpos++
			return true
		} else if ev.Termbox.Ch != 0 {
			StatusLine("Rest: \"" + string(ev.Termbox.Ch) + "\"")
			b.Text = append(b.Text, 0)
			copy(b.Text[b.Cursorpos+1:], b.Text[b.Cursorpos:])
			b.Text[b.Cursorpos] = ev.Termbox.Ch

			//b.Text = append(append(b.Text[:b.Cursorpos], ev.Termbox.Ch), b.Text[b.Cursorpos:]...)
			b.Cursorpos++
			return true
		}
		go ev.Control.Draw(ev.Window)
		return false
	}
	b.Box = Box["s"]
	return b
}

func NewTextboxB(base Controlbase, Name string) *Textbox {
	b := &Textbox{}
	b.Name = Name
	b.Base = base
	b.Events = NewEventBase()
	b.events = NewEventBase()
	b.Cursorpos = 0
	b.events.Onkey = func(ev Event) bool {
		if ev.Termbox.Key == termbox.KeyBackspace || ev.Termbox.Key == termbox.KeyBackspace2 {
			StatusLine("Backspace")
			if len(b.Text) > 0 && b.Cursorpos > 0 {
				//b.Text = b.Text[:len(b.Text)-1]

				copy(b.Text[b.Cursorpos-1:], b.Text[b.Cursorpos-1+1:])
				b.Text[len(b.Text)-1] = 0 // or the zero value of T
				b.Text = b.Text[:len(b.Text)-1]

				b.Cursorpos--
			}
			ev.Control.Draw(ev.Window)
			return true
		} else if ev.Termbox.Key == termbox.KeyDelete {
			StatusLine("DEL")
			if len(b.Text) > 0 && b.Cursorpos < len(b.Text) {
				copy(b.Text[b.Cursorpos:], b.Text[b.Cursorpos+1:])
				b.Text[len(b.Text)-1] = 0 // or the zero value of T
				b.Text = b.Text[:len(b.Text)-1]
			}
			ev.Control.Draw(ev.Window)
			return true
		} else if ev.Termbox.Key == termbox.KeyArrowLeft || ev.Termbox.Key == termbox.KeyArrowRight {
			if ev.Termbox.Key == termbox.KeyArrowLeft && b.Cursorpos > 0 {
				StatusLine("Cursor minused")
				b.Cursorpos--
			} else if ev.Termbox.Key == termbox.KeyArrowRight && b.Cursorpos < len(b.Text) {
				StatusLine("Cursor plussed")
				b.Cursorpos++
			}
			ev.Control.Draw(ev.Window)
			return true
		} else if ev.Termbox.Key == termbox.KeySpace {
			StatusLine("Space: \"" + " " + "\"")
			b.Text = append(b.Text, 0)
			copy(b.Text[b.Cursorpos+1:], b.Text[b.Cursorpos:])
			b.Text[b.Cursorpos] = ' '

			//b.Text = append(append(b.Text[:b.Cursorpos], ev.Termbox.Ch), b.Text[b.Cursorpos:]...)
			b.Cursorpos++
			ev.Control.Draw(ev.Window)
			return true
		} else if ev.Termbox.Ch != 0 {
			StatusLine("Rest: \"" + string(ev.Termbox.Ch) + "\"")
			b.Text = append(b.Text, 0)
			copy(b.Text[b.Cursorpos+1:], b.Text[b.Cursorpos:])
			b.Text[b.Cursorpos] = ev.Termbox.Ch

			//b.Text = append(append(b.Text[:b.Cursorpos], evTermbox.Ch), b.Text[b.Cursorpos:]...)
			b.Cursorpos++
			ev.Control.Draw(ev.Window)
			return true
		}

		return false
	}
	b.Box = Box["s"]
	return b
}

func (b *Textbox) Draw(win *Window) {

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
		if b.Maskchar == 0 {
			if b.Cursorpos == k && win.SelectedControl == b {
				win.SetCell(X+k+1, Y+1, v, win.Foreground, win.Background)
			} else {
				win.SetCell(X+k+1, Y+1, v, fg, bg)
			}
		} else {
			if b.Cursorpos == k && win.SelectedControl == b {
				win.SetCell(X+k+1, Y+1, b.Maskchar, win.Foreground, win.Background)
			} else {
				win.SetCell(X+k+1, Y+1, b.Maskchar, fg, bg)
			}
		}
		if timeeffect {
			termbox.Flush()
			time.Sleep(10 * time.Millisecond)
		}

	}

	if b.Cursorpos == len(b.Text) && win.SelectedControl == b {
		//win.SetCell(X+k+1, Y+1, " ", win.Foreground, win.Background)
		win.SetCell(X+1+len(b.Text), Y+1, Box2["shadow"], fg, bg)
	}

	for k, v := range b.Name {
		win.SetCell(X+k+2, Y, v, fg, bg)

		if timeeffect {
			termbox.Flush()
			time.Sleep(10 * time.Millisecond)
		}
	}

	//win.SetCell(X+1+len(b.Text), Y+1, Box2["shadow"], fg, bg)

	termbox.Flush()
}

func (b *Textbox) B() *Controlbase {
	return &b.Base
}

func (b *Textbox) E() *events {
	return &b.Events
}

func (b *Textbox) intevents() *events {
	return &b.events
}
