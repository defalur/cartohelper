package mapstate

import (
    
)

type MapState interface {
    GenerateBlob(x, y, w, h int) (int, int)//find parameters
    GetNode(x, y int) (MapNode, error)//position in integer coordinates?, also return node position
    GetWidth() int
    GetHeight() int
    Seed(seed int)
    //add ctl function if necessary
}
