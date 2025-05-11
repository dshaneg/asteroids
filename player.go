package main

import (
	"math"

	"github.com/dshaneg/asteroids/assets"
	"github.com/hajimehoshi/ebiten/v2"
)

var speed = float64(300 / ebiten.TPS()) // pixels per second

type Player struct {
	position Vector
	sprite   *ebiten.Image
}

func NewPlayer() *Player {
	sprite := assets.PlayerSprite

	bounds := sprite.Bounds()
	halfW := float64(bounds.Dx() / 2)
	halfH := float64(bounds.Dy() / 2)

	pos := Vector{
		X: (float64(ScreenWidth) / 2) - halfW,
		Y: (float64(ScreenHeight) / 2) - halfH,
	}

	return &Player{
		position: pos,
		sprite:   sprite,
	}
}

func (p *Player) Update() {
	var delta Vector
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		delta.Y += speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		delta.Y -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		delta.X -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		delta.X += speed
	}

	if delta.X != 0 && delta.Y != 0 {
		factor := speed / math.Sqrt(delta.X*delta.X+delta.Y*delta.Y)
		delta.X *= factor
		delta.Y *= factor
	}

	p.position.X += delta.X
	p.position.Y += delta.Y
}

func (p *Player) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(p.position.X, p.position.Y)
	screen.DrawImage(p.sprite, op)

	// width := PlayerSprite.Bounds().Dx()
	// height := PlayerSprite.Bounds().Dy()

	// halfW := float64(width / 2)
	// halfH := float64(height / 2)

	// op.GeoM.Translate(-halfW, -halfH)
	// op.GeoM.Rotate(45.0 * math.Pi / 180.0)
	// op.GeoM.Translate(halfW, halfH)

	// op := &colorm.DrawImageOptions{}
	// cm := colorm.ColorM{}
	// cm.Scale(1.0, 1.0, 1.0, 0.5)
}
