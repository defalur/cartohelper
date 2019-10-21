package algorithms

import (
    "cartohelper/mapState"
    //"fmt"
)

//get value on grid, if value is outside grid, returns the value on the border of the grid
func getValue(grid [][]float64, x, y, w, h int) float64 {
    if x < 0 {
        x = 0
    }
    if x >= w {
        x = w - 1
    }
    if y < 0 {
        y = 0
    }
    if y >= h {
        y = h - 1
    }
    
    return grid[y][x]
}

func computeDelta(prev, cur, next float64) float64 {
    return ((next - cur) + (cur - prev)) / 2
}

func updateMap(mapState mapstate.MapState, integral, heightMap [][]float64, rate float64) {
    for i := 0; i < mapState.GetHeight(); i++ {
        for j := 0; j < mapState.GetWidth(); j++ {
            node, err := mapState.GetNode(j, i)
            if err != nil {
                break
            }
            node.SetHeight(int(heightMap[i][j] - rate * integral[i][j]))
            //fmt.Printf("%f ", integral[i][j])
        }
        //fmt.Println("")
    }
}

func SimpleErosion(mapState mapstate.MapState, rate float64) {
    var heightMap [][]float64
    var deltaMap [][]float64
    var integralMap [][]float64

    heightMap = make([][]float64, mapState.GetHeight())
    deltaMap = make([][]float64, mapState.GetHeight())
    integralMap = make([][]float64, mapState.GetHeight())

    width := mapState.GetWidth()
    height := mapState.GetHeight()
    
    for i := 0; i < height; i++ {
        heightMap[i] = make([]float64, mapState.GetWidth())
        deltaMap[i] = make([]float64, mapState.GetWidth())
        integralMap[i] = make([]float64, mapState.GetWidth())
    }
    
    //fill heightmap with map data
    iterator := mapState.GetIterator()
    for {
        x, y := iterator.GetPos()
        node, err := iterator.Next()
        
        //if error, exit becquse we are done copying the heightmap
        if err != nil {
            break
        }
        
        heightMap[y][x] = float64(node.GetHeight())
    }
    
    /*for i := 0; i < height; i++ {
        for j := 0; j < width; j++ {
            fmt.Printf("%f ", heightMap[i][j])
        }
        fmt.Println("")
    }
    
    fmt.Println("============")
    */
    //compute delta for lines
    for i := 0; i < mapState.GetHeight(); i++ {
        for j := 0; j < mapState.GetWidth(); j++ {
            prev := getValue(heightMap, j - 1, i, width, height)
            cur := getValue(heightMap, j, i, width, height)
            next := getValue(heightMap, j + 1, i, width, height)
            delta := computeDelta(prev, cur, next)
            deltaMap[i][j] = delta
            //fmt.Printf("%f ", deltaMap[i][j])
        }
        //fmt.Println("")
    }
    
    //fmt.Println("============")
    
    //compute integral for horizontal component
    for i := 0; i < height; i++ {
        sum := 0.0
        for j := 0; j < width; j++ {
            sum += deltaMap[i][j]
            integralMap[i][j] += sum
        }
    }
    
    //compute delta for columns
    for i := 0; i < mapState.GetHeight(); i++ {
        for j := 0; j < mapState.GetWidth(); j++ {
            prev := getValue(heightMap, j, i - 1, width, height)
            cur := getValue(heightMap, j, i, width, height)
            next := getValue(heightMap, j, i + 1, width, height)
            delta := computeDelta(prev, cur, next)
            deltaMap[i][j] = delta
            //fmt.Printf("%f ", deltaMap[i][j])
        }
        //fmt.Println("")
    }
    
    //compute integral for vertical component
    for j := 0; j < width; j++ {
        sum := 0.0
        for i := 0; i < height; i++ {
            sum += deltaMap[i][j]
            integralMap[i][j] += sum
            integralMap[i][j] /= 2.0
        }
    }
    
    //fmt.Println("============")
    
    updateMap(mapState, integralMap, heightMap, rate)
}
