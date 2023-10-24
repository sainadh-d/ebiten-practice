package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// https://golangprojectstructure.com/creating-cool-games-with-ebiten-in-go/

const (
	screenWidth  = 300
	screenHeight = 300
)

type Game struct{}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	red := color.RGBA{255, 0, 0, 0}
	for x := 100; x < 200; x++ {
		for y := 100; y < 200; y++ {
			if x == 100 || y == 100 || x == 199 || y == 199 {
				screen.Set(x, y, red)
			}
		}
	}
}

func main() {
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Ebiten Hello")
	g := Game{}

	if err := ebiten.RunGame(&g); err != nil {
		panic(err)
	}
}
