package genutils

import (
    "math/rand"
    "math"
)

type Vector2 struct {
    X int
    Y int
}

type DistributionBlob struct {
    X int
    Y int
    
    Radius float64
}

type ContinentData struct {
    //zone
    X int
    Y int
    Width int
    Height int
    
    //density
    NBlob int64
    
    Probs []DistributionBlob
}

//takes dimensions of the map and returns a new continentData instance
func NewContinent(width, height int) ContinentData {
    dimX := rand.Intn(width / 2 - width / 4) + width / 4
    dimY := rand.Intn(height / 2 - height / 4) + height / 4
    
    x := rand.Intn(width - dimX)
    y := rand.Intn(height - dimY)
    
    density := rand.Float32()
    nBlob := int64(float32(dimX * dimY) * density * 0.007)
    
    result := ContinentData{X: x, Y: y, Width: dimX, Height: dimY, NBlob: nBlob}
    
    result.Probs = make([]DistributionBlob, 4)//distribution blobs as a test value
    for i := 0; i < len(result.Probs); i++ {
        result.Probs[i].X = rand.Intn(dimX) + x
        result.Probs[i].Y = rand.Intn(dimY) + y
        
        result.Probs[i].Radius = rand.Float64() * (float64(dimX + dimY) / 4.0)
    }
    
    return result
}

func Magnitude(a, b Vector2) float64 {
    dx := float64(a.X - b.X)
    dy := float64(a.Y - b.Y)
    return math.Sqrt(dx * dx + dy * dy)
}
