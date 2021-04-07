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
	updated       = true
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

func Update() {
	updated = true
}

// Rendering functions

func PrepRender() {
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT)
}

func Render() {
	//Process job queue
	performJobs()

	// If any jobs caused an update then update
	if updated {
		updated = false

		PrepRender()
		for _, obj := range renderObjects {
			obj.Render()
		}
		window.SwapBuffers()
	}

	window.PollInput()
}
