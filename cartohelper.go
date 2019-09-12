package main

import (
    "fyne.io/fyne/widget"
	"fyne.io/fyne/app"
    "cartohelper/mapWidget"
    "cartohelper/ui"
)

func main() {
    app := app.New()

	w := app.NewWindow("Hello")
    c := w.Canvas()
    mapwidget := mapWidget.MapWidget{Scale: c.Scale()}
    ui := ui.NewUi(&mapwidget)
	w.SetContent(widget.NewVBox(
		ui.GetHBox(),
		widget.NewButton("Quit", func() {
			app.Quit()
		})))
    //w.SetContent(&mapWidget)

	w.ShowAndRun()
}
