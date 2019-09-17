package mapState

import (
    
)

IMapNode interface {
    GetHeight() float32
    SetHeight(height float32)
    //getters and setters for humidity and biome
    Neighbours() []MapNode
}

TileMapNode struct {
    height float32
    neighbours [8]MapNode
}

func (n *TileMapNode) GetHeight() float32 {
    return n.height
}

func (n *TileMapNode) SetHeight(height float32) {
    n.height = height
}

func (n *TileMapNode) Neighbours() []MapNode {
    return n.neighbours[:]
}
