package main

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	// "github.com/hajimehoshi/ebiten/v2/colorm"
)

const (
	ScreenWidth  = 800
	ScreenHeight = 600

	asteroidSpawnTime = 1 * time.Second
)

type Game struct {
	player             *Player
	asteroidSpawnTimer *Timer
	asteroids          []*Asteroid
}

var cnt int = 0

func (g *Game) Update() error {
	g.player.Update()

	g.asteroidSpawnTimer.Update()
	if g.asteroidSpawnTimer.IsReady() {
		g.asteroidSpawnTimer.Reset()

		g.asteroids = append(g.asteroids, NewAsteroid())
	}

	for _, a := range g.asteroids {
		a.Update()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.player.Draw(screen)

	for _, a := range g.asteroids {
		a.Draw(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func main() {
	g := &Game{
		player:             NewPlayer(),
		asteroidSpawnTimer: NewTimer(asteroidSpawnTime),
	}

	err := ebiten.RunGame(g)
	if err != nil {
		panic(err)
	}
}
