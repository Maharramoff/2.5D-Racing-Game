package main

import (
	sfml "github.com/manyminds/gosfml"
	"runtime"
)

func init() { runtime.LockOSThread() }

const (
	screenWidth  = 1024
	screenHeight = 768
	bitsPerPixel = 32
	title        = "Racer 2.5D"
)

func drawPolygon(app *sfml.RenderWindow, color sfml.Color, bottomX, bottomY, bottomWidth, topX, topY, topWidth float32) {
	shape, _ := sfml.NewConvexShape()
	shape.SetPointCount(4)
	shape.SetFillColor(color)
	shape.SetPoint(0, sfml.Vector2f{X: bottomX - bottomWidth, Y: bottomY})
	shape.SetPoint(1, sfml.Vector2f{X: topX - topWidth, Y: topY})
	shape.SetPoint(2, sfml.Vector2f{X: topX + topWidth, Y: topY})
	shape.SetPoint(3, sfml.Vector2f{X: bottomX + bottomWidth, Y: bottomY})
	app.Draw(shape, sfml.DefaultRenderStates())
}

func main() {
	app := sfml.NewRenderWindow(
		sfml.VideoMode{Width: screenWidth, Height: screenHeight, BitsPerPixel: bitsPerPixel},
		title,
		sfml.StyleDefault,
		sfml.DefaultContextSettings())
	app.SetFramerateLimit(60)

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

		app.Clear(sfml.ColorBlack())
		drawPolygon(app, sfml.ColorGreen(), 500, 500, 200, 500, 300, 100)
		app.Display()
	}
}
