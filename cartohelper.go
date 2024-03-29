package main

import (
    "fyne.io/fyne/widget"
	"fyne.io/fyne/app"
    "cartohelper/mapWidget"
    "cartohelper/ui"
    "cartohelper/mapState"
    "cartohelper/mapViewer"
)

func main() {
    mapState := mapstate.NewTileMapState(600, 600, -20)
    mapViewer := mapviewer.NewSimpleMapViewer()
    mapViewer.RegisterMap(mapState)
    mapState.Seed(3)
    
    app := app.New()

	w := app.NewWindow("Hello")
    c := w.Canvas()
    mapwidget := mapWidget.MapWidget{Scale: c.Scale(), MapViewer: mapViewer}
    ui := ui.NewUi(&mapwidget)
	w.SetContent(widget.NewVBox(
		ui.GetHBox(),
		widget.NewButton("Quit", func() {
			app.Quit()
		})))
    //w.SetContent(&mapWidget)

	w.ShowAndRun()
}
