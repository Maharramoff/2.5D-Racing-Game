package main

import (
	sfml "github.com/manyminds/gosfml"
	"runtime"
)

//==========================================================
// GAME VARIABLES
//==========================================================

var startPosition, currentPosition, speed = 0, 0, 0
var maxRoadLen = 1600
var camHeight, maxY float64
var currentGrassColor, currentRumbleColor, currentRoadColor, currentBrokenLineColor sfml.Color
var camX, camZ float64 = 0, 0
var pr RoadLine
var roadLines []RoadLine

var SkyColor = sfml.Color{R: 182, G: 240, B: 255, A: 255}
var RoadLightColor = sfml.Color{R: 73, G: 73, B: 73, A: 255}
var RoadDarkColor = sfml.Color{R: 70, G: 70, B: 70, A: 255}
var GrassDarkColor = sfml.Color{R: 16, G: 154, B: 16, A: 255}
var GrassLightColor = sfml.Color{R: 16, G: 170, B: 16, A: 255}
var RumbleDarkColor = sfml.Color{R: 192, G: 78, B: 73, A: 255}
var RumblightColor = sfml.Color{R: 210, G: 210, B: 210, A: 255}
var BrokenLineColor = sfml.Color{R: 210, G: 210, B: 210, A: 255}

//==========================================================
// GAME CONSTANTS
//==========================================================

const (
	ScreenWidth       = 1024
	ScreenHeight      = 768
	BitsPerPixel      = 24
	Title             = "Racer 2.5D"
	RoadWidth         = 2000
	VisibleRoadLength = 300
	SegmentLength     = 200
	CamDepth          = 0.84
	CamInitialHeight  = CamDepth * 1750
)

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
	line.scale = CamDepth / (line._3dz - camZ)
	line.x = (1 + line.scale*(line._3dx-camX)) * ScreenWidth / 2
	line.y = (1 - line.scale*(line._3dy-camY)) * ScreenHeight / 2
	line.width = line.scale * RoadWidth * ScreenWidth / 2
	return line
}

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

func init() { runtime.LockOSThread() }

func main() {

	videoMode := sfml.VideoMode{
		Width:        ScreenWidth,
		Height:       ScreenHeight,
		BitsPerPixel: BitsPerPixel,
	}

	style := sfml.StyleDefault

	contextSettings := sfml.ContextSettings{
		DepthBits:         BitsPerPixel,
		StencilBits:       0,
		AntialiasingLevel: 0,
		MajorVersion:      0,
		MinorVersion:      0,
	}

	app := sfml.NewRenderWindow(videoMode, Title, style, contextSettings)
	app.SetFramerateLimit(60)
	app.SetMouseCursorVisible(false)

	for count := 0; count < maxRoadLen; count++ {
		roadLine := NewRoadLine()
		roadLine._3dz = float64(count * SegmentLength)
		roadLines = append(roadLines, *roadLine)
	}

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

		speed = 0

		if sfml.KeyboardIsKeyPressed(sfml.KeyUp) {
			speed = 150
		}

		if sfml.KeyboardIsKeyPressed(sfml.KeyDown) {
			speed = -150
		}

		app.Clear(SkyColor)

		maxY = ScreenHeight
		var diff = 0
		currentPosition += speed
		startPosition = currentPosition / SegmentLength
		camHeight = roadLines[startPosition]._3dy + CamInitialHeight

		for count := startPosition; count < startPosition+VisibleRoadLength; count++ {

			if count >= maxRoadLen {
				diff = maxRoadLen * SegmentLength
			} else {
				diff = 0
			}

			camZ = float64(startPosition*SegmentLength - diff)

			line := handleCam(roadLines[count%maxRoadLen], camX, camHeight, camZ)

			if line.y >= maxY {
				continue
			}
			maxY = line.y

			currentGrassColor = GrassDarkColor
			currentRumbleColor = RumbleDarkColor
			currentRoadColor = RoadDarkColor
			currentBrokenLineColor = sfml.Color{}

			if (count/6)%2 == 0 {
				currentGrassColor = GrassLightColor
				currentBrokenLineColor = BrokenLineColor
			}

			if (count/3)%2 == 0 {
				currentRumbleColor = RumblightColor
				currentRoadColor = RoadLightColor
			}

			if count == 0 {
				pr = line
			} else {
				currentIdx := (count - 1) % maxRoadLen
				pr = handleCam(roadLines[currentIdx], camX, camHeight, camZ)
			}

			DrawPolygon(app, currentGrassColor, 0, int(pr.y), ScreenWidth, 0, int(line.y), ScreenWidth)
			DrawPolygon(app, currentRumbleColor, int(pr.x), int(pr.y), int(pr.width*1.2), int(line.x), int(line.y), int(line.width*1.2))
			DrawPolygon(app, currentRoadColor, int(pr.x), int(pr.y), int(pr.width), int(line.x), int(line.y), int(line.width))
			DrawPolygon(app, currentBrokenLineColor, int(pr.x), int(pr.y), int(pr.width*0.03), int(line.x), int(line.y), int(line.width*0.03))

		}

		app.Display()
	}
}
