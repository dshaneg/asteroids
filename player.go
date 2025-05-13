package main

import (
	"math"
	"time"

	"github.com/dshaneg/asteroids/assets"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	shotCooldown      time.Duration = time.Millisecond * 75
	bulletSpawnOffset float64       = 50.0
	accelPerSec       float64       = 10
	speedLimitPerSec  float64       = 400
)

type Player struct {
	sprite   *ebiten.Image
	position Vector
	rotation float64
	speedX   float64
	speedY   float64
	plume    *Plume

	shootCooldown *Timer
	bulletAdder   BulletAdder
}

type Plume struct {
	sprite   *ebiten.Image
	position Vector
	rotation float64
	visible  bool
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
		X: (float64(ScreenWidth) / 2) - halfW,
		Y: (float64(ScreenHeight) / 2) - halfH,
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
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		accel := accelPerSec / float64(ebiten.TPS())
		speedLimit := speedLimitPerSec / float64(ebiten.TPS())

		p.speedX += math.Sin(p.rotation) * accel
		if p.speedX > speedLimit {
			p.speedX = speedLimit
		} else if p.speedX < -speedLimit {
			p.speedX = -speedLimit
		}
		p.speedY += math.Cos(p.rotation) * -accel
		if p.speedY > speedLimit {
			p.speedY = speedLimit
		} else if p.speedY < -speedLimit {
			p.speedY = -speedLimit
		}
		p.plume.visible = true
	} else {
		// auto decelerate
		p.speedX *= 0.95
		p.speedY *= 0.95
		p.plume.visible = false
	}

	// set the position of the ship
	p.position.X += p.speedX
	p.position.Y += p.speedY
	if p.position.X < 0 {
		p.position.X = 0
		p.speedX = 0
	} else if p.position.X > float64(ScreenWidth)-float64(p.sprite.Bounds().Dx()) {
		p.position.X = float64(ScreenWidth) - float64(p.sprite.Bounds().Dx())
		p.speedX = 0
	}
	if p.position.Y < 0 {
		p.position.Y = 0
		p.speedY = 0
	} else if p.position.Y > float64(ScreenHeight)-float64(p.sprite.Bounds().Dy()) {
		p.position.Y = float64(ScreenHeight) - float64(p.sprite.Bounds().Dy())
		p.speedY = 0
	}

	if p.plume.visible {
		p.plume.rotation = p.rotation
		bounds := p.sprite.Bounds()
		halfW := float64(bounds.Dx() / 2)

		// offset the plume position based on the ship's rotation
		p.plume.position.X = p.position.X + halfW - float64(p.plume.sprite.Bounds().Dx()/2)
		p.plume.position.Y = p.position.Y + float64(p.sprite.Bounds().Dy())
	}

	p.shootCooldown.Update()
	if p.shootCooldown.IsReady() && ebiten.IsKeyPressed(ebiten.KeySpace) {
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
		halfPlumeH := float64(plumeBounds.Dy() / 2)
		plumeOp := &ebiten.DrawImageOptions{}
		plumeOp.GeoM.Translate(-halfPlumeW, -halfPlumeH)
		plumeOp.GeoM.Rotate(p.rotation)
		plumeOp.GeoM.Translate(halfPlumeW, halfPlumeH)
		dX := halfW*math.Cos(p.rotation) + halfPlumeW*math.Cos(p.plume.rotation)
		dY := halfH*math.Sin(p.rotation) + halfPlumeH*math.Sin(p.plume.rotation)
		p.plume.position.X = p.position.X + dX
		p.plume.position.Y = p.position.Y + dY
		plumeOp.GeoM.Translate(p.plume.position.X, p.plume.position.Y)
		screen.DrawImage(p.plume.sprite, plumeOp)
	}
}
