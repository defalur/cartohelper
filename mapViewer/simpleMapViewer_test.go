package mapviewer

import (
    "testing"
    "cartohelper/mapState"
    "image/color"
)

func TestRegister(t *testing.T) {
    var tmp mapstate.MapState = mapstate.NewTileMapState(10, 10)
    var viewer MapViewer = &SimpleMapViewer{}
    viewer.RegisterMap(tmp)
    
    if tmp != viewer.MapState() {
        t.Errorf("Registering failed, map is not the same.")
    }
}

func TestWaterTile(t *testing.T) {
    var state mapstate.MapState = mapstate.NewTileMapState(10, 10)
    var viewer MapViewer = &SimpleMapViewer{}
    viewer.RegisterMap(state)
    
    c := viewer.GetPixel(5, 5)//should return a blue pixel
    c2 := color.RGBA{17, 148, 218, 255}
    if c != c2 {
        t.Errorf("Did not get the correct color: expected %v, got %v.", c2, c)
    }
}

func TestGroundTile(t *testing.T) {
    var state mapstate.MapState = mapstate.NewTileMapState(10, 10)
    var viewer MapViewer = &SimpleMapViewer{}
    viewer.RegisterMap(state)
    
    node, err := state.GetNode(5, 5)
    if err != nil {
        t.Error("Error while getting node.")
    }
    node.SetHeight(42)
    
    c := viewer.GetPixel(5, 5)
    c2 := color.RGBA{16, 191, 26, 255}
    if c != c2 {
        t.Errorf("Did not get the correct color: expected %v, got %v.", c2, c)
    }
}

func TestOutOfBoundsTile(t *testing.T) {
   var state mapstate.MapState = mapstate.NewTileMapState(10, 10)
   var viewer MapViewer = NewSimpleMapViewer()
   viewer.RegisterMap(state)
   
    c := viewer.GetPixel(11, 5)
    c2 := color.RGBA{0, 0, 0, 255}
    if c != c2 {
        t.Errorf("Did not get the correct color: expected %v, got %v.", c2, c)
    }
}

func TestNoMapRegistered(t *testing.T) {
    viewer := NewSimpleMapViewer()
    
    c := viewer.GetPixel(5, 5)
    if c != nil {
        t.Errorf("Getpixel: expected nil, got %v", c)
    }
}
