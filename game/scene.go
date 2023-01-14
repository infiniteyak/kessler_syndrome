package game

import (
	"fmt"
	"math"
	"math/rand"
	"github.com/infiniteyak/kessler_syndrome/asset"
	"github.com/infiniteyak/kessler_syndrome/entity"
	"github.com/infiniteyak/kessler_syndrome/event"
	"github.com/infiniteyak/kessler_syndrome/utility"
	"sort"
	"strings"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi"
    "net/http"
    "bytes"
    "io/ioutil"
    "encoding/json"
    "log"
)

type SceneEventId int
const (
    Undefined_sceneEvent SceneEventId = iota
    Init_sceneEvent
    Advance_sceneEvent
    GameOver_sceneEvent
    ScreenClear_sceneEvent
)

type SceneId int
const (
    Undefined_sceneId SceneId = iota
    Splash_sceneId
    Playing_sceneId
    EnterInitials_sceneId
    ScoreBoard_sceneId
)

type Scene struct {
    sceneId SceneId
    entities []*donburi.Entity
    cleanupFuncs []func() //These will be called before entities are removed
                          //Use them for unsubscribing from callbacks etc
}

func NewScene(ecs *ecs.ECS) *Scene {
    scene := &Scene{}
    scene.entities = make([]*donburi.Entity, 0)
    scene.cleanupFuncs = make([]func(), 0)

    // Event to handle adding entities
    registerEntityFunc := func(w donburi.World, event event.RegisterEntity) {
        scene.entities = append(scene.entities, event.Entity)
    }
    event.RegisterEntityEvent.Subscribe(ecs.World, registerEntityFunc)
    scene.cleanupFuncs = append(scene.cleanupFuncs, func() {
        event.RegisterEntityEvent.Unsubscribe(ecs.World, registerEntityFunc)
    })

    // Event to handle removing entities
    removeEntityFunc := func(w donburi.World, event event.RemoveEntity) {
        for i, e := range scene.entities {
            if e == event.Entity {
                scene.entities[i] = scene.entities[len(scene.entities)-1]
                scene.entities = scene.entities[:len(scene.entities)-1]
                w.Remove(*e)
                break
            }
        }
    }
    event.RemoveEntityEvent.Subscribe(ecs.World, removeEntityFunc)
    scene.cleanupFuncs = append(scene.cleanupFuncs, func() {
        event.RemoveEntityEvent.Unsubscribe(ecs.World, removeEntityFunc)
    })

    // Event to handle adding cleanup functions which are called when the scene ends
    registerCleanupFunc := func(w donburi.World, event event.RegisterCleanupFunc) {
        scene.cleanupFuncs = append(scene.cleanupFuncs, event.Function)
    }
    event.RegisterCleanupFuncEvent.Subscribe(ecs.World, registerCleanupFunc)
    scene.cleanupFuncs = append(scene.cleanupFuncs, func() {
        event.RegisterCleanupFuncEvent.Unsubscribe(ecs.World, registerCleanupFunc)
    })

    return scene
}

