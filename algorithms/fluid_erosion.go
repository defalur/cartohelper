package algorithms

import (
    "errors"
    "cartohelper/mapState"
    "math"
    "fmt"
)

var startingWater float64 = 10.0
var waterLimit float64 = 64
var deltaTime float64 = 0.1
var A float64 = 1
var g float64 = 9.81
var l float64 = 1
var lx float64 = 1
var ly float64 = 1
var Kc float64 = 5.0
var Ks float64 = 5.0
var Kd float64 = 5.0

type tileData struct {
    totalFlow float64//total quantity of water that passed through the tile
    height float64//current height of the tile
    water float64//current amount of water on the tile
    sediment float64//amount of suspended sediment
    outTop float64//output flux top
    outBot float64//output flux bottom
    outLeft float64//output flux left
    outRight float64//output flux right
    velx float64//velocity vector x
    vely float64//velocity vector y
}

//============== Iteration functions ====================

func computeFlow(curOut, deltaHeight float64) float64 {
    return math.Max(0, curOut + deltaTime * A * (g * deltaHeight) / l)
}

func computeOutFlux(i, j int, mapData [][]tileData, h, w int) {
    //up
    neighbour, err := getTile(i - 1, j, h, w, mapData)
    if err == nil {
        deltaHeight := mapData[i][j].height + mapData[i][j].water - neighbour.height - neighbour.water
        mapData[i][j].outTop = computeFlow(mapData[i][j].outTop, deltaHeight)
    }
    //down
    neighbour, err = getTile(i + 1, j, h, w, mapData)
    if err == nil {
        deltaHeight := mapData[i][j].height + mapData[i][j].water - neighbour.height -
                            neighbour.water
        mapData[i][j].outBot = computeFlow(mapData[i][j].outBot, deltaHeight)
    }
    //left
    neighbour, err = getTile(i, j - 1, h, w, mapData)
    if err == nil {
        deltaHeight := mapData[i][j].height + mapData[i][j].water - neighbour.height - neighbour.water
        mapData[i][j].outLeft = computeFlow(mapData[i][j].outLeft, deltaHeight)
    }
    //right
    neighbour, err = getTile(i, j + 1, h, w, mapData)
    if err == nil {
        deltaHeight := mapData[i][j].height + mapData[i][j].water - neighbour.height - neighbour.water
        mapData[i][j].outRight = computeFlow(mapData[i][j].outRight, deltaHeight)
    }
    //scale down if necessary
    n := mapData[i][j]
    scale := math.Min(1, (n.water * lx * ly) /
                ((n,outTop + n.outBot + n.outLeft + n.outRight) * deltaTime))
    mapData[i][j].outTop *= scale
    mapData[i][j].outBot *= scale
    mapData[i][j].outLeft *= scale
    mapData[i][j].outRight *= scale
}

func updateWaterLvl(i, j int, mapData [][]tileData, h, w int) {
    oldLevel := mapData[i][j].water
    flowOut := mapData[i][j].outTop + mapData[i][j].outBot + mapData[i][j].outLeft +
                    mapData[i][j].outRight
    flowIn := 0.0
    //top
    var inTop float64
    node, err := getTile(i - 1, j, h, w, mapData)
    if err == nil {
        flowIn += node.outBot
        inTop = node.outBot
    }
    //bottom
    var inBot float64
    node, err := getTile(i + 1, j, h, w, mapData)
    if err == nil {
        flowIn += node.outTop
        inBot = node.outTop
    }
    //left
    var inLeft float64
    node, err := getTile(i, j - 1, h, w, mapData)
    if err == nil {
        flowIn += node.outRight
        inLeft = node.outRight
    }
    //right
    var inRight float64
    node, err := getTile(i, j + 1, h, w, mapData)
    if err == nil {
        flowIn += node.outLeft
        inRight = node.outLeft
    }
    
    //update water level
    mapData[i][j].water += (deltatime * (flowIn - flowOut)) / (lx * ly)
    
    //compute velocity vector
    levelAvg := (oldlevel + mapData[i][j].water) / 2
    deltaWx := (inLeft - outLeft + outRight - inRight) / 2
    deltaWy := (inTop - outTop + outBot - inBot) / 2
    mapData[i][j].velX = deltaWx / (levelAvg * ly)
    mapData[i][j].velY = deltaWy / (levelAvg * lx)
}

