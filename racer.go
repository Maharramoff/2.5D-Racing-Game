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

func RoundtoFloat(num int) float32 {
	return float32(num)
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

func DrawPolygon(app *sfml.RenderWindow, color sfml.Color, bottomX, bottomY, bottomWidth, topX, topY, topWidth int) {
	shape, _ := sfml.NewConvexShape()
	shape.SetPointCount(4)
	shape.SetFillColor(color)
	shape.SetPoint(0, sfml.Vector2f{X: RoundtoFloat(bottomX - bottomWidth), Y: RoundtoFloat(bottomY)})
	shape.SetPoint(1, sfml.Vector2f{X: RoundtoFloat(topX - topWidth), Y: RoundtoFloat(topY)})
	shape.SetPoint(2, sfml.Vector2f{X: RoundtoFloat(topX + topWidth), Y: RoundtoFloat(topY)})
	shape.SetPoint(3, sfml.Vector2f{X: RoundtoFloat(bottomX + bottomWidth), Y: RoundtoFloat(bottomY)})
	app.Draw(shape, sfml.DefaultRenderStates())
}

func main() {

	videoMode := sfml.VideoMode{
		Width:        SCREENWIDTH,
		Height:       SCREENHEIGHT,
		BitsPerPixel: BITSPERPIXEL,
	}

	style := sfml.StyleDefault

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

	var roadLines []RoadLine

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

		for count := 0; count < 1200; count++ {
			linee := NewRoadLine()
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
		var rumbleColor sfml.Color
		var roadColor sfml.Color
		var diff = 0
		var camX, camZ float64 = 0, 0
		var pr RoadLine

		for count := startPosition; count < startPosition+300; count++ {

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

			grassColor = sfml.Color{G: 154, A: 255}
			if (count/2)%2 == 0 {
				grassColor = sfml.Color{R: 16, G: 200, B: 16, A: 255}
			}

			rumbleColor = sfml.Color{R: 226, G: 53, B: 0, A: 255}
			if (count/2)%2 == 0 {
				rumbleColor = sfml.Color{R: 255, G: 255, B: 255, A: 255}
			}

			roadColor = sfml.Color{R: 91, G: 91, B: 91, A: 255}

			if count == 0 {
				pr = line
			} else {
				currentIdx := (count - 1) % n
				pr = handleCam(roadLines[currentIdx], camX, camHeight, camZ)
			}

			//fmt.Printf("%v ", int(line.width))
			DrawPolygon(app, grassColor, 0, int(pr.y), SCREENWIDTH, 0, int(line.y), SCREENWIDTH)
			DrawPolygon(app, rumbleColor, int(pr.x), int(pr.y), int(pr.width*1.2), int(line.x), int(line.y), int(line.width*1.2))
			DrawPolygon(app, roadColor, int(pr.x), int(pr.y), int(pr.width), int(line.x), int(line.y), int(line.width))

		}

		app.Display()
	}
}
