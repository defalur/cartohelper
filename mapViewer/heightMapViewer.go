package mapviewer

import (
    "cartohelper/mapState"
    "image/color"
    "math"
)

type HeightMapViewer struct {
    mapState mapstate.MapState
    lowGround color.RGBA//close to sea level
    highGround color.RGBA//high altutide
    lowSea color.RGBA//close to surface
    highSea color.RGBA//very deep
}

func NewHeightMapViewer() MapViewer {
    res := &HeightMapViewer{}
    res.lowGround = color.RGBA{1, 113, 0, 255}
    res.highGround = color.RGBA{137, 237, 137, 255}
    res.lowSea = color.RGBA{85, 96, 234, 255}
    res.highSea = color.RGBA{18, 75, 104, 255}
    
    return res
}

func (viewer *HeightMapViewer) RegisterMap(mapState mapstate.MapState) {
    viewer.mapState = mapState
}

func interpolate(colorStart, colorEnd uint8, heightEnd, heightCur int) uint8 {
    ratio := math.Abs(float64(heightCur) / float64(heightEnd))
    
    deltaColor := float64(colorEnd) - float64(colorStart)
    
    return colorStart + uint8(deltaColor * ratio)
}

func (viewer *HeightMapViewer) GetPixel(x, y int) color.Color {
    if viewer.mapState == nil {
        return nil
    }
    
    minHeight := viewer.mapState.MinHeight()
    maxHeight := viewer.mapState.MaxHeight()
    
    node, err := viewer.mapState.GetNode(x, y)
    if err != nil {
        return color.RGBA{0, 0, 0, 255}//black for all pixels outside of map
    }
    curHeight := node.GetHeight()
    if curHeight > 0 {
        //ground
        c1 := viewer.lowGround
        c2 := viewer.highGround
        r := interpolate(c1.R, c2.R, maxHeight, curHeight)
        g := interpolate(c1.G, c2.G, maxHeight, curHeight)
        b := interpolate(c1.B, c2.B, maxHeight, curHeight)
        
        return color.RGBA{r, g, b, 255}
    } else {
        //sea
        c1 := viewer.lowSea
        c2 := viewer.highSea
        r := interpolate(c1.R, c2.R, minHeight, curHeight)
        g := interpolate(c1.G, c2.G, minHeight, curHeight)
        b := interpolate(c1.B, c2.B, minHeight, curHeight)
        
        return color.RGBA{r, g, b, 255}
    }
}

func (viewer *HeightMapViewer) MapState() mapstate.MapState {
    return viewer.mapState
}