func (this *Game) LoadSplashScene() {
    println("LoadSplashScene")
    this.curScene.sceneId = Splash_sceneId 

    //TODO is there a better way to pass this data between scenes?
    this.curScore = 0
    this.curWave = 1
    this.curShips = 3
    this.playerInitials = ""

    // Add star field background
    entity.AddStarField(
        this.ecs, 
        float64(this.screenView.Area.Max.X / 2), 
        float64(this.screenView.Area.Max.Y / 2),
        this.screenView,
    )

    // Add spash screen text displaying title
    splashText := entity.StringData{
        String: "KESSLER SYNDROME",
        XAlign: entity.Center_fontalignx,
        YAlign: entity.Middle_fontaligny,
        Kerning: 2,
        Font: "BlackFont",
    }
    entity.AddSpriteText(
        this.ecs, 
        float64(this.screenView.Area.Max.X / 2), 
        float64(this.screenView.Area.Max.Y / 2),
        this.screenView,
        &splashText,
    )

    // Advance to the next state when you hit space
    entity.AddInputTrigger(
        this.ecs, 
        func() {
            this.Transition(Advance_sceneEvent)
        },
    )

    distance := 80.0 //avarage distance to spawn asteroids
    dVarience := 30.0 //variance in how far to spawn asteroids
    baseAsteroidCount := 5
    for i := 0; i < baseAsteroidCount; i++ { //spawns an additional asteroid every wave
        //TODO make a function for calculating this 
        min := distance - dVarience 
        max := distance + dVarience 
        radius := min + rand.Float64() * (max - min)
        angle := rand.Float64() * math.Pi * 2.0
        x := math.Cos(angle) * radius
        y := math.Sin(angle) * radius
        x += float64(this.screenView.Area.Max.X / 2)
        y += float64(this.screenView.Area.Max.Y / 2)
        entity.AddLargeAsteroid(
            this.ecs, 
            x,
            y,
            this.screenView,
            false,
        )
    }
    baseAsteroidCount--
    for i := 0; i < baseAsteroidCount; i++ { //spawns an additional asteroid every wave
        min := distance - dVarience 
        max := distance + dVarience 
        radius := min + rand.Float64() * (max - min)
        angle := rand.Float64() * math.Pi * 2.0
        x := math.Cos(angle) * radius
        y := math.Sin(angle) * radius
        x += float64(this.screenView.Area.Max.X / 2)
        y += float64(this.screenView.Area.Max.Y / 2)
        entity.AddMediumAsteroid(
            this.ecs, 
            x,
            y,
            this.screenView,
            false,
        )
    }
    baseAsteroidCount--
    for i := 0; i < baseAsteroidCount; i++ { //spawns an additional asteroid every wave
        min := distance - dVarience 
        max := distance + dVarience 
        radius := min + rand.Float64() * (max - min)
        angle := rand.Float64() * math.Pi * 2.0
        x := math.Cos(angle) * radius
        y := math.Sin(angle) * radius
        x += float64(this.screenView.Area.Max.X / 2)
        y += float64(this.screenView.Area.Max.Y / 2)
        entity.AddSmallAsteroid(
            this.ecs, 
            x,
            y,
            this.screenView,
            false,
        )
    }
}

