package main

import (
	"math"

	"github.com/dshaneg/asteroids/assets"
	"github.com/hajimehoshi/ebiten/v2"
)

var speed = float64(300 / ebiten.TPS())        // pixels per second
var rotSpeed = math.Pi / float64(ebiten.TPS()) // half way around per second (2 seconds for full rotation)

type Player struct {
	sprite   *ebiten.Image
	position Vector
	rotation float64
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
		sprite:   sprite,
		position: pos,
	}
}

func (p *Player) Update() {
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		p.rotation -= rotSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		p.rotation += rotSpeed
	}
	// var delta Vector
	// if ebiten.IsKeyPressed(ebiten.KeyDown) {
	// 	delta.Y += speed
	// }
	// if ebiten.IsKeyPressed(ebiten.KeyUp) {
	// 	delta.Y -= speed
	// }
	// if ebiten.IsKeyPressed(ebiten.KeyLeft) {
	// 	delta.X -= speed
	// }
	// if ebiten.IsKeyPressed(ebiten.KeyRight) {
	// 	delta.X += speed
	// }

	// if delta.X != 0 && delta.Y != 0 {
	// 	factor := speed / math.Sqrt(delta.X*delta.X+delta.Y*delta.Y)
	// 	delta.X *= factor
	// 	delta.Y *= factor
	// }

	// p.position.X += delta.X
	// p.position.Y += delta.Y
}

func (p *Player) Draw(screen *ebiten.Image) {
	bounds := p.sprite.Bounds()
	halfW := float64(bounds.Dx() / 2)
	halfH := float64(bounds.Dy() / 2)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-halfW, -halfH)
	op.GeoM.Rotate(p.rotation)
	op.GeoM.Translate(halfW, halfH)

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
