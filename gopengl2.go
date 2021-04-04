package main

import (
	"Gopengl2/graphics"
	"Gopengl2/graphics/opengl"
	"runtime"
)

func init() {
	runtime.LockOSThread()
}

func main() {
	opengl.GlInit()
	window := opengl.CreateWindow(800, 600, "test")
	graphics.Init(window)

	ro := graphics.CreateDefaultRenderObject("./resources/sprites/tiles.png", 24)
	ro.CreateRect(0, 0, 32, 32, 0, 0, 1, 1)
	ro.CreateRect(32, 0, 32, 32, 0, 0, 1, 1)

	ro.UpdateBuffers()
	// ro.ShouldRender = true

	for !window.ShouldClose() {
		graphics.Render()
	}
}
