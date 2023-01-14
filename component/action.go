package component

import (
    "github.com/yohamta/donburi"
)

type ActionId int

const (
    Undefined_actionid ActionId = iota
    Upkeep_actionid //TODO this is kind of hacky, no?
    MoveLeft_actionid
    MoveRight_actionid
    MoveDown_actionid
    MoveUp_actionid
    RotateCW_actionid
    RotateCCW_actionid
    Accelerate_actionid
    Shoot_actionid
    TriggerFunction_actionid
    SelfDestruct_actionid
    DestroySilent_actionid
    Destroy_actionid
    Shield_actionid
    Blink_actionid
)

type ActionsData struct {
    TriggerMap map[ActionId]bool
    CooldownMap map[ActionId]Cooldown
    ActionMap map[ActionId]func()
}

var Actions = donburi.NewComponentType[ActionsData]()

type Cooldown struct {
    Cur int
    Max int
}
