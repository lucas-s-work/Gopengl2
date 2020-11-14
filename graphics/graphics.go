package graphics

import (
	"Gopengl2/graphics/opengl"

	"github.com/go-gl/gl/v4.1-core/gl"
)

var (
	window        *opengl.Window
	renderObjects []*RenderObject
)

type RenderObject struct {
	vao                        *opengl.DefaultVAO
	freeVert                   int
	vBuff                      *opengl.Buffer
	tBuff                      *opengl.Buffer
	ShouldRender               bool
	updated, autoUpdate, async bool
}

func Init(w *opengl.Window) {
	opengl.GlInit()
	window = w
}

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
		false,
		false,
	}

	renderObjects = append(renderObjects, &ro)

	return &ro
}

func (ro *RenderObject) CreateRect(x, y, width, height, texX, texY, texWidth, texHeight int) int {
	index := ro.freeVert
	ro.freeVert += 6
	ro.ModifyRect(index, x, y, width, height, texX, texY, texWidth, texHeight)

	return index
}

func (ro *RenderObject) ModifyRect(index, x, y, width, height, texX, texY, texWidth, texHeight int) {
	ro.SetVertex(index, x, y, texX, texY+texHeight)
	ro.SetVertex(index+1, x+width, y, texX+texWidth, texY+texHeight)
	ro.SetVertex(index+2, x, y+height, texX, texY)
	ro.SetVertex(index+3, x, y+height, texX, texY)
	ro.SetVertex(index+4, x+width, y+height, texX+texWidth, texY)
	ro.SetVertex(index+5, x+width, y, texX+texWidth, texY+texHeight)
}

func (ro *RenderObject) RemoveSquare(index int) {
	ro.ModifyRect(index, 0, 0, 0, 0, 0, 0, 0, 0)
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

func (ro *RenderObject) ModifySquare(index, x, y, width, texX, texY, texWidth int) {
	ro.ModifyRect(index, x, y, width, width, texX, texY, texWidth, texWidth)
}

func (ro *RenderObject) SetAutoUpdate(update bool) {
	if ro.async {
		panic("cannot auto-update async render object")
	}
	ro.autoUpdate = update
}

func (ro *RenderObject) UpdateBuffers() {
	ro.vao.BindVao()
	ro.vao.UpdateBuffers()
	ro.updated = false
}

func (ro *RenderObject) PrepRender() {
	if ro.updated && ro.autoUpdate {
		ro.vao.BindVao()
		ro.vao.UpdateBuffers()
		ro.updated = false
	}
	ro.vao.PrepRender()
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

func (ro *RenderObject) GetVAO() *opengl.DefaultVAO {
	return ro.vao
}

func DeleteRenderObjects() {
	for _, ro := range renderObjects {
		ro.Delete()
	}
}
