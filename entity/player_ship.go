package entity

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/hajimehoshi/ebiten/v2"
    "github.com/infiniteyak/kessler_syndrome/layer"
	"github.com/infiniteyak/kessler_syndrome/component"
	"github.com/infiniteyak/kessler_syndrome/utility"
	"github.com/infiniteyak/kessler_syndrome/event"
	"github.com/infiniteyak/kessler_syndrome/asset"
    dmath "github.com/yohamta/donburi/features/math"
    "math"
	"github.com/hajimehoshi/ebiten/v2/audio"
    "log"
)

//TODO should this move to asset? there's probably a better way to do this
func generateShipVertices() []ebiten.Vertex {
    shape := []ebiten.Vertex{}
    var cr float32 = 1.0
    var cg float32 = 1.0
    var cb float32 = 1.0
    var scale float32 = 1.15
    shape = append(shape, ebiten.Vertex{
        DstX: 0.0 * scale,
        DstY: 6.0 * scale,
        SrcX: 0.0,
        SrcY: 0.0,
        ColorR: cr,
        ColorG: cg,
        ColorB: cb,
        ColorA: 1,
    })
    shape = append(shape, ebiten.Vertex{
        DstX: -4.0 * scale,
        DstY: -4.0 * scale,
        SrcX: 0.0,
        SrcY: 0.0,
        ColorR: cr,
        ColorG: cg,
        ColorB: cb,
        ColorA: 1,
    })
    shape = append(shape, ebiten.Vertex{
        DstX: 4.0 * scale,
        DstY: -4.0 * scale,
        SrcX: 0.0,
        SrcY: 0.0,
        ColorR: cr,
        ColorG: cg,
        ColorB: cb,
        ColorA: 1,
    })
    shape = append(shape, ebiten.Vertex{ //center
        DstX: 0.0 * scale,
        DstY: 0.0 * scale,
        SrcX: 0.0,
        SrcY: 0.0,
        ColorR: cr,
        ColorG: cg,
        ColorB: cb,
        ColorA: 1,
    })
    return shape
}

//TODO should this move to asset?
func generateThrusterVertices() []ebiten.Vertex {
    shape := []ebiten.Vertex{}
    var cr float32 = 1.0
    var cg float32 = 1.0
    var cb float32 = 1.0
    var scale float32 = 1.15

    shape = append(shape, ebiten.Vertex{
        DstX: 1.0 * scale,
        DstY: -4.0 * scale,
        SrcX: 0.0,
        SrcY: 0.0,
        ColorR: cr,
        ColorG: cg,
        ColorB: cb,
        ColorA: 1,
    })
    shape = append(shape, ebiten.Vertex{
        DstX: 1.5 * scale,
        DstY: -5.0 * scale,
        SrcX: 0.0,
        SrcY: 0.0,
        ColorR: cr,
        ColorG: cg,
        ColorB: cb,
        ColorA: 1,
    })
    shape = append(shape, ebiten.Vertex{
        DstX: 0.0 * scale,
        DstY: -8.0 * scale,
        SrcX: 0.0,
        SrcY: 0.0,
        ColorR: cr,
        ColorG: cg,
        ColorB: cb,
        ColorA: 1,
    })
    shape = append(shape, ebiten.Vertex{
        DstX: -1.5 * scale,
        DstY: -5.0 * scale,
        SrcX: 0.0,
        SrcY: 0.0,
        ColorR: cr,
        ColorG: cg,
        ColorB: cb,
        ColorA: 1,
    })
    shape = append(shape, ebiten.Vertex{
        DstX: -1.0 * scale,
        DstY: -4.0 * scale,
        SrcX: 0.0,
        SrcY: 0.0,
        ColorR: cr,
        ColorG: cg,
        ColorB: cb,
        ColorA: 1,
    })
    shape = append(shape, ebiten.Vertex{ //center
        DstX: 0.0 * scale,
        DstY: -5.0 * scale,
        SrcX: 0.0,
        SrcY: 0.0,
        ColorR: cr,
        ColorG: cg,
        ColorB: cb,
        ColorA: 1,
    })
    return shape
}

