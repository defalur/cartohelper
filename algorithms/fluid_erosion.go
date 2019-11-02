package algorithms

import (
    "errors"
    "cartohelper/mapState"
    "math"
    "fmt"
)

var startingWater float64 = 1.0
var waterLimit float64 = 64
var deltaTime float64 = 0.1
//pipe cross section
var A float64 = 5
//gravity
var g float64 = 9.81
//length of pipes
var l float64 = 1
//width of tiles
var lx float64 = 1
//height of tiles
var ly float64 = 1
//sediment capacity constant
var Kc float64 = 0.1
//dissolving constant
var Ks float64 = 0.1
//sediment deposition constant
var Kd float64 = 0.1
//minimum terrain tilt constant
var minTilt float64 = 0.0001
//evaporation constant
var Ke float64 = 5.0

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
    tempHeight float64
    tempSediment float64
}

//============== Iteration functions ====================

func computeFlow(curOut, deltaHeight float64) float64 {
    ret := math.Max(0, curOut + deltaTime * A * (g * deltaHeight) / l)

    //if math.IsNaN(ret) {
    //    fmt.Println("out: ", curOut, " deltaHeight: ", deltaHeight)
    //}
    return ret
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
    sumOut := (n.outTop + n.outBot + n.outLeft + n.outRight) * deltaTime
    scale := 1.0
    if sumOut != 0.0 {
        scale = math.Min(1, (n.water * lx * ly) / sumOut)
    }
    
    //fmt.Println("Height: ", n.height)
    //fmt.Println("Sum: ", sumOut)            
    //fmt.Println("Scale: ", scale)
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
    node, err = getTile(i + 1, j, h, w, mapData)
    if err == nil {
        flowIn += node.outTop
        inBot = node.outTop
    }
    //left
    var inLeft float64
    node, err = getTile(i, j - 1, h, w, mapData)
    if err == nil {
        flowIn += node.outRight
        inLeft = node.outRight
    }
    //right
    var inRight float64
    node, err = getTile(i, j + 1, h, w, mapData)
    if err == nil {
        flowIn += node.outLeft
        inRight = node.outLeft
    }
    
    //update water level
    //fmt.Println("in: ", flowIn, " out: ", flowOut)
    mapData[i][j].water += (deltaTime * (flowIn - flowOut)) / (lx * ly)
    
    //compute velocity vector
    levelAvg := (oldLevel + mapData[i][j].water) / 2
    if levelAvg == 0.0 {
        levelAvg = 0.1
    }
    deltaWx := (inLeft - mapData[i][j].outLeft + mapData[i][j].outRight - inRight) / 2
    deltaWy := (inTop - mapData[i][j].outTop + mapData[i][j].outBot - inBot) / 2
    mapData[i][j].velx = deltaWx / (levelAvg * ly)
    mapData[i][j].vely = deltaWy / (levelAvg * lx)
}

func computeGrad(a, b, d float64) float64 {
    ret := (b - a) / (2 * d)
    
    //if math.IsNaN(ret) {
    //    fmt.Println("a: ", a, " b: ", b, " d: ", d)
    //}
    return ret
}

func tiltAngleSin(x0, x1, y0, y1 float64) float64 {
    gradX := computeGrad(x0, x1, lx)
    gradY := computeGrad(y0, y1, ly)
    gradXsq := gradX * gradX
    gradYsq := gradY * gradY
    
    sumSq := gradXsq + gradYsq

    ret := (math.Sqrt(sumSq)) / (math.Sqrt(1 + sumSq))
    
    //if math.IsNaN(ret) {
    //    fmt.Println("x0: ", x0, " x1: ", x1, " y0: ", y0, " y1: ", y1)
    //}
    
    return ret
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
    C := Kc * math.Max(minTilt, tiltAngleSin(leftHeight, rightHeight, upHeight, downHeight)) *
                mag(n.velx, n.vely)

    //fmt.Println("C: ", C)
    //update terrain height and suspended soil
    curSediment := n.sediment
    //if math.IsNaN(curSediment) {
    //    fmt.Println("sediment")
    //}
    if C > curSediment {
        dif := C - curSediment
        mapData[i][j].tempHeight = mapData[i][j].height - Ks * dif
        mapData[i][j].tempSediment = mapData[i][j].sediment + Ks * dif
    } else {
        dif := curSediment - C
        mapData[i][j].tempHeight = mapData[i][j].height + Kd * dif
        mapData[i][j].tempSediment = mapData[i][j].sediment - Kd * dif
    }
    //if math.IsNaN(mapData[i][j].sediment) {
    //    fmt.Println("updated sediment")
    //}
}

