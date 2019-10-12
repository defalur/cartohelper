package mapviewer

import (
    "cartohelper/mapState"
    "image/color"
)

type MapViewer interface {
    RegisterMap(state mapstate.MapState)
    GetPixel(x, y int) color.Color
    MapState() mapstate.MapState
}
