package main

import (
	"fmt"
	"os"
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
	bullets            []*Bullet
}

var cnt int = 0

func (g *Game) AddBullet(b *Bullet) {
	g.bullets = append(g.bullets, b)
}

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

	for _, b := range g.bullets {
		b.Update()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.player.Draw(screen)

	for _, a := range g.asteroids {
		a.Draw(screen)
	}

	for _, b := range g.bullets {
		b.Draw(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintln(os.Stdout, "Recovered from panic:", r)
			os.Stdout.Sync()
		}
	}()

	g := &Game{
		asteroidSpawnTimer: NewTimer(asteroidSpawnTime),
	}

	g.player = NewPlayer(g)

	err := ebiten.RunGame(g)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Game exited with error:", err)
		panic(err)
	}
}
