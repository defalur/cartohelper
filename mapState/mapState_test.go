package mapState

import (
    "testing"
)

func TestDimensions(t *testing.T) {
    var tmp MapState = NewTileMapState(10, 10)
    if !(tmp.GetHeight() == 10) {
        t.Errorf("GetHeight() = %d; want 10", tmp.GetHeight())
    }
    if !(tmp.GetWidth() == 10) {
        t.Errorf("GetWidth() = %d; want 10", tmp.GetWidth())
    }
}

func TestGetNodeNormal(t *testing.T) {
    var tmp MapState = NewTileMapState(10, 10)
    node, err := tmp.GetNode(5, 5)
    if err != nil || node == nil {
        t.Errorf("GetNode(), got error, expected no error.")
    }
}

func TestGetNodeError(t *testing.T) {
    var tmp MapState = NewTileMapState(10, 10)
    _, err := tmp.GetNode(11, 10)
    if err == nil {
        t.Errorf("GetNode(), did not get error, expectect an error.")
    }
}

func TestNodeReference(t *testing.T) {
    var state MapState = NewTileMapState(10, 10)
    node1, err := state.GetNode(5, 5)
    if err != nil {
        t.Errorf("GetNode() returned an error, expected no error.")
    }
    node1.SetHeight(42)
    node2, err := state.GetNode(5, 5)
    if err != nil {
        t.Errorf("GetNode() returned an error, expected no error.")
    }
    if node2.GetHeight() != 42 {
        t.Errorf("node2.GetHeight() returned %d, expected 42.", node2.GetHeight())
    }
}

func TestNodeNeighbours(t *testing.T) {
   var state MapState = NewTileMapState(10, 10)
   node1, err := state.GetNode(5, 5)
    if err != nil {
        t.Errorf("GetNode() returned an error, expected no error.")
    }
    
    neighbours := node1.Neighbours()
    if len(neighbours) == 0 {
       t.Errorf("Neighbour length is null.") 
    }
    for _, n := range neighbours {
        n.SetHeight(42)
    }
    
    node2, err := state.GetNode(5, 6)
    if err != nil {
        t.Errorf("GetNode() returned an error, expected no error.")
    }
    if node2.GetHeight() != 42 {
        t.Errorf("node2.GetHeight() returned %d, expected 42.", node2.GetHeight())
    }
}
