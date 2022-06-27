package component

import (
	"fyne.io/fyne/v2"
)

type AppView struct {
	Title, Intro string
	View         func(w fyne.Window) fyne.CanvasObject
}

var (
	AppViews = map[string]AppView{
		"welcome": {"Welcome", "", welcomeScreen},
		"canvas": {"Canvas",
			"See the canvas capabilities.",
			canvasScreen,
		},
	}

	//index tree
	AppViewsIndex = map[string][]string{
		"": {"welcome", "canvas"},
	}
)
