package main

import (
	"fmt"
	sfml "github.com/manyminds/gosfml"
	"math"
	"runtime"
	"time"
)

//==========================================================
// GAME VARIABLES
//==========================================================

var startPosition = 0
var currentPosition = 0
var currentLap = 1
var camHeight float32
var maxY float32 = 0
var speed float32 = 0.0
var maxSpeed float32 = 300.0
var accelerationRate = maxSpeed / 5
var speedPercent float32 = 0
var deltaTime float32
var camX float32 = 0
var camDx float32 = 0.1
var camZ float32 = 0
var playerZ float32 = 0
var camDepth float32 = 0
var pr, playerSegment RoadLine
var roadMap []RoadLine
var playerRiding = false
var carPos, carScale = sfml.Vector2f{X: 0, Y: 0}, sfml.Vector2f{X: 4.0, Y: 4.0}
var carDim = sfml.IntRect{Left: 136, Top: 89, Width: 52, Height: 31}
var carRideDim = sfml.IntRect{Left: 8, Top: 9, Width: 52, Height: 31}
var carDimLeft = sfml.IntRect{Left: 72, Top: 9, Width: 57, Height: 31}
var carDimRight = sfml.IntRect{Left: 879, Top: 9, Width: 57, Height: 31}

var currentGrassColor, currentRumbleColor, currentRoadColor, currentBrokenLineColor sfml.Color
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
	ScreenWidth            = 1024
	ScreenHeight           = 768
	BitsPerPixel           = 24
	Title                  = "Racer 2.5D"
	RoadWidth              = 2000
	MaxRoadLen             = 1600
	VisibleRoadLength      = 300
	SegmentLength          = 200
	CamInitialHeight       = 1000
	MaxLaps                = 3
	PlayerCentrifugalForce = 0.3
)

func RoundtoFloat(num int) float32 {
	return float32(num)
}

func RoundtoDec(num float64, dec int) float64 {
	return math.Round(num*math.Pow10(dec)) / math.Pow10(dec)
}

type RoadLine struct {
	_3dx, _3dy, _3dz float32
	camX, camY, camZ float32
	x, y, width      float32
	scale, curve     float32
}

func NewRoadLine() *RoadLine {
	return &RoadLine{_3dx: 0, _3dy: 0, _3dz: 0, camX: 0, camY: 0, camZ: 0, x: 0, y: 0, width: 0, scale: 0, curve: 0}
}

func handleCam(line RoadLine, camX, camY, camZ, camDepth float32) RoadLine {
	line.camX = line._3dx - camX
	line.camY = line._3dy - camY
	line.camZ = line._3dz - camZ
	line.scale = camDepth / line.camZ
	line.x = (1 + line.scale*line.camX) * ScreenWidth / 2
	line.y = (1 - line.scale*line.camY) * ScreenHeight / 2
	line.width = line.scale * RoadWidth * ScreenWidth / 2
	return line
}

