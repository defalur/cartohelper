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
    
    return result
}

func (ui *SimpleUi) GetHBox() *widget.Box {
    return ui.hbox
}
