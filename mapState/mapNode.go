package mapstate

import (
    
)

type MapNode interface {
    GetHeight() int
    SetHeight(height int)
    //getters and setters for humidity and biome
    Neighbours() []MapNode
    Modified() bool
}

type TileMapNode struct {
    height int
    neighbours []MapNode
    Mod bool
}

func NewTileMapNode(height int) TileMapNode {
    result := TileMapNode{height: height}
    result.neighbours = make([]MapNode, 0, 8)
    result.Mod = false
    
    return result
}

func (n *TileMapNode) Modified() bool {
    mod := n.Mod
    n.Mod = true
    return mod
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
