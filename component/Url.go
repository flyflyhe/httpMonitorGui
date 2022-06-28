package component

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"log"
)

func buttonFocusLost(buttons ...*widget.Button) {
	for _, v := range buttons {
		v.FocusLost()
	}
}

func urlScreen(_ fyne.Window) fyne.CanvasObject {
	vBox := container.NewVBox()

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
			},
			SubmitText: "保存",
		}

		vBox.Objects = []fyne.CanvasObject{container.NewVBox(form)}
		vBox.Refresh()
	})

	deleteButton = widget.NewButton("删除", func() {
		buttonFocusLost(addButton, deleteButton, showButton)
		deleteButton.FocusGained()
	})

	showButton = widget.NewButton("列表", func() {
		buttonFocusLost(addButton, deleteButton, showButton)
		showButton.FocusGained()
	})
	return container.NewVBox(container.NewHBox(showButton, addButton, deleteButton), widget.NewSeparator(), vBox)
}
