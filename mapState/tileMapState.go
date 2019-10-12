package mapstate

import (
    "errors"
    "math/rand"
    "github.com/golang-collections/go-datastructures/queue"
)

type TileMapState struct {
    nodes [][]TileMapNode
    width int
    height int
}

func NewTileMapState(width, height, baseHeight int) MapState {
    result := &TileMapState{width:width, height:height}
    result.nodes = make([][]TileMapNode, height)
    for i := 0; i < height; i++ {
        result.nodes[i] = make([]TileMapNode, width)
        
        for j := 0; j < width; j++ {
            result.nodes[i][j] = NewTileMapNode(baseHeight)//TODO: baseHeight broke all tests
        }
    }
    
    for i := 0; i < height; i++ {
        for j := 0; j < width; j++ {
            result.fillNeighbours(&result.nodes[i][j], j, i)
        }
    }
    
    rand.Seed(0)
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

func (state *TileMapState) Seed(seed int) {
    rand.Seed(int64(seed))
}

func (state *TileMapState) EndBlob() {
    for i := 0; i < state.height; i++ {
        for j := 0; j < state.width; j++ {
            state.nodes[i][j].Mod = false
        }
    }
}

func (state *TileMapState) GenerateBlob() (height, width int) {    
    width = rand.Intn(40 - 5) + 5//random in [3; 15[ range
    height = rand.Intn(30 - 5) + 5//random in [5; 20[
    curHeight := float32(height)
    slope := curHeight / float32(width)
    
    posx := rand.Intn(state.width)
    posy := rand.Intn(state.height)
    
    node, _ := state.GetNode(posx, posy)//position is always valid
    q := queue.New(0)
    q.Put(node)
    q.Put(nil)
    
    for ;q.Len() > 1; {
        tmp, _ := q.Get(1)
        curNode, _ := tmp[0].(MapNode)
        
        if curNode == nil { //go one step down
            curHeight = curHeight - slope
            if curHeight <= 0.01 {
                break;
            }
            q.Put(nil)
        } else {
            curNode.SetHeight(curNode.GetHeight() + int(curHeight))
            for _, n := range curNode.Neighbours() {
                if !n.Modified() {
                    q.Put(n)
                }
            }
        }
    }
    
    state.EndBlob()
    q.Dispose()
    return
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
