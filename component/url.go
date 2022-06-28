package component

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/flyflyhe/httpMonitorGui/layouts"
	"github.com/flyflyhe/httpMonitorGui/services/rpc"
	"log"
	"strconv"
)

func buttonFocusLost(buttons ...*widget.Button) {
	for _, v := range buttons {
		v.FocusLost()
	}
}

func urlScreen(w fyne.Window) fyne.CanvasObject {
	vBox := container.New(layouts.NewVBoxLayout())

	var addButton *widget.Button
	var deleteButton *widget.Button
	var showButton *widget.Button

	addButton = widget.NewButton("添加", func() {
		buttonFocusLost(addButton, deleteButton, showButton)
		addButton.FocusGained()
		urlEntry := widget.NewEntry()
		intervalEntry := widget.NewEntry()

		form := &widget.Form{
			Items: []*widget.FormItem{ // we can specify items in the constructor
				{Text: "http地址", Widget: urlEntry},
				{Text: "间隔时间毫秒", Widget: intervalEntry},
			},
			OnCancel: func() {
				urlEntry.SetText("")
				intervalEntry.SetText("")
			},
			CancelText: "重置",
			OnSubmit: func() { // optional, handle form submission
				log.Println("Form submitted:", urlEntry.Text, intervalEntry.Text)
				interval, err := strconv.Atoi(intervalEntry.Text)
				if err != nil {
					dialog.ShowError(err, w)
					return
				}
				if err = rpc.SetUrl(urlEntry.Text, int32(interval)); err != nil {
					dialog.ShowError(err, w)
				} else {
					dialog.ShowInformation("提示", "保存成功", w)
				}

			},
			SubmitText: "保存",
		}

		vBox.Objects = []fyne.CanvasObject{container.NewVBox(form)}
		vBox.Refresh()
	})

	deleteButton = widget.NewButton("删除", func() {
		buttonFocusLost(addButton, deleteButton, showButton)
		deleteButton.FocusGained()
		urlEntry := widget.NewEntry()

		form := &widget.Form{
			Items: []*widget.FormItem{ // we can specify items in the constructor
				{Text: "http地址", Widget: urlEntry},
			},
			OnCancel: func() {
				urlEntry.SetText("")
			},
			CancelText: "重置",
			OnSubmit: func() { // optional, handle form submission
				log.Println("Form submitted:", urlEntry.Text)
				if err := rpc.DeleteUrl(urlEntry.Text); err != nil {
					dialog.ShowError(err, w)
				} else {
					dialog.ShowInformation("提示", "删除成功", w)
				}

			},
			SubmitText: "保存",
		}

		vBox.Objects = []fyne.CanvasObject{container.NewVBox(form)}
		vBox.Refresh()
	})

	var showButtonFunc func()
	showButtonFunc = func() {
		buttonFocusLost(addButton, deleteButton, showButton)
		showButton.FocusGained()
		if urls, err := rpc.ListUrl(); err != nil {
			dialog.ShowError(err, w)
		} else {
			list := widget.NewList(
				func() int {
					return len(urls)
				},
				func() fyne.CanvasObject {
					return widget.NewLabel("template")
				},
				func(i widget.ListItemID, o fyne.CanvasObject) {
					o.(*widget.Label).SetText(urls[i])
				})
			list.OnSelected = func(id widget.ListItemID) {
				dialog.ShowConfirm("操作", "是否删除", func(b bool) {
					if b {
						if err := rpc.DeleteUrl(urls[id]); err != nil {
							dialog.ShowError(err, w)
						} else {
							dialog.ShowInformation("提示", "删除成功", w)
							showButtonFunc()
						}
					}
				}, w)
			}

			c := container.New(layouts.NewVBoxLayout(), list)
			layouts.SetObjConfigMap(list, &layouts.Size{Height: 400, Width: 200})
			vBox.Objects = []fyne.CanvasObject{c}
			vBox.Refresh()
		}
	}
	showButton = widget.NewButton("列表", showButtonFunc)
	return container.NewVBox(container.NewHBox(showButton, addButton, deleteButton), widget.NewSeparator(), vBox)
}
