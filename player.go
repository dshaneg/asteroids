package main

import (
	"math"
	"time"

	"github.com/dshaneg/asteroids/assets"
	"github.com/dshaneg/asteroids/system"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	shotCooldown      time.Duration = time.Millisecond * 200
	bulletSpawnOffset float64       = 30.0
	accelPerSec       float64       = 10
	speedLimitPerSec  float64       = 400
)

type Player struct {
	sprite   *ebiten.Image
	position Vector
	rotation float64
	speed    Vector
	plume    *Plume

	shootCooldown *Timer
	bulletAdder   BulletAdder
}

type Plume struct {
	sprite  *ebiten.Image
	visible bool
}

type BulletAdder interface {
	AddBullet(b *Bullet)
}

func NewPlayer(bulletAdder BulletAdder) *Player {
	sprite := assets.PlayerSprite

	bounds := sprite.Bounds()
	halfW := float64(bounds.Dx() / 2)
	halfH := float64(bounds.Dy() / 2)

	pos := Vector{
		X: (float64(system.ScreenWidth) / 2) - halfW,
		Y: (float64(system.ScreenHeight) / 2) - halfH,
	}

	return &Player{
		sprite:        sprite,
		position:      pos,
		shootCooldown: NewTimer(shotCooldown),
		bulletAdder:   bulletAdder,
		plume: &Plume{
			sprite: assets.PlumeSprite,
		},
	}
}

func (p *Player) Update() {
	// point the ship at the mouse cursor
	mouseX, mouseY := ebiten.CursorPosition()
	dx := float64(mouseX) - (p.position.X + float64(p.sprite.Bounds().Dx())/2)
	dy := float64(mouseY) - (p.position.Y + float64(p.sprite.Bounds().Dy())/2)
	p.rotation = math.Atan2(dy, dx) + math.Pi/2

	// set the speed of the ship
	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		accel := accelPerSec / float64(ebiten.TPS())
		speedLimit := speedLimitPerSec / float64(ebiten.TPS())

		p.speed.X += math.Sin(p.rotation) * accel
		if p.speed.X > speedLimit {
			p.speed.X = speedLimit
		} else if p.speed.X < -speedLimit {
			p.speed.X = -speedLimit
		}
		p.speed.Y += math.Cos(p.rotation) * -accel
		if p.speed.Y > speedLimit {
			p.speed.Y = speedLimit
		} else if p.speed.Y < -speedLimit {
			p.speed.Y = -speedLimit
		}
		p.plume.visible = true
	} else {
		// auto decelerate
		p.speed.X *= 0.95
		p.speed.Y *= 0.95
		p.plume.visible = false
	}

	// set the position of the ship
	p.position.X += p.speed.X
	p.position.Y += p.speed.Y
	if p.position.X < 0 {
		p.position.X = 0
		p.speed.X = 0
	} else if p.position.X > float64(system.ScreenWidth)-float64(p.sprite.Bounds().Dx()) {
		p.position.X = float64(system.ScreenWidth) - float64(p.sprite.Bounds().Dx())
		p.speed.X = 0
	}
	if p.position.Y < 0 {
		p.position.Y = 0
		p.speed.Y = 0
	} else if p.position.Y > float64(system.ScreenHeight)-float64(p.sprite.Bounds().Dy()) {
		p.position.Y = float64(system.ScreenHeight) - float64(p.sprite.Bounds().Dy())
		p.speed.Y = 0
	}

	p.shootCooldown.Update()
	if p.shootCooldown.IsReady() && (ebiten.IsKeyPressed(ebiten.KeySpace) || ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)) {
		p.shootCooldown.Reset()

		// bullet should come from the middle of the ship
		bounds := p.sprite.Bounds()
		halfW := float64(bounds.Dx() / 2)
		halfH := float64(bounds.Dy() / 2)

		spawnPos := Vector{
			X: p.position.X + halfW + math.Sin(p.rotation)*bulletSpawnOffset,
			Y: p.position.Y + halfH + math.Cos(p.rotation)*-bulletSpawnOffset,
		}

		bullet := NewBullet(spawnPos, p.rotation)

		p.bulletAdder.AddBullet(bullet)
	}
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

	if p.plume.visible {
		plumeBounds := p.plume.sprite.Bounds()
		halfPlumeW := float64(plumeBounds.Dx() / 2)

		plumeOp := &ebiten.DrawImageOptions{}
		plumeOp.GeoM.Translate(-halfPlumeW, halfH)
		plumeOp.GeoM.Rotate(p.rotation)
		plumeOp.GeoM.Translate(p.position.X+halfW, p.position.Y+halfH)
		screen.DrawImage(p.plume.sprite, plumeOp)
	}
}
