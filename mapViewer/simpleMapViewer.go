// Simple map viewer that displays green for ground and blue for water

package mapviewer

import (
    "cartohelper/mapState"
    "image/color"
)

type SimpleMapViewer struct {
    mapState mapstate.MapState
}

func NewSimpleMapViewer() MapViewer {
    return &SimpleMapViewer{}
}

func (viewer *SimpleMapViewer) RegisterMap(mapState mapstate.MapState) {
    viewer.mapState = mapState
}

func (viewer *SimpleMapViewer) GetPixel(x, y int) color.Color {
    if viewer.mapState != nil {
        node, err := viewer.mapState.GetNode(x, y)
        if err != nil {
            return color.RGBA{0, 0, 0, 255}//black for all pixels outside of map
        }
        
        if node.GetHeight() > 0 {
            return color.RGBA{16, 191, 26, 255}//return green
        } else {
            return color.RGBA{17, 148, 218, 255}//return blue
        }
    }
    
    return nil
}

func (viewer *SimpleMapViewer) MapState() mapstate.MapState {
    return viewer.mapState
}
