package cui

import (
	//"encoding/json"
	"fmt"
	"github.com/nsf/termbox-go"
	"os"
	"os/exec"
	"time"
)

var endchan chan bool = make(chan bool, 1)
var pausechan chan bool = make(chan bool, 1)
var CurrentWindow *Window
var WindowStack []*Window
var timeeffect bool = false
var statusline bool = false
var evchan chan termbox.Event = make(chan termbox.Event, 10)

func init() {
	WindowStack = make([]*Window, 0, 0)
}

var Box map[string]*Boxchars = map[string]*Boxchars{
	"d": &Boxchars{
		tl: '╔',
		tr: '╗',
		bl: '╚',
		br: '╝',
		v:  '║',
		h:  '═',
	},
	"s": &Boxchars{
		tl: '┌',
		tr: '┐',
		bl: '└',
		br: '┘',
		v:  '│',
		h:  '─',
	},
}

var Box2 map[string]rune = map[string]rune{
	"shadow": '▒',
}

type Boxchars struct {
	v, h, tl, tr, bl, br rune
}

type events struct {
	Onstart  func(Event) bool
	Onclick  func(Event) bool
	Onselect func(Event) bool
	//Onchange func(Event) bool
	Onkey func(Event) bool
}

func NewEventBase() events {
	ev := events{}
	//ev.Onstart = make([]func(), 0, 0)
	//ev.Onclick = make([]func(), 0, 0)
	//ev.Onselect = make([]func(), 0, 0)
	//ev.Onchange = make([]func(), 0, 0)
	return ev
}

type Event struct {
	Termbox termbox.Event
	Value   interface{}
	Window  *Window
	Control Control
}

type Controlbase struct {
	X      int
	Y      int
	Width  int
	Height int
	Z      int
}

type Control interface {
	B() *Controlbase
	Draw(*Window)
}

type Eventcontrol interface {
	Control
	E() *events
}

type inteventcontrol interface {
	intevents() *events
}

/*func (c *Control) Clear() {

}*/

//Pauses the UI. Set blank to true to blank the console out (black)
//Set blank to false to keep the current screen
func PauseTermbox(blank bool) {
	pausechan <- true

	if blank {
		termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)
		termbox.SetCursor(0, 0)
	}

	termbox.Flush()
	termbox.Close()
}

//Resumes the UI from where it left off.
func ResumeTermbox() {
	pausechan <- false
	termbox.Init()
	CurrentWindow.Draw()
}

//Starts with a window and returns instantly.
func StartTermbox(window *Window) error {
	err := termbox.Init()

	if err != nil {
		return err
	}
	termbox.Clear(window.Foreground, window.Background)

	Navigate(window)

	go inputloop1()

	go inputloop2()

	return nil
}

//Starts with a window and blocks
func StartTermboxBlocking(window *Window) error {
	err := termbox.Init()

	if err != nil {
		return err
	}
	termbox.Clear(window.Foreground, window.Background)

	Navigate(window)

	go inputloop1()

	inputloop2()

	termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)
	termbox.SetCursor(0, 0)
	termbox.Flush()

	return nil
}

func External(name string, arg ...string) {
	PauseTermbox(true)
	cmd := exec.Command(name, arg...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
	time.Sleep(200 * time.Millisecond)

	ResumeTermbox()
}

func inputloop1() {
	for {

		select {

		case <-pausechan:
			<-pausechan
		default:
			//case evchan <- termbox.PollEvent():
			evchan <- termbox.PollEvent()
		}
	}
}

func inputloop2() {
	defer termbox.Close()
outer:
	for {
		select {
		case ev := <-evchan:
			if CurrentWindow != nil {
				switch ev.Type {
				case termbox.EventKey:

					if CurrentWindow.onkey == nil {
						if CurrentWindow.Events.Onkey == nil {
						} else {
							CurrentWindow.Events.Onkey(Event{Window: CurrentWindow, Termbox: ev})
						}
					} else {
						CurrentWindow.onkey(Event{Window: CurrentWindow, Termbox: ev})
					}
				case termbox.EventResize:
					termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)
					CurrentWindow.Draw()
				}
			} else {
				for k, v := range "No window opened... (You popped more Windows than were on the stack)" {
					termbox.SetCell(k, 0, v, termbox.ColorWhite, termbox.ColorBlack)
				}
				for k, v := range "Or you blanked out CurrentWindow." {
					termbox.SetCell(k, 1, v, termbox.ColorWhite, termbox.ColorBlack)
				}
				termbox.Flush()
			}
		case <-endchan:
			StatusLine("Break!")
			break outer
		}
	}

	termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)
	termbox.SetCursor(0, 0)
	termbox.Flush()
}

//Navigates to a new window, completeley erases the window stack in the process
func Navigate(newwin *Window) {
	CurrentWindow = newwin
	WindowStack = make([]*Window, 0, 0)
	WindowStack = append(WindowStack, newwin)
	termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)
	CurrentWindow.Draw()
	if CurrentWindow.Events.Onstart != nil {
		CurrentWindow.Events.Onstart(Event{Window: CurrentWindow})
	}
}

//Push Window onto the window stack. This allows showing windows on top of each
//other and going backwards again. (Like dialog windows)
func PushWindow(win *Window) {
	WindowStack = append(WindowStack, win)
	CurrentWindow = win
	CurrentWindow.Draw()
	if CurrentWindow.Events.Onstart != nil {
		CurrentWindow.Events.Onstart(Event{Window: CurrentWindow})
	}
}

//Pops a window from the window stack. Focuses to the next window on the stack
func PopWindow() *Window {
	if len(WindowStack) > 1 {
		CurrentWindow = WindowStack[len(WindowStack)-2]
	} else {
		CurrentWindow = nil
		termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)
		termbox.Flush()
	}
	s := WindowStack[len(WindowStack)-1]

	WindowStack = WindowStack[:len(WindowStack)-1]

	termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)

	for _, v := range WindowStack {
		v.Draw()
	}
	return s
}

//Turns off termbox and breaks out of the input loop
func End() {
	endchan <- true
}

func setc(x, y int, ch rune) {
	//Cells := termbox.CellBuffer()
	//Cells[x].Ch = rune(ch[0])
	termbox.SetCell(x, y, ch, termbox.ColorBlack, termbox.ColorWhite)
	termbox.Flush()
}

//Unused
func StatusLine(text string) {
	if statusline {
		_, h := termbox.Size()
		termbox.SetCursor(1, h-2)
		fmt.Print(text)
		termbox.HideCursor()
	}
}
