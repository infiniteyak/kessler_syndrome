package entity

import (
    "github.com/infiniteyak/kessler_syndrome/layer"
	"github.com/infiniteyak/kessler_syndrome/component"
	"github.com/infiniteyak/kessler_syndrome/utility"
	"github.com/infiniteyak/kessler_syndrome/event"
	"github.com/yohamta/donburi"
    "github.com/yohamta/donburi/ecs"
)

func AddSmallExplosion(ecs *ecs.ECS, x, y float64, view *utility.View) *donburi.Entity {
    entity := ecs.Create(
        layer.Foreground, 
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
    nsd := component.SpriteData{}
    nsd.Load("SmallExplosion", nil)
    nsd.SetLoopCallback(func() {
        event.RemoveEntityEvent.Publish(
            ecs.World, 
            event.RemoveEntity{Entity:&entity},
        )
    })
    nsd.SetPlaySpeed(3)
    gobj.Renderables = append(gobj.Renderables, &nsd)
    donburi.SetValue(entry, component.GraphicObject, gobj)

    // View
    donburi.SetValue(entry, component.View, component.ViewData{View:view})

    return &entity
}
