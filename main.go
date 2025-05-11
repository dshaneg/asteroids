package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	// "github.com/hajimehoshi/ebiten/v2/colorm"
)

const (
	ScreenWidth  = 800
	ScreenHeight = 600
)

type Vector struct {
	X float64
	Y float64
}

type Game struct {
	player *Player
	// attackTimer       *Timer
	// attackActiveTimer *Timer
	// attackMessage     string
}

var cnt int = 0

func (g *Game) Update() error {
	// g.attackTimer.Update()
	// if g.attackTimer.IsReady() {
	// 	g.attackTimer.Reset()

	// 	g.attackMessage = fmt.Sprintf("Attack %d!", cnt)
	// 	cnt++
	// 	g.attackActiveTimer.Reset()
	// }
	// if len(g.attackMessage) > 0 && g.attackActiveTimer.IsReady() {
	// 	g.attackMessage = ""
	// }
	g.player.Update()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.player.Draw(screen)
	// colorm.DrawImage(screen, PlayerSprite, cm, op)
	// ebitenutil.DebugPrint(screen, g.attackMessage)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func main() {
	g := &Game{
		player: NewPlayer(),
		// attackTimer:       NewTimer(5 * time.Second),
		// attackActiveTimer: NewTimer(1 * time.Second),
	}

	err := ebiten.RunGame(g)
	if err != nil {
		panic(err)
	}
}