func (this *Game) LoadPlayingScene() {
    println("LoadPlayingScene")
    this.curScene.sceneId = Playing_sceneId 

    // HUD
    hudView := utility.NewView(0.0, 0.0, this.screenView.Area.Max.X, asset.FontHeight)

    entity.AddBlackBar(
        this.ecs, 
        float64(hudView.Area.Max.X / 2),
        float64(hudView.Area.Max.Y / 2),
        hudView,
    )

    curScoreText := entity.StringData{
        String: fmt.Sprintf("%06d", this.curScore),
        XAlign: entity.Left_fontalignx,
        YAlign: entity.Top_fontaligny,
        Kerning: 0,
        Font: "WhiteFont",
    }
    scoreFunction := func(w donburi.World, event event.Score) {
        this.curScore += event.Value
        maxScore := 999999
        if this.curScore > maxScore {
            this.curScore = maxScore 
        }
        curScoreText.String = fmt.Sprintf("%06d", this.curScore)
    }
    event.ScoreEvent.Subscribe(this.ecs.World, scoreFunction)
    event.RegisterCleanupFuncEvent.Publish(
        this.ecs.World, 
        event.RegisterCleanupFunc{
            Function: func() {
                event.ScoreEvent.Unsubscribe(this.ecs.World, scoreFunction)
            },
        },
    )
    entity.AddSpriteText(
        this.ecs, 
        0,
        0,
        hudView,
        &curScoreText,
    )

    waveText := entity.StringData{
        String: fmt.Sprintf("%03d", this.curWave),
        XAlign: entity.Center_fontalignx,
        YAlign: entity.Top_fontaligny,
        Kerning: 0,
        Font: "WhiteFont",
    }
    entity.AddSpriteText(
        this.ecs, 
        float64(hudView.Area.Max.X / 2), 
        0,
        hudView,
        &waveText,
    )

    asteroidsCount := 0
    asteroidCountUpdateFunc := func(w donburi.World, event event.AsteroidsCountUpdate) {
        asteroidsCount += event.Value
        if asteroidsCount <= 0 {
            this.curWave++
            this.Transition(ScreenClear_sceneEvent)
        }
    }
    event.AsteroidsCountUpdateEvent.Subscribe(this.ecs.World, asteroidCountUpdateFunc)
    event.RegisterCleanupFuncEvent.Publish(
        this.ecs.World, 
        event.RegisterCleanupFunc{
            Function: func() {
                event.AsteroidsCountUpdateEvent.Unsubscribe(this.ecs.World, asteroidCountUpdateFunc)
            },
        },
    )

    //TODO Use some sprite to represent lives/ships
    shipsText := entity.StringData{
        String: strings.Repeat("^", this.curShips),
        XAlign: entity.Right_fontalignx,
        YAlign: entity.Top_fontaligny,
        Kerning: 0,
        Font: "WhiteFont",
    }
    entity.AddSpriteText(
        this.ecs, 
        float64(hudView.Area.Max.X), 
        0,
        hudView,
        &shipsText,
    )

    // GAME
    gameView := utility.NewView(
        0.0, 
        hudView.Area.Max.Y,
        this.screenView.Area.Max.X, 
        this.screenView.Area.Max.Y - hudView.Area.Max.Y,
    )

    shipDestFunc := func(w donburi.World, event event.ShipDestroyed) {
        this.curShips--
        if this.curShips < 0 {
            println("game over")
            gameOverText := entity.StringData{
                String: "GAME OVER",
                XAlign: entity.Center_fontalignx,
                YAlign: entity.Middle_fontaligny,
                Kerning: 2,
                Font: "BlackFont",
            }
            entity.AddSpriteText(
                this.ecs, 
                float64(gameView.Area.Max.X / 2), 
                float64(gameView.Area.Max.Y / 2), 
                gameView,
                &gameOverText,
            )
            entity.AddInputTrigger(
                this.ecs, 
                func() {
                    this.Transition(GameOver_sceneEvent)
                },
            )
        } else {
            shipsText.String = strings.Repeat("^", this.curShips) 
            entity.AddPlayerShip(
                this.ecs, 
                float64(gameView.Area.Max.X / 2), 
                float64(gameView.Area.Max.Y / 2), 
                gameView,
                this.audioContext,
            )
        }
    }
    event.ShipDestroyedEvent.Subscribe(this.ecs.World, shipDestFunc)
    event.RegisterCleanupFuncEvent.Publish(
        this.ecs.World, 
        event.RegisterCleanupFunc{
            Function: func() {
                event.ShipDestroyedEvent.Unsubscribe(this.ecs.World, shipDestFunc)
            },
        },
    )

    entity.AddStarField(
        this.ecs, 
        float64(gameView.Area.Max.X / 2), 
        float64(gameView.Area.Max.Y / 2),
        gameView,
    )

    entity.AddPlayerShip(
        this.ecs, 
        float64(gameView.Area.Max.X / 2), 
        float64(gameView.Area.Max.Y / 2), 
        gameView,
        this.audioContext,
    )

    distance := 80.0 //avarage distance to spawn asteroids
    dVarience := 30.0 //variance in how far to spawn asteroids
    baseAsteroidCount := 1
    for i := 0; i < this.curWave + baseAsteroidCount; i++ { //spawns an additional asteroid every wave
        //TODO make a function for calculating this 
        min := distance - dVarience 
        max := distance + dVarience 
        radius := min + rand.Float64() * (max - min)
        angle := rand.Float64() * math.Pi * 2.0
        x := math.Cos(angle) * radius
        y := math.Sin(angle) * radius
        x += float64(gameView.Area.Max.X / 2)
        y += float64(gameView.Area.Max.Y / 2)
        entity.AddLargeAsteroid(
            this.ecs, 
            x,
            y,
            gameView,
            true,
        )
    }
    baseAsteroidCount--
    for i := 0; i < this.curWave + baseAsteroidCount; i++ { //spawns an additional asteroid every wave
        min := distance - dVarience 
        max := distance + dVarience 
        radius := min + rand.Float64() * (max - min)
        angle := rand.Float64() * math.Pi * 2.0
        x := math.Cos(angle) * radius
        y := math.Sin(angle) * radius
        x += float64(gameView.Area.Max.X / 2)
        y += float64(gameView.Area.Max.Y / 2)
        entity.AddMediumAsteroid(
            this.ecs, 
            x,
            y,
            gameView,
            true,
        )
    }
    baseAsteroidCount--
    for i := 0; i < this.curWave + baseAsteroidCount; i++ { //spawns an additional asteroid every wave
        min := distance - dVarience 
        max := distance + dVarience 
        radius := min + rand.Float64() * (max - min)
        angle := rand.Float64() * math.Pi * 2.0
        x := math.Cos(angle) * radius
        y := math.Sin(angle) * radius
        x += float64(gameView.Area.Max.X / 2)
        y += float64(gameView.Area.Max.Y / 2)
        entity.AddSmallAsteroid(
            this.ecs, 
            x,
            y,
            gameView,
            true,
        )
    }

    waveDcopy := *asset.WaveD
    wavePlayer, err := this.audioContext.NewPlayer(&waveDcopy)
    if err != nil {
        log.Fatal(err)
    }

    wavePlayer.Rewind()
    wavePlayer.Play()
}

