package main

import (
	"math"
	"math/rand"

	"github.com/dshaneg/asteroids/assets"
	"github.com/dshaneg/asteroids/system"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	rotationSpeedMin = -0.02
	rotationSpeedMax = 0.02
)

type Asteroid struct {
	position      Vector
	movement      Vector
	rotation      float64
	rotationSpeed float64
	sprite        *ebiten.Image
}

func NewAsteroid() *Asteroid {
	// where the asteroid is headed
	target := Vector{
		X: system.ScreenWidth / 2,
		Y: system.ScreenHeight / 2,
	}

	// the distance from the center of the screen to where the asteroid will spawn
	radius := system.ScreenWidth / 2.0

	// pick a random angle - 2Pi is 360 degrees - so this returns 0 to 360
	angle := rand.Float64() * 2 * math.Pi

	// figure out the spawn position by moving radius pixels from the target at the chosen angle
	pos := Vector{
		X: target.X + radius*math.Cos(angle),
		Y: target.Y + radius*math.Sin(angle),
	}

	// randomized velocity
	velocity := 0.25 + rand.Float64()*1.5 // between 0.25 and 1.75

	// direction is the target minus the current position
	direction := Vector{
		X: target.X - pos.X,
		Y: target.Y - pos.Y,
	}

	normalizedDirection := direction.Normalize()

	movement := Vector{
		X: normalizedDirection.X * velocity,
		Y: normalizedDirection.Y * velocity,
	}

	sprite := assets.AsteroidSprites[rand.Intn(len(assets.AsteroidSprites))]

	return &Asteroid{
		position:      pos,
		movement:      movement,
		rotationSpeed: rotationSpeedMin + rand.Float64()*(rotationSpeedMax-rotationSpeedMin),
		sprite:        sprite,
	}
}

func (a *Asteroid) Update() {
	b := a.sprite.Bounds()
	w := float64(b.Dx())
	h := float64(b.Dy())

	a.position.X += a.movement.X
	if a.position.X < -w {
		a.position.X = system.ScreenWidth
	} else if a.position.X > system.ScreenWidth {
		a.position.X = -w
	}

	a.position.Y += a.movement.Y
	if a.position.Y < -h {
		a.position.Y = system.ScreenHeight
	} else if a.position.Y > system.ScreenHeight {
		a.position.Y = -h
	}

	a.rotation += a.rotationSpeed
}

func (a *Asteroid) Draw(screen *ebiten.Image) {
	bounds := a.sprite.Bounds()
	halfW := float64(bounds.Dx() / 2)
	halfH := float64(bounds.Dy() / 2)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-halfW, -halfH)
	op.GeoM.Rotate(a.rotation)
	op.GeoM.Translate(a.position.X+halfW, a.position.Y+halfH)

	screen.DrawImage(a.sprite, op)
}

func (a *Asteroid) Collider() Rect {
	bounds := a.sprite.Bounds()

	return NewRect(
		a.position.X,
		a.position.Y,
		float64(bounds.Dx()),
		float64(bounds.Dy()),
	)
}
