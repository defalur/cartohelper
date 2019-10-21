package mapstate

import (
    "errors"
    "math/rand"
    "github.com/golang-collections/go-datastructures/queue"
    "cartohelper/genutils"
    "math"
    "fmt"
)

type TileMapState struct {
    nodes [][]TileMapNode
    probMap [][]float64 //probability map for 
    width int
    height int
    maxHeight int
    minHeight int
}

type TileMapIterator struct {
    x int
    y int
    state *TileMapState
}

func (it *TileMapIterator) GetPos() (int, int) {
    return it.x, it.y
}

func newTileMapIterator(state *TileMapState) MapIterator {
    ret := &TileMapIterator{x: 0, y: 0, state: state}
    
    return ret
}

func (it *TileMapIterator) Next() (MapNode, error) {
    ret, err := it.state.GetNode(it.x, it.y)
    
    if err == nil {
        it.x++
        if it.x >= it.state.GetWidth() {
            it.x = 0
            it.y++
        }
    }
    
    return ret, err
}

func NewTileMapState(width, height, baseHeight int) MapState {
    result := &TileMapState{width:width, height:height}
    result.nodes = make([][]TileMapNode, height)
    result.probMap = make([][]float64, height)
    for i := 0; i < height; i++ {
        result.nodes[i] = make([]TileMapNode, width)
        result.probMap[i] = make([]float64, width)
        
        for j := 0; j < width; j++ {
            result.nodes[i][j] = NewTileMapNode(baseHeight)//TODO: baseHeight broke all tests
            result.probMap[i][j] = 0.0
        }
    }
    
    for i := 0; i < height; i++ {
        for j := 0; j < width; j++ {
            result.fillNeighbours(&result.nodes[i][j], j, i)
        }
    }
    
    result.maxHeight = baseHeight
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

func (state *TileMapState) AddDistributionBlob(x, y int, radius float64) {
    modMap := make(map[string]bool)
    
    origin := genutils.Vector2{X: x, Y: y}
    curPoint := genutils.Vector2{X: x, Y: y}
    q := queue.New(0)
    q.Put(curPoint)
    
    pointAdd := func(x, y int, curPoint genutils.Vector2) {
        p := genutils.Vector2{X: curPoint.X + x, Y: curPoint.Y + y}
        pString := fmt.Sprintf("%d:%d", p.X, p.Y)
        mag := genutils.Magnitude(origin, p)
        //fmt.Println("Add: ", p, "pString: ", pString)
        if mag < radius && !modMap[pString] {
            modMap[pString] = true
            state.probMap[p.Y][p.X] = math.Max(1 - mag / radius, state.probMap[p.Y][p.X])
            q.Put(p)
        }
    }
    
    for ;!q.Empty(); {
        tmp, _ := q.Get(1)
        curPoint, _ = tmp[0].(genutils.Vector2)
        
        if curPoint.X > 0 {
           pointAdd(-1, 0, curPoint)
        }
        if curPoint.X < state.width - 1 {
            pointAdd(1, 0, curPoint)
        }
        if curPoint.Y > 0 {
            pointAdd(0, -1, curPoint)
        }
        if curPoint.Y < state.height - 1 {
            pointAdd(0, 1, curPoint)
        }
    }
}

func (state *TileMapState) GetDistribution(x, y int) float64 {
    if x < 0 || y < 0 || x >= state.width || y >= state.height {
        return 0.0
    }
    return state.probMap[y][x]
}

func (state *TileMapState) EndBlob() {
    for i := 0; i < state.height; i++ {
        for j := 0; j < state.width; j++ {
            state.nodes[i][j].Mod = false
        }
    }
}

func (state *TileMapState) GenerateBlob(x, y, w, h int) (posx, posy int) {    
    width := rand.Intn(40 - 5) + 5//random in [3; 15[ range
    height := rand.Intn(30 - 5) + 5//random in [5; 20[
    curHeight := float32(height)
    slope := curHeight / float32(width)
    
    posRetries := 3
    
    posx = rand.Intn(w) + x
    posy = rand.Intn(h) + y
    lastProb := state.probMap[posy][posx]
    
    for i := 0; i < posRetries; i++ {
        tmpx := rand.Intn(w) + x
        tmpy := rand.Intn(h) + y
        curProb := state.probMap[tmpy][tmpx]
        
        if curProb > lastProb {
            posx = tmpx
            posy = tmpy
            lastProb = curProb
        }
    }
    
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
            if curNode.GetHeight() > state.maxHeight {
                state.maxHeight = curNode.GetHeight()
            }
            if curNode.GetHeight() < state.minHeight {
                state.minHeight = curNode.GetHeight()
            }
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

func (state *TileMapState) MaxHeight() int {
    return state.maxHeight
}

func (state *TileMapState) MinHeight() int {
    return state.minHeight
}

func (state *TileMapState) GetIterator() MapIterator {
    return newTileMapIterator(state)
}
