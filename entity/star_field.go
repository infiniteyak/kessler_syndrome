package entity

import (
    "github.com/infiniteyak/kessler_syndrome/layer"
	"github.com/infiniteyak/kessler_syndrome/component"
	"github.com/infiniteyak/kessler_syndrome/utility"
	"github.com/infiniteyak/kessler_syndrome/event"
	"github.com/yohamta/donburi"
    "github.com/yohamta/donburi/ecs"
)

func AddStarField(ecs *ecs.ECS, x, y float64, view *utility.View) *donburi.Entity {
    entity := ecs.Create(
        layer.Background, 
        component.Position, 
        component.GraphicObject,
        component.View,
    )
    event.RegisterEntityEvent.Publish(ecs.World, event.RegisterEntity{Entity:&entity})

    entry := ecs.World.Entry(entity)

    // Position
    pd := component.NewPositionData(x, y)
    donburi.SetValue(entry, component.Position, pd)

    // Graphic Object
    gobj := component.NewGraphicObjectData()
    nsfd := component.StarFieldData{}
    nsfd.Init()
    nsfd.Generate(int(view.Area.Max.X), int(view.Area.Max.Y))
    gobj.Renderables = append(gobj.Renderables, &nsfd)
    donburi.SetValue(entry, component.GraphicObject, gobj)

    // View
    donburi.SetValue(entry, component.View, component.ViewData{View:view})

    return &entity
}
