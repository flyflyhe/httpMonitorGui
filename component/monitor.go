package component

import (
	"encoding/json"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/flyflyhe/httpMonitorGui/layouts"
	"github.com/flyflyhe/httpMonitorGui/services/rpc"
	"github.com/rs/zerolog/log"
	"runtime/debug"
	"sync"
)

func monitorScreen(w fyne.Window) fyne.CanvasObject {
	vBox := container.New(layouts.NewVBoxLayout())

	var startButton *widget.Button
	var startButtonLock sync.Mutex
	var stopButton *widget.Button
	startFunc := func() {
		startButton.FocusGained()
		go func() {
			defer func() {
				if err := recover(); err != nil {
					debug.PrintStack()
					errJson, _ := json.Marshal(err)
					log.Error().Caller().Msg(string(errJson))
				}
			}()
			entry := widget.NewMultiLineEntry()
			layouts.SetObjConfigMap(entry, &layouts.Size{Height: 400, Width: 200})
			for {
				select {
				case res := <-rpc.GetMonitorQueue().Queue:
					entry.Text = entry.Text + "\n" + res.Url
					for proxy, v := range res.Result {
						entry.Text += "\n" + proxy + "<=>" + v
						if v != "success" {
							fyne.NewNotification(res.Url+"监控异常", "代理"+proxy+"信息+"+v)
						}
					}

					vBox.Objects = []fyne.CanvasObject{entry}
					vBox.Refresh()
				default:
				}
			}
		}()
	}
	startButton = widget.NewButton("启动", func() {
		startButtonLock.Lock() //防止重复点击
		defer startButtonLock.Unlock()
		buttonFocusLost(startButton, stopButton)
		startButton.FocusGained()

		dialog.ShowConfirm("url监控", "确认启动", func(b bool) {
			if b {
				if err := rpc.StartMonitor(rpc.GetMonitorQueue()); err != nil {
					dialog.ShowError(err, w)
				} else {
					dialog.ShowInformation("启动结果", "成功", w)
					startFunc()
				}
			}
		}, w)
	})

	stopButton = widget.NewButton("停止", func() {
		buttonFocusLost(startButton, stopButton)
		//stopButton.FocusGained()

		dialog.ShowConfirm("url监控", "确认停止", func(b bool) {
			if b {
				if err := rpc.StopMonitor(rpc.GetMonitorQueue()); err != nil {
					dialog.ShowError(err, w)
				} else {
					dialog.ShowInformation("确认停止", "停止成功", w)
				}
			}
		}, w)
	})

	if rpc.GetMonitorQueue().Running {
		startFunc()
	}

	return container.NewVBox(container.NewHBox(startButton, stopButton), widget.NewSeparator(), vBox)
}
