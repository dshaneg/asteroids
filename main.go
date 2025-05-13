package main

import (
	"fmt"
	"image/color"
	"os"
	"time"

	"github.com/dshaneg/asteroids/assets"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const (
	ScreenWidth  = 1024
	ScreenHeight = 768

	asteroidSpawnTime = 10 * time.Second
)

type Game struct {
	player             *Player
	asteroidSpawnTimer *Timer
	asteroids          []*Asteroid
	bullets            []*Bullet
	score              int
}

func (g *Game) AddBullet(b *Bullet) {
	g.bullets = append(g.bullets, b)
}

func (g *Game) Update() error {
	g.player.Update()

	// g.asteroidSpawnTimer.Update()
	// if g.asteroidSpawnTimer.IsReady() {
	// 	g.asteroidSpawnTimer.Reset()

	// 	g.asteroids = append(g.asteroids, NewAsteroid())
	// }

	for _, a := range g.asteroids {
		a.Update()
	}

	for _, b := range g.bullets {
		b.Update()
	}

	for i, a := range g.asteroids {
		for j, b := range g.bullets {
			if a.Collider().Intersects(b.Collider()) {
				g.score++
				// Remove the bullet and asteroid from their respective slices
				g.bullets = append(g.bullets[:j], g.bullets[j+1:]...)
				g.asteroids = append(g.asteroids[:i], g.asteroids[i+1:]...)
			}
		}
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

	score := fmt.Sprintf("%06d", g.score)
	face := &text.GoTextFace{Source: assets.ScoreFont, Size: 48}
	op := &text.DrawOptions{}
	op.LayoutOptions = text.LayoutOptions{
		PrimaryAlign:   text.AlignCenter,
		SecondaryAlign: text.AlignStart,
	}
	op.GeoM.Translate(ScreenWidth/2, 0)
	op.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, score, face, op)

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Asteroids: %v", len(g.asteroids)), 0, 0)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("TPS: %.2f", ebiten.ActualTPS()), 0, 20)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("FPS: %.2f", ebiten.ActualFPS()), 0, 40)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func main() {
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)

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
