package mapState

import (
    
)

type MapState interface {
    GenerateBlob()//find parameters
    GetNode(x, y int) (MapNode, error)//position in integer coordinates?, also return node position
    GetWidth() int
    GetHeight() int
    //add ctl function if necessary
}