func computeGrad(a, b, d float64) float64 {
    return (b - a) / (2 * d)
}

func tiltAngleSin(x0, x1, y0, y1 float64) float64 {
    gradX := computeGrad(x0, x1, lx)
    gradY := computeGrad(y0, y1, ly)
    gradXsq := gradX * gradX
    gradYsq := gradY * gradY
    
    sumSq := gradXsq + gradYsq
    
    return (math.Sqrt(sumSq)) / (math.Sqrt(1 + sumSq))
}

func mag(x, y float64) float64 {
    return math.Sqrt(x * x + y * y)
}

func computeErosion(i, j int, mapData [][]tileData, h, w int) {
    //up
    up, err := getTile(i - 1, j, h, w, mapData)
    upHeight := mapData[i][j].height
    if err == nil {
        upHeight = up.height
    }
    //bottom
    down, err := getTile(i + 1, j, h, w, mapData)
    downHeight := mapData[i][j].height
    if err == nil {
        downHeight = down.height
    }
    //left
    left, err := getTile(i, j - 1, h, w, mapData)
    leftHeight := mapData[i][j].height
    if err == nil {
        leftHeight = left.height
    }
    //right
    right, err := getTile(i, j + 1, h, w, mapData)
    rightHeight := mapData[i][j].height
    if err == nil {
        rightHeight = right.height
    }
    
    //compute tilt angle
    n := mapData[i][j]
    //Compute capacity
    C := Kc * tiltAngleSin(leftHeight, rightHeight, upHeight, downHeight) * mag(n.velx, n.vely)
    //update terrain height and suspended soil
    curSediment := n.sediment
    if C > curSediment {
        dif := C - curSediment
        mapData[i][j].height = mapData[i][j].height - Ks * dif
        mapData[i][j].sediment = mapData[i][j].sediment + Ks * dif
    } else {
        dif := curSediment - C
        mapData[i][j].height = mapData[i][j].height + Kd * dif
        mapData[i][j].sediment = mapData[i][j].sediment - Kd * dif
    }
}

//============== Utility functions ======================

func getTile(i, j, h, w int, mapData [][]tileData) (*tileData, error) {
    if i < 0 || j < 0 || i >= h || j >= w {
        return nil, errors.New("Pos outside of map.")
    }
    
    return &mapData[i][j], nil
}

func iterate(mapData [][]tileData, h, w int,
             modifier func(int, int, [][]tileData, int, int)) float64 {
    waterCount := 0.0
    for i := 0; i < h; i++ {
        for j := 0; j < w; j++ {
            modifier(i, j, mapData, h, w)
            waterCount += mapData[i][j].water
        }
    }
    
    return waterCount
}

func initData(mapState mapstate.MapState) [][]tileData {
    result :=  make([][]tileData, mapState.GetHeight())
    for i := 0; i < mapState.GetHeight(); i++ {
        result[i] = make([]tileData, mapState.GetWidth())
        
        for j := 0; j < mapState.GetWidth(); j++ {
            node, _ := mapState.GetNode(j, i)
            result[i][j].height = float64(node.GetHeight())
            result[i][j].water = startingWater
        }
    }
    
    return result
}

//============ Main function =======================

func FluidErosion(mapState mapstate.MapState) {
    mapData := initData(mapState)
    //loop while there is water on the map{
    waterAmount := waterLimit + 1
    for waterAmount > waterLimit {
    //compute output flux for all tiles
        waterAmount = iterate(mapData, mapState.GetHeight(), mapState.GetWidth(), computeOutFlux)
    //update water level on all tiles
    //and compute velocity vector on all tiles
        iterate(mapData, mapState.GetHeight(), mapState.GetWidth(), updateWaterLvl)
    //compute sediment deposition and erosion
        iterate(mapData, mapState.GetHeight(), mapState.GetWidth(), computeErosion)
    //compute sediment transportation from velicoty vector
    //update water level with evaporation
        fmt.Println("Water: ", waterAmount)
    //}end loop
    }
}
