package graphics

import (
	"Gopengl2/graphics/opengl"
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
)

var (
	window        *opengl.Window
	renderObjects []*RenderObject
	jobs          = make(chan RenderJob)
)

type RenderObject struct {
	vao          *opengl.VAO
	freeVert     int
	vBuff        *opengl.Buffer
	tBuff        *opengl.Buffer
	ShouldRender bool
	updated      bool
}

func Init(w *opengl.Window) {
	opengl.GlInit()
	window = w
}

func PrepRender() {
	//Process job queue

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

/*
Render object functions
*/

func CreateRenderObject(texture string, elements int) *RenderObject {
	vao := opengl.CreateDefaultVao(window, texture, elements)

	ro := RenderObject{
		vao,
		0,
		vao.GetBuffer("vert"),
		vao.GetBuffer("verttexcoord"),
		true,
		false,
	}

	renderObjects = append(renderObjects, &ro)

	return &ro
}

func (ro *RenderObject) CreateRect(x, y, width, height, texX, texY, texWidth, texHeight int) int {
	index := ro.freeVert

	ro.SetVertex(index, x, y, texX, texY+texHeight)
	ro.SetVertex(index+1, x+width, y, texX+texWidth, texY+texHeight)
	ro.SetVertex(index+2, x, y+height, texX, texY)
	ro.SetVertex(index+3, x, y+height, texX, texY)
	ro.SetVertex(index+4, x+width, y+height, texX+texWidth, texY)
	ro.SetVertex(index+5, x+width, y, texX+texWidth, texY+texHeight)

	fmt.Println(ro.vao.GetBuffer("verttexcoord"))

	return index
}

func (ro *RenderObject) SetVertex(index int, x, y, texX, texY int) {
	i := index * 2

	tX, tY := ro.vao.PixToTex(texX, texY)

	ro.vBuff.Elements[i] = float32(x)
	ro.vBuff.Elements[i+1] = float32(y)
	ro.tBuff.Elements[i] = tX
	ro.tBuff.Elements[i+1] = tY

	ro.updated = true
}

func (ro *RenderObject) CreateSquare(x, y, width, texX, texY, texWidth int) int {
	return ro.CreateRect(x, y, width, width, texX, texY, texWidth, texWidth)
}

func (ro *RenderObject) PrepRender() {
	if ro.updated {
		ro.vao.BindVao()
		ro.vao.UpdateBuffers()
		ro.updated = false
	}
}

func (ro *RenderObject) Render() {
	if ro.ShouldRender {
		ro.PrepRender()
		ro.vao.Render()
	}
}

func (ro *RenderObject) Delete() {
	ro.vao.Delete()
}

func DeleteRenderObjects() {
	for _, ro := range renderObjects {
		ro.Delete()
	}
}