func generateRoadMap(maxLength int) []RoadLine {
	for count := 0; count < maxLength; count++ {
		roadLine := NewRoadLine()
		roadLine._3dz = float32(count * SegmentLength)

		if count > 350 && count < 550 {
			roadLine.curve = 0.3
		}
		if count > 600 && count < 950 {
			roadLine.curve = -0.5
		}

		roadMap = append(roadMap, *roadLine)
	}

	return roadMap
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

func DrawStats(app *sfml.RenderWindow, txt string, x, y float32) {
	font, _ := sfml.NewFontFromFile("assets/fonts/faster_one/regular.ttf")
	statsText, _ := sfml.NewText(font)
	statsText.SetCharacterSize(20)
	statsText.SetPosition(sfml.Vector2f{X: x, Y: y})
	statsText.SetColor(sfml.ColorBlack())
	statsText.SetString(txt)
	app.Draw(statsText, sfml.DefaultRenderStates())
}

func init() {
	runtime.LockOSThread()
}

func main() {

	runtime.UnlockOSThread()
	musicBuffer, err := sfml.NewSoundBufferFromFile("assets/music/boxcat_games_-_tricks.ogg")
	music := sfml.NewSound(musicBuffer)
	if err != nil {
		panic(err)
	}

	music.SetLoop(true)
	music.Play()
	runtime.LockOSThread()

	videoMode := sfml.VideoMode{
		Width:        ScreenWidth,
		Height:       ScreenHeight,
		BitsPerPixel: BitsPerPixel,
	}

	style := sfml.StyleDefault

	contextSettings := sfml.ContextSettings{
		DepthBits:         BitsPerPixel,
		StencilBits:       8,
		AntialiasingLevel: 0,
		MajorVersion:      0,
		MinorVersion:      0,
	}

	app := sfml.NewRenderWindow(videoMode, Title, style, contextSettings)
	app.SetMouseCursorVisible(false)
	app.SetVSyncEnabled(true)
	app.SetActive(false)
	icon, _ := sfml.NewImageFromFile("assets/images/game_icon.png")

	err = app.SetIcon(128, 128, icon.GetPixelData())
	if err != nil {
		panic(err)
	}

	texture, err := sfml.NewTextureFromFile("assets/images/spritesheet.png", nil)
	if err != nil {
		panic(err)
	}

	// Car sprite
	carSprite, _ := sfml.NewSprite(texture)
	if err != nil {
		panic(err)
	}

	carSprite.SetScale(carScale)
	carSprite.SetTextureRect(carDim)
	carPos.X = ScreenWidth/2 - carSprite.GetGlobalBounds().Width + carSprite.GetGlobalBounds().Width/2
	carPos.Y = ScreenHeight - carSprite.GetGlobalBounds().Height - 10

	roadMap = generateRoadMap(MaxRoadLen)

	lastTime := time.Now()
	speed = 0.0
	breaking := -maxSpeed

	for app.IsOpen() {
		app.SetActive(true)
		for event := app.PollEvent(); event != nil; event = app.PollEvent() {
			switch eventType := event.(type) {
			case sfml.EventKeyReleased:
				switch eventType.Code {
				case sfml.KeyEscape:
					app.Close()
				case sfml.KeyLeft, sfml.KeyRight, sfml.KeyUp, sfml.KeyDown:
					carSprite.SetTextureRect(carDim)
				}
			case sfml.EventClosed:
				app.Close()
			}
		}

		now := time.Now()
		deltaTime = float32(time.Since(lastTime).Seconds())
		music.GetStatus()
		playerRiding = false

		if app.HasFocus() {

			if sfml.KeyboardIsKeyPressed(sfml.KeyUp) {
				speed += accelerationRate*deltaTime - (speed * .003)
				carSprite.SetTextureRect(carRideDim)
				playerRiding = true
			}

			if sfml.KeyboardIsKeyPressed(sfml.KeyDown) {
				speed += breaking * deltaTime
				carSprite.SetTextureRect(carRideDim)
				playerRiding = true
			}
			speedPercent = speed / maxSpeed
			camDx = deltaTime * speedPercent * 2.0

			if sfml.KeyboardIsKeyPressed(sfml.KeyRight) {
				camX += camDx
				carSprite.SetTextureRect(carDimRight)
			}

			if sfml.KeyboardIsKeyPressed(sfml.KeyLeft) {
				camX -= camDx
				carSprite.SetTextureRect(carDimLeft)
			}
		}

		// Speed limit
		speed = float32(math.Max(0, math.Min(float64(speed), float64(maxSpeed))))

		maxY = ScreenHeight
		var diff, curveX, curveDx float32 = 0, 0.0, 0.0
		currentPosition += int(speed)

		for currentPosition >= MaxRoadLen*SegmentLength {
			currentPosition -= MaxRoadLen * SegmentLength
			currentLap += 1
		}

		for currentPosition < 0 {
			currentPosition = startPosition
		}

		camDepth = float32(1 / math.Tan((50)*math.Pi/180))
		startPosition = currentPosition / SegmentLength
		camHeight = roadMap[startPosition]._3dy + CamInitialHeight
		playerZ = camHeight * camDepth

		if playerRiding {
			playerSegment = roadMap[(startPosition+int(playerZ/SegmentLength))%len(roadMap)]
			if speed > 0 {
				camX -= camDx * speedPercent * playerSegment.curve * PlayerCentrifugalForce
			} else {
				camX -= -camDx * speedPercent * playerSegment.curve * PlayerCentrifugalForce
			}
		}

		app.Clear(SkyColor)
		for count := startPosition; count < startPosition+VisibleRoadLength; count++ {

			if count >= MaxRoadLen {
				diff = MaxRoadLen*SegmentLength - SegmentLength
			} else {
				diff = 0
			}

			camZ = float32(startPosition*SegmentLength) - diff

			line := handleCam(roadMap[count%MaxRoadLen], camX*RoadWidth-curveX, camHeight, camZ, camDepth)

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
				currentIdx := (count - 1) % MaxRoadLen
				pr = handleCam(roadMap[currentIdx], camX*RoadWidth-curveX, camHeight, camZ, camDepth)
			}

			DrawPolygon(app, currentGrassColor, 0, int(pr.y), ScreenWidth, 0, int(line.y), ScreenWidth)
			DrawPolygon(app, currentRumbleColor, int(pr.x), int(pr.y), int(pr.width*1.2), int(line.x), int(line.y), int(line.width*1.2))
			DrawPolygon(app, currentRoadColor, int(pr.x), int(pr.y), int(pr.width), int(line.x), int(line.y), int(line.width))
			DrawPolygon(app, currentBrokenLineColor, int(pr.x), int(pr.y), int(pr.width*0.03), int(line.x), int(line.y), int(line.width*0.03))

			curveX += curveDx
			curveDx += line.curve
		}
		lastTime = now
		carSprite.SetPosition(carPos)
		app.Draw(carSprite, sfml.DefaultRenderStates())
		DrawStats(
			app,
			fmt.Sprintf(
				"LAP: %d/%d\nFPS: %v\nSPEED:\t%v\nDELTA: %v",
				currentLap,
				MaxLaps,
				RoundtoDec(float64(1/deltaTime), 1),
				RoundtoDec(float64(speed), 1),
				RoundtoDec(float64(deltaTime*1000), 1),
			),
			20,
			20,
		)
		app.Display()
	}
}
