package cui

import (
	"github.com/nsf/termbox-go"
	"time"
)

type Checkbox struct {
	Base       Controlbase
	checked    bool
	Events     events
	Box        *Boxchars
	Checkgroup *Checkgroup
	Text       string
}

func (c *Checkbox) Toggle() {
	c.checked = !c.checked
	if c.Checkgroup != nil {
		if c.checked {
			if c.Checkgroup.Active != nil {
				c.Checkgroup.Active.checked = false
			}
			c.Checkgroup.Active = c
		} else {
			c.Checkgroup.Active = nil
		}
	}
}

func (c *Checkbox) Checked() bool {
	return c.checked
}

type Checkgroup struct {
	Checkboxes []*Checkbox
	Active     *Checkbox
}

func NewCheckgroup() *Checkgroup {
	cg := &Checkgroup{}
	cg.Checkboxes = make([]*Checkbox, 0, 0)
	return cg
}

func (c *Checkgroup) AddControl(chk ...*Checkbox) {
	c.Checkboxes = append(c.Checkboxes, chk...)

	activewasset := false

	for _, v := range chk {
		v.Checkgroup = c
		if v.checked {
			if activewasset {
				v.checked = false
			} else {
				c.Active = v
				activewasset = true
			}
		}
	}
}

func NewCheckbox() *Checkbox {
	b := &Checkbox{}
	b.Base = Controlbase{}
	b.Events = NewEventBase()
	b.Events.Onclick = func(e Event) bool {

		if b.Checkgroup != nil {
			if b.Checkgroup.Active != nil {
				if b.Checkgroup.Active != b {
					b.Checkgroup.Active.checked = false
					b.checked = true
					b.Checkgroup.Active.Draw(e.Window)
					b.Checkgroup.Active = b
					b.Draw(e.Window)
					return false
				} else {
					return false
				}
			} else {
				b.Checkgroup.Active = b
				b.Draw(e.Window)
				b.checked = true
			}
		} else {
			b.checked = !b.checked
			b.Draw(e.Window)
		}

		return false
	}
	b.Box = Box["s"]
	return b
}

func NewCheckboxB(base Controlbase, text string) *Checkbox {
	b := &Checkbox{}
	b.Base = base
	b.Text = text
	b.Events.Onclick = func(e Event) bool {

		if b.Checkgroup != nil {
			if b.Checkgroup.Active != nil {
				if b.Checkgroup.Active != b {
					b.Checkgroup.Active.checked = false
					b.checked = true
					b.Checkgroup.Active.Draw(e.Window)
					b.Checkgroup.Active = b
					b.Draw(e.Window)
					return false
				} else {
					return false
				}
			} else {
				b.Checkgroup.Active = b
				b.checked = true
				b.Draw(e.Window)
			}
		} else {
			b.checked = !b.checked
			b.Draw(e.Window)
		}

		return false
	}
	b.Box = Box["s"]
	return b
}

func (b *Checkbox) Draw(win *Window) {

	var fg termbox.Attribute
	var bg termbox.Attribute

	var (
		X = b.Base.X
		Y = b.Base.Y
		W = 3
		H = 3
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
		win.SetCell(X+3+k, Y+1, v, fg, bg)
	}

	if b.checked {
		win.SetCell(X+1, Y+1, 'X', fg, bg)
	} else {
		win.SetCell(X+1, Y+1, ' ', fg, bg)
	}

	termbox.Flush()
}

func (b *Checkbox) B() *Controlbase {
	return &b.Base
}

func (b *Checkbox) E() *events {
	return &b.Events
}
