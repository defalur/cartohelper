package main

import (
    "fyne.io/fyne/widget"
	"fyne.io/fyne/app"
    "cartohelper/mapWidget"
    "cartohelper/ui"
)

func main() {
    app := app.New()
    mapwidget := mapWidget.MapWidget{}

	w := app.NewWindow("Hello")
    ui := ui.CreateUi(&mapwidget)
	w.SetContent(widget.NewVBox(
		ui.GetHBox(),
		widget.NewButton("Quit", func() {
			app.Quit()
		})))
    //w.SetContent(&mapWidget)

	w.ShowAndRun()
}
