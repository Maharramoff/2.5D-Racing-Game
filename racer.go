package main

import (
	sfml "github.com/manyminds/gosfml"
	"math"
	"runtime"
)

const (
	SCREENWIDTH  = 1024
	SCREENHEIGHT = 768
	BITSPERPIXEL = 24
	TITLE        = "Racer 2.5D"
	ROADWIDTH    = 2000
	SEGMENTLEN   = 200
	CAMDEPTH     = 0.84
)

func Round(num float64) float32 {
	return float32(int(num))
}

type RoadLine struct {
	_3dx, _3dy, _3dz float64
	x, y, width      float64
	scale            float64
}

func NewRoadLine() *RoadLine {
	return &RoadLine{_3dx: 0, _3dy: 0, _3dz: 0, x: 0, y: 0, width: 0, scale: 0}
}

func handleCam(line RoadLine, camX, camY, camZ float64) RoadLine {
	line.scale = CAMDEPTH / (line._3dz - camZ)
	line.x = (1 + line.scale*(line._3dx-camX)) * SCREENWIDTH / 2
	line.y = (1 - line.scale*(line._3dy-camY)) * SCREENHEIGHT / 2
	line.width = line.scale * ROADWIDTH * SCREENWIDTH / 2
	return line
}

func init() { runtime.LockOSThread() }

func DrawPolygon(app *sfml.RenderWindow, color sfml.Color, bottomX, bottomY, bottomWidth, topX, topY, topWidth float64) {
	shape, _ := sfml.NewConvexShape()
	shape.SetPointCount(4)
	shape.SetFillColor(color)
	shape.SetPoint(0, sfml.Vector2f{X: Round(bottomX - bottomWidth), Y: Round(bottomY)})
	shape.SetPoint(1, sfml.Vector2f{X: Round(topX - topWidth), Y: Round(topY)})
	shape.SetPoint(2, sfml.Vector2f{X: Round(topX + topWidth), Y: Round(topY)})
	shape.SetPoint(3, sfml.Vector2f{X: Round(bottomX + bottomWidth), Y: Round(bottomY)})
	app.Draw(shape, sfml.DefaultRenderStates())
}

func main() {

	videoMode := sfml.VideoMode{
		Width:        SCREENWIDTH,
		Height:       SCREENHEIGHT,
		BitsPerPixel: BITSPERPIXEL,
	}

	style := sfml.StyleNone

	contextSettings := sfml.ContextSettings{
		DepthBits:         BITSPERPIXEL,
		StencilBits:       0,
		AntialiasingLevel: 0,
		MajorVersion:      0,
		MinorVersion:      0,
	}

	app := sfml.NewRenderWindow(videoMode, TITLE, style, contextSettings)
	app.SetFramerateLimit(60)
	app.SetMouseCursorVisible(false)

	for app.IsOpen() {
		for event := app.PollEvent(); event != nil; event = app.PollEvent() {
			switch eventType := event.(type) {
			case sfml.EventKeyReleased:
				switch eventType.Code {
				case sfml.KeyEscape:
					app.Close()
				}
			case sfml.EventClosed:
				app.Close()
			}
		}

		app.Clear(sfml.ColorWhite())

		var roadLines []RoadLine
		linee := NewRoadLine()
		for count := 0; count < 1600; count++ {
			linee._3dz = float64(count * SEGMENTLEN)
			if count > 750 {
				linee._3dy = math.Sin(float64(count/30.0)) * 1500
			}
			roadLines = append(roadLines, *linee)
		}

		var n = len(roadLines)
		var startPosition = 0
		var camHeight = roadLines[startPosition]._3dy + 1500
		var maxy float64 = SCREENHEIGHT
		var grassColor sfml.Color
		var diff = 0
		var camX, camZ float64 = 0, 0

		for count := startPosition + 1; count < startPosition+300; count++ {

			if count >= n {
				diff = n * SEGMENTLEN
			} else {
				diff = 0
			}

			camZ = float64(startPosition*SEGMENTLEN - diff)

			line := handleCam(roadLines[count%n], camX, camHeight, camZ)

			if line.y >= maxy {
				continue
			}
			maxy = line.y

			grassColor = sfml.Color{R: 16, G: 200, B: 16, A: 255}
			if (count/3)%2 != 0 {
				grassColor = sfml.Color{G: 154, A: 255}
			}

			currentIdx := (count - 1) % n
			pr := handleCam(roadLines[currentIdx], camX, camHeight, camZ)

			DrawPolygon(app, grassColor, 0, pr.y, SCREENWIDTH, 0, line.y, SCREENWIDTH)
		}

		app.Display()
	}
}
