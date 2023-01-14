package game

import (
	"os"
	"github.com/infiniteyak/kessler_syndrome/asset"
    "github.com/infiniteyak/kessler_syndrome/layer"
	"github.com/infiniteyak/kessler_syndrome/system"
	"github.com/infiniteyak/kessler_syndrome/utility"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/events"
)

type scoreEntry struct {
    initials string
    score int
    id int
}

type Game struct {
    screenView *utility.View //view equiv of the full screen
    curScene *Scene
    ecs *ecs.ECS
    states map[SceneId]map[SceneEventId]func() 
    curScore int
    curWave int
    curShips int
    playerInitials string
    highScores []scoreEntry
    curScoreId int
    audioContext *audio.Context
}

func NewGame(width, height float64) *Game {
	world := donburi.NewWorld()
	ecs := ecs.NewECS(world)
    this := &Game{
        screenView: utility.NewView(0.0, 0.0, width, height),
        ecs: ecs,
    }

    this.audioContext = audio.NewContext(48000)
    
    this.states = map[SceneId]map[SceneEventId]func() {
        Undefined_sceneId: { 
            Init_sceneEvent: this.LoadSplashScene,
        },
        Splash_sceneId: { 
            Advance_sceneEvent: this.LoadPlayingScene,
        },
        Playing_sceneId: {
            GameOver_sceneEvent: this.LoadInitialsScene,
            ScreenClear_sceneEvent: this.LoadPlayingScene,
        },
        EnterInitials_sceneId: {
            Advance_sceneEvent: this.LoadScoreBoardScene,
        },
        ScoreBoard_sceneId: {
            Advance_sceneEvent: this.LoadSplashScene,
        },
    }

    this.highScores = []scoreEntry{}
    for i := 0; i < 10; i++ {
        this.highScores = append(this.highScores, scoreEntry{initials:"XXX", score: 0, id: -1})
    }

    asset.InitSpriteAssets()
    asset.InitAudioAssets()

    this.curScene = NewScene(this.ecs)
    this.Transition(Init_sceneEvent)
    
    this.ecs.AddSystem(system.Velocity.Update)
    this.ecs.AddSystem(system.Wrap.Update)
    this.ecs.AddSystem(system.Collisions.Update)

    this.ecs.AddSystem(system.AnimateGraphicObjects.Update)
    this.ecs.AddSystem(system.Input.Update)
    this.ecs.AddSystem(system.TextInput.Update)
    this.ecs.AddSystem(system.Damage.Update)
    this.ecs.AddSystem(system.Health.Update)

    this.ecs.AddSystem(system.Action.Update)

    this.ecs.AddRenderer(layer.Foreground, system.DrawGraphicObjectsBG.Draw)
    this.ecs.AddRenderer(layer.Foreground, system.DrawGraphicObjectsFG.Draw)
    this.ecs.AddRenderer(layer.Foreground, system.DrawGraphicObjectsHudBG.Draw)
    this.ecs.AddRenderer(layer.Foreground, system.DrawGraphicObjectsHudFG.Draw)
    //this.ecs.AddRenderer(layer.Foreground, system.DrawColliders.Draw)
    return this
}

func (this *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
    return int(this.screenView.Area.Max.X), int(this.screenView.Area.Max.Y)
}

func (this *Game) Update() error {
    if ebiten.IsWindowBeingClosed() {
        this.Exit()
        return nil
    }
    events.ProcessAllEvents(this.ecs.World)
	this.ecs.Update()
	return nil
}

func (this *Game) Draw(screen *ebiten.Image) {
	screen.Clear()
	this.ecs.DrawLayer(layer.Background, screen)
	this.ecs.DrawLayer(layer.Foreground, screen)
	this.ecs.DrawLayer(layer.HudBackground, screen)
	this.ecs.DrawLayer(layer.HudForeground, screen)
}

func (this *Game) Exit() {
    os.Exit(0)
}
