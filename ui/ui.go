package ui

import (
    //"fyne.io/fyne"
    //"fyne.io/fyne/canvas"
    //"fyne.io/fyne/theme"
    "fyne.io/fyne/widget"
    "cartohelper/mapWidget"
    "cartohelper/genutils"
    "cartohelper/mapViewer"
    "fmt"
    "strconv"
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
    
    mapViewerSimple := mapWidget.MapViewer
    
    distriViewer := mapviewer.NewDistributionMapViewer()
    distriViewer.RegisterMap(mapWidget.MapViewer.MapState())
    
    heightViewer := mapviewer.NewHeightMapViewer()
    heightViewer.RegisterMap(mapWidget.MapViewer.MapState())
    
    w := mapWidget.MapViewer.MapState().GetWidth()
    h := mapWidget.MapViewer.MapState().GetHeight()
    
    imgCount := 0
    
    result.menu.Append(widget.NewButton("Useless Button", func() {
        fmt.Println("Useless")
    }))
    result.menu.Append(widget.NewButton("Normal viewer", func() {
        mapWidget.MapViewer = mapViewerSimple
        widget.Refresh(mapWidget)
    }))
    result.menu.Append(widget.NewButton("Distribution Viewer", func() {
        mapWidget.MapViewer = distriViewer
        widget.Refresh(mapWidget)
    }))
    result.menu.Append(widget.NewButton("Height Viewer", func() {
        mapWidget.MapViewer = heightViewer
        widget.Refresh(mapWidget)
    }))
    result.menu.Append(widget.NewButton("Save", func() {
        mapWidget.SaveImg("image_" + strconv.Itoa(imgCount))
        fmt.Println("Saved: image_", imgCount)
        imgCount++
    }))
    result.menu.Append(widget.NewButton("Add blob", func() {
        mapWidget.MapViewer.MapState().GenerateBlob(0, 0, w, h)
        widget.Refresh(mapWidget)
    }))
    result.menu.Append(widget.NewButton("Add 10 blob", func() {
        for i := 0; i < 10; i++ {
            mapWidget.MapViewer.MapState().GenerateBlob(0, 0, w, h)
        }
        widget.Refresh(mapWidget)
    }))
    result.menu.Append(widget.NewButton("Add 100 blob", func() {
        for i := 0; i < 100; i++ {
            mapWidget.MapViewer.MapState().GenerateBlob(0, 0, w, h)
        }
        widget.Refresh(mapWidget)
    }))
    result.menu.Append(widget.NewButton("Add continent", func() {
        continent := genutils.NewContinent(w, h)
        
        for _, p := range continent.Probs {
            mapWidget.MapViewer.MapState().AddDistributionBlob(p.X, p.Y, p.Radius)
        }
        
        for i := continent.NBlob; i > 0; i-- {
            mapWidget.MapViewer.MapState().GenerateBlob(continent.X,
                                                        continent.Y,
                                                        continent.Width,
                                                        continent.Height)
        }
        
        widget.Refresh(mapWidget)
    }))
    result.menu.Append(widget.NewButton("Add continent prob only", func() {
        continent := genutils.NewContinent(w, h)
        
        for _, p := range continent.Probs {
            mapWidget.MapViewer.MapState().AddDistributionBlob(p.X, p.Y, p.Radius)
        }
        
        widget.Refresh(mapWidget)
    }))
    
    return result
}

func (ui *SimpleUi) GetHBox() *widget.Box {
    return ui.hbox
}
