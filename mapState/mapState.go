package mapState

import (
    
)

IMapState interface {
    GenerateBlob()//find parameters
    GetNode(x, y int)//position in integer coordinates?, also return node position
    GetWidth()
    GetHeight()
    //add ctl function if necessary
}
