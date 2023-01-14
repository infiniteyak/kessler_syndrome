package main

import (
    "github.com/infiniteyak/kessler_syndrome/game"
	"github.com/hajimehoshi/ebiten/v2"
    "log"
	"os"
	"os/signal"
	"syscall"
    "math/rand"
	"time"
)

const (
    ScreenWidth = 288
    ScreenHeight = 224
    Title = "Kessler Syndrome"
)

func main() {
    rand.Seed(time.Now().UnixNano())

    ebiten.SetWindowTitle(Title)
    ebiten.SetWindowResizable(true)
    ebiten.SetMaxTPS(120)
    ebiten.SetWindowClosingHandled(true)

    g := game.NewGame(ScreenWidth, ScreenHeight)

    sigc := make(chan os.Signal, 1)
    signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
    go func() {
        <-sigc
        g.Exit()
    }()

    if err := ebiten.RunGame(g); err != nil {
        log.Fatal(err)
    }
}
