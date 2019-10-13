package genutils

import (
    "math/rand"
)

type ContinentData struct {
    //zone
    X int
    Y int
    Width int
    Height int
    
    //density
    NBlob int64
}

//takes dimensions of the map and returns a new continentData instance
func NewContinent(width, height int) ContinentData {
    dimX := rand.Intn(width / 2 - width / 4) + width / 4
    dimY := rand.Intn(height / 2 - height / 4) + height / 4
    
    x := rand.Intn(width - dimX)
    y := rand.Intn(height - dimY)
    
    density := rand.Float32()
    nBlob := int64(float32(dimX * dimY) * density * 0.007)
    
    return ContinentData{X: x, Y: y, Width: dimX, Height: dimY, NBlob: nBlob}
}
