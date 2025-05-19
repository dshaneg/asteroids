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
	booms              []*Boom
	score              int
}

func (g *Game) AddBullet(b *Bullet) {
	g.bullets = append(g.bullets, b)
}

func (g *Game) AddAsteroid(a *Asteroid) {
	g.asteroids = append(g.asteroids, a)
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

	for _, boom := range g.booms {
		boom.Update()
	}

	for i := len(g.asteroids) - 1; i >= 0; i-- {
		a := g.asteroids[i]
		for j := len(g.bullets) - 1; j >= 0; j-- {
			b := g.bullets[j]
			if a.Collider().Intersects(b.Collider()) {
				g.score++
				g.discardBullet(j)
				a.Hit()
				if a.health <= 0 {
					g.booms = append(g.booms, NewBoom(a))
					g.discardAsteroid(i)
					g.asteroids = append(g.asteroids, a.Split()...)
				}
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

func (g *Game) discardAsteroid(index int) {
	if index < 0 || index >= len(g.asteroids) {
		return
	}
	g.asteroids = append(g.asteroids[:index], g.asteroids[index+1:]...)
}

func (g *Game) discardBoom(index int) {
	if index < 0 || index >= len(g.booms) {
		return
	}
	g.booms = append(g.booms[:index], g.booms[index+1:]...)
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, b := range g.bullets {
		b.Draw(screen)
	}

	g.player.Draw(screen)

	for _, a := range g.asteroids {
		a.Draw(screen)
	}

	for i, boom := range g.booms {
		if boom.visibleTimer.IsReady() {
			g.discardBoom(i)
			continue
		}
		boom.Draw(screen)
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
