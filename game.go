package main

import (
	"fmt"
	"image/color"

	"github.com/dshaneg/asteroids/assets"
	"github.com/dshaneg/asteroids/system"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
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

	g.asteroidSpawnTimer.Update()
	if g.asteroidSpawnTimer.IsReady() {
		g.asteroidSpawnTimer.Reset()

		g.asteroids = append(g.asteroids, NewAsteroid())
	}

	for _, a := range g.asteroids {
		a.Update()
	}

	for i, b := range g.bullets {
		if b.IsOffScreen() {
			g.discardBullet(i)
			continue
		}
		b.Update()
	}

	for i, a := range g.asteroids {
		for j, b := range g.bullets {
			if a.Collider().Intersects(b.Collider()) {
				g.score++
				// Remove the bullet and asteroid from their respective slices
				g.discardBullet(j)
				g.asteroids = append(g.asteroids[:i], g.asteroids[i+1:]...)
			}
		}
	}

	return nil
}

func (g *Game) discardBullet(index int) {
	if index < 0 || index >= len(g.bullets) {
		return
	}
	g.bullets = append(g.bullets[:index], g.bullets[index+1:]...)
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, b := range g.bullets {
		b.Draw(screen)
	}

	g.player.Draw(screen)

	for _, a := range g.asteroids {
		a.Draw(screen)
	}

	score := fmt.Sprintf("%06d", g.score)
	face := &text.GoTextFace{Source: assets.ScoreFont, Size: 48}
	op := &text.DrawOptions{}
	op.LayoutOptions = text.LayoutOptions{
		PrimaryAlign:   text.AlignCenter,
		SecondaryAlign: text.AlignStart,
	}
	op.GeoM.Translate(system.ScreenWidth/2, 0)
	op.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, score, face, op)

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Asteroids: %v", len(g.asteroids)), 0, 0)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Bullets: %v", len(g.bullets)), 0, 20)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("TPS: %.2f", ebiten.ActualTPS()), 0, 40)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("FPS: %.2f", ebiten.ActualFPS()), 0, 60)
	// ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Loc: %.2f", g.player.position), 0, 80)
	// ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Rot: %.2f", g.player.rotation), 0, 100)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return system.ScreenWidth, system.ScreenHeight
}
