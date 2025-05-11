package main

import (
	"math"
	"time"

	"github.com/dshaneg/asteroids/assets"
	"github.com/hajimehoshi/ebiten/v2"
)

var speed = float64(300 / ebiten.TPS())        // pixels per second
var rotSpeed = math.Pi / float64(ebiten.TPS()) // half way around per second (2 seconds for full rotation)

const (
	shotCooldown      = time.Millisecond * 500
	bulletSpawnOffset = 50.0
)

type Player struct {
	sprite   *ebiten.Image
	position Vector
	rotation float64

	shootCooldown *Timer
	bulletAdder   BulletAdder
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
	}
}

func (p *Player) Update() {
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		p.rotation -= rotSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		p.rotation += rotSpeed
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
}
