package mapState

import (
    
)

type MapNode interface {
    GetHeight() int
    SetHeight(height int)
    //getters and setters for humidity and biome
    Neighbours() []MapNode
}

type TileMapNode struct {
    height int
    neighbours []MapNode
}

func NewTileMapNode(height int) TileMapNode {
    result := TileMapNode{height: height}
    result.neighbours = make([]MapNode, 0, 8)
    
    return result
}

func (n *TileMapNode) GetHeight() int {
    return n.height
}

func (n *TileMapNode) SetHeight(height int) {
    n.height = height
}

func (n *TileMapNode) Neighbours() []MapNode {
    return n.neighbours[:]
}
