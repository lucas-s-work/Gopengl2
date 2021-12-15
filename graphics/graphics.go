package graphics

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/lucas-s-work/gopengl2/graphics/opengl"
)

/*
All Render Objects are stored here and rendered here, on creation they are added and rendered
without additional work.
*/

var (
	window        *opengl.Window
	renderObjects []RenderObject
)

func Init(w *opengl.Window) {
	opengl.GlInit()
	window = w
}

func DeleteRenderObjects() {
	for _, ro := range renderObjects {
		ro.Delete()
	}
}

// Rendering functions

func PrepRender() {
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT)
}

func Render() {
	//Process job queue
	performJobs()

	PrepRender()
	for _, obj := range renderObjects {
		obj.Render()
	}
	window.SwapBuffers()

	window.PollInput()
}
