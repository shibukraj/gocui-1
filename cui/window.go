package cui

import (
	"github.com/nsf/termbox-go"
	"time"
)

type Window struct {
	Controls             []Control
	ControlMap           map[string]Control
	base                 *Controlbase
	DisplayControls      []Control
	SelectedControlindex int
	SelectedControl      Control
	Background           termbox.Attribute
	Foreground           termbox.Attribute
	SelectedBackground   termbox.Attribute
	SelectedForeground   termbox.Attribute
	Events               events
	events               events
	Caption              string
	onkey                func(Event) bool
}

func (win *Window) SetBase(c *Controlbase) {
	win.base = c
}

func (w *Window) SetCell(x, y int, rn rune, fg, bg termbox.Attribute) {
	if w.base != nil {
		termbox.SetCell(x+w.base.X+1, y+w.base.Y+1, rn, fg, bg)
	} else {
		termbox.SetCell(x, y, rn, fg, bg)
	}

}

func NewWindow() *Window {
	wi := &Window{}
	wi.Background = termbox.ColorWhite
	wi.Foreground = termbox.ColorBlack
	wi.SelectedBackground = termbox.ColorBlack
	wi.SelectedForeground = termbox.ColorWhite
	wi.Controls = make([]Control, 0, 0)
	wi.ControlMap = make(map[string]Control)
	wi.DisplayControls = make([]Control, 0, 0)
	wi.Events = NewEventBase()
	wi.events = NewEventBase()
	wi.onkey = func(e Event) bool {
		ev := e.Termbox

		keyhandled := false

		if ev.Key != termbox.KeyTab {
			if wi.ControlSelected() {
				if val, ok := wi.SelectedControl.(inteventcontrol); ok {
					if val.intevents().Onkey != nil {
						keyhandled = val.intevents().Onkey(Event{Window: wi, Control: wi.SelectedControl, Termbox: ev})
					}
				}

				if val, ok := wi.SelectedControl.(Eventcontrol); ok {
					if val.E().Onkey != nil {
						keyhandled = keyhandled || !val.E().Onkey(Event{Window: wi, Control: wi.SelectedControl, Termbox: ev})
					}
				}

				if !keyhandled {
					if wi.events.Onkey != nil {
						keyhandled = wi.events.Onkey(Event{Window: wi, Control: wi.SelectedControl, Termbox: ev})
					}

					if wi.Events.Onkey != nil {
						keyhandled = keyhandled || wi.Events.Onkey(Event{Window: wi, Control: wi.SelectedControl, Termbox: ev})
					}
				}

			}
		}

		if !keyhandled {
			switch ev.Key {
			case termbox.KeyF5:
				wi.Draw()
				return true
			case termbox.KeyArrowDown, termbox.KeyArrowUp, termbox.KeyTab:
				modifier := 0
				if ev.Key == termbox.KeyArrowUp {
					modifier = -1
				} else {
					modifier = +1
				}

				selecthandled := false
				if wi.ControlSelected() {

					if val, ok := wi.SelectedControl.(inteventcontrol); ok {
						if val.intevents().Onselect != nil {
							selecthandled = val.intevents().Onselect(Event{Window: wi, Control: wi.SelectedControl, Termbox: ev})
						}
					}

					if val, ok := wi.SelectedControl.(Eventcontrol); ok {
						if val.E().Onselect != nil {
							selecthandled = selecthandled || !val.E().Onselect(Event{Window: wi, Control: wi.SelectedControl, Termbox: ev})
						}
					}

					if modifier == -1 && wi.SelectedControlindex+modifier >= 0 || modifier == +1 && wi.SelectedControlindex+modifier < len(wi.Controls) {
						wi.SelectedControlindex = wi.SelectedControlindex + modifier
						c := wi.SelectedControl
						wi.SelectedControl = wi.Controls[wi.SelectedControlindex]
						c.Draw(wi)
						wi.SelectedControl.Draw(wi)
					}
				}

			case termbox.KeyEnter:
				StatusLine("Enter!")

				if val, ok := wi.SelectedControl.(inteventcontrol); ok {
					if val.intevents().Onclick != nil {
						val.intevents().Onclick(Event{Window: wi, Control: wi.SelectedControl, Termbox: ev})
					}
				}

				if val, ok := wi.SelectedControl.(Eventcontrol); ok {
					if val.E().Onclick != nil {
						val.E().Onclick(Event{Window: wi, Control: wi.SelectedControl, Termbox: ev})
					}
				}
			case termbox.KeyCtrlC:
				StatusLine("Break!")
				End()
			default:
			}
		}
		return true
	}

	return wi
}

func (win *Window) AddControl(c ...Control) {
	win.Controls = append(win.Controls, c...)
}

func (win *Window) AddDisplayControl(c ...Control) {
	win.DisplayControls = append(win.DisplayControls, c...)
}

func (win *Window) ControlSelected() bool {
	return win.SelectedControl != nil
}

func (win *Window) intevents() *events {
	return &win.events
}

func (win *Window) Draw() {

	var (
		fg = win.Foreground
		bg = win.Background
	)

	termbox.Flush()

	w, h := 0, 0
	x, y := 0, 0

	if win.base != nil {
		w, h = win.base.Width, win.base.Height
		x, y = win.base.X, win.base.Y
	} else {
		w, h = termbox.Size()
		x, y = 0, 0
	}

	for X := 0; X < w; X++ {
		for Y := 0; Y < h; Y++ {
			termbox.SetCell(X+x, Y+y, ' ', fg, bg)
		}
	}

	termbox.SetCell(x, y, Box["d"].tl, fg, bg)
	termbox.SetCell(x+w-1, y, Box["d"].tr, fg, bg)
	termbox.SetCell(x, y+h-1, Box["d"].bl, fg, bg)
	termbox.SetCell(x+w-1, y+h-1, Box["d"].br, fg, bg)

	for X := 0; X < int(w/2); X++ {
		termbox.SetCell(x+X+1, y, Box["d"].h, fg, bg)
		termbox.SetCell(x+w-X-2, y, Box["d"].h, fg, bg)
		termbox.SetCell(x+X+1, y+h-1, Box["d"].h, fg, bg)
		termbox.SetCell(x+w-X-2, y+h-1, Box["d"].h, fg, bg)

		if timeeffect {
			termbox.Flush()
			time.Sleep(10 * time.Millisecond)
		}
	}

	for Y := 1; Y < int(h/2); Y++ {
		termbox.SetCell(x, y+Y, Box["d"].v, fg, bg)
		termbox.SetCell(x+w-1, y+Y, Box["d"].v, fg, bg)
		termbox.SetCell(x, y+h-Y-1, Box["d"].v, fg, bg)
		termbox.SetCell(x+w-1, y+h-Y-1, Box["d"].v, fg, bg)

		if timeeffect {
			termbox.Flush()
			time.Sleep(10 * time.Millisecond)
		}
	}

	for k, v := range win.Caption {
		termbox.SetCell(x+2+k, y, v, win.Foreground, win.Background)

		if timeeffect {
			termbox.Flush()
			time.Sleep(10 * time.Millisecond)
		}
	}

	for _, c := range win.Controls {
		c.Draw(win)
	}

	for _, c := range win.DisplayControls {
		c.Draw(win)
	}

	termbox.Flush()
}
