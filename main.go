package main

import (
	"fmt"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/dshaneg/asteroids/system"
)

const (
	asteroidSpawnTime = 1 * time.Second
)

func main() {
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowSize(system.ScreenWidth, system.ScreenHeight)

	g := &Game{
		asteroidSpawnTimer: NewTimer(asteroidSpawnTime),
	}

	g.player = NewPlayer(g)

	err := ebiten.RunGame(g)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Game exited with error:", err)
		os.Exit(1)
	}
}
