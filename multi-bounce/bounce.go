package main

import (
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// https://golangprojectstructure.com/creating-cool-games-with-ebiten-in-go/

const (
	screenWidth  = 300
	screenHeight = 300
	ballCount    = 2
	ballRadius   = 15
)

type Ball struct {
	x, y   int // Center
	radius float64

	positionX float64
	positionY float64

	xSpeed float64
	ySpeed float64
}

var (
	simpleShader *ebiten.Shader

	ballPositionX = float64(screenWidth) / 2
	ballPositionY = float64(screenHeight) / 2

	ballMovementX = float64(0.00000006)
	ballMovementY = float64(0.00000004)

	lastUpdatedTime = time.Now()

	balls []*Ball
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

	// Initialize Balls
	for i := 0; i < ballCount; i++ {
		x, y := 50+rand.Intn(100), 50+rand.Intn(100)
		balls = append(balls, &Ball{
			x:         x,
			y:         y,
			positionX: float64(x),
			positionY: float64(y),
			radius:    ballRadius,
			xSpeed:    ballMovementX,
			ySpeed:    ballMovementY,
		})
	}
}

type Game struct{}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) Update() error {
	timeDelta := float64(time.Since(lastUpdatedTime))
	lastUpdatedTime = time.Now()

	for i := range balls {
		ball := balls[i]

		ball.positionX += ball.xSpeed * timeDelta
		ball.positionY += ball.ySpeed * timeDelta

		minX := ball.radius
		minY := ball.radius
		maxX := screenWidth - ball.radius
		maxY := screenHeight - ball.radius

		if ball.positionX >= maxX || ball.positionX <= minX {
			if ball.positionX > maxX {
				ball.positionX = maxX
			} else if ball.positionX < minX {
				ball.positionX = minX
			}
			ball.xSpeed *= -1
		}

		if ball.positionY >= maxY || ball.positionY <= minY {
			if ball.positionY > maxY {
				ball.positionY = maxY
			} else if ball.positionY < minY {
				ball.positionY = minY
			}

			ball.ySpeed *= -1
		}
	}
	return nil
}

// drawCircle - Draws a Circle with x, y as Center
func (g *Game) drawBalls(screen *ebiten.Image, balls []*Ball, clr color.RGBA) error {
	var path vector.Path
	for i := range balls {
		ball := balls[i]
		x := ball.positionX
		y := ball.positionY
		radius := ball.radius

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

	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	red := color.RGBA{255, 0, 0, 0}
	g.drawBalls(screen, balls, red)
}

func main() {
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Ebiten Hello")
	g := Game{}

	if err := ebiten.RunGame(&g); err != nil {
		panic(err)
	}
}
