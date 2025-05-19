package main

import (
	"math"
	"math/rand"

	"github.com/dshaneg/asteroids/assets"
	"github.com/dshaneg/asteroids/system"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
)

const (
	rotationSpeedMin = -0.02
	rotationSpeedMax = 0.02
)

var asteroidMaxHealth = map[asteroidClass]int{
	asteroidClassBig:    15,
	asteroidClassMedium: 7,
	asteroidClassSmall:  3,
	asteroidClassTiny:   1,
}

type asteroidClass int

const (
	asteroidClassBig asteroidClass = iota
	asteroidClassMedium
	asteroidClassSmall
	asteroidClassTiny
)

type asteroidColor int

const (
	asteroidColorBrown asteroidColor = iota
	asteroidColorGray
)

type Asteroid struct {
	class         asteroidClass
	color         asteroidColor
	position      Vector
	speed         Vector
	rotation      float64
	rotationSpeed float64
	sprite        *ebiten.Image
	health        int
	hot           bool
}

func NewAsteroid() *Asteroid {
	// where the asteroid is headed
	target := Vector{
		X: system.ScreenWidth / 2,
		Y: system.ScreenHeight / 2,
	}

	// the distance from the center of the screen to where the asteroid will spawn
	radius := system.ScreenWidth/2.0 + 150

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

	class := asteroidClassBig
	color := asteroidColor(rand.Intn(2))
	return &Asteroid{
		class:         class,
		color:         color,
		position:      pos,
		speed:         movement,
		rotationSpeed: rotationSpeedMin + rand.Float64()*(rotationSpeedMax-rotationSpeedMin),
		sprite:        getAsteroidSprite(class, color),
		health:        asteroidMaxHealth[class],
	}
}

func (a *Asteroid) newChildAsteroid() *Asteroid {
	angle := rand.Float64() * 2 * math.Pi
	velocity := 0.25 + rand.Float64()*1.5 // between 0.25 and 1.75

	movement := Vector{
		X: math.Cos(angle) * velocity,
		Y: math.Sin(angle) * velocity,
	}

	class := a.class + 1

	return &Asteroid{
		class:         class,
		color:         a.color,
		position:      a.position,
		speed:         movement,
		rotationSpeed: rotationSpeedMin + rand.Float64()*(rotationSpeedMax-rotationSpeedMin),
		sprite:        getAsteroidSprite(class, a.color),
		health:        asteroidMaxHealth[class],
	}
}

func getAsteroidSprite(class asteroidClass, color asteroidColor) *ebiten.Image {
	switch {
	case class == asteroidClassTiny && color == asteroidColorBrown:
		return assets.AsteroidTinyBrownSprites[rand.Intn(len(assets.AsteroidTinyBrownSprites))]
	case class == asteroidClassTiny && color == asteroidColorGray:
		return assets.AsteroidTinyGraySprites[rand.Intn(len(assets.AsteroidTinyGraySprites))]
	case class == asteroidClassSmall && color == asteroidColorBrown:
		return assets.AsteroidSmallBrownSprites[rand.Intn(len(assets.AsteroidSmallBrownSprites))]
	case class == asteroidClassSmall && color == asteroidColorGray:
		return assets.AsteroidSmallGraySprites[rand.Intn(len(assets.AsteroidSmallGraySprites))]
	case class == asteroidClassMedium && color == asteroidColorBrown:
		return assets.AsteroidMediumBrownSprites[rand.Intn(len(assets.AsteroidMediumBrownSprites))]
	case class == asteroidClassMedium && color == asteroidColorGray:
		return assets.AsteroidMediumGraySprites[rand.Intn(len(assets.AsteroidMediumGraySprites))]
	case class == asteroidClassBig && color == asteroidColorBrown:
		return assets.AsteroidBigBrownSprites[rand.Intn(len(assets.AsteroidBigBrownSprites))]
	default:
		return assets.AsteroidBigGraySprites[rand.Intn(len(assets.AsteroidBigGraySprites))]
	}
}

func (a *Asteroid) Split() []*Asteroid {
	// remove the check for small to get nuts quick
	if a.class == asteroidClassTiny || a.class == asteroidClassSmall {
		return nil
	}

	splitCount := rand.Intn(3) + 2
	children := make([]*Asteroid, splitCount)
	for i := 0; i < splitCount; i++ {
		children[i] = a.newChildAsteroid()
	}
	return children
}

func (a *Asteroid) Update() {
	a.hot = false

	b := a.sprite.Bounds()
	w := float64(b.Dx())
	h := float64(b.Dy())

	a.position.X += a.speed.X
	if a.position.X < -w {
		a.position.X = system.ScreenWidth
	} else if a.position.X > system.ScreenWidth {
		a.position.X = -w
	}

	a.position.Y += a.speed.Y
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

	if a.hot {
		op := &colorm.DrawImageOptions{}
		op.GeoM.Translate(-halfW, -halfH)
		op.GeoM.Rotate(a.rotation)
		op.GeoM.Translate(a.position.X+halfW, a.position.Y+halfH)
		cm := colorm.ColorM{}
		cm.Translate(1.0, 0.0, 0.0, 0.0)
		colorm.DrawImage(screen, a.sprite, cm, op)
	} else {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-halfW, -halfH)
		op.GeoM.Rotate(a.rotation)
		op.GeoM.Translate(a.position.X+halfW, a.position.Y+halfH)

		screen.DrawImage(a.sprite, op)
	}
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

func (a *Asteroid) Hit() {
	a.health--
	a.hot = true
}
