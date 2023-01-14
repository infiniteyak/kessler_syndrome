package entity

import (
	"image"
	"github.com/infiniteyak/kessler_syndrome/asset"
    "github.com/infiniteyak/kessler_syndrome/layer"
	"github.com/infiniteyak/kessler_syndrome/component"
	"github.com/infiniteyak/kessler_syndrome/utility"
	"github.com/infiniteyak/kessler_syndrome/event"

	"github.com/yohamta/donburi"
    "github.com/yohamta/donburi/ecs"
)

//ABCDEFGHIJKLM
//NOPQRSTUVWXYZ
//0123456789
//TODO should this be a const? in assets maybe
var letterMasks = map[string]image.Rectangle{
    "A": image.Rect(0,0,8,8),
    "B": image.Rect(8,0,16,8),
    "C": image.Rect(16,0,24,8),
    "D": image.Rect(24,0,32,8),
    "E": image.Rect(32,0,40,8),
    "F": image.Rect(40,0,48,8),
    "G": image.Rect(48,0,56,8),
    "H": image.Rect(56,0,64,8),
    "I": image.Rect(64,0,72,8),
    "J": image.Rect(72,0,80,8),
    "K": image.Rect(80,0,88,8),
    "L": image.Rect(88,0,96,8),
    "M": image.Rect(96,0,104,8),
    "N": image.Rect(0,8,8,16),
    "O": image.Rect(8,8,16,16),
    "P": image.Rect(16,8,24,16),
    "Q": image.Rect(24,8,32,16),
    "R": image.Rect(32,8,40,16),
    "S": image.Rect(40,8,48,16),
    "T": image.Rect(48,8,56,16),
    "U": image.Rect(56,8,64,16),
    "V": image.Rect(64,8,72,16),
    "W": image.Rect(72,8,80,16),
    "X": image.Rect(80,8,88,16),
    "Y": image.Rect(88,8,96,16),
    "Z": image.Rect(96,8,104,16),
    "0": image.Rect(0,16,8,24),
    "1": image.Rect(8,16,16,24),
    "2": image.Rect(16,16,24,24),
    "3": image.Rect(24,16,32,24),
    "4": image.Rect(32,16,40,24),
    "5": image.Rect(40,16,48,24),
    "6": image.Rect(48,16,56,24),
    "7": image.Rect(56,16,64,24),
    "8": image.Rect(64,16,72,24),
    "9": image.Rect(72,16,80,24),
    "_": image.Rect(80,16,88,24),
    ".": image.Rect(88,16,96,24),
    "^": image.Rect(96,16,104,24),
}

type FontAlignY int
const (
    Top_fontaligny FontAlignY = iota
    Middle_fontaligny
    Bottom_fontaligny
)
type FontAlignX int
const (
    Left_fontalignx FontAlignX = iota
    Center_fontalignx
    Right_fontalignx
)

type StringData struct {
    String string
    XAlign FontAlignX
    YAlign FontAlignY
    Font string
    Kerning int
    Blink bool
}

//TODO support multiple lines of text
func writeText(curString *StringData, gobj *component.GraphicObjectData) {
    var textY float64
    switch curString.YAlign {
    case Top_fontaligny:
        textY = float64(asset.FontHeight/2)
    case Middle_fontaligny:
        textY = 0.0
    case Bottom_fontaligny:
        textY = -1.0 * float64(asset.FontHeight/2)
    }

    var textX float64
    switch curString.XAlign {
    case Left_fontalignx:
        textX = float64(asset.FontWidth/2) //Because each sprite pos is char center
    case Center_fontalignx:
        textX = float64(
(asset.FontWidth * (1 - len(curString.String)) + -1 * curString.Kerning * (len(curString.String) - 1)) / 2)
    case Right_fontalignx:
        textX = float64(asset.FontWidth * (0 - (len(curString.String) + curString.Kerning * (len(curString.String) - 1))))
    }

    gobj.Renderables = []component.Renderable{}
    for _, c := range curString.String {
        nsd := component.SpriteData{}
        m := letterMasks[string(c)]
        nsd.Load(curString.Font, &m)
        tinfo := nsd.RenderableData.GetTransInfo()
        tinfo.Offset.X = textX
        tinfo.Offset.Y = textY
        nsd.RenderableData.SetTransInfo(tinfo)
        gobj.Renderables = append(gobj.Renderables, &nsd)
        textX += float64(asset.FontWidth) + float64(curString.Kerning)
    }
}

func AddSpriteText(ecs *ecs.ECS, x, y float64, view *utility.View, str *StringData) *donburi.Entity {
    entity := ecs.Create(
        layer.HudForeground, //TODO make these all take layers IDs as args?
        component.Position, 
        component.GraphicObject,
        component.View,
        component.Actions,
    )
    event.RegisterEntityEvent.Publish(ecs.World, event.RegisterEntity{Entity:&entity})

    entry := ecs.World.Entry(entity)

    curString := *str

    pd := component.NewPositionData(x, y)
    donburi.SetValue(entry, component.Position, pd)

    gobj := component.NewGraphicObjectData()
    writeText(&curString, &gobj)
    donburi.SetValue(entry, component.GraphicObject, gobj)

    donburi.SetValue(entry, component.View, component.ViewData{View:view})

    // Actions
    tm := make(map[component.ActionId]bool)
    cdm := make(map[component.ActionId]component.Cooldown)
    am := make(map[component.ActionId]func())
    blinkCounter := 0
    am[component.Upkeep_actionid] = func() { 
        if curString != *str {
            curString = *str
            g := component.GraphicObject.Get(entry)
            writeText(&curString, g)
        }
        if curString.Blink {
            blinkCounter++
            if (blinkCounter / 20) % 2 == 0 {
                *gobj.TransInfo.Hide = true
            } else {
                *gobj.TransInfo.Hide = false
            }
        } else {
            *gobj.TransInfo.Hide = false
        }
    }

    donburi.SetValue(entry, component.Actions, component.ActionsData{
        TriggerMap: tm,
        CooldownMap: cdm,
        ActionMap: am,
    })

    return &entity
}
