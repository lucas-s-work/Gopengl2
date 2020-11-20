package graphics

import (
	"Gopengl2/graphics/opengl"

	"github.com/go-gl/gl/v4.1-core/gl"
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
	//Process job queue
	performJobs()

	gl.ClearColor(0.0, 0.0, 0.0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT)
}

func FinishRender() {
	window.SwapBuffers()
	window.PollInput()
}

func Render() {
	PrepRender()

	for _, obj := range renderObjects {
		obj.Render()
	}

	FinishRender()
}
