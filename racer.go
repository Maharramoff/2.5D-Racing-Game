package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	screenWidth  = 600
	screenHeight = 600
)

func draw(renderer *sdl.Renderer) {
	renderer.Clear()
	renderer.SetDrawColor(255, 255, 255, 255)
	renderer.Present()
}

func eventHandler(quit *bool) {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch eType := event.(type) {
		case *sdl.QuitEvent:
			*quit = true
		case *sdl.KeyboardEvent:
			if eType.Keysym.Sym == sdl.K_ESCAPE {
				*quit = true
			}
		}
	}
}

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		fmt.Println("Init SDL: ", err)
		return
	}

	window, err := sdl.CreateWindow(
		"2.5D Racing Demo",
		sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		screenWidth, screenHeight,
		sdl.WINDOW_OPENGL)
	if err != nil {
		fmt.Println("Init Window: ", err)
		return
	}

	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Println("Init Renderer: ", err)
		return
	}

	defer renderer.Destroy()

	var quit = false

	for !quit {
		eventHandler(&quit)
		draw(renderer)
	}
}
