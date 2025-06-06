package main

import (
	"math"

	"github.com/dshaneg/asteroids/assets"
	"github.com/dshaneg/asteroids/system"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	bulletSpeedPerSecond = 700.0
)

type Bullet struct {
	position Vector
	rotation float64
	sprite   *ebiten.Image
}

func NewBullet(pos Vector, rotation float64) *Bullet {
	sprite := assets.LaserSprite

	bounds := sprite.Bounds()
	halfW := float64(bounds.Dx()) / 2
	halfH := float64(bounds.Dy()) / 2

	pos.X -= halfW
	pos.Y -= halfH

	return &Bullet{
		position: pos,
		rotation: rotation,
		sprite:   sprite,
	}
}

func (b *Bullet) Update() {
	speed := bulletSpeedPerSecond / float64(ebiten.TPS())

	b.position.X += math.Sin(b.rotation) * speed
	b.position.Y += math.Cos(b.rotation) * -speed
}

func (b *Bullet) Draw(screen *ebiten.Image) {
	bounds := b.sprite.Bounds()
	halfW := float64(bounds.Dx()) / 2
	halfH := float64(bounds.Dy()) / 2

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-halfW, -halfH)
	op.GeoM.Rotate(b.rotation)
	op.GeoM.Translate(b.position.X+halfW, b.position.Y+halfH)

	screen.DrawImage(b.sprite, op)
}

func (b *Bullet) IsOffScreen() bool {
	bounds := b.sprite.Bounds()
	if b.position.X < -float64(bounds.Dx()) ||
		b.position.X > float64(system.ScreenWidth) ||
		b.position.Y < -float64(bounds.Dy()) ||
		b.position.Y > float64(system.ScreenHeight) {
		return true
	}
	return false
}

func (b *Bullet) Collider() Rect {
	bounds := b.sprite.Bounds()

	return NewRect(
		b.position.X,
		b.position.Y,
		float64(bounds.Dx()),
		float64(bounds.Dy()),
	)
}
