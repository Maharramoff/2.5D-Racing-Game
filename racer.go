package main

import (
	"gopkg.in/teh-cmc/go-sfml.v24/graphics"
	"gopkg.in/teh-cmc/go-sfml.v24/window"
	"runtime"
)

func init() { runtime.LockOSThread() }

func main() {
	videoMode := window.NewSfVideoMode()
	defer window.DeleteSfVideoMode(videoMode)
	videoMode.SetWidth(1024)
	videoMode.SetHeight(768)
	videoMode.SetBitsPerPixel(32)

	/* Create the main window */
	contextSettings := window.NewSfContextSettings()
	defer window.DeleteSfContextSettings(contextSettings)
	app := graphics.SfRenderWindow_create(videoMode, "SFML window", uint(window.SfResize|window.SfClose), contextSettings)
	defer window.SfWindow_destroy(app)

	ev := window.NewSfEvent()
	defer window.DeleteSfEvent(ev)

	/* Start the game loop */
	for window.SfWindow_isOpen(app) > 0 {
		/* Process events */
		for window.SfWindow_pollEvent(app, ev) > 0 {
			/* Close window: exit */
			if ev.GetXtype() == window.SfEventType(window.SfEvtClosed) {
				return
			}
		}
		graphics.SfRenderWindow_clear(app, graphics.GetSfRed())
		graphics.SfRenderWindow_display(app)
	}
}
