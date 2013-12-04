package cui

import (
	"github.com/nsf/termbox-go"
	"math"
	"time"
)

type Label struct {
	Base       Controlbase
	Text       string
	Pageoffset int
	events     events
}

func NewLabel() Label {
	b := Label{}
	b.Base = Controlbase{}
	b.Pageoffset = 0
	return b
}

func NewLabelB(base Controlbase, text string) *Label {
	b := &Label{}
	b.Base = base
	b.Text = text
	b.Pageoffset = 0
	b.events.Onkey = func(ev Event) bool {
		pagecount := math.Ceil(math.Ceil(float64(len(b.Text)/b.Base.Width)) / float64(b.Base.Height))
		if ev.Termbox.Key == termbox.KeyPgdn {
			if float64(b.Pageoffset) < pagecount {
				b.Pageoffset++
			}
		} else if ev.Termbox.Key == termbox.KeyPgup {
			if b.Pageoffset > 0 {
				b.Pageoffset--
			}
		}

		b.Draw(ev.Window)
		return true
	}

	return b
}

func (l *Label) Draw(win *Window) {

	var fg termbox.Attribute
	var bg termbox.Attribute

	var (
		X = l.Base.X
		Y = l.Base.Y
		//W = l.Base.Width
		//H = l.Base.Height
	)

	if l == win.SelectedControl {
		fg = win.SelectedForeground
		bg = win.SelectedBackground
	} else {
		fg = win.Foreground
		bg = win.Background
	}

	lineoffset := 0
	k := 0
	for _, v := range l.Text {
		if v == '\n' {
			lineoffset++
			k = 0
			continue
		}

		if k == l.Base.Width {
			lineoffset++
			k = 0
		}

		if k == l.Base.Height {

		}

		win.SetCell(X+k, Y+lineoffset, v, fg, bg)
		k++
		if timeeffect {
			termbox.Flush()
			time.Sleep(10 * time.Millisecond)
		}
	}

	termbox.Flush()
}

func (b *Label) B() *Controlbase {
	return &b.Base
}
