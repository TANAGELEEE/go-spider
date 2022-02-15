package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "mario")
}
func (g *Game) Layout(outsideWidth, outsideHeigh int) (screenWidth, screenHeigh int) {
	return 380, 160
}

func main() {
	ebiten.SetWindowSize(400, 200)
	ebiten.SetWindowTitle("mario")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