func (this *Game) LoadInitialsScene() {
    println("LoadInitialsScene")
    this.curScene.sceneId = EnterInitials_sceneId 

    menuDcopy := *asset.MenuD
    menuPlayer, err := this.audioContext.NewPlayer(&menuDcopy)
    if err != nil {
        log.Fatal(err)
    }

    menuPlayer.Rewind()
    menuPlayer.Play()

    scoreText := entity.StringData{
        String: "SCORE " + fmt.Sprintf("%06d", this.curScore),
        XAlign: entity.Center_fontalignx,
        YAlign: entity.Middle_fontaligny,
        Kerning: 0,
        Font: "WhiteFont",
    }
    entity.AddSpriteText(
        this.ecs, 
        float64(this.screenView.Area.Max.X / 2), 
        float64(this.screenView.Area.Max.Y / 2 - 36),
        this.screenView,
        &scoreText,
    )

    explanationText := entity.StringData{
        String: "ENTER YOUR INITIALS",
        XAlign: entity.Center_fontalignx,
        YAlign: entity.Middle_fontaligny,
        Kerning: 0,
        Font: "WhiteFont",
    }
    entity.AddSpriteText(
        this.ecs, 
        float64(this.screenView.Area.Max.X / 2), 
        float64(this.screenView.Area.Max.Y / 2 - 22),
        this.screenView,
        &explanationText,
    )

    initialsText := entity.StringData{
        String: "___",
        XAlign: entity.Center_fontalignx,
        YAlign: entity.Middle_fontaligny,
        Kerning: 0,
        Font: "WhiteFont",
    }
    entity.AddSpriteText(
        this.ecs, 
        float64(this.screenView.Area.Max.X / 2), 
        float64(this.screenView.Area.Max.Y / 2),
        this.screenView,
        &initialsText,
    )
    entity.AddTextInput(
        this.ecs, 
        &initialsText.String,
        3,
        func(){
            this.playerInitials = initialsText.String
        },
    )

    confirmText := entity.StringData{
        String: "FIRE TO CONFIRM",
        XAlign: entity.Center_fontalignx,
        YAlign: entity.Middle_fontaligny,
        Kerning: 0,
        Font: "WhiteFont",
        Blink: true,
    }
    entity.AddSpriteText(
        this.ecs, 
        float64(this.screenView.Area.Max.X / 2), 
        float64(this.screenView.Area.Max.Y / 2 + 22),
        this.screenView,
        &confirmText,
    )

    entity.AddInputTrigger(
        this.ecs, 
        func() {
            // If we've fully filled in the initals advance
            done := true
            for _, c := range initialsText.String {
                if string(c) == "_" {
                    done = false
                }
            }
            if done {
                this.Transition(Advance_sceneEvent)
            }
        },
    )
}