func AddPlayerShip(ecs *ecs.ECS, x, y float64, view *utility.View, audioContext *audio.Context) *donburi.Entity {
    entity := ecs.Create(
        layer.Foreground, 
        component.Position, 
        component.View,
        component.GraphicObject,
        component.Inputs,
        component.Actions,
        component.Velocity,
        component.Wrap,
        component.Collider,
        component.Health,
        component.Factions,
        )
    event.RegisterEntityEvent.Publish(ecs.World, event.RegisterEntity{Entity:&entity})

    entry := ecs.World.Entry(entity)

    // Factions
    factions := []component.FactionId{component.Player_factionid}
    donburi.SetValue(entry, component.Factions, factions)

    // Health
    donburi.SetValue(entry, component.Health, component.HealthData{Value:1.0})

    // Collider
    collider := component.NewColliderData()
    collider.Hitboxes = append(collider.Hitboxes, component.NewHitbox(5, 0, 0))
    donburi.SetValue(entry, component.Collider, collider)

    // Wrap
    wrap := component.WrapData{Distance: new(float64)}
    *wrap.Distance = 5.0
    donburi.SetValue(entry, component.Wrap, wrap)

    // Velocity
    vd := component.VelocityData{Velocity: &dmath.Vec2{}}
    donburi.SetValue(entry, component.Velocity, vd)

    // Position
    pd := component.NewPositionData(x, y)
    donburi.SetValue(entry, component.Position, pd)

    // View
    donburi.SetValue(entry, component.View, component.ViewData{View:view})

    // Graphic Object
    gobj := component.NewGraphicObjectData()

    shipPd := component.PolygonData{}
    shipPd.Load(generateShipVertices())
    gobj.Renderables = append(gobj.Renderables, &shipPd)

    thrustPd := component.PolygonData{}
    thrustPd.Load(generateThrusterVertices())
    *thrustPd.RenderableData.GetTransInfo().Hide = true //only show when applying thrust
    gobj.Renderables = append(gobj.Renderables, &thrustPd)

    donburi.SetValue(entry, component.GraphicObject, gobj)

    // Inputs
    im := make(map[component.ActionId]ebiten.Key)
    im[component.RotateCCW_actionid] = ebiten.KeyLeft
    im[component.RotateCW_actionid] = ebiten.KeyRight
    im[component.Accelerate_actionid] = ebiten.KeyUp
    im[component.Shoot_actionid] = ebiten.KeySpace
    donburi.SetValue(entry, component.Inputs, component.InputData{Mapping: im})

    // Actions
    tm := make(map[component.ActionId]bool)
    tm[component.Shield_actionid] = true //start out invulnerable
    cdm := make(map[component.ActionId]component.Cooldown)
    cdm[component.Shoot_actionid] = component.Cooldown{Cur:50, Max:50}
    cdm[component.Shield_actionid] = component.Cooldown{Cur:300, Max:300}
    am := make(map[component.ActionId]func())

    // Shoot
    bulletVelocity := dmath.Vec2{X:0, Y:1.3}
    am[component.Shoot_actionid] = func() {
        max := cdm[component.Shoot_actionid].Max
        cooldown := component.Cooldown{Cur:max, Max:max}
        cdm[component.Shoot_actionid] = cooldown

        ti := gobj.TransInfo
        bulletVector := bulletVelocity.Rotate(*ti.Rotation)
        bulletVector = vd.Velocity.Add(bulletVector)

        // Make the bullet spawn at the front of the ship, not the middle
        radius := 8.0
        angleCorrection := math.Pi / 2.0
        spawnX := math.Cos(*ti.Rotation + angleCorrection) * radius
        spawnY := math.Sin(*ti.Rotation + angleCorrection) * radius
        spawnX += pd.Point.X
        spawnY += pd.Point.Y

        AddBullet(ecs, spawnX, spawnY, *wrap.Distance, bulletVector, view, audioContext)
    }

    // Shield - actually turns off shield
    am[component.Shield_actionid] = func() {
        tm[component.Shield_actionid] = false
    }

    // Rotate CW
    rotationSpeed := 0.05
    am[component.RotateCW_actionid] = func() {
        ti := gobj.TransInfo
        *ti.Rotation += rotationSpeed
    }

    // Rotate CCW
    am[component.RotateCCW_actionid] = func() {
        ti := gobj.TransInfo
        *ti.Rotation -= rotationSpeed
    }

    // Apply thrust
    thrusterPower := dmath.Vec2{X:0, Y:0.01}
    thrustStarted := false
    thrustSoundSource := audio.NewInfiniteLoop(asset.ThrusterD, asset.ThrusterD.Length())
    thrustPlayer, err := audioContext.NewPlayer(thrustSoundSource)
    if err != nil {
        log.Fatal(err)
    }

    event.RegisterCleanupFuncEvent.Publish(
        ecs.World, 
        event.RegisterCleanupFunc{
            Function: func() {
                thrustPlayer.Close()
            },
        },
    )
    am[component.Accelerate_actionid] = func() {
        ti := gobj.TransInfo
        thrustVector := thrusterPower.Rotate(*ti.Rotation)
        *vd.Velocity = vd.Velocity.Add(thrustVector)
        *thrustPd.RenderableData.GetTransInfo().Hide = false

        if !thrustStarted {
            thrustStarted = true
            thrustPlayer.Rewind()
            thrustPlayer.Play()
        }
    }

    am[component.Destroy_actionid] = func() {
        event.RemoveEntityEvent.Publish(
            ecs.World, 
            event.RemoveEntity{Entity:&entity},
        )
        shipDestroyedDcopy := *asset.DestroyedD
        destroyedPlayer, err := audioContext.NewPlayer(&shipDestroyedDcopy)
        if err != nil {
            log.Fatal(err)
        }

        destroyedPlayer.Rewind()
        destroyedPlayer.Play()

        thrustPlayer.Pause()
        thrustPlayer.Close()

        event.ShipDestroyedEvent.Publish(ecs.World, event.ShipDestroyed{})
    }

    blinkCounter := 0
    am[component.Upkeep_actionid] = func() {
        // do a blinking effect if we are shielded
        if tm[component.Shield_actionid] {
            blinkCounter++
            if (blinkCounter / 10) % 2 == 0 {
                *thrustPd.RenderableData.GetTransInfo().Hide = true
                *shipPd.RenderableData.GetTransInfo().Hide = true
            } else {
                *shipPd.RenderableData.GetTransInfo().Hide = false
            }
        } else {
            *shipPd.RenderableData.GetTransInfo().Hide = false
        }
        if !tm[component.Accelerate_actionid] {
            //only show when applying thrust
            *thrustPd.RenderableData.GetTransInfo().Hide = true 
            thrustPlayer.Pause()
            thrustStarted = false
        }
    }

    donburi.SetValue(entry, component.Actions, component.ActionsData{
        TriggerMap: tm,
        CooldownMap: cdm,
        ActionMap: am,
    })

    return &entity
}
