package cui

import (
	"container/list"
	"github.com/nsf/termbox-go"
	"time"
)

type Listbox struct {
	Base            Controlbase
	Caption         string
	SelectedElement *list.Element
	SelectedIndex   int
	Events          events
	Box             *Boxchars
	Stuff           *list.List
	Shadow          bool
}

type Listelement struct {
	Text   string
	Value  interface{}
	Parent *Listbox
}

func NewListbox() *Listbox {
	b := &Listbox{}
	b.Stuff = list.New()
	b.Base = Controlbase{}
	b.Events = NewEventBase()
	b.Box = Box["s"]
	b.SelectedIndex = -1
	return b
}

func NewListboxB(base Controlbase, Caption string) *Listbox {
	b := NewListbox()
	b.Caption = Caption
	b.Base = base

	b.Events.Onkey = func(e Event) bool {
		if e.Termbox.Key == termbox.KeyArrowLeft {
			if b.SelectedIndex > 0 {
				b.SelectedIndex--
				b.SelectedElement = b.SelectedElement.Prev()
			}
		} else if e.Termbox.Key == termbox.KeyArrowRight {
			if b.SelectedIndex < b.Stuff.Len()-1 {
				b.SelectedIndex++
				b.SelectedElement = b.SelectedElement.Next()

			}
		}
		b.Draw(e.Window)
		return true
	}
	return b
}

func (b *Listbox) Add(text string, value interface{}) {
	elemn := &Listelement{text, value, b}

	if b.Stuff.Len() == 0 {
		b.Stuff.PushBack(elemn)
		b.SelectedElement = b.Stuff.Front()
		b.SelectedIndex = 0
	} else {
		b.Stuff.PushBack(elemn)
	}

}

func (b *Listbox) Draw(win *Window) {

	var fg, bg, cfg, cbg termbox.Attribute

	var (
		X = b.Base.X
		Y = b.Base.Y
		W = b.Base.Width
		H = b.Base.Height
	)

	if b == win.SelectedControl {
		fg = win.SelectedForeground
		bg = win.SelectedBackground
		cfg = win.Foreground
		cbg = win.Background
	} else {
		fg = win.Foreground
		bg = win.Background
		cfg = win.SelectedForeground
		cbg = win.SelectedBackground
	}

	_, _ = cfg, cbg

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
			win.SetCell(X+x, Y+y, ' ', win.Foreground, win.Background)
			if timeeffect {
				termbox.Flush()
				time.Sleep(10 * time.Millisecond)
			}
		}
	}

	e := b.Stuff.Front()
	l := 0
	for {
		if e != nil {
			for k, v := range e.Value.(*Listelement).Text {
				if e.Value.(*Listelement).Parent.SelectedElement == e {
					win.SetCell(X+k+1, Y+1+l, v, win.SelectedForeground, win.SelectedBackground)
				} else {
					win.SetCell(X+k+1, Y+1+l, v, win.Foreground, win.Background)
				}

				if timeeffect {
					termbox.Flush()
					time.Sleep(10 * time.Millisecond)
				}

			}
		} else {
			break
		}
		l++
		e = e.Next()
	}

	for k, v := range b.Caption {
		win.SetCell(X+k+1, Y, v, fg, bg)
		if timeeffect {
			termbox.Flush()
			time.Sleep(10 * time.Millisecond)
		}
	}

	termbox.Flush()
}

func (b *Listbox) B() *Controlbase {
	return &b.Base
}

func (b *Listbox) E() *events {
	return &b.Events
}
