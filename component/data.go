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
		"url":     {Title: "地址管理", View: urlScreen},
		"proxy":   {Title: "代理管理", View: proxyScreen},
		"monitor": {Title: "监控管理", View: monitorScreen},
	}

	//index tree

	AppViewsIndex = map[string][]string{
		"": {"welcome", "url", "proxy", "monitor"},
	}
)
