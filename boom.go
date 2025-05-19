package main

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/dshaneg/asteroids/assets"
)

const (
	visibleMs = 150
)

type Boom struct {
	sprite       *ebiten.Image
	scale        float64
	targetSprite *ebiten.Image
	targetCenter Vector
	position     Vector
	speed        Vector
	visibleTimer *Timer
}

func NewBoom(a *Asteroid) *Boom {
	var scale float64

	switch a.class {
	case asteroidClassTiny:
		scale = 0.20
	case asteroidClassSmall:
		scale = 0.35
	case asteroidClassMedium:
		scale = 0.50
	default:
		scale = 1.0
	}

	// targetBounds := a.sprite.Bounds()
	// targetCenter := Vector{
	// 	X: float64(targetBounds.Dx()) / 2,
	// 	Y: float64(targetBounds.Dy()) / 2,
	// }
	// spriteCenter := Vector{
	// 	X: float64(a.sprite.Bounds().Dx()) / 2,
	// 	Y: float64(a.sprite.Bounds().Dy()) / 2,
	// }
	// dx := targetCenter.X - spriteCenter.X
	// dy := targetCenter.Y - spriteCenter.Y

	return &Boom{
		sprite: assets.BoomSprite,
		scale:  scale,
		position: Vector{
			X: a.position.X,
			Y: a.position.Y,
		},
		targetSprite: a.sprite,
		// targetCenter: targetCenter,
		speed:        a.speed,
		visibleTimer: NewTimer(time.Millisecond * visibleMs),
	}
}

func (b *Boom) Update() {
	b.visibleTimer.Update()
	b.position.X += b.speed.X
	b.position.Y += b.speed.Y
}

func (b *Boom) Draw(screen *ebiten.Image) {
	targetBounds := b.targetSprite.Bounds()
	targetCenter := Vector{
		X: float64(targetBounds.Dx()) / 2,
		Y: float64(targetBounds.Dy()) / 2,
	}

	spriteBounds := b.sprite.Bounds()
	spriteCenter := Vector{
		X: float64(spriteBounds.Dx()) / 2,
		Y: float64(spriteBounds.Dy()) / 2,
	}
	dx := targetCenter.X - spriteCenter.X*b.scale
	dy := targetCenter.Y - spriteCenter.Y*b.scale

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(b.scale, b.scale)
	op.GeoM.Translate(b.position.X+dx, b.position.Y+dy)
	screen.DrawImage(b.sprite, op)
}
