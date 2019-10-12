package ui

import (
    //"fyne.io/fyne"
    //"fyne.io/fyne/canvas"
    //"fyne.io/fyne/theme"
    "fyne.io/fyne/widget"
    "cartohelper/mapWidget"
    
    "fmt"
)

type Ui interface {
    
}

type SimpleUi struct {
    hbox *widget.Box
    menu *widget.Box
    mapWidget *mapWidget.MapWidget
}

func NewUi(mapWidget *mapWidget.MapWidget) *SimpleUi {
    result := &SimpleUi{mapWidget: mapWidget}
    hbox := widget.NewHBox()
    result.menu = widget.NewVBox()
    hbox.Append(result.menu)
    hbox.Append(mapWidget)
    result.hbox = hbox
    
    result.menu.Append(widget.NewButton("Useless Button", func() {
        fmt.Println("Useless")
    }))
    result.menu.Append(widget.NewButton("Add blob", func() {
        mapWidget.MapViewer.MapState().GenerateBlob()
        widget.Refresh(mapWidget)
    }))
    result.menu.Append(widget.NewButton("Add 10 blob", func() {
        for i := 0; i < 10; i++ {
            mapWidget.MapViewer.MapState().GenerateBlob()
        }
        widget.Refresh(mapWidget)
    }))
    result.menu.Append(widget.NewButton("Add 100 blob", func() {
        for i := 0; i < 100; i++ {
            mapWidget.MapViewer.MapState().GenerateBlob()
        }
        widget.Refresh(mapWidget)
    }))
    
    return result
}

func (ui *SimpleUi) GetHBox() *widget.Box {
    return ui.hbox
}
