package main

import (
	"image/color"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// https://golangprojectstructure.com/creating-cool-games-with-ebiten-in-go/

const (
	screenWidth  = 300
	screenHeight = 300
	ballRadius   = 15
)

var (
	simpleShader *ebiten.Shader

	ballPositionX = float64(screenWidth) / 2
	ballPositionY = float64(screenHeight) / 2

	ballMovementX = float64(0.00000006)
	ballMovementY = float64(0.00000004)

	lastUpdatedTime = time.Now()
)

func init() {
	var err error
	simpleShader, err = ebiten.NewShader([]byte(`
        package main

        func Fragment(position vec4, textCoord vec2, color vec4) vec4 {
            return color
        }
    `))
	if err != nil {
		panic(err)
	}
}

type Game struct{}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) Update() error {
	timeDelta := float64(time.Since(lastUpdatedTime))
	lastUpdatedTime = time.Now()

	ballPositionX += ballMovementX * timeDelta
	ballPositionY += ballMovementY * timeDelta

	const minX = ballRadius
	const minY = ballRadius
	const maxX = screenWidth - ballRadius
	const maxY = screenHeight - ballRadius

	if ballPositionX >= maxX || ballPositionX <= minX {
		if ballPositionX > maxX {
			ballPositionX = maxX
		} else if ballPositionX < minX {
			ballPositionX = minX
		}
		ballMovementX *= -1
	}

	if ballPositionY >= maxY || ballPositionY <= minY {
		if ballPositionY > maxY {
			ballPositionY = maxY
		} else if ballPositionY < minY {
			ballPositionY = minY
		}

		ballMovementY *= -1
	}
	return nil
}

// drawCircle - Draws a Circle with x, y as Center
func (g *Game) drawCircle(screen *ebiten.Image, x, y, radius int, clr color.RGBA) error {
	var path vector.Path

	path.MoveTo(float32(x), float32(y))
	path.Arc(float32(x), float32(y), float32(radius), 0, math.Pi*2, vector.Clockwise)
	vertices, indices := path.AppendVerticesAndIndicesForFilling(nil, nil)

	rScaled := float32(clr.R) / 255
	gScaled := float32(clr.G) / 255
	bScaled := float32(clr.B) / 255
	aScaled := float32(clr.A) / 255

	for i := range vertices {
		v := &vertices[i]
		v.ColorR = rScaled
		v.ColorG = gScaled
		v.ColorB = bScaled
		v.ColorA = aScaled
	}

	screen.DrawTrianglesShader(vertices, indices, simpleShader, &ebiten.DrawTrianglesShaderOptions{
		FillRule: ebiten.EvenOdd,
	})

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	red := color.RGBA{255, 0, 0, 0}
	x := int(math.Round(ballPositionX))
	y := int(math.Round(ballPositionY))
	g.drawCircle(screen, x, y, ballRadius, red)
}

func main() {
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Ebiten Hello")
	g := Game{}

	if err := ebiten.RunGame(&g); err != nil {
		panic(err)
	}
}
