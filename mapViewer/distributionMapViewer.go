package mapviewer

import (
    "cartohelper/mapState"
    "image/color"
)

type DistributionMapViewer struct {
    mapState mapstate.MapState
}

func NewDistributionMapViewer() MapViewer {
    return &DistributionMapViewer{}
}

func (viewer *DistributionMapViewer) RegisterMap(mapState mapstate.MapState) {
    viewer.mapState = mapState
}

func (viewer *DistributionMapViewer) GetPixel(x, y int) color.Color {
    if viewer.mapState != nil {
        value := viewer.mapState.GetDistribution(x, y)
        valueConvert := uint8(value * 255)
        return color.RGBA{valueConvert, valueConvert, valueConvert, 255}
    }
    
    return nil
}

func (viewer *DistributionMapViewer) MapState() mapstate.MapState {
    return viewer.mapState
}
