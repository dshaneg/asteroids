package main

import (
	"embed"
	"image"
	_ "image/png"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	// "github.com/hajimehoshi/ebiten/v2/colorm"
)

//go:embed assets
var assets embed.FS

var PlayerSprite = mustLoadImage("assets/PNG/Sprites/Ships/spaceShips_003.png")

func mustLoadImage(name string) *ebiten.Image {
	f, err := assets.Open(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

	return ebiten.NewImageFromImage(img)
}

type Vector struct {
	X float64
	Y float64
}

type Game struct {
	playerPosition Vector
	attackTimer    *Timer
}

func (g *Game) Update() error {
	g.attackTimer.Update()
	if g.attackTimer.IsReady() {
		g.attackTimer.Reset()

		// g..Debug("Attack!")
	}
	// speed := 5.0 // pixels per tick
	speed := float64(300 / ebiten.TPS()) // pixels per second
	// g.playerPosition.X += speed

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
	g.playerPosition.X += delta.X
	g.playerPosition.Y += delta.Y

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// width := PlayerSprite.Bounds().Dx()
	// height := PlayerSprite.Bounds().Dy()

	// halfW := float64(width / 2)
	// halfH := float64(height / 2)

	op := &ebiten.DrawImageOptions{}
	// op.GeoM.Translate(-halfW, -halfH)
	// op.GeoM.Rotate(45.0 * math.Pi / 180.0)
	// op.GeoM.Translate(halfW, halfH)
	op.GeoM.Translate(g.playerPosition.X, g.playerPosition.Y)
	screen.DrawImage(PlayerSprite, op)

	// op := &colorm.DrawImageOptions{}
	// cm := colorm.ColorM{}
	// cm.Scale(1.0, 1.0, 1.0, 0.5)

	// colorm.DrawImage(screen, PlayerSprite, cm, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	g := &Game{
		playerPosition: Vector{X: 100, Y: 100},
		attackTimer:    NewTimer(5 * time.Second),
	}

	err := ebiten.RunGame(g)
	if err != nil {
		panic(err)
	}
}
