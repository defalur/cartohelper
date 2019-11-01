package algorithms

import (
    "errors"
    "cartohelper/mapState"
    "math"
)

var startingWater float64 = 10.0
var waterLimit float64 = 64
var deltaTime float64 = 0.1
var A float64 = 1
var g float64 = 9.81
var l float64 = 1

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
        deltaHeight := mapData[i][j].height + mapData[i][j].water - neighbour.height - neighbour.water
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
    //compute velocity vector on all tiles
    //compute sediment deposition and erosion
    //compute sediment transportation from velicoty vector
    //update water level with evaporation
    //}end loop
    }
}
