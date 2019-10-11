package mapState

import (
    "errors"
)

type TileMapState struct {
    nodes [][]TileMapNode
    width int
    height int
}

func NewTileMapState(width, height int) MapState {
    result := &TileMapState{width:width, height:height}
    result.nodes = make([][]TileMapNode, height)
    for i := 0; i < height; i++ {
        result.nodes[i] = make([]TileMapNode, width)
        
        for j := 0; j < width; j++ {
            result.nodes[i][j] = NewTileMapNode(0)
        }
    }
    
    for i := 0; i < height; i++ {
        for j := 0; j < width; j++ {
            result.fillNeighbours(&result.nodes[i][j], j, i)
        }
    }
    return result
}

func (state *TileMapState) fillNeighbours(node *TileMapNode, x, y int) {
    for i := -1; i < 2; i++ {
        for j := -1; j < 2; j++ {
            tmpNode, err := state.GetNode(j + x, i + y)
            if err == nil &&  (i == 0) != (j == 0) {
                node.neighbours = append(node.neighbours, tmpNode)
            }
        }
    }
}

func (state *TileMapState) GenerateBlob() {
    
}

func (state *TileMapState) GetNode(x, y int) (MapNode, error) {
    if (x >= 0 && y >= 0 && x < state.width && y < state.height) {
        return &state.nodes[y][x], nil
    }
    return nil, errors.New("Position outside of map.")
}

func (state *TileMapState) GetWidth() int {
    return state.width
}

func (state *TileMapState) GetHeight() int {
    return state.height
}
