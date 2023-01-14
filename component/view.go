package component

import (
    "github.com/yohamta/donburi"
    "github.com/infiniteyak/kessler_syndrome/utility"
)

type ViewData struct {
    View *utility.View
}

var View = donburi.NewComponentType[ViewData]()
