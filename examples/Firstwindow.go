// Firstwindow.go
package main

import (
	//"encoding/json"
	"github.com/nsf/termbox-go"
	"github.com/sinni800/gocui/cui"
	"runtime"
	"time"
)

func main() {
	cui.StartTermboxBlocking(WindowLogin())
}

type Window1 struct {
	cui.Window

	Button1   cui.Button
	Button2   cui.Button
	Textbox1  cui.Textbox
	Checkbox1 cui.Checkbox
	Label1    cui.Label
	Label2    cui.Label
}

func NewWindow1() *cui.Window {
	windowdata := `
	{"Controls":[
		{"Base":{
			"X":5, "Y":5, "Width":20, "Height":3, "Z":0 
		}, "Text":"Button??"},
		{"Base":{
			"X":5,"Y":9,"Width":20,"Height":3,"Z":0
		},"Text":"Hau!"}
	],
	"SelectedControl":null,"Background":8,"Foreground":1,"Events":{"Onstart":[],"Onclick":[],"Onselect":[],"Onchange":[]}}`
	_ = windowdata

	//ret := &Window1{}
	ret := cui.NewWindow()
	//ret.SetBase(&cui.Controlbase{3, 3, 46, 24, 0})
	ret.Background = termbox.ColorWhite

	Button1 := cui.NewButtonB(cui.Controlbase{3, 2, 20, 3, 0}, "Mass Tagger?")
	Button2 := cui.NewButtonB(cui.Controlbase{3, 6, 20, 3, 0}, "Exitme!")
	Textbox1 := cui.NewTextboxB(cui.Controlbase{3, 10, 40, 3, 0}, "")
	Checkbox1 := cui.NewCheckboxB(cui.Controlbase{3, 14, 3, 3, 0}, "Checkboxtext!")
	Label1 := cui.NewLabelB(cui.Controlbase{3, 18, 40, 3, 0}, "Check out the text field above! I also have a checkbox.")
	Label2 := cui.NewLabelB(cui.Controlbase{3, 20, 40, 3, 0}, "I am also a label but you can select me!")

	List := cui.NewListboxB(cui.Controlbase{44, 2, 30, 10, 0}, "ProgramRunner")
	List.Add("Gelmass", "C:\\Telnet\\gelmass.exe")
	List.Add("Hurp", "goconsole")
	List.Add("Lol", "cmd")

	List.Events.Onclick = func(e cui.Event) bool {
		v := List.SelectedElement.Value.(*cui.Listelement).Value.(string)
		cui.External(v)

		return true
	}

	Button1.Events.Onclick = func(e cui.Event) bool {
		cui.External("C:\\Telnet\\gelmass.exe", "localhost")
		return true
	}

	Button2.Events.Onclick = func(e cui.Event) bool {
		cui.End()
		return true
	}

	ret.AddControl(Button1, Button2, Textbox1, Checkbox1, Label2, List)
	ret.AddDisplayControl(Label1)
	ret.SelectedControl = ret.Controls[0]
	ret.SelectedControlindex = 0
	ret.Caption = "Honigkuchenpwnzer"

	return ret

	//Firstwindow = &Window{}
	//err := json.Unmarshal([]byte(windowdata), Firstwindow)
	//if err != nil {
	//	panic(err.Error())
	//}
}

func WindowLogin() *cui.Window {
	Login := cui.NewWindow()
	Login.SetBase(&cui.Controlbase{2, 1, 70, 20, 0})
	Login.Caption = "TAC Login"
	Username := cui.NewTextboxB(cui.Controlbase{3, 2, 30, 3, 0}, "User Name")
	Password := cui.NewTextboxB(cui.Controlbase{35, 2, 30, 3, 0}, "Password")
	Password.Maskchar = '*'
	Submit := cui.NewButtonB(cui.Controlbase{3, 14, 15, 3, 0}, "Submit")
	lblLoginProvider := cui.NewLabelB(cui.Controlbase{3, 6, 50, 1, 0}, "Login provider")

	chkLoginNormal := cui.NewCheckboxB(cui.Controlbase{3, 8, 20, 3, 0}, "Regular")
	chkLoginNormal.Toggle()
	chkLoginGelbooru := cui.NewCheckboxB(cui.Controlbase{15, 8, 20, 3, 0}, "Gelbooru")
	chkGroup := cui.NewCheckgroup()

	chkGroup.AddControl(chkLoginNormal, chkLoginGelbooru)

	Login.SelectedControl = Username

	Submit.Events.Onclick = func(e cui.Event) bool {
		if string(Username.Text) == "hurgh" {
			cui.PushWindow(WindowPasswordWrong())
		} else {
			cui.PushWindow(WindowSplash())
		}

		//cui.Navigate(WindowSplash())
		return true
	}

	Login.AddControl(Username, Password, chkLoginNormal, chkLoginGelbooru, Submit)
	Login.AddDisplayControl(lblLoginProvider)

	return Login
}

func WindowPasswordWrong() *cui.Window {
	PwdWindow := cui.NewWindow()
	PwdWindow.Background = termbox.ColorRed
	PwdWindow.Foreground = termbox.ColorWhite
	PwdWindow.SetBase(&cui.Controlbase{20, 8, 32, 9, 0})
	Textbox := cui.NewLabelB(cui.Controlbase{2, 1, 48, 2, 0}, "Password wrong. Try again.")
	OkButton := cui.NewButtonB(cui.Controlbase{10, 3, 6, 3, 0}, " OK")
	OkButton.Events.Onclick = func(e cui.Event) bool {
		cui.PopWindow()
		return true
	}

	PwdWindow.AddDisplayControl(Textbox)
	PwdWindow.AddControl(OkButton)
	PwdWindow.SelectedControl = OkButton
	return PwdWindow
}

func WindowSplash() *cui.Window {
	Splash := cui.NewWindow()
	Splash.SetBase(&cui.Controlbase{10, 3, 50, 19, 0})
	Textbox := cui.NewLabelB(cui.Controlbase{0, 0, 48, 2, 0}, `Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet.`+"\n\nProbably some login text?")

	Splash.AddDisplayControl(Textbox)
	Splash.Events.Onstart = func(e cui.Event) bool {
		//<-time.After(1e9)
		time.AfterFunc(1e9, func() {
			cui.Navigate(NewWindow1())

		})
		return true
	}
	return Splash
}
