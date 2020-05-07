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

func drawRoad(app *sfml.RenderWindow, color sfml.Color, x1, y1, w1, x2, y2, w2 float32) {
	shape, _ := sfml.NewConvexShape()
	shape.SetPointCount(4)
	shape.SetFillColor(color)
	shape.SetPoint(0, sfml.Vector2f{X: x1 - w1, Y: y1})
	shape.SetPoint(1, sfml.Vector2f{X: x2 - w2, Y: y2})
	shape.SetPoint(2, sfml.Vector2f{X: x2 + w2, Y: y2})
	shape.SetPoint(3, sfml.Vector2f{X: x1 + w1, Y: y1})
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
		drawRoad(app, sfml.ColorGreen(), 500, 500, 200, 500, 300, 100)
		app.Display()
	}
}