func sedimentInterpolation(i, j float64, h, w int, mapData [][]tileData) float64 {
    //interger coordinates
    fi := math.Floor(i)
    ci := math.Ceil(i)
    fj := math.Floor(j)
    cj := math.Ceil(j)
    //a
    sA := 0.0
    a, err := getTile(int(fi), int(fj), h, w, mapData)
    if err == nil {
        sA = a.sediment
    }
    //b
    sB := 0.0
    b, err := getTile(int(fi), int(cj), h, w, mapData)
    if err == nil {
        sB = b.sediment
    }
    //c
    sC := 0.0
    c, err := getTile(int(ci), int(fj), h, w, mapData)
    if err == nil {
        sC = c.sediment
    }
    //d
    sD := 0.0
    d, err := getTile(int(ci), int(cj), h, w, mapData)
    if err == nil {
        sD = d.sediment
    }
    
    //interpolation
    //move window in [0, 1] interval to simplify calculations
    di := i - fi
    dj := j - fj
    ret := sA * (1 - dj) * (1 - di) + sB * dj * (1 - di) + sC * (1 - dj) * i + sD * dj * di
    
    //if math.IsNaN(ret) {
    //    fmt.Println("interpolation: (", i, ",", j, ")")
    //}
    return ret
}

func sedimentTransport(i, j int, mapData [][]tileData, h, w int) {
    newI := float64(i) - mapData[i][j].vely * deltaTime
    newJ := float64(j) - mapData[i][j].velx * deltaTime
    mapData[i][j].tempSediment = sedimentInterpolation(newI, newJ, h, w, mapData)
}

func updateHeight(i, j int, mapData [][]tileData, h, w int) {
    mapData[i][j].height = mapData[i][j].tempHeight
}

func updateSediment(i, j int, mapData [][]tileData, h, w int) {
    mapData[i][j].sediment = mapData[i][j].tempSediment
}

func evaporate(i, j int, mapData [][]tileData, h, w int) {
    mapData[i][j].water = math.Max(0.0, mapData[i][j].water * (1 - Ke * deltaTime))
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

func updateMapFluid(mapState mapstate.MapState, mapData [][]tileData) {
    for i := 0; i < mapState.GetHeight(); i++ {
        for j := 0; j < mapState.GetWidth(); j++ {
            node, err := mapState.GetNode(j, i)
            if err != nil {
                break
            }
            node.SetHeight(int(math.Round(mapData[i][j].height)))
        }
    }
    mapState.UpdateExtrema()
}

//============ Main function =======================

func FluidErosion(mapState mapstate.MapState) {
    mapData := initData(mapState)
    //loop while there is water on the map{
    waterAmount := waterLimit + 1
    for waterAmount > waterLimit {
    //compute output flux for all tiles
        //fmt.Println("Ouput flux.")
        waterAmount = iterate(mapData, mapState.GetHeight(), mapState.GetWidth(), computeOutFlux)
        //fmt.Println("Water: ", waterAmount)
    //update water level on all tiles
    //and compute velocity vector on all tiles
        //fmt.Println("WaterLvl and speed vectors.")
        waterAmount = iterate(mapData, mapState.GetHeight(), mapState.GetWidth(), updateWaterLvl)
        //fmt.Println("Water: ", waterAmount)
    //compute sediment deposition and erosion
        //fmt.Println("Erosion.")
        waterAmount = iterate(mapData, mapState.GetHeight(), mapState.GetWidth(), computeErosion)
        //fmt.Println("Water: ", waterAmount)
    //compute sediment transportation from velicoty vector
        //fmt.Println("Sediment transport.")
        waterAmount = iterate(mapData, mapState.GetHeight(), mapState.GetWidth(), sedimentTransport)
        //fmt.Println("Water: ", waterAmount)
        iterate(mapData, mapState.GetHeight(), mapState.GetWidth(), updateHeight)
        iterate(mapData, mapState.GetHeight(), mapState.GetWidth(), updateSediment)
    //update water level with evaporation
        //fmt.Println("Evaporation.")
        waterAmount = iterate(mapData, mapState.GetHeight(), mapState.GetWidth(), evaporate)
        fmt.Println("Water: ", waterAmount)
    //}end loop
    }
    //update mapState
    updateMapFluid(mapState, mapData)
}
