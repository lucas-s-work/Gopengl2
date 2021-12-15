package main

import (
	"runtime"
	"time"

	"github.com/lucas-s-work/gopengl2/graphics"
	"github.com/lucas-s-work/gopengl2/graphics/opengl"
)

func init() {
	runtime.LockOSThread()
}

func main() {
	window := opengl.CreateWindow(1920, 1200, "test")
	opengl.GlInit()

	graphics.Init(window)

	ro := graphics.CreateDefaultRenderObject("./resources/sprites/tiles.png", 1000*10000*2)

	for x := 0; x < 1000; x++ {
		for y := 0; y < 10000; y++ {
			ro.CreateSquare(x*33, y*33, 32, 0, 0, 1)
		}
	}

	ro.UpdateBuffers()
	ro.ShouldRender = true
	var x, y float32
	ro.SetTranslation(&x, &y)
	ch := ro.SetWait(true)

	go tick(ch)

	ticker := time.NewTicker(1000 / 60 * time.Millisecond)
	for !window.ShouldClose() {
		<-ticker.C
		x--
		ro.UpdatePointers()
		graphics.Render()
	}
}

func tick(ch chan graphics.WaitSignal) {
	ticker := time.NewTicker(1 * time.Second)
	for {
		<-ticker.C
		ch <- graphics.WaitSignal{}
	}
}