type scoreAPIResponse struct {
    Items []scoreEntryAPI `json:"items"`
}

type scoreEntryAPI struct {
    Initials string `json:"initials"`
    Score int `json:"score"`
}

func getScores(pInitials string, pScore int) ([]scoreEntryAPI, bool) {
    jsonData, err := json.Marshal(scoreEntryAPI{Initials: pInitials, Score: pScore})
    if err != nil {
        log.Fatal(err)
    }

    response, err := http.Post(
        "https://www.infiniteyak.com/api/collections/kessler_syndrome_scores/records",
        "application/json",
        bytes.NewBuffer(jsonData),
        )
    if err != nil {
        print(err)
        return nil, true
    }

    response, err = http.Get("https://www.infiniteyak.com/api/collections/kessler_syndrome_scores/records?sort=-score&perPage=10")
    if err != nil {
        print(err)
        return nil, true
    }
    responseData, err := ioutil.ReadAll(response.Body)
    if err != nil {
        print(err)
        return nil, true
    }

    var responseObject scoreAPIResponse
    if err := json.Unmarshal(responseData, &responseObject); err != nil {
        panic(err)
    }

    return responseObject.Items, false
}

func (this *Game) LoadScoreBoardScene() {
    println("LoadScoreBoardScene")
    this.curScene.sceneId = ScoreBoard_sceneId 

    titleText := entity.StringData{
        String: "HIGH SCORES",
        XAlign: entity.Center_fontalignx,
        YAlign: entity.Top_fontaligny,
        Kerning: 2,
        Font: "BlackFont",
    }
    entity.AddSpriteText(
        this.ecs, 
        float64(this.screenView.Area.Max.X / 2), 
        10,
        this.screenView,
        &titleText,
    )

    scores, useLocal := getScores(this.playerInitials, this.curScore)

    if useLocal {
        this.highScores = append(this.highScores, scoreEntry{
            initials: this.playerInitials, 
            score: this.curScore,
            id: this.curScoreId,
        })
        sort.Slice(this.highScores, func(i, j int) bool {
            return this.highScores[i].score > this.highScores[j].score
        })
    } else {
        this.highScores = []scoreEntry{}
        for _, s := range scores {
            this.highScores = append(this.highScores, scoreEntry{
                initials:s.Initials, 
                score: s.Score, 
                id: -1,
            })
        }
    }
    scoreCount := len(this.highScores)
    if len(this.highScores) > 10 {
        scoreCount = 10
    }
    blinked := false
    for i := 0; i < scoreCount; i++ {
        match := this.highScores[i].score == this.curScore && 
                 this.highScores[i].initials == this.playerInitials 
        scoreText := entity.StringData{
            String: fmt.Sprintf("%02d", i+1) + 
                    ". " + 
                    this.highScores[i].initials + 
                    " " + 
                    fmt.Sprintf("%06d", this.highScores[i].score),
            XAlign: entity.Center_fontalignx,
            YAlign: entity.Top_fontaligny,
            Kerning: 0,
            Font: "WhiteFont",
            Blink: match && !blinked,
        }
        blinked = blinked || match
        entity.AddSpriteText(
            this.ecs, 
            float64(this.screenView.Area.Max.X / 2), 
            float64(40 + i * 15),
            this.screenView,
            &scoreText,
        )
    }
    this.curScoreId++

    // Advance to the next state when you hit space
    entity.AddInputTrigger(
        this.ecs, 
        func() {
            this.Transition(Advance_sceneEvent)
        },
    )
}

func (this *Game) Transition(event SceneEventId) {
    if this.states[this.curScene.sceneId][event] != nil {
        for _, foo := range this.curScene.cleanupFuncs {
            foo()
        }
        for _, e := range this.curScene.entities {
            this.ecs.World.Remove(*e)
        }
        sid := this.curScene.sceneId
        this.curScene = NewScene(this.ecs)
        this.states[sid][event]()
    } else {
        println("states map miss")
    }
}

